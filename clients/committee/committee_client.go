package committee

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
	"github.com/RiemaLabs/modular-indexer-light/clients/http"
	"github.com/RiemaLabs/modular-indexer-light/constant"
)

// TODO: Medium. Distinguish indexer and committee indexer.
type CommitteeIndexerClient struct {
	ctx      context.Context
	endpoint string
	name     string
	*http.Client
}

func NewCommitteeIndexerClient(ctx context.Context, name, endpoint string) *CommitteeIndexerClient {
	return &CommitteeIndexerClient{ctx, endpoint, name, http.NewClient()}
}

func (c *CommitteeIndexerClient) JoinPath(subURL string) (string, error) {
	path, err := url.JoinPath(c.endpoint, subURL)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *CommitteeIndexerClient) LatestStateProof() (*apis.Brc20VerifiableLatestStateProofResponse, error) {
	var data *apis.Brc20VerifiableLatestStateProofResponse

	path, err := c.JoinPath(constant.LatestStateProof)
	if err != nil {
		return nil, err
	}

	resp, err := c.Get(c.ctx, path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *CommitteeIndexerClient) BlockHeight() (uint, error) {
	var data uint
	path, err := c.JoinPath(constant.LatestStateProof)
	if err != nil {
		return 0, err
	}
	resp, err := c.Get(c.ctx, path)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return 0, err
	}
	return data, nil
}

func (c *CommitteeIndexerClient) CurrentBalanceOfWallet(tick, wallet string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error) {
	var data *apis.Brc20VerifiableCurrentBalanceOfWalletResponse
	path, err := c.JoinPath(constant.CurrentBalanceOfWallet)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("tick", tick)
	q.Set("wallet", wallet)
	u.RawQuery = q.Encode()

	resp, err := c.Get(c.ctx, u.String())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *CommitteeIndexerClient) CurrentBalanceOfPkscript(tick, pkscript string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error) {
	var data *apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse
	path, err := c.JoinPath(constant.CurrentBalanceOfWallet)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("tick", tick)
	q.Set("pkscript", pkscript)
	u.RawQuery = q.Encode()

	resp, err := c.Get(c.ctx, u.String())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
