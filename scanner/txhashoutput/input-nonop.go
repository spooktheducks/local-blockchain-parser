package txhashoutput

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type InputScriptNonOP struct {
	OutDir   string
	Filename string
	data     []inputScriptNonOPOutputLine
}

type inputScriptNonOPOutputLine struct {
	txHash chainhash.Hash
	data   []byte
}

func (o *InputScriptNonOP) OutputTx(tx *Tx) error {
	data, err := tx.ConcatNonOPDataFromTxIns()
	if err != nil {
		return err
	}

	o.data = append(o.data, inputScriptNonOPOutputLine{txHash: *tx.Hash(), data: data})
	return nil
}

func (o *InputScriptNonOP) Close() error {
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
