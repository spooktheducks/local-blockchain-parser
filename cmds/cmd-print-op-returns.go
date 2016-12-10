package cmds

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcd/txscript"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/utils"
)

type csvLine struct {
	blockHash  string
	txHash     string
	scriptData []byte
}

var maxFiles = 256
var fileSemaphore = make(chan bool, maxFiles)

func PrintBlockScriptsOpReturns(startBlock, endBlock uint64, inDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "op-returns")

	err := os.MkdirAll(outSubdir, 0777)
	if err != nil {
		return err
	}

	chErr := make(chan error)
	go func() {
		for err := range chErr {
			fmt.Println("error:", err)
		}
	}()

	// fill up our file semaphore so we can obtain tokens from it
	for i := 0; i < maxFiles; i++ {
		fileSemaphore <- true
	}

	chCSVData := make(chan csvLine)
	chCSVDone := make(chan bool)

	// start a goroutine to write lines to the CSV file
	go writeCSV(outSubdir, chCSVData, chCSVDone, chErr)

	// start a goroutine for each .dat file being parsed
	chDones := []chan bool{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		chDone := make(chan bool)
		go opReturnsParseBlock(inDir, outSubdir, i, chCSVData, chErr, chDone)
		chDones = append(chDones, chDone)

	}

	// wait for all ops to complete
	for _, chDone := range chDones {
		<-chDone
	}

	fmt.Println(".dat files written to", outSubdir)

	// close CSV writer channel
	close(chCSVData)

	// wait for CSV writing to finish
	<-chCSVDone

	// close error channel
	close(chErr)

	return nil
}

func writeCSV(outSubdir string, chCSVData chan csvLine, chCSVDone chan bool, chErr chan error) {
	defer close(chCSVDone)

	csvFilepath := filepath.Join(outSubdir, "all-blocks.csv")
	csvFile, err := createFile(csvFilepath)
	if err != nil {
		chErr <- err
		return
	}
	defer closeFile(csvFile)

	_, err = csvFile.WriteString(fmt.Sprintf("blockHash,txHash,scriptData\n"))
	if err != nil {
		chErr <- err
		return
	}

	for line := range chCSVData {
		_, err := csvFile.WriteString(fmt.Sprintf("%s,%s,%s\n", line.blockHash, line.txHash, string(line.scriptData)))
		if err != nil {
			chErr <- err
			return
		}
	}

	fmt.Println(csvFilepath, "written.")
}

func opReturnsParseBlock(inDir string, outDir string, blockFileNum int, chCSVData chan csvLine, chErr chan error, chDone chan bool) {
	defer close(chDone)

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	<-fileSemaphore
	blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
	fileSemaphore <- true
	if err != nil {
		chErr <- err
		return
	}

	for _, bl := range blocks {
		blockHash := bl.Hash().String()

		blockDir := filepath.Join(outDir, blockHash)

		err = os.MkdirAll(blockDir, 0777)
		if err != nil {
			chErr <- err
			return
		}

		for _, tx := range bl.Transactions() {
			txHash := tx.Hash().String()

			allTxOutData := []byte{}

			for _, txout := range tx.MsgTx().TxOut {
				scriptStr, err := txscript.DisasmString(txout.PkScript)
				if err != nil {
					if err.Error() == "execute past end of script" {
						continue
					} else {
						chErr <- fmt.Errorf("error in txscript.DisasmString: %v", err)
						return
					}
				}

				data, err := getNonOPBytes(scriptStr)
				if err != nil {
					if err.Error() == "encoding/hex: odd length hex string" {
						continue
					} else {
						chErr <- fmt.Errorf("error in getNonOPBytes: %v", err)
						return
					}
				} else if data == nil {
					continue
				}

				allTxOutData = append(allTxOutData, data...)

				chCSVData <- csvLine{blockHash, txHash, data}
			}

			headerMatches, footerMatches := searchDataForKnownFileBits(allTxOutData)
			if len(headerMatches) > 0 {
				for _, match := range headerMatches {
					fmt.Printf("- file header match (type: %v) (block hash: %v) (tx hash: %v)\n", match.filetype, blockHash, txHash)
				}
			}
			if len(footerMatches) > 0 {
				for _, match := range footerMatches {
					fmt.Printf("- file footer match (type: %v) (block hash: %v) (tx hash: %v)\n", match.filetype, blockHash, txHash)
				}
			}

			allTxOutFilename := filepath.Join(blockDir, fmt.Sprintf("txouts-combined-%v.dat", txHash))
			err = createAndWriteFile(allTxOutFilename, allTxOutData)
			if err != nil {
				chErr <- err
				return
			}
		}
	}
}

