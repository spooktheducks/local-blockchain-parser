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
		strs[i] = fmt.Sprintf("%+v", p)
	}
	return strs
}

func FindPGPPackets(data []byte) PGPPacketResult {
	packets := []packet.Packet{}

	reader := packet.NewReader(bytes.NewReader(data))
	for {
		packet, err := reader.Next()
		if err != nil {
			break
		}
		packets = append(packets, packet)
		// if isSatoshi {
		// fmt.Printf("  - GPG packet (satoshi-encoded): %+v\n", packet)
		// } else {
		// fmt.Printf("  - GPG packet: %+v\n", packet)
		// }
	}

	return PGPPacketResult{Packets: packets}
}
