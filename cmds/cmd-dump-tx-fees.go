package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

type DumpTxFeesCommand struct {
	startBlock, endBlock uint64
	datFileDir, outDir   string
	dbFile               string

	db *BlockDB
}

func NewDumpTxFeesCommand(startBlock, endBlock uint64, datFileDir, dbFile, outDir string) *DumpTxFeesCommand {
	return &DumpTxFeesCommand{
		startBlock: startBlock,
		endBlock:   endBlock,
		datFileDir: datFileDir,
		dbFile:     dbFile,
		outDir:     filepath.Join(".", outDir, "dump-tx-fees"),
	}
}

func (cmd *DumpTxFeesCommand) RunCommand() error {
	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	cmd.db = db

	err = os.MkdirAll(cmd.outDir, 0777)
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

func (cmd *DumpTxFeesCommand) parseBlock(blockFileNum int, chErr chan error, chDone chan bool, procLimiter chan bool) {
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

	outFile := utils.NewConditionalFile(filepath.Join(cmd.outDir, fmt.Sprintf("blk%05d-tx-fees.csv", blockFileNum)))
	defer outFile.Close()

	// write CSV header
	_, err = outFile.WriteString("block,tx,fee\n", false)
	if err != nil {
		chErr <- err
		return
	}

	numBlocks := len(blocks)
	for blIdx, bl := range blocks {
		blockHash := bl.Hash().String()

		// numTxs := len(bl.Transactions())
		for _, btctx := range bl.Transactions() {
			tx := Tx{Tx: btctx}
			tx.SetDB(cmd.db)

			txHash := tx.Hash().String()
			fee, err := tx.Fee()
			if err != nil {
				chErr <- err
				return
			}

			outFile.WriteString(fmt.Sprintf("%v,%v,%v\n", blockHash, txHash, fee), true)
		}

		fmt.Printf("finished block %s (%d/%d)\n", blockHash, blIdx, numBlocks)
	}

	if err != nil {
		chErr <- err
		return
	}
}
