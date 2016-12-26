package scanner

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type ConsoleOutput struct {
	Prefix string
}

// ensure ConsoleOutput conforms to IOutput
var _ IOutput = &ConsoleOutput{}

func (o *ConsoleOutput) PrintOutput(txHash chainhash.Hash, txDataSource ITxDataSource, detector IDataDetector, data []byte, result IDetectionResult) error {
	for _, str := range result.DescriptionStrings() {
		fmt.Printf("%s%s: %s: %s\n", o.Prefix, txDataSource.Name(), detector.Name(), str)
	}
	return nil
}

func (o *ConsoleOutput) Close() error {
	return nil
}
