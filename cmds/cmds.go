package cmds

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcd/txscript"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/utils"
)

func PrintBlockScripts(startBlock, endBlock uint64, inDir, outDir string) error {
	dir := filepath.Join(".", outDir, "scripts")

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	for i := int(startBlock); i < int(endBlock)+1; i++ {
		filename := fmt.Sprintf("blk%05d.dat", i)

		blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
		if err != nil {
			panic(err)
		}

		for _, bl := range blocks {
			blockHash := bl.Hash().String()

			f, err := os.Create(filepath.Join(dir, blockHash+".txt"))
			if err != nil {
				return err
			}
			defer f.Close()

			fmt.Println("===== BLOCK " + blockHash + " =====")
			f.WriteString("[BLOCK " + blockHash + "]\n")

			for _, tx := range bl.Transactions() {
				txHash := tx.Hash().String()
				fmt.Println("-   TX " + txHash)
				_, err := f.WriteString("TX: " + txHash + "\n")
				if err != nil {
					return err
				}

				for _, txout := range tx.MsgTx().TxOut {
					scriptStr, err := txscript.DisasmString(txout.PkScript)
					if err != nil {
						if err.Error() == "execute past end of script" {
							continue
						} else {
							return fmt.Errorf("error in txscript.DisasmString: %v", err)
						}
					}

					fmt.Println("        " + scriptStr)
					_, err = f.WriteString("  - " + scriptStr + "\n")
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func PrintBlockScriptsOpReturns(startBlock, endBlock uint64, inDir, outDir string) error {
	dir := filepath.Join(".", outDir, "op-returns")

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	allData, err := os.Create(filepath.Join(dir, "all-blocks.csv"))
	if err != nil {
		return err
	}
	defer allData.Close()

	_, err = allData.WriteString(fmt.Sprintf("blockHash,txHash,scriptData\n"))
	if err != nil {
		return err
	}

	for i := int(startBlock); i < int(endBlock)+1; i++ {
		filename := fmt.Sprintf("blk%05d.dat", i)

		blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
		if err != nil {
			panic(err)
		}

		for _, bl := range blocks {
			blockHash := bl.Hash().String()

			blockDir := filepath.Join(dir, blockHash)

			err = os.MkdirAll(blockDir, 0777)
			if err != nil {
				return err
			}

			for _, tx := range bl.Transactions() {
				// for _, txin := range tx.TxIns {
				//  data, err := blkparser.Pkscript(txin.ScriptSig).DecodeToString()
				//  if err != nil {
				//      return err
				//  }
				//  fmt.Println(data)
				// }

				txHash := tx.Hash().String()

				for txoutIdx, txout := range tx.MsgTx().TxOut {
					scriptStr, err := txscript.DisasmString(txout.PkScript)
					if err != nil {
						if err.Error() == "execute past end of script" {
							continue
						} else {
							return fmt.Errorf("error in txscript.DisasmString: %v", err)
						}
					}

					data, err := getOpReturnBytes(scriptStr)
					if err != nil {
						if err.Error() == "encoding/hex: odd length hex string" {
							continue
						} else {
							return fmt.Errorf("error in getOpReturnBytes: %v", err)
						}
					} else if data == nil {
						continue
					}

					fileHeaderMatches := searchDataForFileHeaders(data)
					if len(fileHeaderMatches) > 0 {
						for _, match := range fileHeaderMatches {
							fmt.Printf("- file header match (type: %v) (block hash: %v) (tx hash: %v)\n", match.filetype, blockHash, txHash)
						}
					}

					f, err := os.Create(filepath.Join(blockDir, fmt.Sprintf("%v-%v.dat", txHash, txoutIdx)))
					if err != nil {
						return err
					}
					// defer f.Close()

					_, err = f.Write(data)
					if err != nil {
						f.Close()
						return err
					}
					f.Close()

					_, err = allData.WriteString(fmt.Sprintf("%s,%s,%s\n", blockHash, txHash, string(data)))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
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
	{"zip", []byte{0x50, 0x4B, 0x03, 0x04, 0x14}},
	{"jpg", []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01}},
	{"gif", []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x4E, 0x01, 0x53, 0x00, 0xC4}},
	{"pdf", []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E}},
	{"7zip", []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}},
	{"Torrent", []byte{0x64, 0x38, 0x3A, 0x61, 0x6E, 0x6E, 0x6F, 0x75, 0x6E, 0x63, 0x65}},
	{"AVI", []byte{0x52, 0x49, 0x46, 0x46}},
}

func searchDataForFileHeaders(data []byte) []fileHeaderDefinition {
	if data == nil {
		return []fileHeaderDefinition{}
	}

	matches := []fileHeaderDefinition{}
	for _, header := range fileHeaders {
		if bytes.Contains(data, header.headerData) {
			fmt.Println("possible match (" + header.filetype + ")")
			matches = append(matches, header)
		}
	}

	return matches
}

func getOpReturnBytes(scriptStr string) ([]byte, error) {
	toks := strings.Split(scriptStr, " ")

	for i := range toks {
		if toks[i] == "OP_RETURN" {
			if len(toks) >= i+2 {
				return hex.DecodeString(toks[i+1])
			} else {
				return nil, errors.New("empty OP_RETURN data")
			}
		}
	}

	return nil, nil
}

func PrintBlockData(startBlock, endBlock uint64, inDir, outDir string) error {
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		filename := fmt.Sprintf("blk%05d.dat", i)

		blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
		if err != nil {
			panic(err)
		}

		for _, bl := range blocks {

			// Basic block info
			fmt.Printf("Block hash: %v\n", bl.Hash().String())
			fmt.Printf("Block time: %v\n", bl.MsgBlock().Header.Timestamp)
			fmt.Printf("Block version: %v\n", bl.MsgBlock().Header.Version)
			// fmt.Printf("Block parent: %v\n", btc.NewUint256(bl.ParentHash()).String())
			fmt.Printf("Block merkle root: %v\n", bl.MsgBlock().Header.MerkleRoot.String())
			fmt.Printf("Block bits: %v\n", bl.MsgBlock().Header.Bits)
			// fmt.Printf("Block size: %v\n", len(bl.Raw))

			// Fetch TXs and iterate over them
			for _, tx := range bl.Transactions() {
				fmt.Printf("TxId: %v\n", tx.Hash().String())
				// fmt.Printf("Tx Size: %v\n", tx.Size)
				// fmt.Printf("Tx Lock time: %v\n", tx.LockTime)
				// fmt.Printf("Tx Version: %v\n", tx.Version)

				fmt.Println("TxIns:")

				// if tx.IsCoinBase() {
				//  fmt.Printf("TxIn coinbase, newly generated coins")
				// } else {
				for txin_index := range tx.MsgTx().TxIn {
					fmt.Printf("TxIn index: %v\n", txin_index)
					// fmt.Printf("TxIn Input hash: %v\n", txin.InputHash)
					// fmt.Printf("TxIn Input vout: %v\n", txin.InputVout)
					// fmt.Printf("TxIn ScriptSig: %v\n", hex.EncodeToString(txin.ScriptSig))
					// fmt.Printf("TxIn Sequence: %v\n", txin.Sequence)
				}
				// }

				fmt.Println("TxOuts:")

				for txo_index := range tx.MsgTx().TxOut {
					fmt.Printf("TxOut index: %v\n", txo_index)
					// txout.PkScript
					// fmt.Printf("TxOut value: %v\n", txout.Value)
					// fmt.Printf("TxOut script: %s\n", hex.EncodeToString(txout.Pkscript))
					// fmt.Printf("TxOut address: %v\n", txout.Addr)
				}
			}
		}
	}

	return nil
}
