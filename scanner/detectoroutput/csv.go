package detectoroutput

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type CSV struct {
	OutDir   string
	csvFiles map[string]*utils.ConditionalFile
}

// ensure CSV conforms to scanner.IDetectorOutput
var _ scanner.IDetectorOutput = &CSV{}

func (o *CSV) PrintOutput(txHash chainhash.Hash, txDataSource scanner.ITxDataSource, dataResult scanner.ITxDataSourceResult, detector scanner.IDetector, result scanner.IDetectionResult) error {
	csvFile, exists := o.csvFiles[detector.SafeName()]
	if !exists {
		if o.csvFiles == nil {
			o.csvFiles = make(map[string]*utils.ConditionalFile)
		}
		csvFile = utils.NewConditionalFile(filepath.Join(o.OutDir, detector.SafeName()+".csv"))
		_, err := csvFile.WriteString("tx hash,data source,description\n", false)
		if err != nil {
			return err
		}

		o.csvFiles[detector.SafeName()] = csvFile
	}

	for _, str := range result.DescriptionStrings() {
		_, err := csvFile.WriteString(fmt.Sprintf("%s,%s,%s\n", txHash.String(), dataResult.SourceName(), str), true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *CSV) Close() error {
	for _, file := range o.csvFiles {
		file.Close()
	}
	return nil
}
