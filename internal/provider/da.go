package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"

	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	sdk "github.com/RiemaLabs/nubit-da-sdk"
	"github.com/RiemaLabs/nubit-da-sdk/constant"
	"github.com/RiemaLabs/nubit-da-sdk/types"
)

type DA struct {
	Config       *configs.SourceDA
	MetaProtocol string
}

func NewProviderDA(sourceDA *configs.SourceDA, metaProtocol string) *DA {
	return &DA{Config: sourceDA, MetaProtocol: metaProtocol}
}

func (p *DA) Get(ctx context.Context, height uint, hash string) (*configs.CheckpointExport, error) {
	nid, net, name, mp := p.Config.NamespaceID, p.Config.Network, p.Config.Name, p.MetaProtocol

	if net == constant.PreAlphaTestNet {
		sdk.SetNet(constant.PreAlphaTestNet)
	} else if net == constant.TestNet {
		sdk.SetNet(constant.TestNet)
	} else {
		return nil, fmt.Errorf("unknown network: %s", net)
	}

	clientDA := sdk.NewNubit(sdk.WithCtx(ctx)).Client
	resp, err := clientDA.GetTotalDataIDsInNamesapce(ctx, &types.GetTotalDataIDsInNamesapceReq{NID: nid})
	if err != nil {
		return nil, fmt.Errorf("failed to get data IDs in namespace: namespaceID=%s, err=%v", nid, err)
	}
	count := int(resp.Count)
	if count == 0 {
		return nil, fmt.Errorf("empty data in namespace %q", nid)
	}

	for i := 0; i < DefaultRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ck, err := p.doDownload(nid, name, mp, strconv.Itoa(int(height)), hash, count-1)
			if err != nil {
				time.Sleep(20 * time.Second)
				logs.Error.Println("Download checkpoint from DA err:", err)
				continue
			}
			return &configs.CheckpointExport{Checkpoint: ck, SourceDA: p.Config}, nil
		}
	}
	return nil, fmt.Errorf("get DA checkpoint max retires exceeded: height=%d, hash=%s", height, hash)
}

func (p *DA) doDownload(namespaceID, name, metaProtocol, height, hash string, offset int) (*checkpoint.Checkpoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	clientDA := sdk.NewNubit(sdk.WithCtx(ctx)).Client
	resp, err := clientDA.GetDataInNamespace(
		ctx,
		&types.GetDataInNamespaceReq{NID: namespaceID, Limit: 100, Offset: offset},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get namespace data error: offset=%d, namespaceID=%s, err=%v",
			offset,
			namespaceID,
			err,
		)
	}

	dataIDs := resp.DataIDs
	if len(dataIDs) == 0 {
		return nil, fmt.Errorf(
			"empty data IDs: offset=%d, namespaceID=%s",
			offset,
			namespaceID,
		)
	}

	for _, dataID := range dataIDs {
		data, err := clientDA.GetData(ctx, &types.GetDataReq{DAID: dataID})
		if err != nil {
			return nil, fmt.Errorf(
				"get DA data error: offset=%d, namespaceID=%s, err=%v",
				offset,
				namespaceID,
				err,
			)
		}

		rawBytes, err := base64.StdEncoding.DecodeString(data.RawData)
		if err != nil {
			return nil, fmt.Errorf(
				"decode checkpoint error: offset=%d, namespaceID=%s, err=%v",
				offset,
				namespaceID,
				err,
			)
		}

		var c checkpoint.Checkpoint
		if err := json.Unmarshal(rawBytes, &c); err != nil {
			return nil, fmt.Errorf(
				"unmarshal checkpoint error: offset=%d, namespaceID=%s, err=%v",
				offset,
				namespaceID,
				err,
			)
		}
		if strings.EqualFold(c.Name, name) &&
			strings.EqualFold(c.MetaProtocol, metaProtocol) &&
			strings.EqualFold(c.Height, height) &&
			strings.EqualFold(c.Hash, hash) {
			return &c, nil
		}
	}

	return nil, fmt.Errorf(
		"checkpoint not found on DA: offset=%d, dataIDs=%v, namespaceID=%s",
		offset,
		dataIDs,
		namespaceID,
	)
}
