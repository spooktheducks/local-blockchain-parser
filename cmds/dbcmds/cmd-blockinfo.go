package dbcmds

import (
	"fmt"

	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
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
	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	blockHash, err := utils.HashFromString(cmd.blockHash)
	if err != nil {
		return err
	}

	blockIndexRow, err := db.GetBlockIndexRow(blockHash)
	if err != nil {
		return err
	}

	block, err := db.GetBlock(blockHash)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", block.Hash().String())
	fmt.Printf("  - %v\n", block.MsgBlock().Header.Timestamp)
	fmt.Printf("  - %v\n", blockIndexRow.DATFilename())

	return nil
}
