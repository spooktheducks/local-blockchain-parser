package blockdb

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func HashFromBytes(bs []byte) (chainhash.Hash, error) {
	hash := &chainhash.Hash{}
	err := hash.SetBytes(bs)
	if err != nil {
		return chainhash.Hash{}, err
	}
	return *hash, nil
}

func HashFromString(s string) (chainhash.Hash, error) {
	h, err := chainhash.NewHashFromStr(s)
	if err != nil {
		return chainhash.Hash{}, err
	}
	return *h, nil
}
