package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

func SearchForPlaintext(startBlock, endBlock uint64, inDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "search-plaintext")

	err := os.MkdirAll(outSubdir, 0777)
	if err != nil {
		return err
	}

	// start a goroutine to log errors
	chErr := make(chan error)
	go func() {
		for err := range chErr {
			fmt.Println("error:", err)
		}
	}()

	// start a goroutine for each .dat file being parsed
	chDones := []chan bool{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		chDone := make(chan bool)
		go searchForPlaintextParseBlock(inDir, outSubdir, i, chErr, chDone)
		chDones = append(chDones, chDone)
	}

	// wait for all ops to complete
	for _, chDone := range chDones {
		<-chDone
	}

	// close error channel
	close(chErr)

	return nil
}

func searchForPlaintextParseBlock(inDir string, outDir string, blockFileNum int, chErr chan error, chDone chan bool) {
	defer close(chDone)

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(inDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	outFile, err := utils.CreateFile(filepath.Join(outDir, fmt.Sprintf("blk%05d-plaintext.txt", blockFileNum)))
	if err != nil {
		chErr <- err
		return
	}
	defer utils.CloseFile(outFile)

	for _, bl := range blocks {
		fmt.Println("==================", bl.MsgBlock().Header.Timestamp, "==================")
		blockHash := bl.Hash().String()

		for _, tx := range bl.Transactions() {
			txHash := tx.Hash().String()

			// extract text from each TxIn scriptSig
			for txinIdx, txin := range tx.MsgTx().TxIn {
				txt, isText := extractText(txin.SignatureScript)
				if !isText || len(txt) < 8 {
					continue
				}

				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "in", txinIdx, string(txt)))
				if err != nil {
					chErr <- err
					return
				}
			}

			// extract text from each TxOut PkScript
			for txoutIdx, txout := range tx.MsgTx().TxOut {
				txt, isText := extractText(txout.PkScript)
				if !isText || len(txt) < 8 {
					continue
				}

				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", txoutIdx, string(txt)))
				if err != nil {
					chErr <- err
					return
				}
			}

			// extract text from concatenated TxOut hex tokens

			parsedScriptData, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
			if err != nil {
				chErr <- err
				return
			}

			parsedScriptText, isText := extractText(parsedScriptData)
			if err != nil {
				chErr <- err
				return
			}

			if isText && len(parsedScriptText) > 8 {
				fmt.Println(string(parsedScriptText))
				_, err := outFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", blockHash, txHash, "out", -1, string(parsedScriptText)))
				if err != nil {
					chErr <- err
					return
				}
			}
		}
	}

	if err != nil {
		chErr <- err
		return
	}
}

func extractText(bs []byte) ([]byte, bool) {
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
