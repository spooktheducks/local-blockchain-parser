package txdatasource

import (
	"fmt"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type InputScriptFirstPushdata struct{}

type InputScriptFirstPushdataResult struct {
	rawData []byte
	index   int
}

// ensure that InputScriptFirstPushdata conforms to ITxDataSource
var _ scanner.ITxDataSource = &InputScriptFirstPushdata{}

// ensure that InputScriptFirstPushdataResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = InputScriptFirstPushdataResult{}

func (ds *InputScriptFirstPushdata) Name() string {
	return "txin-script-first-pushdata"
}

func (ds *InputScriptFirstPushdata) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for i, txin := range tx.MsgTx().TxIn {
		data, err := utils.GetFirstPushdataBytes(txin.SignatureScript)
		if err != nil {
			return nil, err
		}
		results = append(results, InputScriptFirstPushdataResult{rawData: data, index: i})
	}

	return results, nil
}

func (r InputScriptFirstPushdataResult) SourceName() string {
	return fmt.Sprintf("txin-script-pushdata-%d", r.index)
}

func (r InputScriptFirstPushdataResult) RawData() []byte {
	return r.rawData
}
