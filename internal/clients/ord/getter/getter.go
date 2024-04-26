package getter

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/jsonrpc"
)

const OK = 0

type (
	Response[T any] jsonrpc.Response[T, ErrorResponse]

	ErrorResponse struct {
		Code    int
		Message string
	}
)

var _ jsonrpc.HasErr = (*Response[struct{}])(nil)

func (r *Response[T]) Err() error {
	if e := r.Error; e.Code != OK {
		return fmt.Errorf("RPC error: code=%d, msg=%s", e.Code, e.Message)
	}
	return nil
}

type Client struct{ cl jsonrpc.Client }

func New(bitcoinRPC string) (*Client, error) {
	cl, err := jsonrpc.New(bitcoinRPC)
	if err != nil {
		return nil, err
	}
	return &Client{cl: cl}, nil
}

func (r *Client) GetLatestBlockHeight(ctx context.Context) (uint, error) {
	var rsp Response[uint]
	if err := r.cl.Call(ctx, "getblockcount", nil, &rsp); err != nil {
		return 0, fmt.Errorf("get latest block height error: err=%v", err)
	}
	return rsp.Result, nil
}

func (r *Client) GetBlockHash(ctx context.Context, height uint) (string, error) {
	var rsp Response[string]
	if err := r.cl.Call(ctx, "getblockhash", []uint{height}, &rsp); err != nil {
		return "", fmt.Errorf("get block hash error: height=%d, err=%v", height, err)
	}
	return rsp.Result, nil
}

func (r *Client) GetOrdTransfers(ctx context.Context, blockHeight uint) error {
	// TODO: High. Use satpoint mapping to get the transfer from past transactions.
	_ = ctx
	_ = blockHeight
	return errors.New("not implemented")
}

func (r *Client) GetRawTransaction(ctx context.Context, txID string) (*btcjson.TxRawResult, error) {
	var rsp Response[*btcjson.TxRawResult]
	if err := r.cl.Call(ctx, "getrawtransaction", []interface{}{txID, true}, &rsp); err != nil {
		return nil, fmt.Errorf("get raw transaction error: txID=%s, err=%v", txID, err)
	}
	return rsp.Result, nil
}

func (r *Client) GetOutput(ctx context.Context, txID string, index int) (*btcjson.Vout, error) {
	var rsp Response[*btcjson.TxRawResult]
	if err := r.cl.Call(ctx, "getrawtransaction", []interface{}{txID, true}, &rsp); err != nil {
		return nil, fmt.Errorf("get raw transaction error: txID=%s, index=%d, err=%v", txID, index, err)
	}
	if l := len(rsp.Result.Vout); l < index+1 {
		return nil, fmt.Errorf("raw transactions out of index: len=%d, index=%d", l, index)
	}
	return &rsp.Result.Vout[index], nil
}

func (r *Client) GetBlock(ctx context.Context, hash string) (*btcjson.GetBlockVerboseResult, error) {
	var rsp Response[*btcjson.GetBlockVerboseResult]
	if err := r.cl.Call(ctx, "getblock", []interface{}{hash, 1}, &rsp); err != nil {
		return nil, fmt.Errorf("get block error: hash=%s, err=%v", hash, err)
	}
	return rsp.Result, nil
}

func (r *Client) GetBlockDetail(ctx context.Context, hash string) (*btcjson.GetBlockVerboseTxResult, error) {
	var rsp Response[*btcjson.GetBlockVerboseTxResult]
	if err := r.cl.Call(ctx, "getblock", []interface{}{hash, 2}, &rsp); err != nil {
		return nil, fmt.Errorf("get block detail error: hash=%s, err=%v", hash, err)
	}
	return rsp.Result, nil
}

func (r *Client) GetAllInscriptions(ctx context.Context, txID string) (map[string]*parser.TransactionInscription, error) {
	rawTx, err := r.GetRawTransaction(ctx, txID)
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

	res := make(map[string]*parser.TransactionInscription)
	inscriptions := parser.ParseInscriptionsFromTransaction(msgTx)
	idCnt := 0
	for index := range msgTx.TxIn {
		for _, inscription := range inscriptions {
			if int(inscription.TxInIndex) == index {
				res[fmt.Sprintf("%si%d", txID, idCnt)] = inscription
				idCnt++
			}

		}
	}
	return res, nil
}
