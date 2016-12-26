package scanner

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type (
	Scanner struct {
		TxHashSource  TxHashSource
		TxDataSources []ITxDataSource
		DataDetectors []IDataDetector
		Outputs       []IOutput

		DB *blockdb.BlockDB

		// results chan IScannerResult
	}

	// IScannerResult interface {
	// 	Description() string
	// }

	ITxDataSource interface {
		Name() string
		GetData(tx *btcutil.Tx) ([]byte, error)
	}

	IDataDetector interface {
		DetectData([]byte) (IDetectionResult, error)
		Name() string
		SafeName() string
	}

	IDetectionResult interface {
		DescriptionStrings() []string
		IsEmpty() bool
		// RawData() []byte
	}

	IOutput interface {
		PrintOutput(txHash chainhash.Hash, txDataSource ITxDataSource, detector IDataDetector, data []byte, result IDetectionResult) error
		Close() error
	}
)

func (s *Scanner) Close() error {
	for _, out := range s.Outputs {
		err := out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scanner) Run() error {
	for {
		txHash, exists := s.TxHashSource.NextHash()
		if !exists {
			break
		}

		tx, err := s.DB.GetTx(txHash)
		if err != nil {
			return err
		}

		for _, txDataSource := range s.TxDataSources {
			data, err := txDataSource.GetData(tx)
			if err != nil {
				return err
			}

			for _, d := range s.DataDetectors {
				detectionResult, err := d.DetectData(data)
				if err != nil {
					return err
				}

				for _, out := range s.Outputs {
					err := out.PrintOutput(txHash, txDataSource, d, data, detectionResult)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// func (s *Scanner) Results() chan IScannerResult {
// 	return s.results
// }
