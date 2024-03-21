package provide

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RiemaLabs/indexer-committee/checkpoint"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/log"
	"github.com/RiemaLabs/indexer-light/types"
	sdk "github.com/RiemaLabs/nubit-da-sdk"
	sdkconstant "github.com/RiemaLabs/nubit-da-sdk/constant"
	sdktypes "github.com/RiemaLabs/nubit-da-sdk/types"
)

func NewDA(config *types.SourceDa) *ProviderDa {
	ctx := context.Background()
	sdk.SetNet(sdkconstant.TestNet)
	opt := []sdk.Opt{}
	if config.Rpc != "" {
		opt = append(opt, sdk.WithRpc(config.Rpc))
	}
	opt = append(opt, sdk.WithCtx(ctx))
	client := sdk.NewNubit(opt...)
	//sdk.WithRpc("http://middleware.nubit.xyz"),
	//sdk.WithInviteCode("7mkEPWPBBrMr12WKNsL2UALvqYfbox"),
	//sdk.WithPrivateKey("9541ea760acc451684d28033566379a95bfe5a1b4da4a56a7df6055e4fa93eac"),
	return &ProviderDa{
		ctx:    ctx,
		Name:   constant.ProvideDaName,
		Config: config,
		Client: client,
	}
}
func (p *ProviderDa) GetCheckpoint(ctx context.Context, height uint, hash string) *types.CheckPointObject {
	var obj *types.CheckPointObject
	for {
		select {
		case <-ctx.Done():
			log.Error("ProviderDa", "context", ctx.Err())
			return obj
		default:
			count, err := p.Client.Client.GetTotalDataIDsInNamesapce(ctx, &sdktypes.GetTotalDataIDsInNamesapceReq{NID: p.Config.NamespaceID})
			if err != nil {
				log.Error("ProviderDa", "GetTotalDataIDsInNamesapce", err)
				return obj
			}
			dataids, err := p.Client.Client.GetDataInNamespace(ctx, &sdktypes.GetDataInNamespaceReq{
				NID:    p.Config.NamespaceID,
				Limit:  1,
				Offset: int(count.Count - 1),
			})
			if err != nil {
				log.Error("ProviderDa", "GetDataInNamespace", err)
				return obj
			}
			if dataids == nil || len(dataids.DataIDs) == 0 {
				log.Debug("ProviderDa", "GetDataInNamespace", "DataIDs do not exist or are of 0 length")
				continue
			}

			datas, err := p.Client.Client.GetDataByDAID(ctx, &sdktypes.GetDataReq{
				DAID: dataids.DataIDs[0],
			})
			if err != nil {
				log.Error("ProviderDa", "GetDataByDAID", err)
				return obj
			}

			var ck *checkpoint.Checkpoint
			decodeString, err := base64.StdEncoding.DecodeString(datas.CallData)
			if err != nil {
				log.Error("ProviderDa", "DecodeString.CallData", err)
				return obj
			}
			err = json.Unmarshal(decodeString, &ck)
			if err != nil {
				log.Error("ProviderDa", "Unmarshal.rawData", err)
				return obj
			}
			if strings.EqualFold(ck.Height, fmt.Sprintf("%d", height)) && strings.EqualFold(ck.Hash, hash) {
				obj = &types.CheckPointObject{
					CheckPoint: ck,
					Name:       p.Name,
					Type:       constant.ProvideDaName,
					Source: &types.Source{
						SourceDa: p.Config,
					},
				}
				return obj
			}
			log.Debug("ProviderDa", "GetCheckpoint", fmt.Sprintf("No CheckPoint data with %d was obtained", height))
		}
	}
}
