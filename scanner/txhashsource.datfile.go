package scanner

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

func NewDATFileTxHashSource(datFileDir string, startBlock, endBlock int) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		for datIdx := startBlock; datIdx <= endBlock; datIdx++ {
			blocks, err := utils.LoadBlocksFromDAT(filepath.Join(datFileDir, fmt.Sprintf("blk%05d.dat")))
			if err != nil {
				// @@TODO
				panic(err)
			}

			for _, bl := range blocks {
				for _, tx := range bl.Transactions() {
					ch <- *tx.Hash()
				}
			}
		}
	}()

	return TxHashSource(ch)
}
