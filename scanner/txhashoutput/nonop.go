package txhashoutput

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type NonOp struct {
	OutDir   string
	Filename string
	data     []nonOpOutputLine
}

type nonOpOutputLine struct {
	txHash chainhash.Hash
	data   []byte
}

func (o *NonOp) OutputTx(tx *Tx) error {
	data, err := tx.ConcatNonOPDataFromTxOuts()
	if err != nil {
		return err
	}

	o.data = append(o.data, nonOpOutputLine{txHash: *tx.Hash(), data: data})
	return nil
}

func (o *NonOp) Close() error {
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
