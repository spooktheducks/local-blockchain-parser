package txhashsource

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type TxHashSource <-chan chainhash.Hash

func (hs TxHashSource) NextHash() (chainhash.Hash, bool) {
	if hash, ok := <-hs; ok {
		return hash, true
	}
	return chainhash.Hash{}, false
}
