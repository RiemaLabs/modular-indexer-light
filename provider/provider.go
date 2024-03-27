package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/RiemaLabs/nubit-da-sdk/constant"
	"github.com/RiemaLabs/nubit-da-sdk/types"

	sdk "github.com/RiemaLabs/nubit-da-sdk"
)

func GetCheckpoints(providers []CheckpointProvider, height uint, hash string, timeout time.Duration) []*config.CheckpointExport {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	result := make([]*config.CheckpointExport, 0, len(providers))
	var wg sync.WaitGroup
	for _, p := range providers {
		p := p
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					ck, err := p.GetCheckpoint(ctx, height, hash)
					if err != nil || ck == nil {
						continue
					}
					result = append(result, ck)
					return
				}
			}
		}()
	}
	wg.Wait()
	return result
}

func RecordBlacklist(correct, fraud *config.CheckpointExport) {
	h, _ := strconv.ParseUint(correct.Checkpoint.Height, 10, 64)
	b := config.Blacklist{
		Evidence: &config.Evidence{
			Height:            uint(h),
			Hash:              correct.Checkpoint.Hash,
			CorrectCommitment: correct.Checkpoint.Commitment,
			FraudCommitment:   fraud.Checkpoint.Commitment,
		},
	}
	if fraud.SourceDA != nil {
		b.SourceDA = fraud.SourceDA
	}

	if fraud.SourceS3 != nil {
		b.SourceS3 = fraud.SourceS3
	}

	config.AppendBlacklist(&b)
}

// Find the first inconsistent checkpoint pair among multiple ones.
func CheckpointsInconsist(checkpoints []*config.CheckpointExport) (int, int, bool) {
	for i := 0; i < len(checkpoints)-1; i++ {
		if !CheckPointEqual(checkpoints[i].Checkpoint, checkpoints[i+1].Checkpoint) {
			return i, i + 1, true
		}
	}
	return 0, 0, false
}

func CheckPointEqual(a, b *checkpoint.Checkpoint) bool {
	if a.Commitment != b.Commitment ||
		a.Hash != b.Hash ||
		a.Height != b.Height {
		return false
	}
	return true
}

func DownloadCheckpointByDA(namespaceID, network string, name, metaProtocol, height, hash string, runtimeOffset int, timeout time.Duration) (*checkpoint.Checkpoint, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if network == "Pre-Alpha Testnet" {
		sdk.SetNet(constant.PreAlphaTestNet)
	} else if network == "Testnet" {
		sdk.SetNet(constant.TestNet)
	} else {
		return nil, 0, fmt.Errorf("unknown network: %s", network)
	}

	clientDA := sdk.NewNubit(sdk.WithCtx(ctx)).Client

	resDataIDs, err := clientDA.GetDataInNamespace(ctx, &types.GetDataInNamespaceReq{
		NID:    namespaceID,
		Limit:  100,
		Offset: runtimeOffset,
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data with offset %d, in namespace %s, error: %v", runtimeOffset, namespaceID, err)
	}

	dataIDs := resDataIDs.DataIDs

	if len(dataIDs) == 0 {
		return nil, 0, fmt.Errorf("the count of data with offset %d, in namespace %s, error: %v", runtimeOffset, namespaceID, err)
	}

	var c checkpoint.Checkpoint

	for _, dataID := range dataIDs {
		runtimeOffset += 1
		datas, err := clientDA.GetData(ctx, &types.GetDataReq{
			DAID: dataID,
		})
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get data with offset %d, in namespace %s, error: %v", runtimeOffset, namespaceID, err)
		}

		decodeString, err := base64.StdEncoding.DecodeString(datas.RawData)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse data with offset %d, in namespace %s, error: %v", runtimeOffset, namespaceID, err)
		}
		err = json.Unmarshal(decodeString, &c)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse data with offset %d, in namespace %s, error: %v", runtimeOffset, namespaceID, err)
		}
		if strings.EqualFold(c.Name, name) && strings.EqualFold(c.MetaProtocol, metaProtocol) && strings.EqualFold(c.Height, height) && strings.EqualFold(c.Hash, hash) {
			return &c, runtimeOffset, nil
		}
	}
	return nil, 0, fmt.Errorf("failed to find valid checkpoint with offset %d, in namespace %s", runtimeOffset, namespaceID)
}

func DownloadCheckpointByS3(region, bucket string, name, metaProtocol, height, hash string, timeout time.Duration) (*checkpoint.Checkpoint, error) {

	objectKey := fmt.Sprintf("checkpoint-%s-%s-%s-%s.json", name, metaProtocol, height, hash)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, objectKey)

	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c checkpoint.Checkpoint
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file, error: %v", err)
	}

	return &c, nil
}
