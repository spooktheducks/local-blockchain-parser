package blockdb

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"

	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

type Tx struct {
	*btcutil.Tx

	db *BlockDB

	DATFileIdx          uint16
	BlockTimestamp      int64
	BlockIndexInDATFile uint32

	BlockHash    chainhash.Hash
	IndexInBlock uint64
}

func (tx *Tx) DATFilename() string {
	return fmt.Sprintf("blk%05d.dat", tx.DATFileIdx)
}

func (tx *Tx) GetOPReturnDataFromTxOut(txoutIdx int) ([]byte, error) {
	return utils.GetOPReturnBytes(tx.MsgTx().TxOut[txoutIdx].PkScript)
}

func (tx *Tx) ConcatOPReturnDataFromTxOuts() ([]byte, error) {
	allBytes := []byte{}

	for _, txout := range tx.MsgTx().TxOut {
		bs, err := utils.GetOPReturnBytes(txout.PkScript)
		if err != nil {
			continue
		}

		allBytes = append(allBytes, bs...)
	}

	return allBytes, nil
}

func (tx *Tx) GetNonOPDataFromTxOut(txoutIdx int) ([]byte, error) {
	return utils.GetNonOPBytes(tx.MsgTx().TxOut[txoutIdx].PkScript)
}

// func (tx *Tx) IsSpent(txoutIdx int) (bool, error) {
// 	spentTxOut, err := tx.db.GetSpentTxOut(SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(txoutIdx)})
//     if err != nil {
//         return false, err
//     }
//     spentTxOut.
// }

func (tx *Tx) ConcatNonOPDataFromTxOuts() ([]byte, error) {
	allBytes := []byte{}

	for _, txout := range tx.MsgTx().TxOut {
		bs, err := utils.GetNonOPBytes(txout.PkScript)
		if err != nil {
			continue
		}

		allBytes = append(allBytes, bs...)
	}

	return allBytes, nil
}

func (tx *Tx) ConcatSatoshiDataFromTxOuts() ([]byte, error) {
	data, err := tx.ConcatNonOPDataFromTxOuts()
	if err != nil {
		return nil, err
	}

	return utils.GetSatoshiEncodedData(data)
}

func (tx *Tx) ConcatTxInScripts() ([]byte, error) {
	allBytes := []byte{}

	for _, txin := range tx.MsgTx().TxIn {
		allBytes = append(allBytes, txin.SignatureScript...)
	}

	return allBytes, nil
}

func (tx *Tx) GetTxOutAddress(txoutIdx int) ([]btcutil.Address, error) {
	txout := tx.MsgTx().TxOut[txoutIdx]

	_, addresses, _, err := txscript.ExtractPkScriptAddrs(txout.PkScript, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (tx *Tx) GetTxOutAddresses() ([][]btcutil.Address, error) {
	addrs := make([][]btcutil.Address, len(tx.MsgTx().TxOut))

	for i, txout := range tx.MsgTx().TxOut {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(txout.PkScript, &chaincfg.MainNetParams)
		if err != nil {
			return nil, err
		}
		addrs[i] = addresses
	}

	return addrs, nil
}

func (tx *Tx) FindMaxValueTxOut() int {
	var maxValue int64
	var maxValueIdx int
	for txoutIdx, txout := range tx.MsgTx().TxOut {
		if txout.Value > maxValue {
			maxValue = txout.Value
			maxValueIdx = txoutIdx
		}
	}
	return maxValueIdx
}

func (tx *Tx) HasSuspiciousOutputValues() bool {
	numTinyValues := 0
	for _, txout := range tx.MsgTx().TxOut {
		if Satoshis(txout.Value).ToBTC() == 0.00000001 {
			numTinyValues++
		}
	}

	if numTinyValues > 0 && numTinyValues == len(tx.MsgTx().TxOut)-1 {
		return true
	}
	return false
}

func (tx *Tx) Fee() (BTC, error) {
	return 0, nil
	var outValues int64
	for _, txout := range tx.MsgTx().TxOut {
		outValues += txout.Value
	}
	var inValues int64
	for _, txin := range tx.MsgTx().TxIn {
		prevTx, err := tx.db.GetTx(txin.PreviousOutPoint.Hash)
		if err != nil {
			return 0, err
		}
		inValues += prevTx.MsgTx().TxOut[txin.PreviousOutPoint.Index].Value
	}

	return Satoshis(inValues - outValues).ToBTC(), nil
}

func (tx *Tx) GetSpendingTx(txoutIdx int) (*Tx, error) {
	spentTxOut, err := tx.db.GetSpentTxOut(SpentTxOutKey{TxHash: *tx.Hash(), TxOutIndex: uint32(txoutIdx)})
	if err != nil {
		return nil, err
	}

	return tx.db.GetTx(spentTxOut.InputTxHash)
}
