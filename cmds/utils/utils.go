package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func LoadBlocksFromDAT(file string) (blocks []*btcutil.Block, err error) {
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

func LoadBlockFromDAT(file string, height uint32) (*btcutil.Block, error) {
	<-fileSemaphore
	defer func() { fileSemaphore <- true }()

	var network = wire.MainNet

	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fi.Close()

	for i := uint32(0); i <= height; i++ {
		var networkBits uint32
		err := binary.Read(fi, binary.LittleEndian, &networkBits)
		if err == io.EOF {
			// hit end of file at expected offset: no warning
			return nil, fmt.Errorf("block %v not found in DAT file", height)
		}
		if err != nil {
			return nil, err
		}
		if networkBits != uint32(network) {
			return nil, fmt.Errorf("block has bad network bits")
		}

		var blocklen uint32
		err = binary.Read(fi, binary.LittleEndian, &blocklen)

		// read block
		if i == height {
			rbytes := make([]byte, blocklen)
			fi.Read(rbytes)
			return btcutil.NewBlockFromBytes(rbytes)
		} else {
			fi.Seek(int64(blocklen), 1)
		}
	}

	return nil, fmt.Errorf("block %v not found in DAT file", height)
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

func ReverseBytes(bs []byte) []byte {
	length := len(bs)
	reversed := make([]byte, length)

	for i := range bs {
		reversed[length-1-i] = bs[i]
	}

	return reversed
}

func DATFilename(idx uint16) string {
	return fmt.Sprintf("blk%05d.dat", idx)
}

func HashFromBytes(bs []byte) (chainhash.Hash, error) {
	hash := &chainhash.Hash{}
	err := hash.SetBytes(bs)
	if err != nil {
		return chainhash.Hash{}, err
	}
	return *hash, nil
}

func HashFromString(s string) (chainhash.Hash, error) {
	h, err := chainhash.NewHashFromStr(s)
	if err != nil {
		return chainhash.Hash{}, err
	}
	return *h, nil
}
