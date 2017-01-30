package utils

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/openpgp/packet"
)

type PGPPacketResult struct {
	Packets []packet.Packet
}

func (r PGPPacketResult) IsEmpty() bool {
	return len(r.Packets) == 0
}

func (r PGPPacketResult) DescriptionStrings() []string {
	strs := make([]string, len(r.Packets))
	for i, p := range r.Packets {
		strs[i] = fmt.Sprintf("%T", p)
	}
	return strs
}

func FindPGPPackets(data []byte) PGPPacketResult {
	packets := []packet.Packet{}

	for i := range data {
		reader := packet.NewReader(bytes.NewReader(data[i:]))
		for {
			packet, err := reader.Next()
			if err != nil {
				break
			}
			packets = append(packets, packet)
		}
	}

	return PGPPacketResult{Packets: packets}
}
