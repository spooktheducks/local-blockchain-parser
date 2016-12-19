package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

func GetSatoshiEncodedData(data []byte) ([]byte, error) {
	length := binary.LittleEndian.Uint32(data[0:4])
	checksum := binary.LittleEndian.Uint32(data[4:8])
	if len(data) < 8+int(length) {
		return nil, fmt.Errorf("GetSatoshiEncodedData: not enough data")
	}

	data = data[8 : 8+length]

	if crc32.ChecksumIEEE(data) != checksum {
		return nil, fmt.Errorf("GetSatoshiEncodedData: crc32 failed")
	}
	return data, nil
}

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
				return nil, errors.New("GetOPReturnBytes: empty OP_RETURN data")
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
	FileMagicBytesDef struct {
		Filetype  string
		MagicData []byte
	}

	FileMagicBytesResult struct {
		Filetype string
		Reversed bool
	}
)

func (f FileMagicBytesResult) Description() string {
	if f.Reversed {
		return f.Filetype + " (reversed)"
	} else {
		return f.Filetype
	}
}

var fileMagicBytes = []FileMagicBytesDef{
	{"DOC Header", []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}},
	{"DOC Footer", []byte{0x57, 0x6f, 0x72, 0x64, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x2e}},
	{"XLS Header", []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}},
	{"XLS Footer", []byte{0xfe, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x57, 0x00, 0x6f, 0x00, 0x72, 0x00, 0x6b, 0x00, 0x62, 0x00, 0x6f, 0x00, 0x6f, 0x00, 0x6b, 0x00}},
	{"PPT Header", []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}},
	{"PPT Footer", []byte{0xa0, 0x46, 0x1d, 0xf0}},
	{"ZIP Header", []byte{0x50, 0x4b, 0x03, 0x04, 0x14}},
	{"ZIP Footer", []byte{0x50, 0x4b, 0x05, 0x06, 0x00}},
	{"ZIPLock Footer", []byte{0x50, 0x4b, 0x03, 0x04, 0x14, 0x00, 0x01, 0x00, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"JPG Header", []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01}},
	{"GIF Header", []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}},
	{"GIF Footer", []byte{0x21, 0x00, 0x00, 0x3b, 0x00}},
	{"PDF Header", []byte{0x25, 0x50, 0x44, 0x46}},
	{"PDF Header (alternate)", []byte{0x26, 0x23, 0x32, 0x30, 0x35}},
	{"PDF Footer", []byte{0x25, 0x25, 0x45, 0x4f, 0x46}},
	{"Torrent Header", []byte{0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65}},
	{"GZ Header", []byte{0x1f, 0x8b, 0x08, 0x08}},
	{"TAR Header", []byte{0x1f, 0x8b, 0x08, 0x00}},
	{"TAR.GZ Header", []byte{0x1f, 0x9d, 0x90, 0x70}},
	{"EPUB Header", []byte{0x50, 0x4b, 0x03, 0x04, 0x0a, 0x00, 0x02, 0x00}},
	{"PNG Header", []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}},
	{"8192 Header", []byte{0x6d, 0x51, 0x51, 0x4e, 0x42}},
	{"4096 Header", []byte{0x6d, 0x51, 0x49, 0x4e, 0x42, 0x46, 0x67, 0x2f}},
	{"2048 Header", []byte{0x95, 0x2e, 0x3e, 0x2e, 0x58, 0x4b, 0x7a}},
	{"Secret Header", []byte{0x52, 0x61, 0x72, 0x21, 0x1a, 0x07, 0x00}},
	{"RAR Header", []byte{0x6d, 0x51, 0x45, 0x4e, 0x42, 0x46, 0x67}},
	{"OGG Header", []byte{0x4f, 0x67, 0x67, 0x53}},
	{"WAV Header", []byte{0x42, 0x49, 0x46, 0x46}},
	{"WAV Header (alternate)", []byte{0x57, 0x41, 0x56, 0x45}},
	{"AVI Header", []byte{0x42, 0x49, 0x46, 0x46}},
	{"AVI Header (alternate)", []byte{0x41, 0x56, 0x49, 0x20}},
	{"MIDI Header", []byte{0x4d, 0x54, 0x68, 0x64}},
	{"7z Header", []byte{0x37, 0x7a, 0xbc, 0xaf, 0x27, 0x1c}},
	{"7z Footer", []byte{0x00, 0x00, 0x00, 0x17, 0x06}},
	{"DMG Header", []byte{0x78, 0x01, 0x73, 0x0d, 0x62, 0x62, 0x60}},
	{"Wikileaks", []byte{0x57, 0x69, 0x6b, 0x69, 0x6c, 0x65, 0x61, 0x6b, 0x73}},
	{"Julian Assange", []byte{0x4a, 0x75, 0x6c, 0x69, 0x61, 0x6e, 0x20, 0x41, 0x73, 0x73, 0x61, 0x6e, 0x67, 0x65}},
	{"Mendax", []byte{0x4d, 0x65, 0x6e, 0x64, 0x61, 0x7}},
}

func SearchDataForMagicFileBytes(data []byte) []FileMagicBytesResult {
	if data == nil {
		return []FileMagicBytesResult{}
	}

	chMatches := make(chan []FileMagicBytesResult)
	go func() {
		matches := []FileMagicBytesResult{}
		for _, header := range fileMagicBytes {
			if bytes.Contains(data, header.MagicData) {
				matches = append(matches, FileMagicBytesResult{Filetype: header.Filetype, Reversed: false})
			}
		}
		chMatches <- matches
	}()

	chMatchesReversed := make(chan []FileMagicBytesResult)
	go func() {
		matches := []FileMagicBytesResult{}
		for _, header := range fileMagicBytes {
			if bytes.Contains(data, ReverseBytes(header.MagicData)) {
				matches = append(matches, FileMagicBytesResult{Filetype: header.Filetype, Reversed: true})
			}
		}
		chMatchesReversed <- matches
	}()

	matches := []FileMagicBytesResult{}
	matches = append(matches, <-chMatches...)
	matches = append(matches, <-chMatchesReversed...)

	return matches
}

func ExtractText(bs []byte) ([]byte, bool) {
	start := 0

	for start < len(bs) {
		if isValidPlaintextByte(bs[start]) {
			break
		}
		start++
	}
	if start == len(bs) {
		return nil, false
	}

	end := start
	for end < len(bs) {
		if !isValidPlaintextByte(bs[end]) {
			break
		}
		end++
	}

	sublen := end - start + 1
	if sublen < 5 {
		return nil, false
	}

	substr := bs[start:end]
	return substr, true
}

func stripNonTextBytes(bs []byte) []byte {
	newBs := make([]byte, len(bs))
	newBsLen := 0
	for i := range bs {
		if isValidPlaintextByte(bs[i]) {
			newBs[newBsLen] = bs[i]
			newBsLen++
		}
	}

	if newBsLen == 0 {
		return nil
	}

	return newBs[0:newBsLen]
}

func isValidPlaintextByte(x byte) bool {
	switch x {
	case '\r', '\n', '\t', ' ':
		return true
	}

	i := int(rune(x))
	if i >= 32 && i < 127 {
		return true
	}

	return false
}
