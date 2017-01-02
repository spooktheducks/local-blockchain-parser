package detector

import (
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

type MagicBytes struct{}

// ensure MagicBytes conforms to scanner.IDetector
var _ scanner.IDetector = &MagicBytes{}

// ensure MagicBytesResult conforms to scanner.IDetectionResult
var _ scanner.IDetectionResult = utils.MagicBytesResult{}

func (d *MagicBytes) DetectData(data []byte) (scanner.IDetectionResult, error) {
	return utils.SearchDataForMagicFileBytes(data), nil
}

func (d *MagicBytes) Name() string {
	return "Magic bytes"
}

func (d *MagicBytes) SafeName() string {
	return "magic-bytes"
}
