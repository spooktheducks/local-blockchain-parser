package dbcmds

import (
	"fmt"
	"strings"
	"time"

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

	txHash, err := blockdb.HashFromString(cmd.txHash)
	if err != nil {
		return err
	}

	txRow, blockRow, err := db.GetTxIndexRow(txHash)
	if err != nil {
		return err
	}

	tx, err := db.GetTx(txHash)
	if err != nil {
		return err
	}

	fmt.Printf("transaction %v\n", tx.Hash().String())
	fmt.Printf("  - Block %v (%v) (%v)\n", txRow.BlockHash, blockRow.DATFilename(), time.Unix(blockRow.Timestamp, 0))

	err = cmd.printOutputsSpentUnspent(db, tx)
	if err != nil {
		return err
	}

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

	return nil
}

func (cmd *TxInfoCommand) printOutputsSpentUnspent(db *blockdb.BlockDB, tx *btcutil.Tx) error {
	for txoutIdx := range tx.MsgTx().TxOut {
		key := blockdb.SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(txoutIdx)}

		spentTxOut, err := db.GetSpentTxOut(key)
		if err != nil {
			if strings.Contains(err.Error(), "can't find SpentTxOut") {
				fmt.Printf("  - TxOut %v: unspent\n", txoutIdx)
				continue
			}
			return err
		}

		fmt.Printf("  - TxOut %v: spent by %v (%v)\n", txoutIdx, spentTxOut.InputTxHash.String(), spentTxOut.TxInIndex)
	}
	return nil
}

func (cmd *TxInfoCommand) findGPGData(tx *btcutil.Tx) error {
	data, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return err
	}

	isSatoshi := false
	satoshiData, err := utils.GetSatoshiEncodedData(data)
	if err != nil {
		// ignore
	} else {
		isSatoshi = true
		data = satoshiData
	}

	result := utils.FindPGPPackets(data)
	for _, packet := range result.Packets {
		if isSatoshi {
			fmt.Printf("  - GPG packet (satoshi-encoded): %+v\n", packet)
		} else {
			fmt.Printf("  - GPG packet: %+v\n", packet)
		}
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
		matches := utils.SearchDataForMagicFileBytes(txin.SignatureScript)
		for _, m := range matches {
			fmt.Printf("  - TxIn %v magic match: %v\n", txinIdx, m.Description())
		}
	}

	// check TxOut scripts for known file headers/footers
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		matches := utils.SearchDataForMagicFileBytes(txout.PkScript)
		for _, m := range matches {
			fmt.Printf("  - TxOut %v magic match: %v\n", txoutIdx, m.Description())
		}
	}

	parsedScriptData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return err
	}

	matches := utils.SearchDataForMagicFileBytes(parsedScriptData)
	for _, m := range matches {
		fmt.Printf("  - Concatenated TxOut magic match: %v\n", m.Description())
	}

	return nil
}

func (cmd *TxInfoCommand) findPlaintext(tx *btcutil.Tx) error {
	// extract text from each TxIn scriptSig
	for txinIdx, txin := range tx.MsgTx().TxIn {
		txt := utils.StripNonTextBytes(txin.SignatureScript)
		if len(txt) < 8 {
			continue
		}

		fmt.Printf("  - TxIn %v plaintext: %v\n", txinIdx, string(txt))
	}

	// extract text from each TxOut PkScript
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		txt := utils.StripNonTextBytes(txout.PkScript)
		if len(txt) < 8 {
			continue
		}

		fmt.Printf("  - TxOut %v plaintext: %v\n", txoutIdx, string(txt))
	}

	// extract text from concatenated TxOut hex tokens

	parsedScriptData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return err
	}

	parsedScriptText := utils.StripNonTextBytes(parsedScriptData)
	if err != nil {
		return err
	}
	if len(parsedScriptText) > 8 {
		fmt.Printf("  - Concatenated TxOut plaintext: %v\n", string(parsedScriptText))
	}

	return nil
}
