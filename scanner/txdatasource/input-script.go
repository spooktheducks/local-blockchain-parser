package txdatasource

import (
	"fmt"

	// "github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type InputScript struct{}

type InputScriptResult struct {
	rawData []byte
	index   int
}

// ensure that InputScript conforms to ITxDataSource
var _ scanner.ITxDataSource = &InputScript{}

// ensure that InputScriptResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = InputScriptResult{}

func (ds *InputScript) Name() string {
	return "txin-script"
}

func (ds *InputScript) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	results := []scanner.ITxDataSourceResult{}
	for txinIdx, txin := range tx.MsgTx().TxIn {
		results = append(results, InputScriptResult{rawData: txin.SignatureScript, index: txinIdx})
	}

	return results, nil
}

func (r InputScriptResult) SourceName() string {
	return fmt.Sprintf("txin-script-%d", r.index)
}

func (r InputScriptResult) RawData() []byte {
	return r.rawData
}

func (r InputScriptResult) InOut() string {
	return "in"
}

func (r InputScriptResult) Index() int {
	return r.index
}
