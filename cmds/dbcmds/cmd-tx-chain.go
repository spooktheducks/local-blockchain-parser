package dbcmds

import (
	"os"
	"path/filepath"

	"github.com/spooktheducks/local-blockchain-parser/scanner"
	"github.com/spooktheducks/local-blockchain-parser/scanner/detector"
	"github.com/spooktheducks/local-blockchain-parser/scanner/detectoroutput"
	"github.com/spooktheducks/local-blockchain-parser/scanner/txdatasource"
	"github.com/spooktheducks/local-blockchain-parser/scanner/txdatasourceoutput"
	"github.com/spooktheducks/local-blockchain-parser/scanner/txhashoutput"
	"github.com/spooktheducks/local-blockchain-parser/scanner/txhashsource"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

type TxChainCommand struct {
	dbFile     string
	datFileDir string
	outDir     string
	direction  string
	txHash     string
	limit      uint

	db *BlockDB
}

func NewTxChainCommand(datFileDir, dbFile, outDir, direction string, limit uint, txHash string) *TxChainCommand {
	if direction != "forward" &&
		direction != "backward" &&
		direction != "both" {
		panic("--direction (-d) must be 'forward', 'backward', or 'both'")
	}

	return &TxChainCommand{
		dbFile:     dbFile,
		datFileDir: datFileDir,
		txHash:     txHash,
		direction:  direction,
		limit:      limit,
		outDir:     filepath.Join(outDir, "tx-chain", txHash),
	}
}

func (cmd *TxChainCommand) RunCommand() error {
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

	startHash, err := utils.HashFromString(cmd.txHash)
	if err != nil {
		return err
	}

	var txHashSource scanner.ITxHashSource
	if cmd.direction == "forward" {
		txHashSource = txhashsource.NewForwardChain(db, startHash, cmd.limit)
	} else if cmd.direction == "backward" {
		txHashSource = txhashsource.NewBackwardChain(db, startHash, cmd.limit)
	} else {
		txHashSource = txhashsource.NewChain(db, startHash, cmd.limit)
	}

	s := &scanner.Scanner{
		DB:           db,
		TxHashSource: txHashSource,
		TxHashOutputs: []scanner.ITxHashOutput{
			&txhashoutput.HashOnly{OutDir: cmd.outDir, Filename: "transactions.txt"},
			&txhashoutput.OpReturn{OutDir: cmd.outDir, Filename: "transactions-opreturn.txt"},
			&txhashoutput.NonOp{OutDir: cmd.outDir, Filename: "transactions-nonop.txt"},
			&txhashoutput.InputScript{OutDir: cmd.outDir, Filename: "transactions-inputscripts.txt"},
			&txhashoutput.InputScriptNonOP{OutDir: cmd.outDir, Filename: "transactions-inputscripts-nonop.txt"},
		},
		TxDataSources: []scanner.ITxDataSource{
			&txdatasource.InputScript{},
			&txdatasource.InputScriptNonOP{},
			&txdatasource.InputScriptPushdata{},
			&txdatasource.InputScriptsConcat{},
			&txdatasource.OutputScript{},
			&txdatasource.OutputScript{OrderByValue: true},
			&txdatasource.OutputScript{SkipMaxValueTxOut: true},
			&txdatasource.OutputScript{SkipMaxValueTxOut: true, OrderByValue: true},
			&txdatasource.OutputScriptsSatoshi{},
			&txdatasource.OutputScriptOpReturn{},
			&txdatasource.OutputScriptsConcat{},
		},
		TxDataSourceOutputs: []scanner.ITxDataSourceOutput{
			&txdatasourceoutput.RawData{OutDir: cmd.outDir},
			&txdatasourceoutput.RawDataEachDataSource{OutDir: cmd.outDir},
		},
		Detectors: []scanner.IDetector{
			// &detector.PGPPackets{},
			&detector.AESKeys{},
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
