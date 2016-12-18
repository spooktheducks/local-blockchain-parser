package dbcmds

import (
	"fmt"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type BuildBlockDBCommand struct {
	dbFile     string
	datFileDir string
	startBlock uint64
	endBlock   uint64
	indexWhat  string
}

func NewBuildBlockDBCommand(startBlock, endBlock uint64, datFileDir, dbFile, indexWhat string) (*BuildBlockDBCommand, error) {
	if indexWhat != "blocks" && indexWhat != "transactions" {
		return nil, fmt.Errorf("must specify either 'blocks' or 'transactions'")
	}

	return &BuildBlockDBCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		startBlock: startBlock,
		endBlock:   endBlock,
		indexWhat:  indexWhat,
	}, nil
}

func (cmd *BuildBlockDBCommand) RunCommand() error {
	db, err := blockdb.NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	switch cmd.indexWhat {
	case "transactions":
		err = db.IndexDATFileTransactions(cmd.startBlock, cmd.endBlock)
		if err != nil {
			return err
		}
		return nil

	case "blocks":
		err = db.IndexDATFileBlocks(cmd.startBlock, cmd.endBlock)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
