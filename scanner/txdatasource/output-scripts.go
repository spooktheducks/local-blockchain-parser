package txdatasource

import (
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type OutputScripts struct{}

// ensure that OutputScripts conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScripts{}

func (ds *OutputScripts) Name() string {
	return "outputs"
}

func (ds *OutputScripts) GetData(tx *btcutil.Tx) ([]byte, error) {
	return utils.ConcatNonOPHexTokensFromTxOuts(tx)
}
