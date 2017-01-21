package txhashsource

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

func NewChain(db *BlockDB, startHash chainhash.Hash) TxHashSource {
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

func NewForwardChain(db *BlockDB, startHash chainhash.Hash) TxHashSource {
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

			key := SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(tx.FindMaxValueTxOut())}
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

func NewBackwardChain(db *BlockDB, startHash chainhash.Hash) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		emptyHash := chainhash.Hash{}

		foundHashesReverse := []chainhash.Hash{}
		currentTxHash := startHash
		for {
			if currentTxHash == emptyHash {
				// this is the coinbase, so we can't follow further backwards
				break
			}

			tx, err := db.GetTx(currentTxHash)
			if err != nil {
				// @@TODO
				panic(err)
			}

			// if utils.TxHasSuspiciousOutputValues(tx) {
			foundHashesReverse = append(foundHashesReverse, currentTxHash)
			if len(tx.MsgTx().TxIn) == 1 {
				currentTxHash = tx.MsgTx().TxIn[0].PreviousOutPoint.Hash
			} else {
				break
			}
			// } else {
			// 	break
			// }
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
