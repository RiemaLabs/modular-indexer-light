package states

import (
	"testing"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"

	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/provider"
)

var correctCheckpoint = &checkpoint.Checkpoint{
	Commitment: "AilOEj54jwAX3quXaqCoZlssmuvpbg+0BNFdmaWezX8=",
	Hash:       "00000000000000000002d0aa70024c05626a0746e401ec33ac6a4c0dff4f4e44",
	Height:     "841221",

	//Commitment: "CD/zpV/e0ukEShxlV2/XKi2NUWGRqKCqfPxBBeCIMF0=",
	//Hash: "00000000000000000000b2a38e98aa32d96d7349e8b772047946db2b7c19f5ae",
	//Height: "841218",

	MetaProtocol: "brc-20",
	Name:         "nubit-official-01",
	URL:          "file:///nubit-official-01.json",
	//URL:          "http://35.88.177.54:8080",
	Version: "v0.1.0-rc.3",
}

var maliciousCheckpointProviders = []provider.CheckpointProvider{
	(*configs.SourceRaw)(correctCheckpoint),
	&configs.SourceRaw{
		Commitment: "CD/zpV/e0ukEShxlV2/XKi2NUWGRqKCqfPxBBeCIMF0=",
		//Commitment:   correctCheckpoint.Commitment,
		Hash:         correctCheckpoint.Hash,
		Height:       correctCheckpoint.Height,
		MetaProtocol: correctCheckpoint.MetaProtocol,
		Name:         "nubit-malicious-01",
		URL:          "file:///nubit-malicious-01.json",
		Version:      correctCheckpoint.Version,
	},
}

func TestState_UpdateCheckpoints(t *testing.T) {
	// 1. Missing transactions
	// 2. Wrong transaction fees
	// 3. Duplicated transactions
	// 4. Random malicious transactions

	state := New(
		"",
		maliciousCheckpointProviders,
		&configs.CheckpointExport{Checkpoint: correctCheckpoint},
		2,
		3*time.Second,
	)
	if err := state.UpdateCheckpoints(
		841221,
		"00000000000000000002d0aa70024c05626a0746e401ec33ac6a4c0dff4f4e44",
	); err != nil {
		t.Fatal(err)
	}
}
