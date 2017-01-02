package scanner

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

type (
	Scanner struct {
		TxHashSource        ITxHashSource
		TxDataSources       []ITxDataSource
		TxDataSourceOutputs []ITxDataSourceOutput
		Detectors           []IDetector
		DetectorOutputs     []IDetectorOutput

		DB *blockdb.BlockDB
	}

	ITxHashSource interface {
		NextHash() (chainhash.Hash, bool)
	}

	ITxDataSource interface {
		Name() string
		GetData(tx *btcutil.Tx) ([]byte, error)
	}

	ITxDataSourceOutput interface {
		PrintOutput(txDataSource ITxDataSource, data []byte) error
		Close() error
	}

	IDetector interface {
		DetectData([]byte) (IDetectionResult, error)
		Name() string
		SafeName() string
	}

	IDetectionResult interface {
		DescriptionStrings() []string
		IsEmpty() bool
	}

	IDetectorOutput interface {
		PrintOutput(txHash chainhash.Hash, txDataSource ITxDataSource, detector IDetector, data []byte, result IDetectionResult) error
		Close() error
	}
)

func (s *Scanner) Close() error {
	for _, out := range s.DetectorOutputs {
		err := out.Close()
		if err != nil {
			return err
		}
	}

	for _, out := range s.TxDataSourceOutputs {
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
				continue
				// return err
			}

			for _, out := range s.TxDataSourceOutputs {
				err := out.PrintOutput(txDataSource, data)
				if err != nil {
					return err
				}
			}

			for _, d := range s.Detectors {
				detectionResult, err := d.DetectData(data)
				if err != nil {
					return err
				}

				for _, out := range s.DetectorOutputs {
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
