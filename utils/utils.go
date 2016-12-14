package utils

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func SatoshisToBTCs(satoshis int64) float64 {
	return float64(satoshis) * 0.00000001
}

func LoadBlockFile(file string) (blocks []*btcutil.Block, err error) {
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
