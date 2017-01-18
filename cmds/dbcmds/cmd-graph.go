package dbcmds

import (
	"fmt"
	"os"
	"path/filepath"

	// "github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/types"
)

type GraphCommand struct {
	dbFile     string
	datFileDir string
	walletAddr string
	outDir     string
	db         *blockdb.BlockDB
}

type graph struct {
	addrs   []string
	txs     []string
	entries map[string][]string
}

func NewGraphCommand(datFileDir, dbFile, outDir, walletAddr string) *GraphCommand {
	return &GraphCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		walletAddr: walletAddr,
		outDir:     filepath.Join(outDir, "graph", walletAddr),
	}
}

func (cmd *GraphCommand) RunCommand() error {
	err := os.MkdirAll(cmd.outDir, 0777)
	if err != nil {
		return err
	}

	db, err := blockdb.NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	cmd.db = db

	txHash, err := blockdb.HashFromString("b0800606ee9f5e73868ed6c61b55b802a7454abf06a0d69a7fcabe4904afb665")
	if err != nil {
		return err
	}
	tx, err := db.GetTx(txHash)
	if err != nil {
		return err
	}

	g := &graph{entries: map[string][]string{}}
	err = cmd.graphTx(g, tx)
	if err != nil {
		return err
	}

	err = cmd.writeGraphFile(g)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *GraphCommand) graphTx(g *graph, tx *Tx) error {
	g.txs = append(g.txs, tx.Hash().String())

	for txoutIdx := range tx.MsgTx().TxOut {
		txoutAddrs, err := tx.GetTxOutAddress(txoutIdx)
		if err != nil {
			return err
		} else if len(txoutAddrs) == 0 {
			continue
		} else if len(txoutAddrs) > 1 {
			fmt.Println("weird txout (multiple addresses):", txoutAddrs, "(only using first address)")
		}
		outAddr := txoutAddrs[0].String()

		g.addrs = append(g.addrs, outAddr)
		g.entries[tx.Hash().String()] = append(g.entries[tx.Hash().String()], outAddr)

		spentTxOut, err := cmd.db.GetSpentTxOut(blockdb.SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(txoutIdx)})
		if err == nil {
			g.txs = append(g.txs, spentTxOut.InputTxHash.String())
			g.entries[txoutAddrs[0].String()] = append(g.entries[txoutAddrs[0].String()], spentTxOut.InputTxHash.String())
		}
	}

	return nil
}

func (cmd *GraphCommand) writeGraphFile(g *graph) error {
	f, err := os.Create(filepath.Join(cmd.outDir, "graph.dot"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("digraph tx {\n")
	if err != nil {
		return err
	}

	for _, addr := range g.addrs {
		_, err := f.WriteString(fmt.Sprintf("    \"%s\" [color=blue];\n", addr))
		if err != nil {
			return err
		}
	}

	for _, tx := range g.txs {
		_, err := f.WriteString(fmt.Sprintf("    \"%s\" [color=red];\n", tx))
		if err != nil {
			return err
		}
	}

	for k, v := range g.entries {
		for _, x := range v {
			_, err := f.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", k, x))
			if err != nil {
				return err
			}
		}
	}

	_, err = f.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}
