package txhashsource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/blockdb"
)

func NewAddressTxHashSource(db *blockdb.BlockDB, addr string) TxHashSource {
	ch := make(chan chainhash.Hash)
	go func() {
		defer close(ch)

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

		type TxResponse struct {
			Hash string `json:"hash"`
		}

		type AddressResponse struct {
			Txs []TxResponse `json:"txs"`
		}

		var addrResp AddressResponse
		err = json.Unmarshal(bs, &addrResp)
		if err != nil {
			return
		}

		reversedTxs := make([]TxResponse, len(addrResp.Txs))

		for i := range addrResp.Txs {
			reversedTxs[len(reversedTxs)-i-1] = addrResp.Txs[i]
		}

		for _, tx := range reversedTxs {
			fmt.Printf("tx %v\n", tx.Hash)
			txHash, err := blockdb.HashFromString(tx.Hash)
			if err != nil {
				fmt.Println("err decoding tx hash:", err)
				return
			}

			ch <- txHash
		}

	}()

	return TxHashSource(ch)
}
