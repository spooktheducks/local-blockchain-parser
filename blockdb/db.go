package blockdb

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
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

		err = db.writeBlockIndexToDB(blocks, datFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BlockDB) IndexDATFileTransactions(startBlock, endBlock uint64) error {
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

		err = db.writeBlockIndexToDB(blocks, datFilename)
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

func (db *BlockDB) writeBlockIndexToDB(blocks []*btcutil.Block, datFilename string) error {
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
				blockHash := bl.Hash().String()

				err = bucket.Put([]byte(blockHash), []byte(fmt.Sprintf("%s:%d:%d", datFilename, (g*groupLen)+blIdx, bl.MsgBlock().Header.Timestamp.Unix())))
				if err != nil {
					return err
				}

				fmt.Printf("finished block %v (%v) (%v/%v)\n", blockHash, bl.MsgBlock().Header.Timestamp, (g*groupLen)+blIdx+1, len(blocks))
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
				blockHash := bl.Hash().String()
				numTxs := len(bl.Transactions())

				for txIdx, tx := range bl.Transactions() {
					txHash := tx.Hash().String()

					err := bucket.Put([]byte(txHash), []byte(fmt.Sprintf("%s:%d", blockHash, txIdx)))
					if err != nil {
						return err
					}

					fmt.Printf("finished tx %v (%v/%v) (%v/%v)\n", txHash, txIdx+1, numTxs, (g*groupLen)+blkIdx+1, numBlocks)
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

type (
	BlockIndexRow struct {
		Filename       string
		Timestamp      int64
		IndexInDATFile int
	}

	TxIndexRow struct {
		BlockHash    string
		IndexInBlock int
	}
)

func NewBlockIndexRowFromBytes(bs []byte) (BlockIndexRow, error) {
	parts := strings.Split(string(bs), ":")

	if len(parts) != 3 {
		return BlockIndexRow{}, fmt.Errorf("badly formatted BlockIndex value")
	}

	idx, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return BlockIndexRow{}, err
	}

	timestamp, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return BlockIndexRow{}, err
	}

	row := BlockIndexRow{
		Filename:       parts[0],
		IndexInDATFile: int(idx),
		Timestamp:      timestamp,
	}

	return row, nil
}

func NewTxIndexRowFromBytes(bs []byte) (TxIndexRow, error) {
	parts := strings.Split(string(bs), ":")

	if len(parts) != 2 {
		return TxIndexRow{}, fmt.Errorf("badly formatted TransactionIndex value")
	}

	txIndex, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return TxIndexRow{}, err
	}

	row := TxIndexRow{
		BlockHash:    parts[0],
		IndexInBlock: int(txIndex),
	}

	return row, nil
}

func (db *BlockDB) GetBlockIndexRow(blockHash string) (BlockIndexRow, error) {
	var err error
	var blockRow BlockIndexRow

	err = db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte("BlockIndex"))
		if bucket == nil {
			return fmt.Errorf("could not find bucket BlockIndex")
		}

		val := bucket.Get([]byte(blockHash))
		if val == nil {
			return fmt.Errorf("could not find block %v", blockHash)
		}

		blockRow, err = NewBlockIndexRowFromBytes(val)
		if err != nil {
			return err
		}

		return nil
	})

	return blockRow, err
}

func (db *BlockDB) GetBlock(blockHash string) (*btcutil.Block, error) {
	blockRow, err := db.GetBlockIndexRow(blockHash)
	if err != nil {
		return nil, err
	}

	block, err := utils.LoadBlockFromDAT(filepath.Join(db.datFileDir, blockRow.Filename), blockRow.IndexInDATFile)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (db *BlockDB) GetTxIndexRow(txHash string) (TxIndexRow, BlockIndexRow, error) {
	var err error
	var txRow TxIndexRow
	var blockRow BlockIndexRow

	err = db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte("TransactionIndex"))
		if bucket == nil {
			return fmt.Errorf("could not find bucket TransactionIndex")
		}

		val := bucket.Get([]byte(txHash))
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

		val = bucket.Get([]byte(txRow.BlockHash))
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

func (db *BlockDB) GetTx(txHash string) (*btcutil.Tx, error) {
	txRow, blockRow, err := db.GetTxIndexRow(txHash)
	if err != nil {
		return nil, err
	}

	block, err := utils.LoadBlockFromDAT(filepath.Join(db.datFileDir, blockRow.Filename), blockRow.IndexInDATFile)
	if err != nil {
		return nil, err
	}

	tx, err := block.Tx(txRow.IndexInBlock)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
