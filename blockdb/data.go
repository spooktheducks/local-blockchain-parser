package blockdb

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type BlockIndexRow struct {
	DATFileIdx     uint16
	Timestamp      int64
	IndexInDATFile uint32
}

func NewBlockIndexRowFromBytes(bs []byte) (BlockIndexRow, error) {
	row := BlockIndexRow{}
	err := binary.Read(bytes.NewReader(bs), binary.LittleEndian, &row)
	if err != nil {
		return BlockIndexRow{}, err
	}
	return row, nil
}

func (r BlockIndexRow) DATFilename() string {
	return fmt.Sprintf("blk%05d.dat", r.DATFileIdx)
}

func (r BlockIndexRow) ToBytes() ([]byte, error) {
	rowData := &bytes.Buffer{}
	err := binary.Write(rowData, binary.LittleEndian, r)
	if err != nil {
		return nil, err
	}
	return rowData.Bytes(), nil
}

type TxIndexRow struct {
	BlockHash    chainhash.Hash
	IndexInBlock uint64
}

func NewTxIndexRowFromBytes(bs []byte) (TxIndexRow, error) {
	row := TxIndexRow{}
	err := binary.Read(bytes.NewReader(bs), binary.LittleEndian, &row)
	if err != nil {
		return TxIndexRow{}, err
	}
	return row, nil
}

func (r TxIndexRow) ToBytes() ([]byte, error) {
	rowData := &bytes.Buffer{}
	err := binary.Write(rowData, binary.LittleEndian, r)
	if err != nil {
		return nil, err
	}
	return rowData.Bytes(), nil
}

type SpentTxOutKey struct {
	TxHash     chainhash.Hash
	TxOutIndex uint32
}

type SpentTxOutRow struct {
	InputTxHash chainhash.Hash
	TxInIndex   uint32
}

func (k SpentTxOutKey) ToBytes() ([]byte, error) {
	data := &bytes.Buffer{}
	err := binary.Write(data, binary.LittleEndian, k)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func NewSpentTxOutRowFromBytes(bs []byte) (SpentTxOutRow, error) {
	row := SpentTxOutRow{}
	err := binary.Read(bytes.NewReader(bs), binary.LittleEndian, &row)
	if err != nil {
		return SpentTxOutRow{}, err
	}
	return row, nil
}

func (r SpentTxOutRow) ToBytes() ([]byte, error) {
	data := &bytes.Buffer{}
	err := binary.Write(data, binary.LittleEndian, r)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}
