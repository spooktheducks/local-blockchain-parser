package scanner

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
)

type (
	Scanner struct {
		TxHashSource        ITxHashSource
		TxHashOutputs       []ITxHashOutput
		TxDataSources       []ITxDataSource
		TxDataSourceOutputs []ITxDataSourceOutput
		Detectors           []IDetector
		DetectorOutputs     []IDetectorOutput

		DB *BlockDB
	}

	ITxHashSource interface {
		NextHash() (chainhash.Hash, bool)
	}

	ITxHashOutput interface {
		OutputTx(*Tx) error
		Close() error
	}

	ITxDataSource interface {
		Name() string
		GetData(tx *Tx) ([]ITxDataSourceResult, error)
	}

	ITxDataSourceResult interface {
		SourceName() string
		RawData() []byte
		InOut() string
		Index() int
	}

	ITxDataSourceOutput interface {
		PrintOutput(*Tx, ITxDataSource, []ITxDataSourceResult) error
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
		PrintOutput(txHash chainhash.Hash, txDataSource ITxDataSource, dataResult ITxDataSourceResult, detector IDetector, result IDetectionResult) error
		Close() error
	}
)

func (s *Scanner) Close() error {
	for _, out := range s.TxHashOutputs {
		err := out.Close()
		if err != nil {
			return err
		}
	}

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
			fmt.Printf("cannot get tx %v\n", txHash)
			return err
		}

		for _, out := range s.TxHashOutputs {
			err := out.OutputTx(tx)
			if err != nil {
				return err
			}
		}

		for _, txDataSource := range s.TxDataSources {
			dataResults, err := txDataSource.GetData(tx)
			if err != nil {
				// if a data source returns an error, we assume that means there's
				// no data of that type, so we just continue to the next source
				continue
			}

			for _, out := range s.TxDataSourceOutputs {
				err := out.PrintOutput(tx, txDataSource, dataResults)
				if err != nil {
					return err
				}
			}

			for _, dataResult := range dataResults {
				for _, detector := range s.Detectors {
					detectionResult, err := detector.DetectData(dataResult.RawData())
					if err != nil {
						return err
					}

					for _, out := range s.DetectorOutputs {
						err := out.PrintOutput(txHash, txDataSource, dataResult, detector, detectionResult)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
