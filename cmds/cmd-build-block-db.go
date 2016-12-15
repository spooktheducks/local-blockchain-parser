package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	// "github.com/btcsuite/btcd/txscript"
	// "github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

func BuildBlockDB(startBlock, endBlock uint64, inDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "blockdb")

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

	db, err := bolt.Open(filepath.Join(outSubdir, "blockchain.db"), 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// start a goroutine for each .dat file being parsed
	chDones := []chan bool{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		chDone := make(chan bool)
		go buildDBParseBlock(inDir, outSubdir, db, i, chErr, chDone)
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

func buildDBParseBlock(inDir string, outDir string, db *bolt.DB, blockFileNum int, chErr chan error, chDone chan bool) {
	defer close(chDone)

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
	if err != nil {
		chErr <- err
		return
	}
	fmt.Println("writing TxOut scripts...")

	groupLen := 10
	blockGroups := utils.GroupBlocks(blocks, groupLen)

	for g, group := range blockGroups {
		err := db.Update(func(boltTx *bolt.Tx) error {
			b, err := boltTx.CreateBucketIfNotExists([]byte("TxOutScripts"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for blIdx, bl := range group {
				blockHash := bl.Hash().String()

				for _, tx := range bl.Transactions() {
					txHash := tx.Hash().String()

					for txoutIdx, txout := range tx.MsgTx().TxOut {
						err = b.Put([]byte(fmt.Sprintf("%v:%v:%v", blockHash, txHash, txoutIdx)), txout.PkScript)
						if err != nil {
							return err
						}
					}
				}
				fmt.Printf("finished block %v (%v/%v)\n", blockHash, (g*groupLen)+blIdx+1, len(blocks))
			}

			return nil
		})

		if err != nil {
			chErr <- err
			return
		}
	}
}
