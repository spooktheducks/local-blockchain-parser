package txdatasource

import (
	"fmt"

	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type OutputScriptOpReturn struct{}

type OutputScriptOpReturnResult struct {
	rawData []byte
	index   int
}

// ensure that OutputScriptOpReturn conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScriptOpReturn{}

// ensure that OutputScriptOpReturnResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = OutputScriptOpReturnResult{}

func (ds *OutputScriptOpReturn) Name() string {
	return "txout-script-opreturn"
}

func (ds *OutputScriptOpReturn) GetData(tx *btcutil.Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		bs, err := utils.GetOPReturnBytes(txout.PkScript)
		if err != nil {
			continue
		}

		results = append(results, OutputScriptOpReturnResult{rawData: bs, index: txoutIdx})
	}

	return results, nil
}

func (r OutputScriptOpReturnResult) SourceName() string {
	return fmt.Sprintf("txout-script-opreturn-%d", r.index)
}

func (r OutputScriptOpReturnResult) RawData() []byte {
	return r.rawData
}
