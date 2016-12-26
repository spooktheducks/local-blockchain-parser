package scanner

// import (
// 	"fmt"

// 	"github.com/btcsuite/btcd/chaincfg/chainhash"

// 	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
// )

// type CSVOutput struct {
// 	csvFile *utils.ConditionalFile
// }

// // ensure CSVOutput conforms to IOutput
// var _ IOutput = &CSVOutput{}

// func NewCSVOutput(filenameFunc func(txHash chainhash.Hash, txDataSource ITxDataSource, detector IDataDetector, data []byte, result IDetectionResult)) *CSVOutput {
// 	csvFile := utils.NewConditionalFile(filename)
// 	return &CSVOutput{csvFile: csvFile}
// }

// func (o *CSVOutput) PrintOutput(txHash chainhash.Hash, txDataSource ITxDataSource, detector IDataDetector, data []byte, result IDetectionResult) error {
// 	for _, str := range result.DescriptionStrings() {
// 		_, err := o.csvFile.WriteString(fmt.Sprintf("%s,%s,%s\n", txHash.String(), txDataSource.Name(), str), true)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (o *CSVOutput) Close() error {
// 	return o.csvFile.Close()
// }
