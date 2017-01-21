package txdatasource

import (
	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type OutputScriptsConcat struct{}

type OutputScriptsConcatResult []byte

// ensure that OutputScriptsConcat conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScriptsConcat{}

// ensure that OutputScriptsConcatResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = &OutputScriptsConcatResult{}

func (ds *OutputScriptsConcat) Name() string {
	return "outputs-concatenated"
}

func (ds *OutputScriptsConcat) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	data, err := tx.ConcatNonOPDataFromTxOuts()
	if err != nil {
		return nil, err
	}

	return []scanner.ITxDataSourceResult{OutputScriptsConcatResult(data)}, nil
}

func (r OutputScriptsConcatResult) SourceName() string {
	return "outputs-concatenated"
}

func (r OutputScriptsConcatResult) RawData() []byte {
	return r
}
