package txhashoutput

import (
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type HashOnly struct {
	OutDir   string
	Filename string
	data     []chainhash.Hash
}

func (o *HashOnly) OutputTx(tx *Tx) error {
	o.data = append(o.data, *tx.Hash())
	return nil
}

func (o *HashOnly) Close() error {
	f, err := os.Create(filepath.Join(o.OutDir, o.Filename))
	if err != nil {
		return err
	}
	defer f.Close()

	for _, txHash := range o.data {
		f.WriteString(txHash.String() + "\n")
	}

	return nil
}
