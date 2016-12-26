package scanner

import (
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type PGPDataDetector struct{}

// ensure PGPDataDetector conforms to IDataDetector
var _ IDataDetector = &PGPDataDetector{}

// ensure PGPPacketResult conforms to IDetectionResult
var _ IDetectionResult = &utils.PGPPacketResult{}

func (d *PGPDataDetector) DetectData(data []byte) (IDetectionResult, error) {
	return utils.FindPGPPackets(data), nil
}

func (d *PGPDataDetector) Name() string {
	return "PGP data"
}

func (d *PGPDataDetector) SafeName() string {
	return "pgp-data"
}
