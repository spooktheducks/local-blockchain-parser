package scanner

import (
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type OutputScriptTxDataSource struct{}

// ensure that OutputScriptTxDataSource conforms to ITxDataSource
var _ ITxDataSource = &OutputScriptTxDataSource{}

func (ds *OutputScriptTxDataSource) Name() string {
	return "outputs"
}

func (ds *OutputScriptTxDataSource) GetData(tx *btcutil.Tx) ([]byte, error) {
	return utils.ConcatNonOPHexTokensFromTxOuts(tx)
}
