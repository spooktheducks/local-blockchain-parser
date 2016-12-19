package dbcmds

import (
	"fmt"

	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type TxChainCommand struct {
	dbFile     string
	datFileDir string
	txHash     string
}

func NewTxChainCommand(datFileDir, dbFile, txHash string) *TxChainCommand {
	return &TxChainCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		txHash:     txHash,
	}
}

func (cmd *TxChainCommand) RunCommand() error {
	db, err := blockdb.NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	foundHashesReverse := []string{}
	currentTxHash := cmd.txHash
	for {
		tx, err := db.GetTx(currentTxHash)
		if err != nil {
			return err
		}

		if cmd.isSuspiciousTx(tx) {
			foundHashesReverse = append(foundHashesReverse, currentTxHash)
			if len(tx.MsgTx().TxIn) == 1 {
				currentTxHash = tx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String()
			} else {
				break
			}
		} else {
			break
		}
	}

	numHashes := len(foundHashesReverse)
	foundHashes := make([]string, numHashes)
	for i := 0; i < numHashes; i++ {
		foundHashes[numHashes-i-1] = foundHashesReverse[i]
	}

	err = cmd.writeDataFromTxs(foundHashes, db)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *TxChainCommand) isSuspiciousTx(tx *btcutil.Tx) bool {
	numTinyValues := 0
	for _, txout := range tx.MsgTx().TxOut {
		if utils.SatoshisToBTCs(txout.Value) == 0.00000001 {
			numTinyValues++
		}
	}

	if numTinyValues == len(tx.MsgTx().TxOut)-1 {
		return true
	}
	return false
}

func (cmd *TxChainCommand) writeDataFromTxs(txHashes []string, db *blockdb.BlockDB) error {
	outFile, err := utils.CreateFile("output/txchain-output")
	if err != nil {
		return err
	}
	defer utils.CloseFile(outFile)

	for _, txHash := range txHashes {
		tx, err := db.GetTx(txHash)
		if err != nil {
			return err
		}

		data := []byte{}
		// we skip the final two TxOuts because one goes to WL and one is used to pass BTC to the next transaction in the chain
		for i := 0; i < len(tx.MsgTx().TxOut)-2; i++ {
			bs, err := utils.GetNonOPBytes(tx.MsgTx().TxOut[i].PkScript)
			if err != nil {
				continue
			}

			data = append(data, bs...)
		}

		data, err = utils.GetSatoshiEncodedData(data)
		if err != nil {
			return err
		}

		matches := utils.SearchDataForKnownFileBits(data)
		for _, m := range matches {
			fmt.Println(m.Description())
		}

		_, err = outFile.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}
