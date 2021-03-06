package txhashsource

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

func NewChain(db *BlockDB, startHash chainhash.Hash, limit uint) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		chBackwards := NewBackwardChain(db, startHash, limit)
		chForwards := NewForwardChain(db, startHash, limit)

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

func NewForwardChain(db *BlockDB, startHash chainhash.Hash, limit uint) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		currentTxHash := startHash
		var i uint
		for {
			if limit > 0 && i >= limit {
				break
			}

			tx, err := db.GetTx(currentTxHash)
			if err != nil {
				// @@TODO
				panic(err)
			}

			// if !tx.HasSuspiciousOutputValues() {
			// 	fmt.Println("no suspicious output values, stopping")
			// 	break
			// }
			ch <- currentTxHash

			key := SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(tx.FindMaxValueTxOut())}
			spentTxOut, err := db.GetSpentTxOut(key)
			if err != nil {
				fmt.Println("err", err)
				// @@TODO
				// panic(err)
				break
			}

			currentTxHash = spentTxOut.InputTxHash
			i++
		}
	}()

	return TxHashSource(ch)
}

func NewBackwardChain(db *BlockDB, startHash chainhash.Hash, limit uint) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		emptyHash := chainhash.Hash{}

		foundHashesReverse := []chainhash.Hash{}
		currentTxHash := startHash
		var i uint
		for {
			if limit > 0 && i >= limit {
				break
			}

			if currentTxHash == emptyHash {
				// this is the coinbase, so we can't follow further backwards
				break
			}

			tx, err := db.GetTx(currentTxHash)
			if err != nil {
				// @@TODO
				panic(err)
			}

			// if tx.HasSuspiciousOutputValues() {
			foundHashesReverse = append(foundHashesReverse, currentTxHash)
			if len(tx.MsgTx().TxIn) == 1 {
				currentTxHash = tx.MsgTx().TxIn[0].PreviousOutPoint.Hash
			} else {
				break
			}
			// } else {
			// 	break
			// }

			i++
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
