package detectoroutput

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type Console struct {
	Prefix string
}

// ensure Console conforms to scanner.IDetectorOutput
var _ scanner.IDetectorOutput = &Console{}

func (o *Console) PrintOutput(txHash chainhash.Hash, txDataSource scanner.ITxDataSource, detector scanner.IDetector, data []byte, result scanner.IDetectionResult) error {
	for _, str := range result.DescriptionStrings() {
		fmt.Printf("%s%s %s: %s: %s\n", o.Prefix, txHash.String(), detector.Name(), txDataSource.Name(), str)
	}
	return nil
}

func (o *Console) Close() error {
	return nil
}
