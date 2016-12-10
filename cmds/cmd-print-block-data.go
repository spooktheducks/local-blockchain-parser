package cmds

import (
	"fmt"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/utils"
)

func PrintBlockData(startBlock, endBlock uint64, inDir, outDir string) error {
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		filename := fmt.Sprintf("blk%05d.dat", i)

		blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
		if err != nil {
			panic(err)
		}

		for _, bl := range blocks {

			// Basic block info
			fmt.Printf("Block hash: %v\n", bl.Hash().String())
			fmt.Printf("Block time: %v\n", bl.MsgBlock().Header.Timestamp)
			fmt.Printf("Block version: %v\n", bl.MsgBlock().Header.Version)
			// fmt.Printf("Block parent: %v\n", btc.NewUint256(bl.ParentHash()).String())
			fmt.Printf("Block merkle root: %v\n", bl.MsgBlock().Header.MerkleRoot.String())
			fmt.Printf("Block bits: %v\n", bl.MsgBlock().Header.Bits)
			// fmt.Printf("Block size: %v\n", len(bl.Raw))

			// Fetch TXs and iterate over them
			for _, tx := range bl.Transactions() {
				fmt.Printf("TxId: %v\n", tx.Hash().String())
				// fmt.Printf("Tx Size: %v\n", tx.Size)
				// fmt.Printf("Tx Lock time: %v\n", tx.LockTime)
				// fmt.Printf("Tx Version: %v\n", tx.Version)

				fmt.Println("TxIns:")

				// if tx.IsCoinBase() {
				//  fmt.Printf("TxIn coinbase, newly generated coins")
				// } else {
				for txin_index := range tx.MsgTx().TxIn {
					fmt.Printf("TxIn index: %v\n", txin_index)
					// fmt.Printf("TxIn Input hash: %v\n", txin.InputHash)
					// fmt.Printf("TxIn Input vout: %v\n", txin.InputVout)
					// fmt.Printf("TxIn ScriptSig: %v\n", hex.EncodeToString(txin.ScriptSig))
					// fmt.Printf("TxIn Sequence: %v\n", txin.Sequence)
				}
				// }

				fmt.Println("TxOuts:")

				for txo_index := range tx.MsgTx().TxOut {
					fmt.Printf("TxOut index: %v\n", txo_index)
					// txout.PkScript
					// fmt.Printf("TxOut value: %v\n", txout.Value)
					// fmt.Printf("TxOut script: %s\n", hex.EncodeToString(txout.Pkscript))
					// fmt.Printf("TxOut address: %v\n", txout.Addr)
				}
			}
		}
	}

	return nil
}
