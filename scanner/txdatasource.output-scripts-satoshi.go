package scanner

import (
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type OutputScriptSatoshiTxDataSource struct{}

// ensure that OutputScriptSatoshiTxDataSource conforms to ITxDataSource
var _ ITxDataSource = &OutputScriptSatoshiTxDataSource{}

func (ds *OutputScriptSatoshiTxDataSource) Name() string {
	return "outputs-satoshi"
}

func (ds *OutputScriptSatoshiTxDataSource) GetData(tx *btcutil.Tx) ([]byte, error) {
	data, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return nil, err
	}

	return utils.GetSatoshiEncodedData(data)
}
