package dbcmds

import (
	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type BuildSpentTxOutIndexCommand struct {
	dbFile     string
	datFileDir string
	startBlock uint64
	endBlock   uint64
	force      bool
}

func NewBuildSpentTxOutIndexCommand(startBlock, endBlock uint64, datFileDir, dbFile string, force bool) *BuildSpentTxOutIndexCommand {
	return &BuildSpentTxOutIndexCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		startBlock: startBlock,
		endBlock:   endBlock,
		force:      force,
	}
}

func (cmd *BuildSpentTxOutIndexCommand) RunCommand() error {
	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.IndexDATFileSpentTxOuts(cmd.startBlock, cmd.endBlock, cmd.force)
	if err != nil {
		return err
	}
	return nil
}
