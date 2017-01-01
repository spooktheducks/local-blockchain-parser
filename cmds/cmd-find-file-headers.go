package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type FindFileHeadersCommand struct {
	startBlock, endBlock uint64
	datFileDir, outDir   string
}

func NewFindFileHeadersCommand(startBlock, endBlock uint64, datFileDir, outDir string) *FindFileHeadersCommand {
	return &FindFileHeadersCommand{
		startBlock: startBlock,
		endBlock:   endBlock,
		datFileDir: datFileDir,
		outDir:     filepath.Join(".", outDir, "find-file-headers"),
	}
}

func (cmd *FindFileHeadersCommand) RunCommand() error {
	err := os.MkdirAll(cmd.outDir, 0777)
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

	// start a goroutine for each .dat file being parsed, limited to 5 at a time
	chDones := []chan bool{}
	procLimiter := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		procLimiter <- true
	}

	for i := int(cmd.startBlock); i < int(cmd.endBlock)+1; i++ {
		chDone := make(chan bool)
		chDones = append(chDones, chDone)
		go cmd.parseBlock(i, chErr, chDone, procLimiter)
	}

	// wait for all goroutines to complete
	for _, chDone := range chDones {
		<-chDone
	}

	// close error channel
	close(chErr)

	return nil
}

func (cmd *FindFileHeadersCommand) parseBlock(blockFileNum int, chErr chan error, chDone chan bool, procLimiter chan bool) {
	defer close(chDone)
	defer func() { procLimiter <- true }()
	<-procLimiter

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(cmd.datFileDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	outFile := utils.NewConditionalFile(filepath.Join(cmd.outDir, fmt.Sprintf("blk%05d-file-headers.txt", blockFileNum)))
	defer outFile.Close()

	// write CSV header
	_, err = outFile.WriteString("block,tx,input or output,index,description\n", false)
	if err != nil {
		chErr <- err
		return
	}

	// numBlocks := len(blocks)
	for _, bl := range blocks {
		blockHash := bl.Hash().String()

		// numTxs := len(bl.Transactions())
		for _, tx := range bl.Transactions() {
			txHash := tx.Hash().String()

			/*
				// check TxIn scripts for known file headers/footers
				for txinIdx, txin := range tx.MsgTx().TxIn {
					matches := utils.SearchDataForMagicFileBytes(txin.SignatureScript)
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
					matches := utils.SearchDataForMagicFileBytes(txout.PkScript)
					for _, m := range matches {
						_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", txoutIdx, m.Description()), true)
						if err != nil {
							chErr <- err
							return
						}
					}
				}*/

			inData, err := utils.ConcatTxInScripts(tx)
			if err != nil {
				chErr <- err
				return
			}

			matches := utils.SearchDataForMagicFileBytes(inData)
			for _, m := range matches {
				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "in", -1, m.Description()), true)
				if err != nil {
					chErr <- err
					return
				}
			}

			outData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
			if err != nil {
				chErr <- err
				return
			}

			matches = utils.SearchDataForMagicFileBytes(outData)
			for _, m := range matches {
				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", -1, m.Description()), true)
				if err != nil {
					chErr <- err
					return
				}
			}

			// fmt.Printf("finished %v (%v/%v) (%v/%v)\n", txHash, txIdx, numTxs, blIdx, numBlocks)
		}
	}

	if err != nil {
		chErr <- err
		return
	}
}
