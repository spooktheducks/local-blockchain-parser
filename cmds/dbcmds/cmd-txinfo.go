package dbcmds

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcutil"
	// "golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"

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

	// fmt.Println("TxOuts:")
	// for txoutIdx, txout := range tx.MsgTx().TxOut {
	// 	fmt.Printf("\n")
	// }

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

	err = cmd.findGPGData(tx)
	if err != nil {
		return err
	}

	// err = cmd.blah(db)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (cmd *TxInfoCommand) blah(db *blockdb.BlockDB) error {
	fmt.Println("blah")
	txs := []string{
		// "7379ab5047b143c0b6cfe5d8d79ad240b4b4f8cced55aa26f86d1d3d370c0d4c",
		// "d3c1cb2cdbf07c25e3c5f513de5ee36081a7c590e621f1f1eab62e8d4b50b635",
		"cce82f3bde0537f82a55f3b8458cb50d632977f85c81dad3e1983a3348638f5c",
	}

	allData := []byte{}
	for _, txHash := range txs {
		tx, err := db.GetTx(txHash)
		if err != nil {
			return err
		}

		data, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
		if err != nil {
			return err
		}

		data, err = utils.GetSatoshiEncodedData(data)
		if err != nil {
			return err
		}

		allData = append(allData, data...)
	}

	reader := packet.NewReader(bytes.NewReader(allData))
	for {
		packet, err := reader.Next()
		if err != nil {
			return err
		}
		fmt.Printf("  - GPG packet: %+v\n", packet)
	}

	return nil
}

func (cmd *TxInfoCommand) findGPGData(tx *btcutil.Tx) error {
	data, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return err
	}

	data, err = utils.GetSatoshiEncodedData(data)
	if err != nil {
		return nil
	}

	reader := packet.NewReader(bytes.NewReader(data))
	for {
		packet, err := reader.Next()
		if err != nil {
			// return err
			break
		}
		fmt.Printf("  - GPG packet: %+v\n", packet)
	}

	return nil
}

func (cmd *TxInfoCommand) findSatoshiEncodedData(tx *btcutil.Tx) error {
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
