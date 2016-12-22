package blockdb

import (
	"fmt"

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

func DecodeHashList(txListBytes []byte) ([]chainhash.Hash, error) {
	if txListBytes == nil || len(txListBytes) == 0 {
		return nil, fmt.Errorf("blockdb.DecodeHashList: empty bytes")
	} else if len(txListBytes)%chainhash.HashSize != 0 {
		return nil, fmt.Errorf("blockdb.DecodeHashList: value is corrupted")
	}

	numTxs := len(txListBytes) / chainhash.HashSize
	txList := make([]chainhash.Hash, numTxs)
	for i := 0; i < numTxs; i++ {
		txHash, err := HashFromBytes(txListBytes[i*chainhash.HashSize : (i+1)*chainhash.HashSize])
		if err != nil {
			return nil, fmt.Errorf("blockdb.DecodeHashList: %v", err.Error())
		}

		txList[i] = txHash
	}

	return txList, nil
}
