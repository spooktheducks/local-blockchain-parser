package blockdb

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

func DecodeHashList(txListBytes []byte) ([]chainhash.Hash, error) {
	if txListBytes == nil || len(txListBytes) == 0 {
		return nil, fmt.Errorf("DecodeHashList: empty bytes")
	} else if len(txListBytes)%chainhash.HashSize != 0 {
		return nil, fmt.Errorf("DecodeHashList: value is corrupted")
	}

	numTxs := len(txListBytes) / chainhash.HashSize
	txList := make([]chainhash.Hash, numTxs)
	for i := 0; i < numTxs; i++ {
		txHash, err := utils.HashFromBytes(txListBytes[i*chainhash.HashSize : (i+1)*chainhash.HashSize])
		if err != nil {
			return nil, fmt.Errorf("DecodeHashList: %v", err.Error())
		}

		txList[i] = txHash
	}

	return txList, nil
}
