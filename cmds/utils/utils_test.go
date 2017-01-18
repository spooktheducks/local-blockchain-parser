package utils

import (
	"testing"
)

func TestSatoshisToBTCs(T *testing.T) {
	s := Satoshis(1234567)
	b := s.ToBTC()

	if b != 0.01234567 {
		T.Fatalf("Expected 0.01234567")
	}
}

func TestBTCsToSatoshis(T *testing.T) {
	b := BTC(0.1234567)
	s := b.ToSatoshis()

	if s != 12345670 {
		T.Fatalf("Expected 12345670")
	}
}
