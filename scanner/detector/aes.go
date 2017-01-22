package detector

import (
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils/aeskeyfind"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type AESKeys struct{}

// ensure AESKeys conforms to scanner.IDetector
var _ scanner.IDetector = &AESKeys{}

// ensure PGPPacketResult conforms to scanner.IDetectionResult
var _ scanner.IDetectionResult = &aeskeyfind.AESResult{}

func (d *AESKeys) DetectData(data []byte) (scanner.IDetectionResult, error) {
	return aeskeyfind.Detect(data), nil
}

func (d *AESKeys) Name() string {
	return "AES key"
}

func (d *AESKeys) SafeName() string {
	return "aes-keys"
}
