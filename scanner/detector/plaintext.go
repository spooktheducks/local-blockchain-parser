package detector

import (
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
	"github.com/spooktheducks/local-blockchain-parser/scanner"
)

type (
	Plaintext struct{}

	PlaintextResult struct {
		TextData []byte
	}
)

// ensure that Plaintext conforms to scanner.IDetector
var _ scanner.IDetector = &Plaintext{}

// ensure that PlaintextResult conforms to scanner.IDetectionResult
var _ scanner.IDetectionResult = PlaintextResult{}

func (d *Plaintext) DetectData(bs []byte) (scanner.IDetectionResult, error) {
	nonText := utils.StripNonTextBytes(bs)
	return PlaintextResult{TextData: nonText}, nil
}

func (d *Plaintext) Name() string {
	return "Plaintext"
}

func (d *Plaintext) SafeName() string {
	return "plaintext"
}

func (r PlaintextResult) DescriptionStrings() []string {
	return []string{string(r.TextData)}
}

func (r PlaintextResult) IsEmpty() bool {
	return r.TextData == nil || len(r.TextData) < 8
}
