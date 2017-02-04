package txdatasource

import (
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/wire"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type OutputScript struct {
	SkipMaxValueTxOut bool
	OrderByValue      bool
}

type OutputScriptResult struct {
	rawData []byte
	index   int
}

// ensure that OutputScript conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScript{}

// ensure that OutputScriptResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = OutputScriptResult{}

func (ds *OutputScript) Name() string {
	if ds.OrderByValue && ds.SkipMaxValueTxOut {
		return "txout-script-byvalue-skipmaxvalue"
	} else if ds.OrderByValue {
		return "txout-script-byvalue"
	} else if ds.SkipMaxValueTxOut {
		return "txout-script-skipmaxvalue"
	} else {
		return "txout-script"
	}
}

func (ds *OutputScript) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	txouts := tx.MsgTx().TxOut

	skipTxoutIdx := tx.FindMaxValueTxOut()
	if ds.OrderByValue {
		txouts = ds.sortTxOuts(txouts)

		var maxValue int64
		for txoutIdx, txout := range txouts {
			if txout.Value > maxValue {
				maxValue = txout.Value
				skipTxoutIdx = txoutIdx
			}
		}
	}

	results := []scanner.ITxDataSourceResult{}
	for i := range txouts {
		if ds.SkipMaxValueTxOut && i == skipTxoutIdx {
			continue
		}
		bs, err := tx.GetNonOPDataFromTxOut(i)
		if err != nil {
			continue
		}

		results = append(results, OutputScriptResult{rawData: bs, index: i})
	}

	return results, nil
}

func (ds *OutputScript) sortTxOuts(txouts []*wire.TxOut) []*wire.TxOut {
	// this copy is absolutely necessary since sort.Sort is in-place.  otherwise, later
	// ITxDataSources will receive shuffled data.
	txoutsCopy := make([]*wire.TxOut, len(txouts))
	copy(txoutsCopy, txouts)
	sort.Sort(sortableTxOuts(txoutsCopy))
	return txoutsCopy
}

func (r OutputScriptResult) SourceName() string {
	return fmt.Sprintf("txout-script-%d", r.index)
}

func (r OutputScriptResult) RawData() []byte {
	return r.rawData
}

type sortableTxOuts []*wire.TxOut

func (sto sortableTxOuts) Len() int {
	return len(sto)
}

func (sto sortableTxOuts) Less(i, j int) bool {
	return sto[i].Value < sto[j].Value
}

func (sto sortableTxOuts) Swap(i, j int) {
	sto[i], sto[j] = sto[j], sto[i]
}
