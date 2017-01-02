package txhashsource

import (
	"fmt"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/scanner"
)

func TestChainTxHashSource(T *testing.T) {
	datFileDir := os.Getenv("DAT_FILE_DIR")
	if datFileDir == "" {
		panic("must specify DAT_FILE_DIR enviroment variable (example: DAT_FILE_DIR=/path/to/dat/files go test)")
	}

	db, err := blockdb.NewBlockDB("../blockchain.db", datFileDir)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	startHash, err := blockdb.HashFromString("691dd277dc0e90a462a3d652a1171686de49cf19067cd33c7df0392833fb986a")
	if err != nil {
		panic(err)
	}

	src := scanner.NewChainTxHashSource(db, startHash)

	received := []chainhash.Hash{}
	for {
		hash, exists := src.NextHash()
		if !exists {
			break
		}

		received = append(received, hash)
	}

	for _, x := range received {
		fmt.Println(x)
	}
}
