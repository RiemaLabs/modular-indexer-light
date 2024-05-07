package utils

import (
	"testing"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
)

func TestToJSObject(t *testing.T) {
	proof := "xxx"
	balance := apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse{
		Error: nil,
		Result: &apis.Brc20VerifiableCurrentBalanceOfPkscriptResult{
			AvailableBalance: "123",
			OverallBalance:   "456",
		},
		Proof: &proof,
	}
	o := ToRawMap(balance)
	r, ok := o["result"]
	if !ok {
		t.Fatal(o)
	}
	if _, ok := r.(RawMap); !ok {
		t.Fatal(r)
	}
}
