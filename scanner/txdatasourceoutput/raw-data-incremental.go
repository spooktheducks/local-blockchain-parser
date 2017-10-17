package txdatasourceoutput

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type RawDataIncremental struct {
	OutDir   string

	initialized bool

	buffers     map[string][]byte
	textResults map[string][]textResult

	currentTx chainhash.Hash
	txIndex   uint
}

type textResult struct {
	txHash string
	method string
	inOut string
	index int
	scriptClass string
	size int
	fee int64
	timestamp int64
}

// ensure RawDataIncremental conforms to ITxDataSourceOutput
var _ scanner.ITxDataSourceOutput = &RawDataIncremental{}

func (o *RawDataIncremental) PrintOutput(tx *Tx, txDataSource scanner.ITxDataSource, dataResults []scanner.ITxDataSourceResult) error {
	txHash := *tx.Hash()

	if !o.initialized {
		o.currentTx = txHash
		o.buffers = make(map[string][]byte)
		o.textResults = make(map[string][]textResult)
		o.initialized = true
	}

	if txHash != o.currentTx {
		err := o.writeBuffers()
		if err != nil {
			return err
		}

		o.txIndex++
		o.currentTx = txHash
	}

	fee, err := tx.Fee()
	if err != nil {
		// no-op
	}

	var timestamp int64
	block, err := tx.GetBlock()
	if err == nil {
		timestamp = block.Timestamp
	}

	buffer := o.buffers[txDataSource.Name()]
	textResults := o.textResults[txDataSource.Name()]

	for _, result := range dataResults {
		buffer = append(buffer, result.RawData()...)

		var scriptClass string
		if result.InOut() == "out" {
			scriptClass = txscript.GetScriptClass(tx.MsgTx().TxOut[result.Index()].PkScript).String()
		}

		textResults = append(textResults, textResult{
			txHash:      txHash.String(),
			method:      txDataSource.Name(),
			inOut:       result.InOut(),
			index:       result.Index(),
			scriptClass: scriptClass,
			size:        len(result.RawData()),
			fee:         int64(fee),
			timestamp:   timestamp,
		})
	}

	o.buffers[txDataSource.Name()] = buffer
	o.textResults[txDataSource.Name()] = textResults

	return nil
}

func (o *RawDataIncremental) writeBuffers() error {
	for dataSource, buffer := range o.buffers {
		if len(buffer) == 0 {
			continue
		}

		filename := fmt.Sprintf("incr-%s-%d.dat", dataSource, o.txIndex)
		fullpath := filepath.Join(o.OutDir, filename)
		err := ioutil.WriteFile(fullpath, buffer, 0644)
		if err != nil {
			return err
		}
	}

	for dataSource, textResults := range o.textResults {
		if len(textResults) == 0 {
			continue
		}

		csvFile, err := os.Create(filepath.Join(o.OutDir, fmt.Sprintf("incr-%s-%d.csv", dataSource, o.txIndex)))
		if err != nil {
			return err
		}
		defer csvFile.Close()

		csvWriter := csv.NewWriter(csvFile)

		// write the columns
		err = csvWriter.Write([]string{"tx", "method", "inout", "index", "script class", "size", "tx fee", "block timestamp"})
		if err != nil {
			return err
		}
		defer csvWriter.Flush()

		for _, r := range textResults {
			csvWriter.Write([]string{
				r.txHash,
				r.method,
				r.inOut,
				fmt.Sprintf("%d", r.index),
				r.scriptClass,
				fmt.Sprintf("%d", r.size),
				fmt.Sprintf("%d", r.fee),
				fmt.Sprintf("%d", r.timestamp),
			})
		}
	}

	return nil
}

func (o *RawDataIncremental) Close() error {
	return o.writeBuffers()
}
