package getter

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
)

type HttpError struct {
	Code    int
	Message string
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
	if nil != err {
		return nil, err
	}

	type Result struct {
		Result *btcjson.TxRawResult
		Error  HttpError
	}
	var ret Result
	if err := json.Unmarshal(body, &ret); nil != err {
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
	if nil != err {
		return nil, err
	}

	type Result struct {
		Result *btcjson.TxRawResult
		Error  HttpError
	}
	var ret Result
	if err := json.Unmarshal(body, &ret); nil != err {
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

func (r *BitcoinOrdGetter) GetBlock1(hash string) (*btcjson.GetBlockVerboseResult, error) {
	type param struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}

	body, err := r.client.Post(context.Background(), r.Endpoint, param{
		Method: "getblock",
		Params: []interface{}{hash, 1},
	}, map[string]string{"Content-Type": "application/json"})
	if nil != err {
		return nil, err
	}

	type result struct {
		Result *btcjson.GetBlockVerboseResult
		Error  HttpError
	}

	var ret result
	if err := json.Unmarshal(body, &ret); nil != err {
		return nil, err
	}

	return ret.Result, nil
}

func (r *BitcoinOrdGetter) GetBlock2(hash string) (*btcjson.GetBlockVerboseTxResult, error) {
	type param struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}

	body, err := r.client.Post(context.Background(), r.Endpoint, param{
		Method: "getblock",
		Params: []interface{}{hash, 2},
	}, map[string]string{"Content-Type": "application/json"})
	if nil != err {
		return nil, err
	}

	type result struct {
		Result *btcjson.GetBlockVerboseTxResult
		Error  HttpError
	}

	var ret result
	if err := json.Unmarshal(body, &ret); nil != err {
		return nil, err
	}

	return ret.Result, nil
}

func (r *BitcoinOrdGetter) GetAllInscriptions(txID string) (map[string]*parser.TransactionInscription, error) {
	rawTx, err := r.GetRawTransaction(txID)
	if nil != err {
		return nil, err
	}

	buf, err := hex.DecodeString(rawTx.Hex)
	if nil != err {
		return nil, err
	}
	msgTx := new(wire.MsgTx)
	if err := msgTx.Deserialize(bytes.NewReader(buf)); nil != err {
		return nil, err
	}

	// inscription -> content
	res := make(map[string]*parser.TransactionInscription)
	inscriptions := parser.ParseInscriptionsFromTransaction(msgTx)
	id_counter := 0
	for index, _ := range msgTx.TxIn {
		for _, ii := range inscriptions {
			if int(ii.TxInIndex) == index {
				res[fmt.Sprintf("%si%d", txID, id_counter)] = ii
				id_counter++
			}

		}
	}
	return res, nil
}
