package scanner

import (
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type (
	PlaintextDetector struct {
	}

	PlaintextResult struct {
		TextData []byte
	}
)

// ensure that PlaintextDetector conforms to IDataDetector
var _ IDataDetector = &PlaintextDetector{}

// ensure that PlaintextResult conforms to IDetectionResult
var _ IDetectionResult = PlaintextResult{}

func (d *PlaintextDetector) DetectData(bs []byte) (IDetectionResult, error) {
	nonText := utils.StripNonTextBytes(bs)
	return PlaintextResult{TextData: nonText}, nil
}

func (d *PlaintextDetector) Name() string {
	return "Plaintext"
}

func (d *PlaintextDetector) SafeName() string {
	return "plaintext"
}

func (r PlaintextResult) DescriptionStrings() []string {
	return []string{string(r.TextData)}
}

func (r PlaintextResult) IsEmpty() bool {
	return r.TextData == nil || len(r.TextData) < 8
}
