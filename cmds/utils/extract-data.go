package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

func GetOPReturnBytes(scriptData []byte) ([]byte, error) {
	scriptStr, err := txscript.DisasmString(scriptData)
	if err != nil {
		return nil, err
	}

	toks := strings.Split(scriptStr, " ")

	for i := range toks {
		if toks[i] == "OP_RETURN" {
			if len(toks) >= i+2 {
				return hex.DecodeString(toks[i+1])
			} else {
				return nil, errors.New("empty OP_RETURN data")
			}
		}
	}
	return nil, nil
}

func GetNonOPBytes(scriptData []byte) ([]byte, error) {
	scriptStr, err := txscript.DisasmString(scriptData)
	if err != nil {
		return nil, err
	}

	toks := strings.Split(scriptStr, " ")

	bs := []byte{}
	for _, tok := range toks {
		if len(tok) <= 3 {
			continue
		}

		if tok[:3] != "OP_" && len(tok) >= 40 {
			decoded, err := hex.DecodeString(tok)
			if err != nil {
				return nil, err
			}
			bs = append(bs, decoded...)
		}
	}

	return bs, nil
}

func ConcatNonOPHexTokensFromTxOuts(tx *btcutil.Tx) ([]byte, error) {
	allBytes := []byte{}

	for _, txout := range tx.MsgTx().TxOut {
		bs, err := GetNonOPBytes(txout.PkScript)
		if err != nil {
			continue
		}

		allBytes = append(allBytes, bs...)
	}

	return allBytes, nil
}

type (
	FileHeaderDef struct {
		Filetype  string
		MagicData []byte
	}

	FileHeaderResult struct {
		Filetype string
		Header   bool
	}
)

func (f FileHeaderResult) Description() string {
	if f.Header {
		return f.Filetype + " header"
	} else {
		return f.Filetype + " footer"
	}
}

var fileHeaders = []FileHeaderDef{
	{"doc", []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	{"xls", []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	{"ppt", []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	{"zip", []byte{0x50, 0x4B, 0x03, 0x04, 0x14}}, // probably wrong
	{"jpg", []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01}},
	{"gif", []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x4E, 0x01, 0x53, 0x00, 0xC4}},
	{"pdf", []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E}}, // probably wrong
	{"7zip", []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}},      // verified at endchan
	{"Torrent", []byte{0x64, 0x38, 0x3A, 0x61, 0x6E, 0x6E, 0x6F, 0x75, 0x6E, 0x63, 0x65}},
	{"AVI", []byte{0x52, 0x49, 0x46, 0x46}},
}

var fileFooters = []FileHeaderDef{
	{"doc", []byte{0x57, 0x6F, 0x72, 0x64, 0x2E, 0x44, 0x6F, 0x63, 0x75, 0x6D, 0x65, 0x6E, 0x74, 0x2E}},
	{"xls", []byte{0xFE, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x57, 0x00, 0x6F, 0x00, 0x72, 0x00, 0x6B, 0x00, 0x62, 0x00, 0x6F, 0x00, 0x6F, 0x00, 0x6B, 0x00}},
	{"ppt", []byte{0x50, 0x00, 0x6F, 0x00, 0x77, 0x00, 0x65, 0x00, 0x72, 0x00, 0x50, 0x00, 0x6F, 0x00, 0x69, 0x00, 0x6E, 0x00, 0x74, 0x00, 0x20, 0x00, 0x44, 0x00, 0x6F, 0x00, 0x63, 0x00, 0x75, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x6E, 0x00, 0x74}},
	{"zip", []byte{0x50, 0x4B, 0x05, 0x06, 0x00}}, // probably wrong
	{"gif", []byte{0x21, 0x00, 0x00, 0x3B, 0x00}},
	{"pdf", []byte{0x25, 0x25, 0x45, 0x4F, 0x46}},  // probably wrong
	{"7zip", []byte{0x00, 0x00, 0x00, 0x17, 0x06}}, // verified with cablegate
}

func SearchDataForKnownFileBits(data []byte) []FileHeaderResult {
	if data == nil {
		return []FileHeaderResult{}
	}

	chHeaderMatches := make(chan FileHeaderDef)
	go func() {
		for _, header := range fileHeaders {
			if bytes.Contains(data, header.MagicData) {
				chHeaderMatches <- header
			}
		}
		close(chHeaderMatches)
	}()

	chFooterMatches := make(chan FileHeaderDef)
	go func() {
		for _, footer := range fileFooters {
			if bytes.Contains(data, footer.MagicData) {
				chFooterMatches <- footer
			}
		}
		close(chFooterMatches)
	}()

	matches := []FileHeaderResult{}
	for match := range chHeaderMatches {
		matches = append(matches, FileHeaderResult{Filetype: match.Filetype, Header: true})
	}

	for match := range chFooterMatches {
		matches = append(matches, FileHeaderResult{Filetype: match.Filetype, Header: false})
	}

	return matches
}
