package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/btcsuite/btcd/txscript"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

func CheckScriptPatterns(startBlock, endBlock uint64, inDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "script-patterns")

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

	patterns := map[string]struct{}{}
	chPatterns := make(chan string)
	go func() {
		for pattern := range chPatterns {
			patterns[pattern] = struct{}{}
		}
	}()

	// start a goroutine for each .dat file being parsed
	chDones := []chan bool{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		chDone := make(chan bool)
		go scriptPatternsParseBlock(inDir, outSubdir, i, chPatterns, chErr, chDone)
		chDones = append(chDones, chDone)
	}

	// wait for all ops to complete
	for _, chDone := range chDones {
		<-chDone
	}

	close(chPatterns)

	for k := range patterns {
		fmt.Println(k)
	}

	// close error channel
	close(chErr)

	return nil
}

func scriptPatternsParseBlock(inDir string, outDir string, blockFileNum int, chPatterns chan string, chErr chan error, chDone chan bool) {
	defer close(chDone)

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	regex := regexp.MustCompile("(?:^| )([^(?:OP_)]+)(?:$| )")

	for _, bl := range blocks {
		// blockHash := bl.Hash().String()

		for _, tx := range bl.Transactions() {
			for _, txout := range tx.MsgTx().TxOut {
				scriptStr, err := txscript.DisasmString(txout.PkScript)
				if err != nil {
					if err.Error() == "execute past end of script" {
						continue
					} else {
						chErr <- fmt.Errorf("error in txscript.DisasmString: %v", err)
						return
					}
				}

				newBs := regex.ReplaceAllFunc([]byte(scriptStr), func(bs []byte) []byte {
					x := fmt.Sprintf("[chunk size=%v]", len(bs))
					if string(bs[0]) == " " {
						x = " " + x
					}
					if string(bs[len(bs)-1]) == " " {
						x = x + " "
					}
					return []byte(x)
				})

				chPatterns <- string(newBs)
			}
		}
	}
}
