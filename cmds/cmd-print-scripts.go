package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/txscript"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

func PrintBlockScripts(startBlock, endBlock uint64, inDir, outDir string) error {
	dir := filepath.Join(".", outDir, "scripts")

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	for i := int(startBlock); i < int(endBlock)+1; i++ {
		filename := fmt.Sprintf("blk%05d.dat", i)

		blocks, err := utils.LoadBlocksFromDAT(filepath.Join(inDir, filename))
		if err != nil {
			panic(err)
		}

		for _, bl := range blocks {
			blockHash := bl.Hash().String()

			f, err := utils.CreateFile(filepath.Join(dir, blockHash+".txt"))
			if err != nil {
				return err
			}
			defer utils.CloseFile(f)

			fmt.Println("===== BLOCK " + blockHash + " =====")
			f.WriteString("[BLOCK " + blockHash + "]\n")

			for _, tx := range bl.Transactions() {
				txHash := tx.Hash().String()
				fmt.Println("-   TX " + txHash)
				_, err := f.WriteString("TX: " + txHash + "\n")
				if err != nil {
					return err
				}

				for _, txout := range tx.MsgTx().TxOut {
					scriptStr, err := txscript.DisasmString(txout.PkScript)
					if err != nil {
						if err.Error() == "execute past end of script" {
							continue
						} else {
							return fmt.Errorf("error in txscript.DisasmString: %v", err)
						}
					}

					fmt.Println("        " + scriptStr)
					_, err = f.WriteString("  - " + scriptStr + "\n")
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
