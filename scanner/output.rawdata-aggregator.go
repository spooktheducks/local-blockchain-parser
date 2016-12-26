package scanner

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type RawDataAggregatorOutput struct {
	OutDir string
	file   *utils.ConditionalFile
}

// ensure RawDataAggregatorOutput conforms to IOutput
var _ IOutput = &RawDataAggregatorOutput{}

func (o *RawDataAggregatorOutput) PrintOutput(txHash chainhash.Hash, txDataSource ITxDataSource, detector IDataDetector, data []byte, result IDetectionResult) error {
	if o.file == nil {
		filename := fmt.Sprintf("%s-%s-aggregated.dat", detector.SafeName(), txDataSource.Name())
		o.file = utils.NewConditionalFile(filepath.Join(o.OutDir, filename))
	}

	// if !result.IsEmpty() {
	_, err := o.file.Write(data, true)
	if err != nil {
		return err
	}
	// }
	return nil
}

func (o *RawDataAggregatorOutput) Close() error {
	return o.file.Close()
}
