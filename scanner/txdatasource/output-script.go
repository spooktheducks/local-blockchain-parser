package txdatasource

import (
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/wire"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
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
	_txouts := tx.MsgTx().TxOut
	txouts := make([]wrappedTxOut, len(_txouts))

	for i := range _txouts {
		txouts[i] = wrappedTxOut{TxOut: _txouts[i], index: i}
	}

	skipTxoutIdx := tx.FindMaxValueTxOut()

	if ds.OrderByValue {
		sort.Sort(valueSortedTxOuts(txouts))
	}

	results := []scanner.ITxDataSourceResult{}
	for _, txout := range txouts {
		if ds.SkipMaxValueTxOut && txout.index == skipTxoutIdx {
			continue
		}
		bs, err := utils.GetNonOPBytesFromOutputScript(txout.PkScript)
		if err != nil {
			continue
		}

		results = append(results, OutputScriptResult{rawData: bs, index: txout.index})
	}

	return results, nil
}

// func (ds *OutputScript) sortTxOutsByValue(txouts []wrappedTxOut) []*wire.TxOut {
// 	// this copy is absolutely necessary since sort.Sort is in-place.  otherwise, later
// 	// ITxDataSources will receive shuffled data.
// 	// txoutsCopy := make([]*wire.TxOut, len(txouts))
// 	// copy(txoutsCopy, txouts)
// 	sort.Sort(valueSortedTxOuts(txoutsCopy))
// 	return txoutsCopy
// }

func (r OutputScriptResult) SourceName() string {
	return fmt.Sprintf("txout-script-%d", r.index)
}

func (r OutputScriptResult) RawData() []byte {
	return r.rawData
}

func (r OutputScriptResult) InOut() string {
	return "out"
}

func (r OutputScriptResult) Index() int {
	return r.index
}

type wrappedTxOut struct {
	*wire.TxOut
	index int
}

type valueSortedTxOuts []wrappedTxOut

func (sto valueSortedTxOuts) Len() int {
	return len(sto)
}

func (sto valueSortedTxOuts) Less(i, j int) bool {
	return sto[i].Value < sto[j].Value
}

func (sto valueSortedTxOuts) Swap(i, j int) {
	sto[i], sto[j] = sto[j], sto[i]
}
