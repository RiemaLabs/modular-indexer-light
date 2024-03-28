package getter

import (
	"context"
	"encoding/json"

	"github.com/RiemaLabs/modular-indexer-committee/ord/getter"
	"github.com/RiemaLabs/modular-indexer-light/clients/http"

	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
)

type GetBlockCountResponse struct {
	Result int         `json:"result"`
	Error  interface{} `json:"error"`
	Id     interface{} `json:"id"`
}

type GetBlockHashResponse struct {
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
	Id     interface{} `json:"id"`
}

type BitcoinOrdGetter struct {
	client   *http.Client
	Endpoint string
}

type HttpError struct {
	Code    int
	Message string
}

func NewBitcoinOrdGetter(bitcoinRPC string) (*BitcoinOrdGetter, error) {
	return &BitcoinOrdGetter{
		client:   http.NewClient(),
		Endpoint: bitcoinRPC,
	}, nil
}

func (r *BitcoinOrdGetter) GetLatestBlockHeight() (uint, error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	param := make(map[string]string)
	param["method"] = "getblockcount"
	post, err := r.client.Post(context.Background(), r.Endpoint, param, header)
	if err != nil {
		return 0, err
	}
	var count GetBlockCountResponse
	err = json.Unmarshal(post, &count)
	if err != nil {
		return 0, err
	}

	return uint(count.Result), err
}

func (r *BitcoinOrdGetter) GetBlockHash(blockHeight uint) (string, error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	param := make(map[string]any)
	param["method"] = "getblockhash"
	if blockHeight != 0 {
		param["params"] = []int{int(blockHeight)}
	}

	post, err := r.client.Post(context.Background(), r.Endpoint, param, header)
	if err != nil {
		return "", err
	}
	var resp GetBlockHashResponse
	err = json.Unmarshal(post, &resp)
	if err != nil {
		return "", err
	}
	return resp.Result, err
}

func (r *BitcoinOrdGetter) GetOrdTransfers(blockHeight uint) ([]getter.OrdTransfer, error) {
	// TODO: High. Use satpoint mapping to get the transfer from past transactions.
	return []getter.OrdTransfer{}, nil
}

func (r *BitcoinOrdGetter) GetRawTransaction(txID string) (*btcjson.TxRawResult, error) {
	type txReq struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}

	body, err := r.client.Post(context.Background(), r.Endpoint, txReq{
		Method: "getrawtransaction",
		Params: []interface{}{txID, true},
	}, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}

	type Result struct {
		Result *btcjson.TxRawResult
		Error  HttpError
	}
	var ret Result
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	if ret.Error.Code != 0 {
		return nil, errors.New(ret.Error.Message)
	}
	return ret.Result, nil
}

func (r *BitcoinOrdGetter) GetOutput(txID string, index int) (*btcjson.Vout, error) {
	type txReq struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}

	body, err := r.client.Post(context.Background(), r.Endpoint, txReq{
		Method: "getrawtransaction",
		Params: []interface{}{txID, true},
	}, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}

	type Result struct {
		Result *btcjson.TxRawResult
		Error  HttpError
	}
	var ret Result
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	if ret.Error.Code != 0 {
		return nil, errors.New(ret.Error.Message)
	}

	if len(ret.Result.Vout) < index+1 {
		return nil, fmt.Errorf("RawTransction not have enough vout, cap %d, need %d", len(ret.Result.Vout), index)
	}
	return &ret.Result.Vout[index], nil
}

// Getblock returns an Object with information about block ‘hash’
func (r *BitcoinOrdGetter) GetBlock(hash string) (*btcjson.GetBlockVerboseResult, error) {
	type param struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}

	body, err := r.client.Post(context.Background(), r.Endpoint, param{
		Method: "getblock",
		Params: []interface{}{hash, 1},
	}, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}

	type result struct {
		Result *btcjson.GetBlockVerboseResult
		Error  HttpError
	}

	var ret result
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return ret.Result, nil
}

// GetBlockDetail returns an Object with information about block ‘hash’ and information about each transaction.
func (r *BitcoinOrdGetter) GetBlockDetail(hash string) (*btcjson.GetBlockVerboseTxResult, error) {
	type param struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}

	body, err := r.client.Post(context.Background(), r.Endpoint, param{
		Method: "getblock",
		Params: []interface{}{hash, 2},
	}, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}

	type result struct {
		Result *btcjson.GetBlockVerboseTxResult
		Error  HttpError
	}

	var ret result
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return ret.Result, nil
}

func (r *BitcoinOrdGetter) GetAllInscriptions(txID string) (map[string]*parser.TransactionInscription, error) {
	rawTx, err := r.GetRawTransaction(txID)
	if err != nil {
		return nil, err
	}

	buf, err := hex.DecodeString(rawTx.Hex)
	if err != nil {
		return nil, err
	}
	msgTx := new(wire.MsgTx)
	if err := msgTx.Deserialize(bytes.NewReader(buf)); err != nil {
		return nil, err
	}

	// inscription -> content
	res := make(map[string]*parser.TransactionInscription)
	inscriptions := parser.ParseInscriptionsFromTransaction(msgTx)
	id_counter := 0
	for index := range msgTx.TxIn {
		for _, ii := range inscriptions {
			if int(ii.TxInIndex) == index {
				res[fmt.Sprintf("%si%d", txID, id_counter)] = ii
				id_counter++
			}

		}
	}
	return res, nil
}
