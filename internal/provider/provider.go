package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	sdk "github.com/RiemaLabs/nubit-da-sdk"
	"github.com/RiemaLabs/nubit-da-sdk/constant"
	"github.com/RiemaLabs/nubit-da-sdk/types"

	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
)

type CheckpointProvider interface {
	Get(ctx context.Context, height uint, hash string) (*configs.CheckpointExport, error)
}

func GetCheckpoints(ctx context.Context, providers []CheckpointProvider, height uint, hash string) ([]*configs.CheckpointExport, error) {
	var (
		wg          sync.WaitGroup
		errs        = make(chan error, len(providers))
		checkpoints = make(chan *configs.CheckpointExport, len(providers))
	)
	for _, p := range providers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				default:
					ck, err := p.Get(ctx, height, hash)
					if err != nil {
						logs.Error.Printf("Get checkpoint error: height=%d, hash=%s, err=%v", height, hash, err)
						continue
					}
					checkpoints <- ck
					return
				}
			}
		}()
	}
	wg.Wait()

	close(errs)
	var retErrs []error
	for err := range errs {
		retErrs = append(retErrs, err)
	}

	close(checkpoints)
	var ret []*configs.CheckpointExport
	for ck := range checkpoints {
		ret = append(ret, ck)
	}

	return ret, errors.Join(retErrs...)
}

func DenyCheckpoint(path string, correct, fraud *configs.CheckpointExport) {
	h, _ := strconv.ParseUint(correct.Checkpoint.Height, 10, 64)
	b := configs.DenyList{
		Evidence: &configs.Evidence{
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

	if err := configs.AppendDenyList(path, &b); err != nil {
		logs.Error.Println("Append to deny list error:", err)
	}
}

// Find the first inconsistent checkpoint pair among multiple ones.
func CheckpointsInconsist(checkpoints []*configs.CheckpointExport) (int, int, bool) {
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
		return nil, 0, fmt.Errorf("the count of data with offset %d is zero, in namespace %s", runtimeOffset, namespaceID)
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