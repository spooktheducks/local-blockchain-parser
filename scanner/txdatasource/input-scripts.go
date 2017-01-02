package txdatasource

import (
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type InputScripts struct{}

// ensure that InputScripts conforms to ITxDataSource
var _ scanner.ITxDataSource = &InputScripts{}

func (ds *InputScripts) Name() string {
	return "inputs"
}

func (ds *InputScripts) GetData(tx *btcutil.Tx) ([]byte, error) {
	return utils.ConcatTxInScripts(tx)
}
