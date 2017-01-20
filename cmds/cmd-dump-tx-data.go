package cmds

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/types"
)

type DumpTxDataCommand struct {
	startBlock, endBlock uint64
	datFileDir, outDir   string
}

func NewDumpTxDataCommand(startBlock, endBlock uint64, datFileDir, outDir string) *DumpTxDataCommand {
	return &DumpTxDataCommand{
		startBlock: startBlock,
		endBlock:   endBlock,
		datFileDir: datFileDir,
		outDir:     filepath.Join(".", outDir, "dump-tx-data"),
	}
}

func (cmd *DumpTxDataCommand) RunCommand() error {
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

func (cmd *DumpTxDataCommand) parseBlock(blockFileNum int, chErr chan error, chDone chan bool, procLimiter chan bool) {
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

	outFile := utils.NewConditionalFile(filepath.Join(cmd.outDir, fmt.Sprintf("blk%05d-tx-data.csv", blockFileNum)))
	defer outFile.Close()

	// write CSV header
	_, err = outFile.WriteString("block,recipient address,tx,output index,data\n", false)
	if err != nil {
		chErr <- err
		return
	}

	// numBlocks := len(blocks)
	for _, bl := range blocks {
		blockHash := bl.Hash().String()

		// numTxs := len(bl.Transactions())
		for _, btctx := range bl.Transactions() {
			tx := Tx{Tx: btctx}

			txHash := tx.Hash().String()

			for txoutIdx := range tx.MsgTx().TxOut {
				addrs, err := tx.GetTxOutAddress(txoutIdx)
				if err != nil {
					chErr <- err
					return
				}

				recipientAddr := ""
				if len(addrs) > 0 {
					recipientAddr = addrs[0].EncodeAddress()
				}

				txoutData, err := tx.GetNonOPDataFromTxOut(txoutIdx)
				if err != nil {
					chErr <- err
					return
				}

				txoutDataHex := hex.EncodeToString(txoutData)

				outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, recipientAddr, txHash, txoutIdx, txoutDataHex), true)
			}
		}
	}

	if err != nil {
		chErr <- err
		return
	}
}
