package txdatasource

import (
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/types"
)

type InputScriptsConcat struct{}

type InputScriptsConcatResult []byte

// ensure that InputScriptsConcat conforms to ITxDataSource
var _ scanner.ITxDataSource = &InputScriptsConcat{}

// ensure that InputScriptsConcatResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = &InputScriptsConcatResult{}

func (ds *InputScriptsConcat) Name() string {
	return "inputs-concatenated"
}

func (ds *InputScriptsConcat) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	data, err := tx.ConcatTxInScripts()
	if err != nil {
		return nil, err
	}

	return []scanner.ITxDataSourceResult{InputScriptsConcatResult(data)}, nil
}

func (r InputScriptsConcatResult) SourceName() string {
	return "inputs-concatenated"
}

func (r InputScriptsConcatResult) RawData() []byte {
	return r
}
