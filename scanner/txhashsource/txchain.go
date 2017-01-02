package txhashsource

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

func NewChain(db *blockdb.BlockDB, startHash chainhash.Hash) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		chBackwards := NewBackwardChain(db, startHash)
		chForwards := NewForwardChain(db, startHash)

		for hash := range chBackwards {
			ch <- hash
		}

		<-chForwards // skip first element from forward source because both backwards + forwards contain startHash
		for hash := range chForwards {
			ch <- hash
		}
	}()

	return TxHashSource(ch)
}

func NewForwardChain(db *blockdb.BlockDB, startHash chainhash.Hash) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		currentTxHash := startHash
		for {
			tx, err := db.GetTx(currentTxHash)
			if err != nil {
				// @@TODO
				panic(err)
			}

			// if !utils.TxHasSuspiciousOutputValues(tx) {
			// 	break
			// }
			ch <- currentTxHash

			maxValueTxoutIdx := utils.FindMaxValueTxOut(tx)

			key := blockdb.SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(maxValueTxoutIdx)}
			spentTxOut, err := db.GetSpentTxOut(key)
			if err != nil {
				// @@TODO
				// panic(err)
				break
			}

			currentTxHash = spentTxOut.InputTxHash
		}
	}()

	return TxHashSource(ch)
}

func NewBackwardChain(db *blockdb.BlockDB, startHash chainhash.Hash) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		foundHashesReverse := []chainhash.Hash{}
		currentTxHash := startHash
		for {
			tx, err := db.GetTx(currentTxHash)
			if err != nil {
				// @@TODO
				panic(err)
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

		for _, hash := range foundHashes {
			ch <- hash
		}
	}()

	return TxHashSource(ch)
}
