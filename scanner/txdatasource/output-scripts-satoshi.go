package txdatasource

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

type OutputScriptsSatoshi struct{}

type OutputScriptsSatoshiResult struct {
	rawData []byte
	index int
}

// ensure that OutputScriptsSatoshi conforms to ITxDataSource
var _ scanner.ITxDataSource = &OutputScriptsSatoshi{}

// ensure that OutputScriptsSatoshiResult conforms to ITxDataSourceResult
var _ scanner.ITxDataSourceResult = &OutputScriptsSatoshiResult{}

func (ds *OutputScriptsSatoshi) Name() string {
	return "outputs-satoshi"
}

func u32len(bs []byte) uint32 {
	return uint32(len(bs))
}

func (ds *OutputScriptsSatoshi) GetData(tx *Tx) ([]scanner.ITxDataSourceResult, error) {
	firstData, err := utils.GetNonOPBytesFromOutputScript(tx.MsgTx().TxOut[0].PkScript)
	if err != nil {
		return nil, err
	}

	if len(firstData) < 8 {
		return nil, fmt.Errorf("GetSatoshiEncodedData: not enough data")
	}

	length := binary.LittleEndian.Uint32(firstData[0:4])
	checksum := binary.LittleEndian.Uint32(firstData[4:8])

	var data []byte
	var results []scanner.ITxDataSourceResult

	for i, txout := range tx.MsgTx().TxOut {
		bs, err := utils.GetNonOPBytesFromOutputScript(txout.PkScript)
		if err != nil {
			continue
		}

		if i == 0 {
			bs = bs[8:]
		}

		// trim the last result if needed
		if u32len(data) + u32len(bs) > length {
			extra := u32len(data) + u32len(bs) - length
			bs = bs[0 : u32len(bs) - extra]
		}

		data = append(data, bs...)
		results = append(results, OutputScriptsSatoshiResult{rawData: bs, index: i})

		if u32len(data) >= length {
			break
		}
	}

	if u32len(data) < length {
		return nil, fmt.Errorf("GetSatoshiEncodedData: not enough data")
	}

	if crc32.ChecksumIEEE(data) != checksum {
		return nil, fmt.Errorf("GetSatoshiEncodedData: crc32 failed")
	}

	return results, nil
}

func (r OutputScriptsSatoshiResult) SourceName() string {
	return fmt.Sprintf("outputs-satoshi-%d", r.index)
}

func (r OutputScriptsSatoshiResult) RawData() []byte {
	return r.rawData
}

func (r OutputScriptsSatoshiResult) InOut() string {
	return "out"
}

func (r OutputScriptsSatoshiResult) Index() int {
	return r.index
}
