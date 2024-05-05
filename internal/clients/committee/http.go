package committee

import (
	"context"
	"net/url"

	"github.com/RiemaLabs/modular-indexer-committee/apis"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/httputl"
)

type endpoint struct{ u *url.URL }

func New(rawURL string) (Client, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "file" {
		return fromFile(u.Path), nil
	}
	return &endpoint{u: u}, nil
}

func (e *endpoint) Params(path string, queries url.Values) *url.URL {
	u := *e.u
	u.Path = path
	if len(queries) > 0 {
		u.RawQuery = queries.Encode()
	}
	return &u
}

func (e *endpoint) LatestStateProof(ctx context.Context) (*apis.Brc20VerifiableLatestStateProofResponse, error) {
	var ret apis.Brc20VerifiableLatestStateProofResponse
	if err := httputl.GetJSON(ctx, e.Params("/v1/brc20_verifiable/latest_state_proof", nil), &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (e *endpoint) BlockHeight(ctx context.Context) (height uint, err error) {
	err = httputl.GetJSON(ctx, e.Params("/v1/brc20_verifiable/block_height", nil), &height)
	return
}

func (e *endpoint) CurrentBalanceOfWallet(ctx context.Context, tick, wallet string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error) {
	var ret apis.Brc20VerifiableCurrentBalanceOfWalletResponse
	q := make(url.Values)
	q.Set("tick", tick)
	q.Set("wallet", wallet)
	if err := httputl.GetJSON(ctx, e.Params("/v1/brc20_verifiable/current_balance_of_wallet", q), &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (e *endpoint) CurrentBalanceOfPkscript(ctx context.Context, tick, pkscript string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error) {
	var ret *apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse
	q := make(url.Values)
	q.Set("tick", tick)
	q.Set("pkscript", pkscript)
	if err := httputl.GetJSON(ctx, e.Params("/v1/brc20_verifiable/current_balance_of_pkscript", q), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
