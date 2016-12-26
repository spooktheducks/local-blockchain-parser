package scanner

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type RawDataOutput struct {
	OutDir string
}

// ensure RawDataOutput conforms to IOutput
var _ IOutput = &RawDataOutput{}

func (o *RawDataOutput) PrintOutput(txHash chainhash.Hash, txDataSource ITxDataSource, detector IDataDetector, data []byte, result IDetectionResult) error {
	if !result.IsEmpty() {
		filename := filepath.Join(o.OutDir, fmt.Sprintf("%s-%s-%s.dat", detector.SafeName(), txHash.String(), txDataSource.Name()))
		return utils.CreateAndWriteFile(filename, data)
	}
	return nil
}

func (o *RawDataOutput) Close() error {
	return nil
}
