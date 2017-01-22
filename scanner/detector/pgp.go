package detector

import (
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type PGPPackets struct{}

// ensure PGPPackets conforms to scanner.IDetector
var _ scanner.IDetector = &PGPPackets{}

// ensure PGPPacketResult conforms to scanner.IDetectionResult
var _ scanner.IDetectionResult = &utils.PGPPacketResult{}

func (d *PGPPackets) DetectData(data []byte) (scanner.IDetectionResult, error) {
	return utils.FindPGPPackets(data), nil
}

func (d *PGPPackets) Name() string {
	return "PGP data"
}

func (d *PGPPackets) SafeName() string {
	return "pgp-data"
}
