package txdatasource

import (
	"fmt"

	// "github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type InputScriptNonOP struct{}

type InputScriptNonOPResult struct {
	rawData []byte
	index   int
}

// ensure that InputScriptNonOP conforms to ITxDataSource
var _ scanner.ITxDataSource = &InputScriptNonOP{}

// ensure that InputScriptNonOPResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = InputScriptNonOPResult{}

func (ds *InputScriptNonOP) Name() string {
	return "txin-script-nonop"
}

func (ds *InputScriptNonOP) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for i := range tx.MsgTx().TxIn {
		data, err := tx.GetNonOPDataFromTxIn(i)
		if err != nil {
			return nil, err
		}
		results = append(results, InputScriptNonOPResult{rawData: data, index: i})
	}

	return results, nil
}

func (r InputScriptNonOPResult) SourceName() string {
	return fmt.Sprintf("txin-script-nonop-%d", r.index)
}

func (r InputScriptNonOPResult) RawData() []byte {
	return r.rawData
}
