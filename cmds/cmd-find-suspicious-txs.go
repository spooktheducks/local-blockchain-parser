package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

func FindSuspiciousTxs(startBlock, endBlock uint64, datFileDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "suspicious-txs")

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
	// chDones := []chan bool{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		// chDone := make(chan bool)
		findSuspiciousTxsParseBlock(datFileDir, outSubdir, i, chErr)
		// chDones = append(chDones, chDone)
	}

	// wait for all ops to complete
	// for _, chDone := range chDones {
	// 	<-chDone
	// }

	// close error channel
	close(chErr)

	return nil
}

func findSuspiciousTxsParseBlock(datFileDir string, outDir string, blockFileNum int, chErr chan error) {
	// defer close(chDone)

	csvFile, err := utils.CreateFile(filepath.Join(outDir, fmt.Sprintf("suspicious-txs-blk%05d.txt", blockFileNum)))
	if err != nil {
		chErr <- err
		return
	}
	defer utils.CloseFile(csvFile)

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(datFileDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	for _, bl := range blocks {
		blockHash := bl.Hash().String()

		blockTimestamp := bl.MsgBlock().Header.Timestamp
		csvFile.WriteString("======= " + blockTimestamp.String() + " =======\n")

		for _, btctx := range bl.Transactions() {
			tx := Tx{Tx: btctx}

			txHash := tx.Hash().String()

			if tx.HasSuspiciousOutputValues() {
				numInputs := len(tx.MsgTx().TxIn)
				numOutputs := len(tx.MsgTx().TxOut)
				csvFile.WriteString(fmt.Sprintf("%v,%v,%v,%v\n", blockHash, txHash, numInputs, numOutputs))
			}
		}
	}
}
