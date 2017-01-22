package dbcmds

import (
	"os"
	"path/filepath"

	"github.com/spooktheducks/local-blockchain-parser/scanner"
	"github.com/spooktheducks/local-blockchain-parser/scanner/detector"
	"github.com/spooktheducks/local-blockchain-parser/scanner/detectoroutput"
	"github.com/spooktheducks/local-blockchain-parser/scanner/txdatasource"
	"github.com/spooktheducks/local-blockchain-parser/scanner/txdatasourceoutput"
	"github.com/spooktheducks/local-blockchain-parser/scanner/txhashsource"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type ScanAddressCommand struct {
	dbFile     string
	datFileDir string
	walletAddr string
	outDir     string
	db         *BlockDB
}

func NewScanAddressCommand(datFileDir, dbFile, outDir, walletAddr string) *ScanAddressCommand {
	return &ScanAddressCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		walletAddr: walletAddr,
		outDir:     filepath.Join(outDir, "address", walletAddr),
	}
}

func (cmd *ScanAddressCommand) RunCommand() error {
	err := os.MkdirAll(cmd.outDir, 0777)
	if err != nil {
		return err
	}

	db, err := NewBlockDB(cmd.dbFile, cmd.datFileDir)
	if err != nil {
		return err
	}
	defer db.Close()

	cmd.db = db

	s := &scanner.Scanner{
		DB:           db,
		TxHashSource: txhashsource.NewAddressTxHashSource(db, cmd.walletAddr),
		TxDataSources: []scanner.ITxDataSource{
			&txdatasource.InputScript{},
			&txdatasource.InputScriptsConcat{},
			&txdatasource.OutputScript{},
			&txdatasource.OutputScriptsConcat{},
			&txdatasource.OutputScriptsSatoshi{},
			&txdatasource.OutputScriptOpReturn{},
		},
		TxDataSourceOutputs: []scanner.ITxDataSourceOutput{
			&txdatasourceoutput.RawData{OutDir: cmd.outDir},
			&txdatasourceoutput.RawDataEachDataSource{OutDir: cmd.outDir},
		},
		Detectors: []scanner.IDetector{
			&detector.PGPPackets{},
			// &detector.AESKeys{},
			&detector.MagicBytes{},
			// &detector.Plaintext{},
		},
		DetectorOutputs: []scanner.IDetectorOutput{
			&detectoroutput.Console{Prefix: "  - "},
			&detectoroutput.RawData{OutDir: cmd.outDir},
			&detectoroutput.CSV{OutDir: cmd.outDir},
			&detectoroutput.CSVTxAnalysis{OutDir: cmd.outDir},
		},
	}

	err = s.Run()
	if err != nil {
		return err
	}

	return s.Close()
}
