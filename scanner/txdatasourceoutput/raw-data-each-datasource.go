package txdatasourceoutput

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type RawDataEachDataSource struct {
	OutDir   string
}

// ensure RawDataEachDataSource conforms to ITxDataSourceOutput
var _ scanner.ITxDataSourceOutput = &RawDataEachDataSource{}

func (o *RawDataEachDataSource) PrintOutput(tx *Tx, txDataSource scanner.ITxDataSource, dataResults []scanner.ITxDataSourceResult) error {
	txHash := tx.Hash()

	for _, result := range dataResults {
		if len(result.RawData()) == 0 {
			continue
		}

		filename := filepath.Join(o.OutDir, fmt.Sprintf("%s-%s.dat", txHash.String(), result.SourceName()))
		err := ioutil.WriteFile(filename, result.RawData(), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *RawDataEachDataSource) Close() error {
	return nil
}
