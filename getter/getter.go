package getter

import (
	"context"
	"encoding/json"

	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/clients/http"
	"github.com/RiemaLabs/indexer-light/types"
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

func NewGetter(config *types.Config) (*BitcoinOrdGetter, error) {
	return &BitcoinOrdGetter{
		client:   http.NewClient(),
		Endpoint: config.BitCoinRpc.Host,
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
	// {"method": "getblockhash", "params": [835525]}

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
	//hash, err := r.client.GetBlockHash(int64(blockHeight))
	//if nil != err || hash == nil {
	//	return []getter.OrdTransfer{}, err
	//}
	//
	//block, err := r.client.GetBlock(hash)
	//if nil != err {
	//	return []getter.OrdTransfer{}, err
	//}
	//block = block

	// TODO fetch tx from  block.txdata
	return []getter.OrdTransfer{}, nil
}
