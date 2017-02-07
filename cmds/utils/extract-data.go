package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"strings"

	// "github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	// "github.com/btcsuite/btcd/wire"
	// "github.com/btcsuite/btcutil"
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

func GetNonOPBytesFromOutputScript(scriptData []byte) ([]byte, error) {
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
