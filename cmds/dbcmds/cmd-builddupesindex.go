package dbcmds

import (
	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type BuildDupesIndexCommand struct {
	dbFile     string
	datFileDir string
	startBlock uint64
	endBlock   uint64
}

func NewBuildDupesIndexCommand(startBlock, endBlock uint64, datFileDir, dbFile string) *BuildDupesIndexCommand {
	return &BuildDupesIndexCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		startBlock: startBlock,
		endBlock:   endBlock,
	}
}

func (cmd *BuildDupesIndexCommand) RunCommand() error {
	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.IndexDATFileTxOutDuplicates(cmd.startBlock, cmd.endBlock)
	if err != nil {
		return err
	}
	return nil
}
