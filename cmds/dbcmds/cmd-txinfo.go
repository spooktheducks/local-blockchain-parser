package dbcmds

import (
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
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
	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	txHash, err := utils.HashFromString(cmd.txHash)
	if err != nil {
		return err
	}

	tx, err := db.GetTx(txHash)
	if err != nil {
		return err
	}

	fmt.Printf("transaction %v\n", tx.Hash().String())
	fmt.Printf("  - Block %v (%v) (%v)\n", tx.BlockHash, tx.DATFilename(), time.Unix(tx.BlockTimestamp, 0))
	fmt.Printf("  - Lock time: %v\n", tx.MsgTx().LockTime)

	fee, err := tx.Fee()
	if err != nil {
		return err
	}
	fmt.Printf("  - Fee: %v BTC\n", fee)

	// txoutAddrs, err := utils.GetTxOutAddresses(tx)
	// if err != nil {
	// 	return err
	// }

	// for txoutIdx, addrs := range txoutAddrs {
	// 	if len(addrs) == 0 {
	// 		fmt.Printf("  - TxOut %v: can't decode address\n", txoutIdx)
	// 	} else {
	// 		addrStrings := make([]string, len(addrs))
	// 		for i := range addrs {
	// 			addrStrings[i] = addrs[i].String()
	// 		}
	// 		fmt.Printf("  - TxOut %v: paid to %v\n", txoutIdx, strings.Join(addrStrings, ", "))
	// 	}
	// }

	// err = cmd.printOutputsSpentUnspent(db, tx)
	// if err != nil {
	// 	return err
	// }

	// err = cmd.findPlaintext(tx)
	// if err != nil {
	// 	return err
	// }

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

	err = os.MkdirAll("output/tx-info/"+tx.Hash().String(), 0777)
	if err != nil {
		return err
	}

	f, err := os.Create("output/tx-info/" + tx.Hash().String() + "/txout-data.dat")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := tx.ConcatNonOPDataFromTxOuts()
	if err != nil {
		return err
	}
	f.Write(data)

	sf, err := os.Create("output/tx-info/" + tx.Hash().String() + "/satoshi-txout-data.dat")
	if err != nil {
		return err
	}
	defer sf.Close()

	sd, err := utils.GetSatoshiEncodedData(data)
	if err == nil {
		sf.Write(sd)
	}

	inf, err := os.Create("output/tx-info/" + tx.Hash().String() + "/txin-data.dat")
	if err != nil {
		return err
	}
	defer sf.Close()

	data, err = tx.ConcatTxInScripts()
	if err == nil {
		inf.Write(data)
	}

	return nil
}

func (cmd *TxInfoCommand) printOutputsSpentUnspent(db *BlockDB, tx *Tx) error {
	for txoutIdx := range tx.MsgTx().TxOut {
		addr, err := tx.GetTxOutAddress(txoutIdx)
		if err != nil {
			return err
		}

		addrString := ""
		if len(addr) == 0 {
			addrString = "unable to decode output address"
		} else {
			addrString = fmt.Sprintf("%v", addr)
		}

		spentString := ""
		spentTxOut, err := db.GetSpentTxOut(SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(txoutIdx)})
		if err != nil {
			if strings.Contains(err.Error(), "can't find SpentTxOut") {
				spentString = "unspent"
			} else {
				return err
			}
		} else {
			spentString = fmt.Sprintf("spent in tx %v (%v)", spentTxOut.InputTxHash.String(), spentTxOut.TxInIndex)
		}

		fmt.Printf("  - TxOut %v: %v (addr: %v)\n", txoutIdx, spentString, addrString)
	}
	return nil
}

func (cmd *TxInfoCommand) findGPGData(tx *Tx) error {
	data, err := tx.ConcatNonOPDataFromTxOuts()
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

func (cmd *TxInfoCommand) findSatoshiEncodedData(tx *Tx) error {
	data, err := tx.ConcatNonOPDataFromTxOuts()
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

func (cmd *TxInfoCommand) findFileHeaders(tx *Tx) error {
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

	parsedScriptData, err := tx.ConcatNonOPDataFromTxOuts()
	if err != nil {
		return err
	}

	matches := utils.SearchDataForMagicFileBytes(parsedScriptData)
	for _, m := range matches {
		fmt.Printf("  - Concatenated TxOut magic match: %v\n", m.Description())
	}

	return nil
}

func (cmd *TxInfoCommand) findPlaintext(tx *Tx) error {
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

	parsedScriptData, err := tx.ConcatNonOPDataFromTxOuts()
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
