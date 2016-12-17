package querycmds

import (
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type BuildBlockDBCommand struct {
	dbFile     string
	datFileDir string
	startBlock uint64
	endBlock   uint64
}

func NewBuildBlockDBCommand(startBlock, endBlock uint64, datFileDir, dbFile string) *BuildBlockDBCommand {
	return &BuildBlockDBCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		startBlock: startBlock,
		endBlock:   endBlock,
	}
}

func (cmd *BuildBlockDBCommand) RunCommand() error {
	db, err := blockdb.NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.IndexDATFiles(cmd.startBlock, cmd.endBlock)
	if err != nil {
		return err
	}

	return nil
}
