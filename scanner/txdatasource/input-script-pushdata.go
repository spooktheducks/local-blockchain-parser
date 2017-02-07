package txdatasource

import (
	"fmt"

	// "github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type InputScriptPushdata struct{}

type InputScriptPushdataResult struct {
	rawData []byte
	index   int
}

// ensure that InputScriptPushdata conforms to ITxDataSource
var _ scanner.ITxDataSource = &InputScriptPushdata{}

// ensure that InputScriptPushdataResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = InputScriptPushdataResult{}

func (ds *InputScriptPushdata) Name() string {
	return "txin-script-pushdata"
}

func (ds *InputScriptPushdata) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for i := range tx.MsgTx().TxIn {
		data, err := tx.GetPushdataFromTxIn(i)
		if err != nil {
			return nil, err
		}
		results = append(results, InputScriptPushdataResult{rawData: data, index: i})
	}

	return results, nil
}

func (r InputScriptPushdataResult) SourceName() string {
	return fmt.Sprintf("txin-script-pushdata-%d", r.index)
}

func (r InputScriptPushdataResult) RawData() []byte {
	return r.rawData
}
