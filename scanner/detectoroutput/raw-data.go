package detectoroutput

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type RawData struct {
	OutDir string
}

// ensure RawData conforms to scanner.IDetectorOutput
var _ scanner.IDetectorOutput = &RawData{}

func (o *RawData) PrintOutput(txHash chainhash.Hash, txDataSource scanner.ITxDataSource, detector scanner.IDetector, data []byte, result scanner.IDetectionResult) error {
	if !result.IsEmpty() {
		filename := filepath.Join(o.OutDir, fmt.Sprintf("%s-%s-%s.dat", detector.SafeName(), txHash.String(), txDataSource.Name()))
		return utils.CreateAndWriteFile(filename, data)
	}
	return nil
}

func (o *RawData) Close() error {
	return nil
}
