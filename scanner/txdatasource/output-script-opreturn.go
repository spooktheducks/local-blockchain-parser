package txdatasource

import (
	"fmt"

	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
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

func (ds *OutputScriptOpReturn) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for i := range tx.MsgTx().TxOut {
		bs, err := tx.GetOPReturnDataFromTxOut(i)
		if err != nil {
			continue
		}

		results = append(results, OutputScriptOpReturnResult{rawData: bs, index: i})
	}

	return results, nil
}

func (r OutputScriptOpReturnResult) SourceName() string {
	return fmt.Sprintf("txout-script-opreturn-%d", r.index)
}

func (r OutputScriptOpReturnResult) RawData() []byte {
	return r.rawData
}
