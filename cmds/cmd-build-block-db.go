package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcd/txscript"
	// "github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/utils"
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

	// fill up our file semaphore so we can obtain tokens from it
	for i := 0; i < maxFiles; i++ {
		fileSemaphore <- true
	}

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

	<-fileSemaphore
	blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
	fileSemaphore <- true
	if err != nil {
		chErr <- err
		return
	}

	err = db.Update(func(boltTx *bolt.Tx) error {
		b, err := boltTx.CreateBucket([]byte("Transactions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for _, bl := range blocks {
			blockHash := bl.Hash().String()

			for _, tx := range bl.Transactions() {
				txHash := tx.Hash().String()

				for txoutIdx, txout := range tx.MsgTx().TxOut {
					scriptStr, err := txscript.DisasmString(txout.PkScript)
					if err != nil {
						return err
					}

					err = b.Put([]byte(fmt.Sprintf("%v:%v:%v", blockHash, txHash, txoutIdx)), []byte(scriptStr))
					if err != nil {
						return err
					}
					fmt.Println("put succeeded:", txHash)
				}
			}
		}

		return nil
	})

	if err != nil {
		chErr <- err
		return
	}
}
