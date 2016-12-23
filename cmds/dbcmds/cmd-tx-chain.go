package dbcmds

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

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

	startHash, err := blockdb.HashFromString(cmd.txHash)
	if err != nil {
		return err
	}

	foundHashes1, err := cmd.CrawlBackwards(startHash, db)
	if err != nil {
		return err
	}

	foundHashes2, err := cmd.CrawlForwards(startHash, db)
	if err != nil {
		return err
	}

	// both foundHashes1 and foundHashes2 contain startHash, so we omit it from one of them
	foundHashes := append(foundHashes1, foundHashes2[1:]...)

	err = cmd.writeDataFromTxs(foundHashes, db)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *TxChainCommand) CrawlBackwards(startHash chainhash.Hash, db *blockdb.BlockDB) ([]chainhash.Hash, error) {
	foundHashesReverse := []chainhash.Hash{}
	currentTxHash := startHash
	for {
		tx, err := db.GetTx(currentTxHash)
		if err != nil {
			return nil, err
		}

		if utils.TxHasSuspiciousOutputValues(tx) {
			foundHashesReverse = append(foundHashesReverse, currentTxHash)
			if len(tx.MsgTx().TxIn) == 1 {
				currentTxHash = tx.MsgTx().TxIn[0].PreviousOutPoint.Hash
			} else {
				break
			}
		} else {
			break
		}
	}

	numHashes := len(foundHashesReverse)
	foundHashes := make([]chainhash.Hash, numHashes)
	for i := 0; i < numHashes; i++ {
		foundHashes[numHashes-i-1] = foundHashesReverse[i]
	}

	return foundHashes, nil
}

func (cmd *TxChainCommand) CrawlForwards(startHash chainhash.Hash, db *blockdb.BlockDB) ([]chainhash.Hash, error) {
	foundHashes := []chainhash.Hash{}
	currentTxHash := startHash
	for {
		tx, err := db.GetTx(currentTxHash)
		if err != nil {
			return nil, err
		}

		if utils.TxHasSuspiciousOutputValues(tx) {
			foundHashes = append(foundHashes, currentTxHash)

			maxValueTxoutIdx := utils.FindMaxValueTxOut(tx)

			key := blockdb.SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(maxValueTxoutIdx)}
			spentTxOut, err := db.GetSpentTxOut(key)
			if err != nil {
				return nil, err
			}

			currentTxHash = spentTxOut.InputTxHash

		} else {
			break
		}
	}
	return foundHashes, nil
}

func (cmd *TxChainCommand) writeDataFromTxs(txHashes []chainhash.Hash, db *blockdb.BlockDB) error {
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

		matches := utils.SearchDataForMagicFileBytes(data)
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
