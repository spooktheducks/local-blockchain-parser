package cmds

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type csvLine struct {
	blockHash  string
	txHash     string
	scriptData []byte
}

func PrintOpReturns(startBlock, endBlock uint64, inDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "op-returns")

	err := os.MkdirAll(outSubdir, 0777)
	if err != nil {
		return err
	}

	// start a goroutine to log errors
	chErr := make(chan error)
	go func() {
		for err := range chErr {
			fmt.Println("error:", err)
		}
	}()

	// start a goroutine to write lines to the CSV file
	chCSVData := make(chan csvLine)
	chCSVDone := make(chan bool)
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
	csvFile, err := utils.CreateFile(csvFilepath)
	if err != nil {
		chErr <- err
		return
	}
	defer utils.CloseFile(csvFile)

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

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(inDir, filename))
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

			txOuts := tx.MsgTx().TxOut

			for _, txout := range txOuts {
				data, err := utils.GetNonOPBytes(txout.PkScript)
				if err != nil {
					if err.Error() == "encoding/hex: odd length hex string" {
						continue
					} else if err.Error() == "execute past end of script" {
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

			if len(allTxOutData) == 0 {
				continue
			}

			matches := utils.SearchDataForMagicFileBytes(allTxOutData)
			for _, match := range matches {
				fmt.Printf("- file magic byte match -> type: %v (block hash: %v) (tx hash: %v)\n", match.Description(), blockHash, txHash)
			}

			length := binary.LittleEndian.Uint32(allTxOutData[0:4])
			expectedChecksum := binary.LittleEndian.Uint32(allTxOutData[4:8])
			if len(allTxOutData) < 8+int(length) {
				continue
			}

			data := allTxOutData[8 : 8+int(length)]

			checksum := crc32.ChecksumIEEE(data)

			if checksum == expectedChecksum {
				fmt.Println("EXPECTED CHECKSUM MATCHED")
			} else {
				continue
			}

			allTxOutFilename := filepath.Join(blockDir, fmt.Sprintf("txouts-combined-%v.dat", txHash))
			err = utils.CreateAndWriteFile(allTxOutFilename, data)
			if err != nil {
				chErr <- err
				return
			}
		}
	}
}
