package dbcmds

import (
	"fmt"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type TxInfoCommand struct {
	dbFile     string
	datFileDir string
	txHash     string
}

func NewTxInfoCommand(datFileDir, dbFile, txHash string) *TxInfoCommand {
	return &TxInfoCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		txHash:     txHash,
	}
}

func (cmd *TxInfoCommand) RunCommand() error {
	db, err := blockdb.NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.GetTx(cmd.txHash)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", tx.Hash().String())
	for _, txin := range tx.MsgTx().TxIn {
		fmt.Printf("  - PrevOutPoint: (%v) %v\n", txin.PreviousOutPoint.Index, txin.PreviousOutPoint.Hash.String())
	}

	// for _, txout := range tx.MsgTx().TxOut {
	// 	fmt.Printf("%v\n", txout.)
	// }

	return nil
}
