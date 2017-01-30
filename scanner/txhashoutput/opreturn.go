package txhashoutput

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type OpReturn struct {
	OutDir   string
	Filename string
	data     []outputLine
}

type outputLine struct {
	txHash chainhash.Hash
	data   []byte
}

func (o *OpReturn) OutputTx(tx *Tx) error {
	data, err := tx.ConcatOPReturnDataFromTxOuts()
	if err != nil {
		return err
	}

	o.data = append(o.data, outputLine{txHash: *tx.Hash(), data: data})
	return nil
}

func (o *OpReturn) Close() error {
	f, err := os.Create(filepath.Join(o.OutDir, o.Filename))
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range o.data {
		f.WriteString(fmt.Sprintf("%s %s\n", line.txHash.String(), string(line.data)))
	}

	return nil
}
