package blkparser

import (
	"crypto/sha256"
	"fmt"
)

// Get the Tx count, decode the variable length integer
// https://en.bitcoin.it/wiki/Protocol_specification#Variable_length_integer
func DecodeVariableLengthInteger(raw []byte) (cnt int, size int) {
	if raw[0] < 0xfd {
		return int(raw[0]), 1
	}
	size = 1 + (2 << (2 - (0xff - raw[0])))
	if len(raw) < 1+size {
		return 0, size
	}

	res := uint64(0)
	for i := 1; i < size; i++ {
		res |= (uint64(raw[i]) << uint64(8*(i-1)))
	}

	return int(res), size
}

func GetShaString(data []byte) string {
	sha := sha256.New()
	sha.Write(data[:])
	tmp := sha.Sum(nil)
	sha.Reset()
	sha.Write(tmp)
	hash := sha.Sum(nil)
	return HashString(hash)
}

func HashString(data []byte) string {
	res := ""
	for i := 0; i < 32; i++ {
		res += fmt.Sprintf("%02x", data[31-i])
	}
	return res
}
