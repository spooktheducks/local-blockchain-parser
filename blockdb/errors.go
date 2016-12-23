package blockdb

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type DataNotIndexedError struct {
	Index string
}

func (e DataNotIndexedError) Error() string {
	return fmt.Sprintf(`It seems that you need to build the %v index before running this command.  Try running the "builddb %v" command on your .dat files.`, e.Index, e.Index)
}

type TxNotFoundError struct {
	TxHash chainhash.Hash
}

func (e TxNotFoundError) Error() string {
	return fmt.Sprintf(`Can't find transaction %v.\n`+
		`Try running the "builddb transactions" command on the .dat file containing this transaction.\n`+
		`You can look up the transaction on blockchain.info to determine which block it's in.  Once\n`+
		`you have the block hash, you can run "querydb blocks [block hash]" to find the filename of\n`+
		`the .dat file you need to index.`, e.TxHash.String())
}

type BlockNotFoundError struct {
	BlockHash chainhash.Hash
}

func (e BlockNotFoundError) Error() string {
	return fmt.Sprintf(`Can't find block %v.\n`+
		`Try running the "builddb blocks" command on the .dat file containing this block, or simply\n`+
		`run the command on the entire set of .dat files.\n`, e.BlockHash.String())
}
