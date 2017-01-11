package blockdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

type BlockchainInfoAPI struct{}

func (api *BlockchainInfoAPI) GetBlockHashForTx(hash chainhash.Hash) (chainhash.Hash, error) {
	var outHash chainhash.Hash

	type RawTxResponse struct {
		BlockHeight uint32 `json:"block_height"`
	}

	var rawTxResponse RawTxResponse
	{
		url := fmt.Sprintf("https://blockchain.info/rawtx/%v", hash.String())
		resp, err := http.Get(url)
		if err != nil {
			return outHash, err
		}
		defer resp.Body.Close()

		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return outHash, err
		}

		err = json.Unmarshal(bs, &rawTxResponse)
		if err != nil {
			return outHash, err
		}
	}

	type RawBlock struct {
		HashStr string `json:"hash"`
		Txs     []struct {
			HashStr string `json:"hash"`
		} `json:"tx"`
	}

	type BlockHeightResponse struct {
		Blocks []RawBlock `json:"blocks"`
	}

	var blockHeightResponse BlockHeightResponse
	{
		url := fmt.Sprintf("https://blockchain.info/block-height/%v?format=json", rawTxResponse.BlockHeight)
		// fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			return outHash, err
		}
		defer resp.Body.Close()

		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return outHash, err
		}

		err = json.Unmarshal(bs, &blockHeightResponse)
		if err != nil {
			return outHash, err
		}
	}

	for _, bl := range blockHeightResponse.Blocks {
		for _, tx := range bl.Txs {
			h, err := HashFromString(tx.HashStr)
			if err != nil {
				return outHash, err
			}

			if h == hash {
				// fmt.Printf("tx found in block %v\n", bl.HashStr)
				return HashFromString(bl.HashStr)
			}
		}
	}

	return outHash, fmt.Errorf("could not find block hash for tx %v", hash.String())
}

var (
	errBlockchainAPINotFound = fmt.Errorf("not found on blockchain API")
)

func (api *BlockchainInfoAPI) GetSpentTxOut(tx *btcutil.Tx, txoutIdx uint32) (SpentTxOutRow, error) {
	txout := tx.MsgTx().TxOut[txoutIdx]

	addrs, err := utils.GetTxOutAddress(txout)
	if err != nil {
		return SpentTxOutRow{}, err
	}

	for _, addr := range addrs {
		// fmt.Printf("searching address %v...\n", addr.EncodeAddress())
		url := fmt.Sprintf("https://blockchain.info/address/%v?format=json", addr.EncodeAddress())
		resp, err := http.Get(url)
		if err != nil {
			return SpentTxOutRow{}, err
		}
		defer resp.Body.Close()

		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return SpentTxOutRow{}, err
		}

		type AddressResponse struct {
			Txs []struct {
				Hash   string `json:"hash"`
				Inputs []struct {
					PrevOut struct {
						TxIndex uint32 `json:"tx_index"`
						Index   uint32 `json:"n"`
					} `json:"prev_out"`
				} `json:"inputs"`
			} `json:"txs"`
		}

		var addrResp AddressResponse
		err = json.Unmarshal(bs, &addrResp)
		if err != nil {
			return SpentTxOutRow{}, err
		}

		for _, txResp := range addrResp.Txs {
			fmt.Printf("searching tx %v...\n", txResp.Hash)
			for txinIdx, txin := range txResp.Inputs {
				prevOutHash := api.getTxHashByTxIndex(txin.PrevOut.TxIndex)

				if prevOutHash == tx.Hash().String() && txin.PrevOut.Index == uint32(txoutIdx) {
					inputTxHash, err := HashFromString(txResp.Hash)
					if err != nil {
						return SpentTxOutRow{}, err
					}
					return SpentTxOutRow{InputTxHash: inputTxHash, TxInIndex: uint32(txinIdx)}, nil
				}
			}
		}
	}
	return SpentTxOutRow{}, errBlockchainAPINotFound
}

func (api *BlockchainInfoAPI) getTxHashByTxIndex(txIndex uint32) string {
	url := fmt.Sprintf("https://blockchain.info/tx-index/%v?format=json", txIndex)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	type TxIndexResponse struct {
		Hash string `json:"hash"`
	}

	var txIndexResp TxIndexResponse
	err = json.Unmarshal(bs, &txIndexResp)
	if err != nil {
		return ""
	}

	return txIndexResp.Hash
}
