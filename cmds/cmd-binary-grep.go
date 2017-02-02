package cmds

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

type BinaryGrepCommand struct {
	startBlock, endBlock uint64
	blocks               []int
	datFileDir           string
	outDir               string
	hexPattern           string
	carveLen             uint64
	carveExt             string
}

func NewBinaryGrepCommand(blocks []int, carveLen uint64, carveExt, outDir, datFileDir, hexPattern string) *BinaryGrepCommand {
	// if endBlock == 0 {
	// 	endBlock = startBlock
	// }

	return &BinaryGrepCommand{
		// startBlock: startBlock,
		// endBlock:   endBlock,
		blocks:     blocks,
		datFileDir: datFileDir,
		outDir:     filepath.Join(".", outDir, "binary-grep"),
		carveLen:   carveLen,
		carveExt:   carveExt,
		hexPattern: hexPattern,
	}
}

func (cmd *BinaryGrepCommand) RunCommand() error {
	pattern, err := hex.DecodeString(cmd.hexPattern)
	if err != nil {
		return err
	}

	err = os.MkdirAll(cmd.outDir, 0777)
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

	// start a goroutine for each .dat file being parsed, limited to 5 at a time
	chDones := []chan bool{}
	procLimiter := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		procLimiter <- true
	}

	chResults := make(chan string)
	chResultsDone := make(chan bool)
	go func() {
		defer close(chResultsDone)
		for x := range chResults {
			os.Stdout.WriteString(x + "\n")
		}
	}()

	// for i := int(cmd.startBlock); i < int(cmd.endBlock)+1; i++ {
	for _, i := range cmd.blocks {
		chDone := make(chan bool)
		chDones = append(chDones, chDone)
		go cmd.parseBlock(i, pattern, chResults, chErr, chDone, procLimiter)
	}

	// wait for all goroutines to complete
	for _, chDone := range chDones {
		<-chDone
	}

	// close results channel
	close(chResults)

	// close error channel
	close(chErr)

	<-chResultsDone

	return nil
}

func (cmd *BinaryGrepCommand) parseBlock(blockFileNum int, pattern []byte, chResults chan string, chErr chan error, chDone chan bool, procLimiter chan bool) {
	defer close(chDone)
	defer func() { procLimiter <- true }()
	<-procLimiter

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	os.Stderr.WriteString("parsing block " + filename + "\n")

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(cmd.datFileDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	for _, bl := range blocks {
		blockHash := bl.Hash().String()

		for _, btctx := range bl.Transactions() {
			tx := Tx{Tx: btctx}

			txHash := tx.Hash().String()

			inData, err := tx.ConcatTxInScripts()
			if err != nil {
				chErr <- err
				return
			}

			// it is unlikely that a given pattern will be found more than once in the scripts,
			// given how short they are, so we don't loop the search

			offset := bytes.Index(inData, pattern)
			if offset > -1 {
				chResults <- fmt.Sprintf("%v,%v,%v,%v,%v", filename, blockHash, txHash, "in", offset)
				if cmd.carveLen > 0 {
					cmd.carve(filename, txHash, "in", offset, inData)
				}
			}

			outData, err := tx.ConcatNonOPDataFromTxOuts()
			if err != nil {
				chErr <- err
				return
			}

			offset = bytes.Index(outData, pattern)
			if offset > -1 {
				chResults <- fmt.Sprintf("%v,%v,%v,%v,%v", filename, blockHash, txHash, "out", offset)
				if cmd.carveLen > 0 {
					cmd.carve(filename, txHash, "out", offset, outData)
				}
			}
		}
	}

	if err != nil {
		chErr <- err
		return
	}
}

func (cmd *BinaryGrepCommand) carve(filename string, txHash string, inOut string, offset int, fileData []byte) {
	if offset >= len(fileData) {
		return
	}

	end := offset + int(cmd.carveLen)
	if end > len(fileData) {
		end = len(fileData)
	}

	outFilename := fmt.Sprintf("%s-%s-%s-%d.%s", filename, txHash, inOut, offset, cmd.carveExt)
	err := ioutil.WriteFile(filepath.Join(cmd.outDir, outFilename), fileData[offset:end], 0666)
	if err != nil {
		panic(err)
	}
	// os.Stderr.WriteString("wrote " + outFilename + "\n")
}
