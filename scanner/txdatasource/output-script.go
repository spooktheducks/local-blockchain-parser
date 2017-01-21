package txdatasource

import (
	"fmt"

	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
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

func (ds *OutputScript) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for i := range tx.MsgTx().TxOut {
		bs, err := tx.GetNonOPDataFromTxOut(i)
		if err != nil {
			continue
		}

		results = append(results, OutputScriptResult{rawData: bs, index: i})
	}

	return results, nil
}

func (r OutputScriptResult) SourceName() string {
	return fmt.Sprintf("txout-script-%d", r.index)
}

func (r OutputScriptResult) RawData() []byte {
	return r.rawData
}
