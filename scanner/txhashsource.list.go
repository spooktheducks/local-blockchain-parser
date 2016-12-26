package scanner

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func NewListTxHashSource(hashes []chainhash.Hash) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		for _, hash := range hashes {
			ch <- hash
		}
	}()

	return TxHashSource(ch)
}
