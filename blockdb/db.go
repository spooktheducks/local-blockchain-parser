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

const (
	BucketBlockIndex       = "BlockIndex"
	BucketTransactionIndex = "TransactionIndex"
	BucketTxOutDupes       = "TxOutDupes"
	BucketSpentTxOuts      = "SpentTxOuts"
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

func (db *BlockDB) DATFilename(datIdx uint16) string {
	return filepath.Join(db.datFileDir, utils.DATFilename(datIdx))
}

func (db *BlockDB) LoadBlocksFromDAT(datIdx uint16) ([]*btcutil.Block, error) {
	return utils.LoadBlocksFromDAT(db.DATFilename(datIdx))
}

func (db *BlockDB) LoadBlockFromDAT(datIdx uint16, blockIdx uint32) (*btcutil.Block, error) {
	return utils.LoadBlockFromDAT(db.DATFilename(datIdx), blockIdx)
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
			bucket, err := boltTx.CreateBucketIfNotExists([]byte(BucketBlockIndex))
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

func (db *BlockDB) putBlockIndexRow(blockHash chainhash.Hash, row BlockIndexRow) error {
	return db.store.Update(func(boltTx *bolt.Tx) error {
		bucket, err := boltTx.CreateBucketIfNotExists([]byte(BucketBlockIndex))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		rowBytes, err := row.ToBytes()
		if err != nil {
			return err
		}

		return bucket.Put(blockHash[:], rowBytes)
	})
}

func (db *BlockDB) writeTxIndexToDB(blocks []*btcutil.Block) error {
	const writesPerBoltTx = 5000 // this value was determined by benchmarking, change at your own peril

	fmt.Println("writing transaction index...")

	keys := []chainhash.Hash{}
	vals := []TxIndexRow{}

	for _, bl := range blocks {
		for txIdx, tx := range bl.Transactions() {
			row := TxIndexRow{
				BlockHash:    *bl.Hash(),
				IndexInBlock: uint64(txIdx),
			}

			keys = append(keys, *tx.Hash())
			vals = append(vals, row)

			if len(keys) > writesPerBoltTx {
				err := db.putTxIndexRows(keys, vals)
				if err != nil {
					return err
				}

				keys = []chainhash.Hash{}
				vals = []TxIndexRow{}
			}
		}

		fmt.Printf("finished block %s\n", bl.Hash().String())
	}

	if len(keys) > 0 {
		err := db.putTxIndexRows(keys, vals)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BlockDB) putTxIndexRows(keys []chainhash.Hash, rows []TxIndexRow) error {
	err := db.store.Update(func(boltTx *bolt.Tx) error {
		bucket, err := boltTx.CreateBucketIfNotExists([]byte(BucketTransactionIndex))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for i := range keys {
			rowBytes, err := rows[i].ToBytes()
			if err != nil {
				return err
			}

			err = bucket.Put(keys[i][:], rowBytes)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (db *BlockDB) GetBlockIndexRow(blockHash chainhash.Hash) (BlockIndexRow, error) {
	if blockHash.String() == "0000000000000000017275d59d5ab479d0df454acad34227abf3d2911e253914" {
		// panic("DAT BLOCK")
		return BlockIndexRow{}, fmt.Errorf("dat block :(")
	}

	hasBytes := false
	for i := range blockHash {
		if blockHash[i] != 0x00 {
			hasBytes = true
			break
		}
	}
	if !hasBytes {
		panic("NO BYTES")
	}

	row, err := db.getBlockIndexRowFromDB(blockHash)
	if err == nil {
		return row, nil
	}

	switch err.(type) {
	case DataNotIndexedError, BlockNotFoundError:
		break
	default:
		return row, err
	}

	row, err = db.getBlockIndexRowFromDATFiles(blockHash, 0)
	if err == nil {
		err = db.putBlockIndexRow(blockHash, row)
		return row, err
	}

	return row, err
}

func (db *BlockDB) getBlockIndexRowFromDB(blockHash chainhash.Hash) (BlockIndexRow, error) {
	var err error
	var blockRow BlockIndexRow

	err = db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte(BucketBlockIndex))
		if bucket == nil {
			return DataNotIndexedError{Index: "blocks"}
		}

		val := bucket.Get(blockHash[:])
		if val == nil {
			return BlockNotFoundError{BlockHash: blockHash}
		}

		blockRow, err = NewBlockIndexRowFromBytes(val)
		if err != nil {
			return err
		}

		return nil
	})

	return blockRow, err
}

func (db *BlockDB) getBlockIndexRowFromDATFiles(blockHash chainhash.Hash, startDatIdx uint16) (BlockIndexRow, error) {
	for i := startDatIdx; ; i++ {
		fmt.Printf("\rsearching for block %v in DAT file %v", blockHash.String(), i)

		blocks, err := db.LoadBlocksFromDAT(i)
		if err != nil {
			return BlockIndexRow{}, err
		}

		for blkIdx, bl := range blocks {
			if *bl.Hash() == blockHash {
				row := BlockIndexRow{
					DATFileIdx:     i,
					Timestamp:      bl.MsgBlock().Header.Timestamp.Unix(),
					IndexInDATFile: uint32(blkIdx),
				}

				fmt.Printf("\r")
				return row, nil
			}
		}
	}

	return BlockIndexRow{}, fmt.Errorf("could not find block %v in DAT files", blockHash.String())
}

func (db *BlockDB) GetBlock(blockHash chainhash.Hash) (*btcutil.Block, error) {
	blockRow, err := db.GetBlockIndexRow(blockHash)
	if err != nil {
		return nil, err
	}

	block, err := db.LoadBlockFromDAT(blockRow.DATFileIdx, blockRow.IndexInDATFile)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (db *BlockDB) GetTxIndexRow(txHash chainhash.Hash) (TxIndexRow, error) {
	row, err := db.getTxIndexRowFromDB(txHash)
	if err == nil {
		return row, nil
	}

	switch err.(type) {
	case DataNotIndexedError, TxNotFoundError:
		break
	default:
		return TxIndexRow{}, err
	}

	row, err = db.getTxIndexRowFromBlockchainInfoAPI(txHash)
	if err == nil {
		err = db.putTxIndexRows([]chainhash.Hash{txHash}, []TxIndexRow{row})
		return row, err
	}
	fmt.Printf("error: tx index row %v not found\n", txHash.String())
	return row, err
}

func (db *BlockDB) getTxIndexRowFromDB(txHash chainhash.Hash) (TxIndexRow, error) {
	var err error
	var txRow TxIndexRow

	err = db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte(BucketTransactionIndex))
		if bucket == nil {
			return DataNotIndexedError{Index: "transactions"}
		}

		val := bucket.Get(txHash[:])
		if val == nil {
			return TxNotFoundError{TxHash: txHash}
		}

		txRow, err = NewTxIndexRowFromBytes(val)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return TxIndexRow{}, err
	}

	return txRow, nil
}

func (db *BlockDB) getTxIndexRowFromBlockchainInfoAPI(txHash chainhash.Hash) (TxIndexRow, error) {
	api := &BlockchainInfoAPI{}

	blockHash, err := api.GetBlockHashForTx(txHash)
	if err != nil {
		return TxIndexRow{}, err
	}

	blockRow, err := db.GetBlockIndexRow(blockHash)
	if err != nil {
		return TxIndexRow{}, err
	}

	bl, err := db.LoadBlockFromDAT(blockRow.DATFileIdx, blockRow.IndexInDATFile)
	if err != nil {
		return TxIndexRow{}, err
	}

	for txIdx, tx := range bl.Transactions() {
		if txHash == *tx.Hash() {
			row := TxIndexRow{
				BlockHash:    blockHash,
				IndexInBlock: uint64(txIdx),
			}
			return row, nil
		}
	}

	return TxIndexRow{}, fmt.Errorf("BlockDB.getTxIndexRowFromBlockchainInfoAPI: could not find transaction %v", txHash.String())
}

func (db *BlockDB) GetTx(txHash chainhash.Hash) (*btcutil.Tx, error) {
	txRow, err := db.GetTxIndexRow(txHash)
	if err != nil {
		return nil, err
	}

	blockRow, err := db.GetBlockIndexRow(txRow.BlockHash)
	if err != nil {
		fmt.Printf("block index row %v not found\n", txRow.BlockHash.String())
		return nil, err
	}

	block, err := db.LoadBlockFromDAT(blockRow.DATFileIdx, blockRow.IndexInDATFile)
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
				err = db.PutTxOutDuplicateData(tx)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (db *BlockDB) PutTxOutDuplicateData(tx *btcutil.Tx) error {
	data, err := utils.ConcatNonOPDataFromTxOuts(tx)
	if err != nil {
		return err
	}

	hasher := sha256.New()
	_, err = hasher.Write(data)
	if err != nil {
		return err
	}
	hashedData := hasher.Sum(nil)

	err = db.store.Update(func(boltTx *bolt.Tx) error {
		bucket, err := boltTx.CreateBucketIfNotExists([]byte(BucketTxOutDupes))
		if err != nil {
			return err
		}

		existing := bucket.Get(hashedData)
		if existing == nil {
			existing = tx.Hash()[:]
		} else {
			existing = append(existing, tx.Hash()[:]...)
		}

		err = bucket.Put(hashedData, existing)
		return err
	})

	return err
}

func (db *BlockDB) GetTxOutDuplicateData(txHash chainhash.Hash) ([]chainhash.Hash, error) {
	tx, err := db.GetTx(txHash)
	if err != nil {
		return nil, err
	}

	data, err := utils.ConcatNonOPDataFromTxOuts(tx)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	_, err = hasher.Write(data)
	if err != nil {
		return nil, err
	}
	hashedData := hasher.Sum(nil)

	var txListBytes []byte
	err = db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte(BucketTxOutDupes))
		if bucket == nil {
			return DataNotIndexedError{Index: "duplicates"}
		}

		txListBytes = bucket.Get(hashedData)
		return nil
	})

	return DecodeHashList(txListBytes)
}

func (db *BlockDB) ScanTxOutDuplicateData() error {
	err := db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte(BucketTxOutDupes))
		if bucket == nil {
			return DataNotIndexedError{Index: "duplicates"}
		}

		err := bucket.ForEach(func(key []byte, val []byte) error {
			txList, err := DecodeHashList(val)
			if err != nil {
				return err
			}

			if len(txList) == 0 {
				return nil
			}

			fmt.Printf("- %v txs sharing data:\n", len(txList))
			for _, txHash := range txList {
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

	keys := []SpentTxOutKey{}
	vals := []SpentTxOutRow{}

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

					keys = append(keys, key)
					vals = append(vals, val)

					if len(keys) == 40 {
						err = db.PutSpentTxOuts(keys, vals)
						if err != nil {
							return err
						}
						keys = []SpentTxOutKey{}
						vals = []SpentTxOutRow{}
					}
				}
			}
		}
	}

	if len(keys) > 0 {
		err := db.PutSpentTxOuts(keys, vals)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BlockDB) PutSpentTxOut(key SpentTxOutKey, val SpentTxOutRow) error {
	err := db.store.Update(func(boltTx *bolt.Tx) error {
		bucket, err := boltTx.CreateBucketIfNotExists([]byte(BucketSpentTxOuts))
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

func (db *BlockDB) PutSpentTxOuts(keys []SpentTxOutKey, vals []SpentTxOutRow) error {
	err := db.store.Update(func(boltTx *bolt.Tx) error {
		bucket, err := boltTx.CreateBucketIfNotExists([]byte(BucketSpentTxOuts))
		if err != nil {
			return err
		}

		for i := range keys {
			keyBytes, err := keys[i].ToBytes()
			if err != nil {
				return err
			}

			valBytes, err := vals[i].ToBytes()
			if err != nil {
				return err
			}

			err = bucket.Put(keyBytes, valBytes)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

var (
	errNotFoundDB = fmt.Errorf("not found in db")
)

func (db *BlockDB) GetSpentTxOut(key SpentTxOutKey) (SpentTxOutRow, error) {
	var row SpentTxOutRow
	err := db.store.View(func(boltTx *bolt.Tx) error {
		bucket := boltTx.Bucket([]byte(BucketSpentTxOuts))
		if bucket == nil {
			return errNotFoundDB //DataNotIndexedError{Index: "spent-txouts"}
		}

		keyBytes, err := key.ToBytes()
		if err != nil {
			return err
		}

		valBytes := bucket.Get(keyBytes)
		if valBytes == nil {
			return errNotFoundDB
		}

		row, err = NewSpentTxOutRowFromBytes(valBytes)
		if err != nil {
			return err
		}

		return nil
	})

	if err == errNotFoundDB {
		return row, err

		tx, err := db.GetTx(key.TxHash)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return row, err
		}

		fmt.Printf("requesting spent txout %v from api...", key)
		row, err = (&BlockchainInfoAPI{}).GetSpentTxOut(tx, key.TxOutIndex)
		if err == errBlockchainAPINotFound {
			fmt.Printf(" not found\n")
			return row, fmt.Errorf("can't find SpentTxOut %+v", key)
		} else if err != nil {
			fmt.Printf("error: %v\n", err)
			return row, err
		}
		// fmt.Printf("FOUND: %+v\n", row)
		// return row, err
	} else if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		err = db.PutSpentTxOut(key, row)
	}

	return row, err
}

func (db *BlockDB) GetSpentTxOutFromDATFiles(key SpentTxOutKey) (SpentTxOutRow, error) {
	datFileStartIndex := 0

	for datIdx := datFileStartIndex; ; datIdx++ {
		filename := filepath.Join(db.datFileDir, fmt.Sprintf("blk%05d.dat", datIdx))
		blocks, err := utils.LoadBlocksFromDAT(filename)
		if err != nil {
			return SpentTxOutRow{}, err
		}

		for _, bl := range blocks {
			for _, tx := range bl.Transactions() {

				for txinIdx, txin := range tx.MsgTx().TxIn {
					if key.TxHash == txin.PreviousOutPoint.Hash && key.TxOutIndex == txin.PreviousOutPoint.Index {
						fmt.Println("found SpentTxOutRow by searching dat files", tx.Hash().String())
						return SpentTxOutRow{InputTxHash: *tx.Hash(), TxInIndex: uint32(txinIdx)}, nil
					}
				}
			}
		}
	}

	return SpentTxOutRow{}, fmt.Errorf("could not find transaction %v", key.TxHash.String())
}
