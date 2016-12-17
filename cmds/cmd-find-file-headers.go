package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

func FindFileHeaders(startBlock, endBlock uint64, inDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "file-headers")

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

	// start a goroutine for each .dat file being parsed
	chDones := []chan bool{}
	procLimiter := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		procLimiter <- true
	}

	for i := int(startBlock); i < int(endBlock)+1; i++ {
		chDone := make(chan bool)
		go findFileHeadersParseBlock(inDir, outSubdir, i, chErr, chDone, procLimiter)
		chDones = append(chDones, chDone)
	}

	// wait for all ops to complete
	for _, chDone := range chDones {
		<-chDone
	}

	// close error channel
	close(chErr)

	return nil
}

func findFileHeadersParseBlock(inDir string, outDir string, blockFileNum int, chErr chan error, chDone chan bool, procLimiter chan bool) {
	defer close(chDone)
	defer func() { procLimiter <- true }()
	<-procLimiter

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(inDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	outFile := utils.NewConditionalFile(filepath.Join(outDir, fmt.Sprintf("blk%05d-file-headers.txt", blockFileNum)))
	defer outFile.Close()

	// write CSV header
	_, err = outFile.WriteString("block,tx,input or output,index,description\n", false)
	if err != nil {
		chErr <- err
		return
	}

	for _, bl := range blocks {
		blockHash := bl.Hash().String()

		for _, tx := range bl.Transactions() {
			txHash := tx.Hash().String()

			// check TxIn scripts for known file headers/footers
			for txinIdx, txin := range tx.MsgTx().TxIn {
				matches := utils.SearchDataForKnownFileBits(txin.SignatureScript)
				for _, m := range matches {
					_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "in", txinIdx, m.Description()), true)
					if err != nil {
						chErr <- err
						return
					}
				}
			}

			// check TxOut scripts for known file headers/footers
			for txoutIdx, txout := range tx.MsgTx().TxOut {
				matches := utils.SearchDataForKnownFileBits(txout.PkScript)
				for _, m := range matches {
					_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", txoutIdx, m.Description()), true)
					if err != nil {
						chErr <- err
						return
					}
				}
			}

			parsedScriptData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
			if err != nil {
				chErr <- err
				return
			}

			matches := utils.SearchDataForKnownFileBits(parsedScriptData)
			for _, m := range matches {
				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", -1, m.Description()), true)
				if err != nil {
					chErr <- err
					return
				}
			}
		}
	}

	if err != nil {
		chErr <- err
		return
	}
}
