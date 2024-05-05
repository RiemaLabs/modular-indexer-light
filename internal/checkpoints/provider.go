package checkpoints

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
)

const DefaultRetries = 3

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

func Deny(path string, correct, fraud *configs.CheckpointExport) {
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

func Inconsistent(checkpoints []*configs.CheckpointExport) bool {
	for i := 0; i < len(checkpoints)-1; i++ {
		if !Equal(checkpoints[i].Checkpoint, checkpoints[i+1].Checkpoint) {
			return true
		}
	}
	return false
}

func Equal(a, b *checkpoint.Checkpoint) bool {
	return a.Commitment == b.Commitment && a.Hash == b.Hash && a.Height == b.Height
}
