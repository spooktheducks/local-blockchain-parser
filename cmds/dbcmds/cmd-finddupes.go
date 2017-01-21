package dbcmds

import (
	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type ScanDupesIndexCommand struct {
	dbFile     string
	datFileDir string
}

func NewScanDupesIndexCommand(datFileDir, dbFile string) *ScanDupesIndexCommand {
	return &ScanDupesIndexCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
	}
}

func (cmd *ScanDupesIndexCommand) RunCommand() error {
	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.ScanTxOutDuplicateData()
	if err != nil {
		return err
	}
	return nil
}
