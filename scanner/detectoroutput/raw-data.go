package detectoroutput

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type RawData struct {
	OutDir string
}

// ensure RawData conforms to scanner.IDetectorOutput
var _ scanner.IDetectorOutput = &RawData{}

func (o *RawData) PrintOutput(txHash chainhash.Hash, txDataSource scanner.ITxDataSource, dataResult scanner.ITxDataSourceResult, detector scanner.IDetector, result scanner.IDetectionResult) error {
	if !result.IsEmpty() {
		filename := filepath.Join(o.OutDir, fmt.Sprintf("%s-%s-%s.dat", detector.SafeName(), txHash.String(), dataResult.SourceName()))
		return utils.CreateAndWriteFile(filename, dataResult.RawData())
	}
	return nil
}

func (o *RawData) Close() error {
	return nil
}
