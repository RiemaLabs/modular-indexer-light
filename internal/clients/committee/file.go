package committee

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
)

type fromFile string

func (f fromFile) LatestStateProof(context.Context) (*apis.Brc20VerifiableLatestStateProofResponse, error) {
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
func (fromFile) BlockHeight(context.Context) (uint, error) { panic("not supported") }
func (fromFile) CurrentBalanceOfWallet(context.Context, string, string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error) {
	panic("not supported")
}
func (fromFile) CurrentBalanceOfPkscript(context.Context, string, string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error) {
	panic("not supported")
}
