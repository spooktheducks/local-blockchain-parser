package txdatasource

import (
	"fmt"

	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type OutputScript struct{}

type OutputScriptResult struct {
	rawData []byte
	index   int
}

// ensure that OutputScript conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScript{}

// ensure that OutputScriptResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = OutputScriptResult{}

func (ds *OutputScript) Name() string {
	return "txout-script"
}

func (ds *OutputScript) GetData(tx *btcutil.Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		bs, err := utils.GetNonOPBytes(txout.PkScript)
		if err != nil {
			continue
		}

		results = append(results, OutputScriptResult{rawData: bs, index: txoutIdx})
	}

	return results, nil
}

func (r OutputScriptResult) SourceName() string {
	return fmt.Sprintf("txout-script-%d", r.index)
}

func (r OutputScriptResult) RawData() []byte {
	return r.rawData
}
