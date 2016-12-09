package cmds

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

func PrintBlockScripts(bl *btcutil.Block, outDir string) error {
	dir := filepath.Join(".", outDir, "scripts")

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	blockHash := bl.Hash().String()

	f, err := os.Create(filepath.Join(dir, blockHash+".txt"))
	if err != nil {
		return err
	}
	defer f.Close()

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

	return nil
}

func PrintBlockScriptsOpReturns(bl *btcutil.Block, outDir string) error {
	dir := filepath.Join(".", outDir, "op-returns")

	// err := os.RemoveAll(dir)
	// if err != nil {
	//  return err
	// }

	blockDir := filepath.Join(dir, bl.Hash().String())

	err := os.MkdirAll(blockDir, 0777)
	if err != nil {
		return err
	}

	for _, tx := range bl.Transactions() {
		// for _, txin := range tx.TxIns {
		//  data, err := blkparser.Pkscript(txin.ScriptSig).DecodeToString()
		//  if err != nil {
		//      return err
		//  }
		//  fmt.Println(data)
		// }

		for txoutIdx, txout := range tx.MsgTx().TxOut {
			scriptStr, err := txscript.DisasmString(txout.PkScript)
			if err != nil {
				if err.Error() == "execute past end of script" {
					continue
				} else {
					return fmt.Errorf("error in txscript.DisasmString: %v", err)
				}
			}

			data, err := getOpReturnBytes(scriptStr)
			if err != nil {
				if err.Error() == "encoding/hex: odd length hex string" {
					continue
				} else {
					return fmt.Errorf("error in getOpReturnBytes: %v", err)
				}
			} else if data == nil {
				continue
			}

			f, err := os.Create(filepath.Join(blockDir, fmt.Sprintf("%v-%v.dat", tx.Hash().String(), txoutIdx)))
			if err != nil {
				return err
			}
			// defer f.Close()

			_, err = f.Write(data)
			if err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	return nil
}

func getOpReturnBytes(scriptStr string) ([]byte, error) {
	toks := strings.Split(scriptStr, " ")

	for i := range toks {
		if toks[i] == "OP_RETURN" {
			if len(toks) >= i+2 {
				return hex.DecodeString(toks[i+1])
			} else {
				return nil, errors.New("empty OP_RETURN data")
			}
		}
	}

	return nil, nil
}

func PrintBlockData(bl *btcutil.Block) error {
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

	return nil
}
