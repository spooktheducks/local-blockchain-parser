package txdatasourceoutput

import (
	"fmt"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type RawData struct {
	OutDir   string
	outFiles map[string]*utils.ConditionalFile
}

// ensure RawData conforms to ITxDataSourceOutput
var _ scanner.ITxDataSourceOutput = &RawData{}

func (o *RawData) PrintOutput(txDataSource scanner.ITxDataSource, data []byte) error {
	outFile, exists := o.outFiles[txDataSource.Name()]
	if !exists {
		if o.outFiles == nil {
			o.outFiles = make(map[string]*utils.ConditionalFile)
		}
		filename := fmt.Sprintf("%s-aggregated.dat", txDataSource.Name())
		outFile = utils.NewConditionalFile(filepath.Join(o.OutDir, filename))
		o.outFiles[txDataSource.Name()] = outFile
	}

	_, err := outFile.Write(data, true)
	if err != nil {
		return err
	}
	return nil
}

func (o *RawData) Close() error {
	for _, f := range o.outFiles {
		f.Close()
	}
	return nil
}
