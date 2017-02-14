package cmds

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

type DumpTxDataCommand struct {
	startBlock, endBlock uint64
	datFileDir, outDir   string
	coalesce             bool
}

func NewDumpTxDataCommand(startBlock, endBlock uint64, datFileDir, outDir string, coalesce bool) *DumpTxDataCommand {
	return &DumpTxDataCommand{
		startBlock: startBlock,
		endBlock:   endBlock,
		datFileDir: datFileDir,
		outDir:     filepath.Join(".", outDir, "dump-tx-data"),
		coalesce:   coalesce,
	}
}

func (cmd *DumpTxDataCommand) RunCommand() error {
	err := os.MkdirAll(cmd.outDir, 0777)
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

	for i := int(cmd.startBlock); i < int(cmd.endBlock)+1; i++ {
		chDone := make(chan bool)
		chDones = append(chDones, chDone)
		go cmd.parseBlock(i, chErr, chDone, procLimiter)
	}

	// wait for all goroutines to complete
	for _, chDone := range chDones {
		<-chDone
	}

	// close error channel
	close(chErr)

	return nil
}

func (cmd *DumpTxDataCommand) parseBlock(blockFileNum int, chErr chan error, chDone chan bool, procLimiter chan bool) {
	defer close(chDone)
	defer func() { procLimiter <- true }()
	<-procLimiter

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlocksFromDAT(filepath.Join(cmd.datFileDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	csvFile := utils.NewConditionalFile(filepath.Join(cmd.outDir, fmt.Sprintf("blk%05d-tx-data.csv", blockFileNum)))
	defer csvFile.Close()

	// write CSV header
	_, err = csvFile.WriteString("block,recipient address,tx,output index,data\n", false)
	if err != nil {
		chErr <- err
		return
	}

	numBlocks := len(blocks)
	for blIdx, bl := range blocks {
		blockHash := bl.Hash().String()

		for _, btctx := range bl.Transactions() {
			tx := Tx{Tx: btctx}

			if cmd.coalesce {
				err := cmd.writeCoalesced(tx)
				if err != nil {
					chErr <- err
					return
				}
			} else {
				err := cmd.writeNonCoalesced(tx, csvFile)
				if err != nil {
					chErr <- err
					return
				}
			}
		}

		fmt.Printf("finished block %s (%d/%d)\n", blockHash, blIdx, numBlocks)
	}

	if err != nil {
		chErr <- err
		return
	}
}

func (cmd *DumpTxDataCommand) writeCoalesced(tx Tx) error {
	txHash := tx.Hash().String()

	allTxinData := []byte{}
	allTxoutData := []byte{}

	for _, txin := range tx.MsgTx().TxIn {
		allTxinData = append(allTxinData, txin.SignatureScript...)
	}

	err := ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("%s-txin.dat", txHash)), allTxinData, 0666)
	if err != nil {
		return err
	}

	for _, txout := range tx.MsgTx().TxOut {
		allTxoutData = append(allTxoutData, txout.PkScript...)
	}

	err = ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("%s-txout.dat", txHash)), allTxoutData, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *DumpTxDataCommand) writeNonCoalesced(tx Tx, csvFile *utils.ConditionalFile) error {
	txHash := tx.Hash().String()

	for txinIdx, txin := range tx.MsgTx().TxIn {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("%s-txin-%d.dat", txHash, txinIdx)), txin.SignatureScript, 0666)
		if err != nil {
			return err
		}
	}

	for txoutIdx, txout := range tx.MsgTx().TxOut {
		err := ioutil.WriteFile(filepath.Join(cmd.outDir, fmt.Sprintf("%s-txout-%d.dat", txHash, txoutIdx)), txout.PkScript, 0666)
		if err != nil {
			return err
		}

		addrs, err := tx.GetTxOutAddress(txoutIdx)
		if err != nil {
			return err
		}

		recipientAddr := ""
		if len(addrs) > 0 {
			recipientAddr = addrs[0].EncodeAddress()
		}

		txoutData, err := tx.GetNonOPDataFromTxOut(txoutIdx)
		if err != nil {
			return err
		}

		txoutDataHex := hex.EncodeToString(txoutData)

		csvFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", tx.BlockHash.String(), recipientAddr, txHash, txoutIdx, txoutDataHex), true)
	}

	return nil
}
