package txhashsource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	. "github.com/spooktheducks/local-blockchain-parser/blockdb"
	"github.com/spooktheducks/local-blockchain-parser/cmds/utils"
)

func NewAddressTxHashSource(db *BlockDB, addr string) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

		var numTxs int
		{
			type AddressNTxResponse struct {
				NumTxs int `json:"n_tx"`
			}

			url := fmt.Sprintf("https://blockchain.info/address/%v?format=json", addr)
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			var ntxResp AddressNTxResponse
			err = json.Unmarshal(bs, &ntxResp)
			if err != nil {
				return
			}

			numTxs = ntxResp.NumTxs
		}

		var err error
		for offset := 0; offset < numTxs; {
			offset, err = getTxs(addr, offset, ch)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
		}
	}()

	return TxHashSource(ch)
}

func getTxs(addr string, offset int, ch chan chainhash.Hash) (int, error) {
	url := fmt.Sprintf("https://blockchain.info/address/%v?format=json&offset=%v&sort=1", addr, offset)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	type TxResponse struct {
		Hash string `json:"hash"`
	}

	type AddressResponse struct {
		Txs []TxResponse `json:"txs"`
	}

	var addrResp AddressResponse
	err = json.Unmarshal(bs, &addrResp)
	if err != nil {
		return 0, err
	}

	for _, tx := range addrResp.Txs {
		txHash, err := utils.HashFromString(tx.Hash)
		if err != nil {
			fmt.Println("err decoding tx hash:", err)
			return 0, err
		}

		ch <- txHash
	}

	return offset + len(addrResp.Txs), nil
}
