package blockdb

import (
	// "github.com/btcsuite/btcd/chaincfg"
	// "github.com/btcsuite/btcd/chaincfg/chainhash"
	// "github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	// "github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type Block struct {
	*btcutil.Block

	DATFileIdx     uint16
	Timestamp      int64
	IndexInDATFile uint32
}
