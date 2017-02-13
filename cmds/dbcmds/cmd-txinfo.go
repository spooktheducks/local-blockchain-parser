package dbcmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	// "sort"
	"strings"
	"time"

	// "github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	// "github.com/spooktheducks/local-blockchain-parser/scanner"
	// "github.com/spooktheducks/local-blockchain-parser/scanner/detector"
	// "github.com/spooktheducks/local-blockchain-parser/scanner/detectoroutput"
	// "github.com/spooktheducks/local-blockchain-parser/scanner/txdatasource"
	// "github.com/spooktheducks/local-blockchain-parser/scanner/txdatasourceoutput"
	// "github.com/spooktheducks/local-blockchain-parser/scanner/txhashsource"
)

type TxInfoCommand struct {
	dbFile     string
	datFileDir string
	outDir     string
	txHash     string
}

func NewTxInfoCommand(datFileDir, dbFile, outDir, txHash string) *TxInfoCommand {
	return &TxInfoCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		outDir:     filepath.Join(outDir, "tx-info", txHash),
		txHash:     txHash,
	}
}

type sortableTxOut struct {
	Index int
	Value uint64
}

type sortableTxOuts []sortableTxOut

func (s sortableTxOuts) Len() int {
	return len(s)
}

func (s sortableTxOuts) Less(i, j int) bool {
	return s[i].Value < s[j].Value
}

func (s sortableTxOuts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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

	// {
	// 	// txList := []string{
	// 	// 	"5970ae129d1141663bd5e441a1555c16fb1c0586dd05f40c1db3d3e81218ee41",
	// 	// 	"6e5c13edf8bb594ad850882173c9b5e906187269595c0154db36668d337b42e1",
	// 	// }

	// 	allData := []byte{}
	// 	hashStr := "5970ae129d1141663bd5e441a1555c16fb1c0586dd05f40c1db3d3e81218ee41"
	// 	h, _ := utils.HashFromString(hashStr)
	// 	tx, err := db.GetTx(h)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	i := 0
	// 	for { //_, hashStr := range txList {
	// 		if i > 20 {
	// 			break
	// 		}
	// 		sortedTxOuts := sortableTxOuts{}
	// 		for i, txout := range tx.MsgTx().TxOut {
	// 			sortedTxOuts = append(sortedTxOuts, sortableTxOut{Index: i, Value: uint64(txout.Value)})
	// 		}
	// 		sort.Sort(sortedTxOuts)

	// 		for _, txout := range sortedTxOuts {
	// 			spendingTx, err := tx.GetSpendingTx(txout.Index)
	// 			if err == nil || spendingTx != nil {
	// 				continue
	// 			}

	// 			data, err := tx.GetNonOPDataFromTxOut(txout.Index)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			allData = append(allData, data...)
	// 		}

	// 		// find next tx
	// 		for i := range tx.MsgTx().TxOut {
	// 			spendingTx, err := tx.GetSpendingTx(i)
	// 			if err != nil || spendingTx == nil {
	// 				continue
	// 			}
	// 			tx = spendingTx
	// 			break
	// 		}
	// 		fmt.Println("next:", tx.Hash().String())
	// 		i++
	// 	}

	// 	fmt.Println(string(allData))
	// 	panic("done")
	// }

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

	err = os.MkdirAll(cmd.outDir, 0777)
	if err != nil {
		return err
	}

	// s := &scanner.Scanner{
	// 	DB:           db,
	// 	TxHashSource: txhashsource.NewListTxHashSource([]chainhash.Hash{*tx.Hash()}),
	// 	TxDataSources: []scanner.ITxDataSource{
	// 		&txdatasource.InputScript{},
	// 		// &txdatasource.InputScriptsConcat{},
	// 		&txdatasource.OutputScript{},
	// 		// &txdatasource.OutputScriptsConcat{},
	// 		&txdatasource.OutputScriptsSatoshi{},
	// 		&txdatasource.OutputScriptOpReturn{},
	// 	},
	// 	TxDataSourceOutputs: []scanner.ITxDataSourceOutput{
	// 		&txdatasourceoutput.RawData{OutDir: cmd.outDir},
	// 		&txdatasourceoutput.RawDataEachDataSource{OutDir: cmd.outDir},
	// 	},
	// 	Detectors: []scanner.IDetector{
	// 		// &detector.PGPPackets{},
	// 		// &detector.AESKeys{},
	// 		&detector.MagicBytes{},
	// 		// &detector.Plaintext{},
	// 	},
	// 	DetectorOutputs: []scanner.IDetectorOutput{
	// 		&detectoroutput.Console{Prefix: "  - "},
	// 		&detectoroutput.RawData{OutDir: cmd.outDir},
	// 		&detectoroutput.CSV{OutDir: cmd.outDir},
	// 		&detectoroutput.CSVTxAnalysis{OutDir: cmd.outDir},
	// 	},
	// }

	// err = s.Run()
	// if err != nil {
	// 	return err
	// }

	// return s.Close()

	err = cmd.findFileHeaders(tx)
	if err != nil {
		return err
	}

	err = cmd.findSatoshiEncodedData(tx)
	if err != nil {
		return err
	}

	// err = cmd.findGPGData(tx)
	// if err != nil {
	// 	return err
	// }

	err = os.MkdirAll("output/tx-info/"+tx.Hash().String(), 0777)
	if err != nil {
		return err
	}

	// write individual txouts
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("txout-%v.dat", txoutIdx)), txout.PkScript, 0666)
		if err != nil {
			return err
		}
	}

	// write concatenated txouts
	data, err := tx.ConcatNonOPDataFromTxOuts()
	if err != nil {
		return err
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(cmd.outDir, "txout-data.dat"), data, 0666)
		if err != nil {
			return err
		}
	}

	// write satoshi-encoded txouts
	sd, err := utils.GetSatoshiEncodedData(data)
	if err == nil {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, "satoshi-txout-data.dat"), sd, 0666)
		if err != nil {
			return err
		}
	}

	// write individual txins
	for txinIdx, txin := range tx.MsgTx().TxIn {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("txin-%v.dat", txinIdx)), txin.SignatureScript, 0666)
		if err != nil {
			return err
		}

		// write txin non-OP data
		nonOPData, err := utils.GetNonOPBytesFromInputScript(txin.SignatureScript)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("txin-nonop-%v.dat", txinIdx)), nonOPData, 0666)
		if err != nil {
			return err
		}

		// write txin OP_PUSHDATA data
		pushdata, err := utils.GetPushdataBytesFromInputScript(txin.SignatureScript)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("txin-pushdata-%v.dat", txinIdx)), pushdata, 0666)
		if err != nil {
			return err
		}
	}

	// write concatenated txins
	data, err = tx.ConcatTxInScripts()
	if err == nil {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, "txin-data.dat"), data, 0666)
		if err != nil {
			return err
		}
	}

	// write concatenated txin non-OP data
	data, err = tx.ConcatNonOPDataFromTxIns()
	if err == nil {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, "txin-nonop-concat.dat"), data, 0666)
		if err != nil {
			return err
		}
	}

	// write concatenated txin OP_PUSHDATA data
	data, err = tx.ConcatPushdataFromTxIns()
	if err == nil {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, "txin-pushdata-concat.dat"), data, 0666)
		if err != nil {
			return err
		}
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
		fmt.Println(string(txin.SignatureScript))
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
