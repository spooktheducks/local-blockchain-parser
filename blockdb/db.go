package blockdb

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type (
	BlockDB struct {
		store      *bolt.DB
		datFileDir string
	}
)

func NewBlockDB(dbFilename string, datFileDir string) (*BlockDB, error) {
	// open the BoltDB file
	store, err := bolt.Open(dbFilename, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BlockDB{store: store, datFileDir: datFileDir}, nil
}

func (db *BlockDB) Close() error {
	return db.store.Close()
}

func (db *BlockDB) IndexDATFileBlocks(startBlock, endBlock uint64) error {
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		datFilename := fmt.Sprintf("blk%05d.dat", i)
		datFilepath := filepath.Join(db.datFileDir, datFilename)

		fmt.Println("parsing block file", datFilepath)

		blocks, err := utils.LoadBlocksFromDAT(datFilepath)
		if err != nil {
			return err
		}

		err = db.writeBlockIndexToDB(blocks, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BlockDB) IndexDATFileTransactions(startBlock, endBlock uint64) error {
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		datFilename := fmt.Sprintf("blk%05d.dat", i)
		datFilepath := filepath.Join(db.datFileDir, datFilename)

		fmt.Println("parsing block file", datFilepath)

		blocks, err := utils.LoadBlocksFromDAT(datFilepath)
		if err != nil {
			return err
		}

		err = db.writeBlockIndexToDB(blocks, i)
		if err != nil {
			return err
		}

		err = db.writeTxIndexToDB(blocks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BlockDB) writeBlockIndexToDB(blocks []*btcutil.Block, datFileIdx int) error {
	fmt.Println("writing block metadata...")

	// we break the blocks into a bunch of smaller groups because BoltDB writes much more quickly this way
	groupLen := 10
	blockGroups := utils.GroupBlocks(blocks, groupLen)

	for g, group := range blockGroups {
		err := db.store.Update(func(boltTx *bolt.Tx) error {
			bucket, err := boltTx.CreateBucketIfNotExists([]byte("BlockIndex"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for blIdx, bl := range group {
				row := BlockIndexRow{
					DATFileIdx:     uint16(datFileIdx),
					Timestamp:      bl.MsgBlock().Header.Timestamp.Unix(),
					IndexInDATFile: uint32((g * groupLen) + blIdx),
				}

				rowBytes, err := row.ToBytes()
				if err != nil {
					return err
				}

				err = bucket.Put(bl.Hash()[:], rowBytes)
				if err != nil {
					return err
				}

				fmt.Printf("finished block %v (%v) (%v/%v)\n", bl.Hash().String(), bl.MsgBlock().Header.Timestamp, (g*groupLen)+blIdx+1, len(blocks))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BlockDB) writeTxIndexToDB(blocks []*btcutil.Block) error {
	fmt.Println("writing transaction index...")

	// we break the blocks into a bunch of smaller groups because BoltDB writes much more quickly this way
	groupLen := 10
	blockGroups := utils.GroupBlocks(blocks, groupLen)
	numBlocks := len(blocks)

	for g, group := range blockGroups {
		err := db.store.Update(func(boltTx *bolt.Tx) error {
			bucket, err := boltTx.CreateBucketIfNotExists([]byte("TransactionIndex"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for blkIdx, bl := range group {
				numTxs := len(bl.Transactions())

				for txIdx, tx := range bl.Transactions() {
					row := TxIndexRow{
						BlockHash:    *bl.Hash(),
						IndexInBlock: uint64(txIdx),
					}

					rowBytes, err := row.ToBytes()
					if err != nil {
						return err
					}

					err = bucket.Put(tx.Hash()[:], rowBytes)
					if err != nil {
						return err
					}

					fmt.Printf("finished tx %v (%v/%v) (%v/%v)\n", tx.Hash().String(), txIdx+1, numTxs, (g*groupLen)+blkIdx+1, numBlocks)
					// fmt.Printf("finished tx (%v/%v) (%v/%v)\n", txIdx+1, numTxs, (g*groupLen)+blkIdx+1, numBlocks)
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BlockDB) GetBlockIndexRow(blockHash chainhash.Hash) (BlockIndexRow, error) {
	var err error
	var blockRow BlockIndexRow

	err = db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte("BlockIndex"))
		if bucket == nil {
			return fmt.Errorf("could not find bucket BlockIndex")
		}

		val := bucket.Get(blockHash[:])
		if val == nil {
			return fmt.Errorf("could not find block %v", blockHash.String())
		}

		blockRow, err = NewBlockIndexRowFromBytes(val)
		if err != nil {
			return err
		}

		return nil
	})

	return blockRow, err
}

func (db *BlockDB) GetBlock(blockHash chainhash.Hash) (*btcutil.Block, error) {
	blockRow, err := db.GetBlockIndexRow(blockHash)
	if err != nil {
		return nil, err
	}

	block, err := utils.LoadBlockFromDAT(filepath.Join(db.datFileDir, blockRow.DATFilename()), blockRow.IndexInDATFile)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (db *BlockDB) GetTxIndexRow(txHash chainhash.Hash) (TxIndexRow, BlockIndexRow, error) {
	var err error
	var txRow TxIndexRow
	var blockRow BlockIndexRow

	err = db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte("TransactionIndex"))
		if bucket == nil {
			return fmt.Errorf("could not find bucket TransactionIndex")
		}

		val := bucket.Get(txHash[:])
		if val == nil {
			return fmt.Errorf("could not find transaction %v", txHash)
		}

		txRow, err = NewTxIndexRowFromBytes(val)
		if err != nil {
			return err
		}

		bucket = boltTx.Bucket([]byte("BlockIndex"))
		if bucket == nil {
			return fmt.Errorf("could not find bucket BlockIndex")
		}

		val = bucket.Get(txRow.BlockHash[:])
		if val == nil {
			return fmt.Errorf("could not find block %v", txRow.BlockHash)
		}

		blockRow, err = NewBlockIndexRowFromBytes(val)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return TxIndexRow{}, BlockIndexRow{}, err
	}

	return txRow, blockRow, nil
}

func (db *BlockDB) GetTx(txHash chainhash.Hash) (*btcutil.Tx, error) {
	txRow, blockRow, err := db.GetTxIndexRow(txHash)
	if err != nil {
		return nil, err
	}

	block, err := utils.LoadBlockFromDAT(filepath.Join(db.datFileDir, blockRow.DATFilename()), blockRow.IndexInDATFile)
	if err != nil {
		return nil, err
	}

	tx, err := block.Tx(int(txRow.IndexInBlock))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (db *BlockDB) IndexDATFileTxOutDuplicates(startBlock, endBlock uint64) error {
	blockDATFiles := []string{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		blockDATFiles = append(blockDATFiles, fmt.Sprintf("blk%05d.dat", i))
	}

	for _, datFilename := range blockDATFiles {
		datFilepath := filepath.Join(db.datFileDir, datFilename)

		fmt.Println("parsing block file", datFilepath)

		blocks, err := utils.LoadBlocksFromDAT(datFilepath)
		if err != nil {
			return err
		}

		for _, bl := range blocks {
			for _, tx := range bl.Transactions() {
				// txHashBytes := tx.Hash().CloneBytes()

				data, err := utils.ConcatNonOPHexTokensFromTxOuts(tx)
				if err != nil {
					return err
				}

				err = db.PutTxOutDuplicateData(*tx.Hash(), data)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (db *BlockDB) PutTxOutDuplicateData(txHash chainhash.Hash, data []byte) error {
	// if len(txHashBytes) != chainhash.HashSize {
	// 	return fmt.Errorf("txHashBytes must be %v bytes long", chainhash.HashSize)
	// }

	err := db.store.Update(func(boltTx *bolt.Tx) error {
		bucket, err := boltTx.CreateBucketIfNotExists([]byte("TxOutDupes"))
		if err != nil {
			return err
		}

		hasher := sha256.New()
		_, err = hasher.Write(data)
		if err != nil {
			return err
		}
		hashedData := hasher.Sum(nil)

		existing := bucket.Get(hashedData)
		if existing == nil {
			existing = txHash[:]
		} else {
			existing = append(existing, txHash[:]...)
		}

		err = bucket.Put(hashedData, existing)
		return err
	})

	return err
}

func (db *BlockDB) ReadTxOutDuplicateData() error {
	err := db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte("TxOutDupes"))
		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(key []byte, val []byte) error {
			if len(val)%chainhash.HashSize != 0 {
				return fmt.Errorf("value is corrupted")
			}

			numTxs := len(val) / chainhash.HashSize

			if numTxs == 1 {
				return nil
			}

			fmt.Printf("- %v txs sharing data:\n", numTxs)
			for i := 0; i < numTxs; i++ {
				txHash, err := HashFromBytes(val[i*chainhash.HashSize : (i+1)*chainhash.HashSize])
				if err != nil {
					return err
				}

				fmt.Printf("  - %v\n", txHash.String())
			}

			return nil
		})

		return err
	})

	return err
}

func (db *BlockDB) IndexDATFileSpentTxOuts(startBlock, endBlock uint64) error {
	blockDATFiles := []string{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		blockDATFiles = append(blockDATFiles, fmt.Sprintf("blk%05d.dat", i))
	}

	for _, datFilename := range blockDATFiles {
		datFilepath := filepath.Join(db.datFileDir, datFilename)

		fmt.Println("parsing block file", datFilepath)

		blocks, err := utils.LoadBlocksFromDAT(datFilepath)
		if err != nil {
			return err
		}

		for _, bl := range blocks {
			for _, tx := range bl.Transactions() {

				for txinIdx, txin := range tx.MsgTx().TxIn {
					key := SpentTxOutKey{TxHash: txin.PreviousOutPoint.Hash, TxOutIndex: txin.PreviousOutPoint.Index}
					val := SpentTxOutRow{InputTxHash: *tx.Hash(), TxInIndex: uint32(txinIdx)}

					err = db.PutSpentTxOut(key, val)
					if err != nil {
						return err
					}
				}
			}
		}

	}

	return nil
}

func (db *BlockDB) PutSpentTxOut(key SpentTxOutKey, val SpentTxOutRow) error {
	err := db.store.Update(func(boltTx *bolt.Tx) error {
		bucket, err := boltTx.CreateBucketIfNotExists([]byte("SpentTxOuts"))
		if err != nil {
			return err
		}

		keyBytes, err := key.ToBytes()
		if err != nil {
			return err
		}

		valBytes, err := val.ToBytes()
		if err != nil {
			return err
		}

		err = bucket.Put(keyBytes, valBytes)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (db *BlockDB) GetSpentTxOut(key SpentTxOutKey) (SpentTxOutRow, error) {
	var row SpentTxOutRow
	err := db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte("SpentTxOuts"))
		if bucket == nil {
			return fmt.Errorf("can't find bucket SpentTxOuts")
		}

		keyBytes, err := key.ToBytes()
		if err != nil {
			return err
		}

		valBytes := bucket.Get(keyBytes)
		if valBytes == nil {
			return fmt.Errorf("can't find SpentTxOut %+v", key)
		}

		row, err = NewSpentTxOutRowFromBytes(valBytes)
		if err != nil {
			return err
		}

		return nil
	})

	return row, err
}
