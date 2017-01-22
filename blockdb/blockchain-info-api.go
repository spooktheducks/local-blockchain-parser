package blockdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
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
			h, err := utils.HashFromString(tx.HashStr)
			if err != nil {
				return outHash, err
			}

			if h == hash {
				// fmt.Printf("tx found in block %v\n", bl.HashStr)
				return utils.HashFromString(bl.HashStr)
			}
		}
	}

	return outHash, fmt.Errorf("could not find block hash for tx %v", hash.String())
}

var (
	errBlockchainAPINotFound = fmt.Errorf("not found on blockchain API")
)

func (api *BlockchainInfoAPI) GetSpentTxOut(tx *Tx, txoutIdx uint32) (SpentTxOutRow, error) {
	addrs, err := tx.GetTxOutAddress(int(txoutIdx))
	if err != nil {
		return SpentTxOutRow{}, err
	}

	{
		type RawTxResponse struct {
			TxOuts []struct {
				Spent    bool   `json:"spent"`
				TxOutIdx uint32 `json:"n"`
			} `json:"out"`
		}

		url := fmt.Sprintf("https://blockchain.info/rawtx/%v", tx.Hash().String())

		var resp RawTxResponse
		err = getJSON(url, &resp)
		if err != nil {
			return SpentTxOutRow{}, err
		}

		for _, txout := range resp.TxOuts {
			if txout.TxOutIdx == txoutIdx && txout.Spent == false {
				return SpentTxOutRow{}, errBlockchainAPINotFound
			}
		}
	}

	for _, addr := range addrs {
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
		err = getJSON(fmt.Sprintf("https://blockchain.info/address/%v?format=json", addr.EncodeAddress()), &addrResp)
		if err != nil {
			return SpentTxOutRow{}, err
		}

		for _, txResp := range addrResp.Txs {
			fmt.Printf("searching tx %v...\n", txResp.Hash)
			for txinIdx, txin := range txResp.Inputs {
				prevOutHash := api.getTxHashByTxIndex(txin.PrevOut.TxIndex)

				if prevOutHash == tx.Hash().String() && txin.PrevOut.Index == uint32(txoutIdx) {
					inputTxHash, err := utils.HashFromString(txResp.Hash)
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

func getJSON(url string, x interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, x)
	if err != nil {
		return err
	}

	return nil
}
