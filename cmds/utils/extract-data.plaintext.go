package utils

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

	// sublen := end - start + 1
	// if sublen < 5 {
	//  return nil, false
	// }
	// fmt.Println("~~~~ EXTRACT TEXT", start, end)

	substr := bs[start:end]
	return substr, true
}

func StripNonTextBytes(bs []byte) []byte {
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
