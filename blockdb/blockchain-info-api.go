package blockdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
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
		fmt.Println(url)
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

	return outHash, nil
}
