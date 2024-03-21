package getter

import (
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/types"
	"github.com/btcsuite/btcd/rpcclient"
)

type BitcoinOrdGetter struct {
	client *rpcclient.Client
}

func NewGetter(config *types.Config) (*BitcoinOrdGetter, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         config.Host,
		User:         config.User,
		Pass:         config.Password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	return &BitcoinOrdGetter{
		client: client,
	}, nil
}

func (r *BitcoinOrdGetter) GetLatestBlockHeight() (uint, error) {
	count, err := r.client.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return uint(count), err
}

func (r *BitcoinOrdGetter) GetBlockHash(blockHeight uint) (string, error) {
	hash, err := r.client.GetBlockHash(int64(blockHeight))
	if nil != err {
		return "", err
	}
	return hash.String(), err
}

func (r *BitcoinOrdGetter) GetOrdTransfers(blockHeight uint) ([]getter.OrdTransfer, error) {
	hash, err := r.client.GetBlockHash(int64(blockHeight))
	if nil != err || hash == nil {
		return []getter.OrdTransfer{}, err
	}

	block, err := r.client.GetBlock(hash)
	if nil != err {
		return []getter.OrdTransfer{}, err
	}
	block = block

	// TODO fetch tx from  block.txdata
	return []getter.OrdTransfer{}, nil
}
