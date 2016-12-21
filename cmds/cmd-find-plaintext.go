package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type (
	FindPlaintextCommand struct {
		startBlock, endBlock uint64
		inDir, outDir        string
	}
)

func NewFindPlaintextCommand(startBlock, endBlock uint64, inDir, outDir string) *FindPlaintextCommand {
	return &FindPlaintextCommand{
		startBlock: startBlock,
		endBlock:   endBlock,
		inDir:      inDir,
		outDir:     filepath.Join(outDir, "search-plaintext"),
	}
}

func (cmd *FindPlaintextCommand) RunCommand() error {
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

	// start a goroutine for each .dat file being parsed
	chDones := []chan bool{}
	procLimiter := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		procLimiter <- true
	}
	for i := int(cmd.startBlock); i < int(cmd.endBlock)+1; i++ {
		chDone := make(chan bool)
		go cmd.parseBlock(i, chErr, chDone, procLimiter)
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

func (cmd *FindPlaintextCommand) parseBlock(blockFileNum int, chErr chan error, chDone chan bool, procLimiter chan bool) {
	defer close(chDone)
	defer func() { procLimiter <- true }()
	<-procLimiter

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(cmd.inDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	outFile, err := utils.CreateFile(filepath.Join(cmd.outDir, fmt.Sprintf("blk%05d-plaintext.txt", blockFileNum)))
	if err != nil {
		chErr <- err
		return
	}
	defer utils.CloseFile(outFile)

	for _, bl := range blocks {
		blockHash := bl.Hash().String()

		for _, tx := range bl.Transactions() {
			txHash := tx.Hash().String()

			// extract text from each TxIn scriptSig
			data := make([]byte, 0)
			for _, txin := range tx.MsgTx().TxIn {
				data = append(data, txin.SignatureScript...)
				// txt, isText := utils.ExtractText(txin.SignatureScript)
				// if !isText || len(txt) < 8 {
				// 	continue
				// }
			}
			txt := utils.StripNonTextBytes(data)
			if len(txt) < 8 {
				continue
			}

			_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "in", -1, string(txt)))
			if err != nil {
				chErr <- err
				return
			}
			// }

			// extract text from each TxOut PkScript
			/*for txoutIdx, txout := range tx.MsgTx().TxOut {
				// txt, isText := utils.ExtractText(txout.PkScript)
				// if !isText || len(txt) < 8 {
				// 	continue
				// }
				txt := utils.StripNonTextBytes(txout.PkScript)
				if len(txt) < 8 {
					continue
				}

				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", txoutIdx, string(txt)))
				if err != nil {
					chErr <- err
					return
				}
			}*/

			// extract text from concatenated TxOut hex tokens

			parsedScriptData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
			if err != nil {
				chErr <- err
				return
			}

			// parsedScriptText, isText := utils.ExtractText(parsedScriptData)
			// if err != nil {
			// 	chErr <- err
			// 	return
			// }
			parsedScriptText := utils.StripNonTextBytes(parsedScriptData)
			// if isText && len(parsedScriptText) > 8 {
			if len(parsedScriptText) > 8 {
				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", -1, string(parsedScriptText)))
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
