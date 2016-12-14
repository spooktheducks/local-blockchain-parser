package cmds

import (
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

var maxFiles = 256
var fileSemaphore = make(chan bool, maxFiles)

func createAndWriteFile(path string, bytes []byte) error {
	f, err := createFile(path)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		closeFile(f)
		return err
	}
	return closeFile(f)
}

func createFile(path string) (*os.File, error) {
	<-fileSemaphore
	f, err := os.Create(path)
	if err != nil {
		fileSemaphore <- true
		return nil, err
	}
	return f, nil
}

func closeFile(file *os.File) error {
	err := file.Close()
	fileSemaphore <- true
	return err
}

func getOPReturnBytes(scriptStr string) ([]byte, error) {
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

func getNonOPBytes(scriptStr string) ([]byte, error) {
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

func concatNonOPHexTokensFromTxOuts(tx *btcutil.Tx) ([]byte, error) {
	allBytes := []byte{}

	for _, txout := range tx.MsgTx().TxOut {
		scriptStr, err := txscript.DisasmString(txout.PkScript)
		if err != nil {
			// if err.Error() == "execute past end of script" {
			// 	continue
			// } else {
			// 	return nil, fmt.Errorf("error in txscript.DisasmString: %v", err)
			// }
			continue
		}

		bs, err := getNonOPBytes(scriptStr)
		if err != nil {
			continue
		}

		allBytes = append(allBytes, bs...)
	}

	return allBytes, nil
}
