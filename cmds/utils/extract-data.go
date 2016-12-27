package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

func GetSatoshiEncodedData(data []byte) ([]byte, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("GetSatoshiEncodedData: not enough data")
	}

	length := binary.LittleEndian.Uint32(data[0:4])
	checksum := binary.LittleEndian.Uint32(data[4:8])
	if len(data) < 8+int(length) {
		return nil, fmt.Errorf("GetSatoshiEncodedData: not enough data")
	}

	data = data[8 : 8+length]

	if crc32.ChecksumIEEE(data) != checksum {
		return nil, fmt.Errorf("GetSatoshiEncodedData: crc32 failed")
	}
	return data, nil
}

func GetOPReturnBytes(scriptData []byte) ([]byte, error) {
	scriptStr, err := txscript.DisasmString(scriptData)
	if err != nil {
		return nil, err
	}

	toks := strings.Split(scriptStr, " ")

	for i := range toks {
		if toks[i] == "OP_RETURN" {
			if len(toks) >= i+2 {
				return hex.DecodeString(toks[i+1])
			} else {
				return nil, errors.New("GetOPReturnBytes: empty OP_RETURN data")
			}
		}
	}
	return nil, nil
}

func GetNonOPBytes(scriptData []byte) ([]byte, error) {
	scriptStr, err := txscript.DisasmString(scriptData)
	if err != nil {
		return nil, err
	}

	toks := strings.Split(scriptStr, " ")

	bs := []byte{}
	for _, tok := range toks {
		if len(tok) <= 3 {
			continue
		}

		if tok[:3] != "OP_" && len(tok) >= 40 {
			decoded, err := hex.DecodeString(tok)
			if err != nil {
				return nil, err
			}
			bs = append(bs, decoded...)
		}
	}

	return bs, nil
}

func ConcatNonOPHexTokensFromTxOuts(tx *btcutil.Tx) ([]byte, error) {
	allBytes := []byte{}

	for _, txout := range tx.MsgTx().TxOut {
		bs, err := GetNonOPBytes(txout.PkScript)
		if err != nil {
			continue
		}

		allBytes = append(allBytes, bs...)
	}

	return allBytes, nil
}

func ConcatSatoshiDataFromTxOuts(tx *btcutil.Tx) ([]byte, error) {
	data, err := ConcatNonOPHexTokensFromTxOuts(tx)
	if err != nil {
		return nil, err
	}

	return GetSatoshiEncodedData(data)
}

func ConcatTxInScripts(tx *btcutil.Tx) ([]byte, error) {
	allBytes := []byte{}

	for _, txin := range tx.MsgTx().TxIn {
		allBytes = append(allBytes, txin.SignatureScript...)
	}

	return allBytes, nil
}

func GetTxOutAddresses(tx *btcutil.Tx) ([][]btcutil.Address, error) {
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

func FindMaxValueTxOut(tx *btcutil.Tx) int {
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

func TxHasSuspiciousOutputValues(tx *btcutil.Tx) bool {
	numTinyValues := 0
	for _, txout := range tx.MsgTx().TxOut {
		if SatoshisToBTCs(txout.Value) == 0.00000001 {
			numTinyValues++
		}
	}

	if numTinyValues > 0 && numTinyValues == len(tx.MsgTx().TxOut)-1 {
		return true
	}
	return false
}
