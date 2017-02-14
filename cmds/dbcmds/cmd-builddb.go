package dbcmds

import (
	"fmt"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type BuildBlockDBCommand struct {
	dbFile     string
	datFileDir string
	startBlock uint64
	endBlock   uint64
	indexWhat  string
	force      bool
}

func NewBuildBlockDBCommand(startBlock, endBlock uint64, datFileDir, dbFile, indexWhat string, force bool) (*BuildBlockDBCommand, error) {
	if indexWhat != "blocks" && indexWhat != "transactions" {
		return nil, fmt.Errorf("must specify either 'blocks' or 'transactions'")
	}

	return &BuildBlockDBCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		startBlock: startBlock,
		endBlock:   endBlock,
		indexWhat:  indexWhat,
		force:      force,
	}, nil
}

func (cmd *BuildBlockDBCommand) RunCommand() error {
	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	switch cmd.indexWhat {
	case "transactions":
		err = db.IndexDATFileTransactions(cmd.startBlock, cmd.endBlock, cmd.force)
		if err != nil {
			return err
		}
		return nil

	case "blocks":
		err = db.IndexDATFileBlocks(cmd.startBlock, cmd.endBlock, cmd.force)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
