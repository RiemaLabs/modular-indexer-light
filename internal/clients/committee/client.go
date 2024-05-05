package committee

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/RiemaLabs/modular-indexer-committee/apis"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/http"
	"github.com/RiemaLabs/modular-indexer-light/internal/constant"
)

// TODO: Medium. Distinguish indexer and committee indexer.

type Client interface {
	LatestStateProof() (*apis.Brc20VerifiableLatestStateProofResponse, error)
	BlockHeight() (uint, error)
	CurrentBalanceOfWallet(tick, wallet string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error)
	CurrentBalanceOfPkscript(tick, pkscript string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error)
}

type fromFile string

func (f fromFile) LatestStateProof() (*apis.Brc20VerifiableLatestStateProofResponse, error) {
	var ret apis.Brc20VerifiableLatestStateProofResponse
	data, err := os.ReadFile(strings.TrimPrefix(string(f), "/"))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
func (fromFile) BlockHeight() (uint, error) { panic("not supported") }
func (fromFile) CurrentBalanceOfWallet(string, string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error) {
	panic("not supported")
}
func (fromFile) CurrentBalanceOfPkscript(string, string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error) {
	panic("not supported")
}

type remoteClient struct {
	ctx      context.Context
	endpoint string
	name     string
	*http.Client
}

func New(ctx context.Context, name, endpoint string) Client {
	if u, err := url.Parse(endpoint); err == nil && u.Scheme == "file" {
		return fromFile(u.Path)
	}
	return &remoteClient{ctx, endpoint, name, http.New()}
}

func (c *remoteClient) joinPath(subURL string) (string, error) {
	path, err := url.JoinPath(c.endpoint, subURL)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *remoteClient) LatestStateProof() (*apis.Brc20VerifiableLatestStateProofResponse, error) {
	var data *apis.Brc20VerifiableLatestStateProofResponse

	path, err := c.joinPath(constant.LatestStateProof)
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

func (c *remoteClient) BlockHeight() (uint, error) {
	var data uint
	path, err := c.joinPath(constant.LatestStateProof)
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

func (c *remoteClient) CurrentBalanceOfWallet(tick, wallet string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error) {
	var data *apis.Brc20VerifiableCurrentBalanceOfWalletResponse
	path, err := c.joinPath(constant.CurrentBalanceOfWallet)
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

func (c *remoteClient) CurrentBalanceOfPkscript(tick, pkscript string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error) {
	var data *apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse
	path, err := c.joinPath(constant.CurrentBalanceOfPkscript)
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
