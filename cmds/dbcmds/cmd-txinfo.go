package dbcmds

import (
	"fmt"

	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type TxInfoCommand struct {
	dbFile     string
	datFileDir string
	txHash     string
}

func NewTxInfoCommand(datFileDir, dbFile, txHash string) *TxInfoCommand {
	return &TxInfoCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		txHash:     txHash,
	}
}

func (cmd *TxInfoCommand) RunCommand() error {
	db, err := blockdb.NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.GetTx(cmd.txHash)
	if err != nil {
		return err
	}

	fmt.Printf("transaction %v\n", tx.Hash().String())

	err = cmd.findPlaintext(tx)
	if err != nil {
		return err
	}

	err = cmd.findFileHeaders(tx)
	if err != nil {
		return err
	}

	err = cmd.findSatoshiEncodedData(tx)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *TxInfoCommand) findSatoshiEncodedData(tx *btcutil.Tx) error {
	// data := []byte{}
	// // we skip the final two TxOuts because one goes to WL and one is used to pass BTC to the next transaction in the chain
	// for i := 0; i < len(tx.MsgTx().TxOut)-2; i++ {
	// 	bs, err := utils.GetNonOPBytes(tx.MsgTx().TxOut[i].PkScript)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	data = append(data, bs...)
	// }
	data, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return err
	}

	_, err = utils.GetSatoshiEncodedData(data)
	if err != nil {
		return nil
	}

	fmt.Printf("  - TxOut Satoshi-encoded data found\n")

	return nil
}

func (cmd *TxInfoCommand) findFileHeaders(tx *btcutil.Tx) error {
	// check TxIn scripts for known file headers/footers
	for txinIdx, txin := range tx.MsgTx().TxIn {
		matches := utils.SearchDataForKnownFileBits(txin.SignatureScript)
		for _, m := range matches {
			fmt.Printf("  - TxIn %v magic match: %v\n", txinIdx, m.Description())
		}
	}

	// check TxOut scripts for known file headers/footers
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		matches := utils.SearchDataForKnownFileBits(txout.PkScript)
		for _, m := range matches {
			fmt.Printf("  - TxOut %v magic match: %v\n", txoutIdx, m.Description())
		}
	}

	parsedScriptData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return err
	}

	matches := utils.SearchDataForKnownFileBits(parsedScriptData)
	for _, m := range matches {
		fmt.Printf("  - Concatenated TxOut magic match: %v\n", m.Description())
	}

	return nil
}

func (cmd *TxInfoCommand) findPlaintext(tx *btcutil.Tx) error {
	// extract text from each TxIn scriptSig
	for txinIdx, txin := range tx.MsgTx().TxIn {
		txt, isText := utils.ExtractText(txin.SignatureScript)
		if !isText || len(txt) < 8 {
			continue
		}

		fmt.Printf("  - TxIn %v plaintext: %v\n", txinIdx, string(txt))
	}

	// extract text from each TxOut PkScript
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		txt, isText := utils.ExtractText(txout.PkScript)
		if !isText || len(txt) < 8 {
			continue
		}

		fmt.Printf("  - TxOut %v plaintext: %v\n", txoutIdx, string(txt))
	}

	// extract text from concatenated TxOut hex tokens

	parsedScriptData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return err
	}

	parsedScriptText, isText := utils.ExtractText(parsedScriptData)
	if err != nil {
		return err
	}
	if isText {
		fmt.Printf("  - Concatenated TxOut plaintext: %v\n", string(parsedScriptText))
	}

	return nil
}
