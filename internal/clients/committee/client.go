package committee

import (
	"context"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
)

// TODO: Medium. Distinguish indexer and committee indexer.

type Client interface {
	LatestStateProof(ctx context.Context) (*apis.Brc20VerifiableLatestStateProofResponse, error)
	BlockHeight(ctx context.Context) (uint, error)
	CurrentBalanceOfWallet(ctx context.Context, tick, wallet string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error)
	CurrentBalanceOfPkscript(ctx context.Context, tick, pkscript string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error)
}
