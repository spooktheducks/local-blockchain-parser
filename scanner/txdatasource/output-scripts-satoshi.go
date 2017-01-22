package txdatasource

import (
	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type OutputScriptsSatoshi struct{}

type OutputScriptsSatoshiResult []byte

// ensure that OutputScriptsSatoshi conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScriptsSatoshi{}

// ensure that OutputScriptsSatoshiResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = &OutputScriptsSatoshiResult{}

func (ds *OutputScriptsSatoshi) Name() string {
	return "all-txs-outputs-satoshi-concatenated"
}

func (ds *OutputScriptsSatoshi) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	data, err := tx.ConcatSatoshiDataFromTxOuts()
	if err != nil {
		return nil, err
	}

	return []scanner.ITxDataSourceResult{OutputScriptsSatoshiResult(data)}, nil
}

func (r OutputScriptsSatoshiResult) SourceName() string {
	return "all-txs-outputs-satoshi-concatenated"
}

func (r OutputScriptsSatoshiResult) RawData() []byte {
	return r
}
