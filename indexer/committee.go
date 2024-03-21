package indexer

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/RiemaLabs/indexer-committee/apis"
	"github.com/RiemaLabs/indexer-light/clients/http"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/log"
)

type CommitteeIndexer struct {
	ctx      context.Context
	endpoint string
	name     string
	*http.Client
}

func NewClient(ctx context.Context, name, endpoint string) *CommitteeIndexer {
	return &CommitteeIndexer{ctx, endpoint, name, http.NewClient()}
}

func (c *CommitteeIndexer) path(str string) string {
	path, err := url.JoinPath(c.endpoint, str)
	if err != nil {
		log.Error("CommitteeIndexer", "method", "JoinPath", "err", err)
		return ""
	}
	return path
}

func (c *CommitteeIndexer) StateDiff() (*apis.Brc20VerifiableLatestStateProofResponse, error) {
	var data *apis.Brc20VerifiableLatestStateProofResponse
	post, err := c.Post(c.ctx, c.path(constant.StateDiff), nil, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(post, &data)
	if err != nil {
		log.Error("CommitteeIndexer", "method", "StateDiff", "Unmarshal", err)
		return nil, err
	}
	return data, err
}

func (c *CommitteeIndexer) BlockHigh() (uint, error) {
	var data uint
	post, err := c.Post(c.ctx, c.path(constant.BlockHigh), nil, nil)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(post, &data)
	if err != nil {
		log.Error("CommitteeIndexer", "method", "BlockHigh", "Unmarshal", err)
	}
	return data, err
}

func (c *CommitteeIndexer) GetBalance(tick, pkscript string) (*apis.Brc20VerifiableGetCurrentBalanceOfWalletResponse, error) {
	var data *apis.Brc20VerifiableGetCurrentBalanceOfWalletResponse
	post, err := c.Post(c.ctx, c.path(constant.Balance), apis.Brc20VerifiableGetCurrentBalanceOfWalletRequest{Tick: tick, Pkscript: pkscript}, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(post, &data)
	if err != nil {
		log.Error("CommitteeIndexer", "method", "GetBalance", "Unmarshal", err)
	}
	return data, err
}
