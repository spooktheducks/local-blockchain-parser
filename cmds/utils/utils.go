package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

var maxFiles = 256
var fileSemaphore = make(chan bool, maxFiles)

func init() {
	for i := 0; i < maxFiles; i++ {
		fileSemaphore <- true
	}
}

func CreateAndWriteFile(path string, bytes []byte) error {
	f, err := CreateFile(path)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		CloseFile(f)
		return err
	}
	return CloseFile(f)
}

func CreateFile(path string) (*os.File, error) {
	<-fileSemaphore
	f, err := os.Create(path)
	if err != nil {
		fileSemaphore <- true
		return nil, err
	}
	return f, nil
}

func CloseFile(file *os.File) error {
	err := file.Close()
	fileSemaphore <- true
	return err
}

func GetOPReturnBytes(scriptStr string) ([]byte, error) {
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

func GetNonOPBytes(scriptStr string) ([]byte, error) {
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
		scriptStr, err := txscript.DisasmString(txout.PkScript)
		if err != nil {
			// if err.Error() == "execute past end of script" {
			//  continue
			// } else {
			//  return nil, fmt.Errorf("error in txscript.DisasmString: %v", err)
			// }
			continue
		}

		bs, err := GetNonOPBytes(scriptStr)
		if err != nil {
			continue
		}

		allBytes = append(allBytes, bs...)
	}

	return allBytes, nil
}

func SatoshisToBTCs(satoshis int64) float64 {
	return float64(satoshis) * 0.00000001
}

func LoadBlockFile(file string) (blocks []*btcutil.Block, err error) {
	<-fileSemaphore
	defer func() { fileSemaphore <- true }()

	var network = wire.MainNet
	var dr io.Reader
	var fi io.ReadCloser

	fi, err = os.Open(file)
	if err != nil {
		return
	}

	dr = fi
	defer fi.Close()

	var block *btcutil.Block

	err = nil
	for height := int64(1); err == nil; height++ {
		var rintbuf uint32
		err = binary.Read(dr, binary.LittleEndian, &rintbuf)
		if err == io.EOF {
			// hit end of file at expected offset: no warning
			height--
			err = nil
			break
		}
		if err != nil {
			break
		}
		if rintbuf != uint32(network) {
			break
		}
		err = binary.Read(dr, binary.LittleEndian, &rintbuf)
		blocklen := rintbuf

		rbytes := make([]byte, blocklen)

		// read block
		dr.Read(rbytes)

		block, err = btcutil.NewBlockFromBytes(rbytes)
		if err != nil {
			return
		}
		blocks = append(blocks, block)
	}

	return
}

func GroupBlocks(blocks []*btcutil.Block, groupLen int) [][]*btcutil.Block {
	extra := len(blocks) % groupLen
	numGroups := ((len(blocks) - extra) / groupLen) + 1

	blockIdx := 0
	groups := make([][]*btcutil.Block, numGroups)
	for g := 0; g < numGroups; g++ {
		groups[g] = []*btcutil.Block{}

		for i := 0; i < groupLen; i++ {
			groups[g] = append(groups[g], blocks[blockIdx])
			blockIdx++

			if blockIdx == len(blocks) {
				return groups
			}
		}
	}

	return groups
}
