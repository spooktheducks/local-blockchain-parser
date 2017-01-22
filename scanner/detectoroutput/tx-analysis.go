package detectoroutput

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type CSVTxAnalysis struct {
	OutDir string
	DB     *BlockDB

	columns []string
	data    map[chainhash.Hash]map[string]bool
}

// ensure CSVTxAnalysis conforms to scanner.IDetectorOutput
var _ scanner.IDetectorOutput = &CSVTxAnalysis{}

func (o *CSVTxAnalysis) PrintOutput(txHash chainhash.Hash, txDataSource scanner.ITxDataSource, dataResult scanner.ITxDataSourceResult, detector scanner.IDetector, result scanner.IDetectionResult) error {
	o.appendColumnIfUnique(detector.Name())

	txData, exists := o.data[txHash]
	if !exists {
		if o.data == nil {
			o.data = make(map[chainhash.Hash]map[string]bool)
		}

		txData = make(map[string]bool)
	}

	if !result.IsEmpty() {
		txData[detector.Name()] = true
	}

	o.data[txHash] = txData

	return nil
}

func (o *CSVTxAnalysis) appendColumnIfUnique(c string) {
	found := false
	for i := range o.columns {
		if o.columns[i] == c {
			found = true
			break
		}
	}

	if !found {
		o.columns = append(o.columns, c)
	}
}

func (o *CSVTxAnalysis) Close() error {
	// add 'tx-hash' as the first column
	o.columns = append([]string{"tx hash"}, o.columns...)
	// o.columns = append(o.columns, "fee")

	// open the csv file
	csvFile, err := os.Create(filepath.Join(o.OutDir, "tx-analysis.csv"))
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	// write the columns
	err = csvWriter.Write(o.columns)
	if err != nil {
		return err
	}
	defer csvWriter.Flush()

	// write the data
	for txHash, txData := range o.data {
		row := []string{}
		for _, col := range o.columns {
			if col == "tx hash" {
				row = append(row, txHash.String())
				continue
			}

			val := ""
			if txData[col] {
				val = "X"
			}
			row = append(row, val)
		}
		csvWriter.Write(row)
	}

	return nil
}
