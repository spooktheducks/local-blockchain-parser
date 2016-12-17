package querycmds

import (
	"fmt"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type BlockInfoCommand struct {
	dbFile     string
	datFileDir string
	blockHash  string
}

func NewBlockInfoCommand(datFileDir, dbFile, blockHash string) *BlockInfoCommand {
	return &BlockInfoCommand{
		datFileDir: datFileDir,
		dbFile:     dbFile,
		blockHash:  blockHash,
	}
}

func (cmd *BlockInfoCommand) RunCommand() error {
	db, err := blockdb.NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	blockIndexRow, err := db.GetBlockIndexRow(cmd.blockHash)
	if err != nil {
		return err
	}

	block, err := db.GetBlock(cmd.blockHash)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", block.Hash().String())
	fmt.Printf("  - %v\n", block.MsgBlock().Header.Timestamp)
	fmt.Printf("  - %v\n", blockIndexRow.Filename)

	return nil
}
