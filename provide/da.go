package provide

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RiemaLabs/indexer-committee/checkpoint"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/indexer"
	"github.com/RiemaLabs/indexer-light/types"
	sdk "github.com/RiemaLabs/nubit-da-sdk"
	sdkconstant "github.com/RiemaLabs/nubit-da-sdk/constant"
	sdktypes "github.com/RiemaLabs/nubit-da-sdk/types"
)

func NewDA(config *types.SourceDa) *ProviderDa {
	ctx := context.Background()
	sdk.SetNet(sdkconstant.TestNet)
	client := sdk.NewNubit(sdk.WithCtx(ctx)) //sdk.WithRpc("http://middleware.nubit.xyz"),
	//sdk.WithInviteCode("7mkEPWPBBrMr12WKNsL2UALvqYfbox"),
	//sdk.WithPrivateKey("9541ea760acc451684d28033566379a95bfe5a1b4da4a56a7df6055e4fa93eac"),
	committee := indexer.NewClient(ctx, config.IndexerName, config.ApiUrl)

	return &ProviderDa{
		ctx:       ctx,
		Name:      constant.ProvideDaName,
		Config:    config,
		Client:    client,
		Committee: committee,
	}
}
func (p *ProviderDa) GetCheckpoint(height uint, hash string) *types.CheckPointObject {
	high, err := p.Committee.BlockHigh() // TODO:: This interface is replaced by an interface that uses height to obtain the da chain height.
	if err != nil {
		return nil
	}
	datas, err := p.Client.Client.GetDatas(p.ctx, &sdktypes.GetDatasReq{
		NID:         []string{p.Config.NamespaceID},
		BlockNumber: int64(high),
	})
	if err != nil {
		return nil
	}
	if datas != nil && len(datas.Datas) > 0 {
		for _, data := range datas.Datas {
			if data.NID == p.Config.NamespaceID {
				var ck *checkpoint.Checkpoint
				decodeString, err := base64.StdEncoding.DecodeString(data.RawData)
				if err != nil {
					return nil
				}
				err = json.Unmarshal(decodeString, &ck)
				if err != nil {
					return nil
				}
				if strings.EqualFold(ck.Height, fmt.Sprintf("%d", height)) && strings.EqualFold(ck.Hash, hash) {
					return &types.CheckPointObject{
						CheckPoint: ck,
						Name:       p.Name,
						Type:       constant.ProvideDaName,
						Source: &types.Source{
							SourceDa: p.Config,
						},
					}
				}
			}
		}
	}
	return nil
}
