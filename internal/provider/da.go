package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"

	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	sdk "github.com/RiemaLabs/nubit-da-sdk"
	"github.com/RiemaLabs/nubit-da-sdk/constant"
	"github.com/RiemaLabs/nubit-da-sdk/types"
)

type DA struct {
	Config               *configs.SourceDA
	MetaProtocol         string
	Retry                int
	LastCheckpointOffset int
}

func NewProviderDA(sourceDA *configs.SourceDA, metaProtocol string, retry int) *DA {
	return &DA{
		Config:       sourceDA,
		MetaProtocol: metaProtocol,
		Retry:        retry,
		// At the beginning, we don't know the offset
		LastCheckpointOffset: 0,
	}
}

func (p *DA) GetCheckpoint(ctx context.Context, height uint, hash string) (*configs.CheckpointExport, error) {

	// We don't use the timeout to limit the single call of DownloadCheckpointByDA.
	maxTimeout := 1000 * time.Second
	nid, net, name, mp := p.Config.NamespaceID, p.Config.Network, p.Config.Name, p.MetaProtocol

	if net == "Pre-Alpha Testnet" {
		sdk.SetNet(constant.PreAlphaTestNet)
	} else if net == "Testnet" {
		sdk.SetNet(constant.TestNet)
	} else {
		return nil, fmt.Errorf("unknown network: %s", net)
	}

	clientDA := sdk.NewNubit(sdk.WithCtx(ctx)).Client
	resCount, err := clientDA.GetTotalDataIDsInNamesapce(ctx, &types.GetTotalDataIDsInNamesapceReq{NID: nid})
	if err != nil {
		return nil, fmt.Errorf("failed to get the count of data in namespace %s, error: %v", nid, err)
	}
	count := int(resCount.Count)
	if count == 0 {
		return nil, fmt.Errorf("the count of data in namespace %s is zero", nid)
	}

	p.LastCheckpointOffset = count - 1
	var ck *checkpoint.Checkpoint
	var offset int = count

OuterLoop:
	for i := 0; i < p.Retry; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ck, offset, err = DownloadCheckpointByDA(nid, net, name, mp, strconv.Itoa(int(height)), hash, p.LastCheckpointOffset, maxTimeout)
			if err != nil {
				time.Sleep(20 * time.Second)
				logs.Error.Printf(err.Error())
				continue
			}
			break OuterLoop
		}
	}
	p.LastCheckpointOffset = offset - 1
	if err != nil {
		return nil, err
	}
	return &configs.CheckpointExport{
		Checkpoint: ck,
		SourceDA:   p.Config,
	}, nil
}
