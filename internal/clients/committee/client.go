package committee

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/RiemaLabs/modular-indexer-committee/apis"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/http"
	"github.com/RiemaLabs/modular-indexer-light/internal/constant"
)

// TODO: Medium. Distinguish indexer and committee indexer.

type Client struct {
	ctx      context.Context
	endpoint string
	name     string
	*http.Client
}

func New(ctx context.Context, name, endpoint string) *Client {
	return &Client{ctx, endpoint, name, http.New()}
}

func (c *Client) JoinPath(subURL string) (string, error) {
	path, err := url.JoinPath(c.endpoint, subURL)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *Client) LatestStateProof() (*apis.Brc20VerifiableLatestStateProofResponse, error) {
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

func (c *Client) BlockHeight() (uint, error) {
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

func (c *Client) CurrentBalanceOfWallet(tick, wallet string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error) {
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

func (c *Client) CurrentBalanceOfPkscript(tick, pkscript string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error) {
	var data *apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse
	path, err := c.JoinPath(constant.CurrentBalanceOfPkscript)
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