func createAndWriteFile(path string, bytes []byte) error {
	f, err := createFile(path)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		closeFile(f)
		return err
	}
	return closeFile(f)
}

func createFile(path string) (*os.File, error) {
	<-fileSemaphore
	f, err := os.Create(path)
	if err != nil {
		fileSemaphore <- true
		return nil, err
	}
	return f, nil
}

func closeFile(file *os.File) error {
	err := file.Close()
	fileSemaphore <- true
	return err
}

type (
	fileHeaderDefinition struct {
		filetype   string
		headerData []byte
	}
)

var fileHeaders = []fileHeaderDefinition{
	{"doc", []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	{"xls", []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	{"ppt", []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	{"zip", []byte{0x50, 0x4B, 0x03, 0x04, 0x14}}, // probably wrong
	{"jpg", []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01}},
	{"gif", []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x4E, 0x01, 0x53, 0x00, 0xC4}},
	{"pdf", []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E}}, // probably wrong
	{"7zip", []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}},      // verified at endchan
	{"Torrent", []byte{0x64, 0x38, 0x3A, 0x61, 0x6E, 0x6E, 0x6F, 0x75, 0x6E, 0x63, 0x65}},
	{"AVI", []byte{0x52, 0x49, 0x46, 0x46}},
}

var fileFooters = []fileHeaderDefinition{
	{"doc", []byte{0x57, 0x6F, 0x72, 0x64, 0x2E, 0x44, 0x6F, 0x63, 0x75, 0x6D, 0x65, 0x6E, 0x74, 0x2E}},
	{"xls", []byte{0xFE, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x57, 0x00, 0x6F, 0x00, 0x72, 0x00, 0x6B, 0x00, 0x62, 0x00, 0x6F, 0x00, 0x6F, 0x00, 0x6B, 0x00}},
	{"ppt", []byte{0x50, 0x00, 0x6F, 0x00, 0x77, 0x00, 0x65, 0x00, 0x72, 0x00, 0x50, 0x00, 0x6F, 0x00, 0x69, 0x00, 0x6E, 0x00, 0x74, 0x00, 0x20, 0x00, 0x44, 0x00, 0x6F, 0x00, 0x63, 0x00, 0x75, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x6E, 0x00, 0x74}},
	{"zip", []byte{0x50, 0x4B, 0x05, 0x06, 0x00}}, // probably wrong
	{"gif", []byte{0x21, 0x00, 0x00, 0x3B, 0x00}},
	{"pdf", []byte{0x25, 0x25, 0x45, 0x4F, 0x46}},  // probably wrong
	{"7zip", []byte{0x00, 0x00, 0x00, 0x17, 0x06}}, // verified with cablegate
}

func searchDataForKnownFileBits(data []byte) ([]fileHeaderDefinition, []fileHeaderDefinition) {
	if data == nil {
		return []fileHeaderDefinition{}, []fileHeaderDefinition{}
	}

	headerMatches := []fileHeaderDefinition{}
	for _, header := range fileHeaders {
		if bytes.Contains(data, header.headerData) {
			headerMatches = append(headerMatches, header)
		}
	}

	footerMatches := []fileHeaderDefinition{}
	for _, header := range fileHeaders {
		if bytes.Contains(data, header.headerData) {
			footerMatches = append(footerMatches, header)
		}
	}

	return headerMatches, footerMatches
}

func getNonOPBytes(scriptStr string) ([]byte, error) {
	toks := strings.Split(scriptStr, " ")

	// for i := range toks {
	// 	if toks[i] == "OP_RETURN" {
	// 		if len(toks) >= i+2 {
	// 			return hex.DecodeString(toks[i+1])
	// 		} else {
	// 			return nil, errors.New("empty OP_RETURN data")
	// 		}
	// 	}
	// }

	bs := []byte{}
	for _, tok := range toks {
		if len(tok) <= 3 {
			continue
		}

		if tok[:3] != "OP_" && len(tok) >= 40 {
			decoded, err := hex.DecodeString(tok)
			if err != nil {
				return nil, err
			}
			bs = append(bs, decoded...)
		}
	}

	return bs, nil
}
