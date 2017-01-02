package txdatasource

import (
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type OutputScriptsSatoshi struct{}

// ensure that OutputScriptsSatoshi conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScriptsSatoshi{}

func (ds *OutputScriptsSatoshi) Name() string {
	return "outputs-satoshi"
}

func (ds *OutputScriptsSatoshi) GetData(tx *btcutil.Tx) ([]byte, error) {
	return utils.ConcatSatoshiDataFromTxOuts(tx)
}
