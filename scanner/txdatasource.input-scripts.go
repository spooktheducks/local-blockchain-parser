package scanner

import (
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type InputScriptTxDataSource struct{}

// ensure that InputScriptTxDataSource conforms to ITxDataSource
var _ ITxDataSource = &InputScriptTxDataSource{}

func (ds *InputScriptTxDataSource) Name() string {
	return "inputs"
}

func (ds *InputScriptTxDataSource) GetData(tx *btcutil.Tx) ([]byte, error) {
	return utils.ConcatTxInScripts(tx)
}
