package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/lib/blkparser"
)

var (
	flagInDir      = flag.String("inDir", "", "The .dat file containing blockchain input data")
	flagStartBlock = flag.Int64("startBlock", 0, "The block number to start from")
	flagEndBlock   = flag.Int64("endBlock", 0, "The block number to end on")
	flagOutDir     = flag.String("outDir", "output", "Output directory")
)

type (
	Cmd int
)

const (
	CmdBlockData Cmd = iota
	CmdScripts
	CmdOpReturns
)

func getCmd(arg string) Cmd {
	switch arg {
	case "blockdata":
		return CmdBlockData
	case "scripts":
		return CmdScripts
	case "opreturns":
		return CmdOpReturns

	case "":
		fallthrough
	default:
		panic("Must specify a command (blockdata, scripts, or opreturns)")
	}
}

func main() {
	flag.Parse()

	if *flagInDir == "" {
		panic("Missing --inDir param")
	} else if *flagEndBlock == 0 {
		panic("Must specify --endBlock param")
	}

	cmd := getCmd(flag.Arg(0))

	startBlock := uint64(*flagStartBlock)
	endBlock := uint64(*flagEndBlock)

	// Set real Bitcoin network
	magic := [4]byte{0xF9, 0xBE, 0xB4, 0xD9}

	// Specify blocks directory
	blockDB, err := blkparser.NewBlockchain(*flagInDir, magic, uint32(startBlock))
	if err != nil {
		panic("error opening file: " + err.Error())
	}

	for i := int(startBlock); i < int(endBlock)+1; i++ {
		dat, err := blockDB.FetchNextBlock()
		if dat == nil || err != nil {
			fmt.Println("END of DB file")
			break
		}
		bl, err := blkparser.NewBlock(dat[:])

		if err != nil {
			fmt.Println("Block inconsistent:", err.Error())
			break
		}

		// Read block till we reach startBlock
		if uint64(i) < startBlock {
			continue
		}

		switch cmd {
		case CmdBlockData:
			err = printBlockData(bl)
		case CmdScripts:
			err = printBlockScripts(bl)
		case CmdOpReturns:
			err = printBlockScriptsOpReturns(bl)
		}

		if err != nil {
			panic(err)
		}
	}
}

func printBlockScripts(bl *blkparser.Block) error {
	dir := filepath.Join(".", *flagOutDir, "scripts")

	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dir, bl.Hash+".txt"))
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Println("===== BLOCK " + bl.Hash + " =====")
	f.WriteString("[BLOCK " + bl.Hash + "]\n")

	for _, tx := range bl.Txs {
		fmt.Println("-   TX " + tx.Hash)
		_, err := f.WriteString("TX: " + tx.Hash + "\n")
		if err != nil {
			return err
		}

		for _, txout := range tx.TxOuts {
			scriptStr, err := txout.Pkscript.DecodeToString()
			if err != nil {
				return err
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

func printBlockScriptsOpReturns(bl *blkparser.Block) error {
	dir := filepath.Join(".", *flagOutDir, "op-returns")

	// err := os.RemoveAll(dir)
	// if err != nil {
	// 	return err
	// }

	for _, tx := range bl.Txs {
		blockDir := filepath.Join(dir, bl.Hash)

		err := os.MkdirAll(blockDir, 0777)
		if err != nil {
			return err
		}

		for txoutIdx, txout := range tx.TxOuts {
			data, err := txout.Pkscript.GetOpReturnBytes()
			if err != nil {
				return err
			} else if data == nil {
				continue
			}

			f, err := os.Create(filepath.Join(blockDir, fmt.Sprintf("%v-%v.dat", tx.Hash, txoutIdx)))
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

func printBlockData(bl *blkparser.Block) error {
	// Basic block info
	fmt.Printf("Block hash: %v\n", bl.Hash)
	fmt.Printf("Block time: %v\n", bl.BlockTime)
	fmt.Printf("Block version: %v\n", bl.Version)
	// fmt.Printf("Block parent: %v\n", btc.NewUint256(bl.ParentHash()).String())
	fmt.Printf("Block merkle root: %v\n", bl.MerkleRoot)
	fmt.Printf("Block bits: %v\n", bl.Bits)
	fmt.Printf("Block size: %v\n", len(bl.Raw))

	// Fetch TXs and iterate over them
	for _, tx := range bl.Txs {
		fmt.Printf("TxId: %v\n", tx.Hash)
		fmt.Printf("Tx Size: %v\n", tx.Size)
		fmt.Printf("Tx Lock time: %v\n", tx.LockTime)
		fmt.Printf("Tx Version: %v\n", tx.Version)

		fmt.Println("TxIns:")

		// if tx.IsCoinBase() {
		//  fmt.Printf("TxIn coinbase, newly generated coins")
		// } else {
		for txin_index, txin := range tx.TxIns {
			fmt.Printf("TxIn index: %v\n", txin_index)
			fmt.Printf("TxIn Input hash: %v\n", txin.InputHash)
			fmt.Printf("TxIn Input vout: %v\n", txin.InputVout)
			fmt.Printf("TxIn ScriptSig: %v\n", hex.EncodeToString(txin.ScriptSig))
			fmt.Printf("TxIn Sequence: %v\n", txin.Sequence)
		}
		// }

		fmt.Println("TxOuts:")

		for txo_index, txout := range tx.TxOuts {
			fmt.Printf("TxOut index: %v\n", txo_index)
			fmt.Printf("TxOut value: %v\n", txout.Value)
			fmt.Printf("TxOut script: %s\n", hex.EncodeToString(txout.Pkscript))
			fmt.Printf("TxOut address: %v\n", txout.Addr)
		}
	}

	return nil
}
