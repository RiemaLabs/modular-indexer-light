package runtime

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-committee/ord/getter"
	"github.com/ethereum/go-verkle"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/committee"
	"github.com/RiemaLabs/modular-indexer-light/internal/clients/ord/transfer"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/constant"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/provider"
)

// TODO: Medium. Uniform the error report.

type State struct {
	denyListPath string

	providers []provider.CheckpointProvider

	// The consistent check point at the current height - 1.
	lastCheckpoint *configs.CheckpointExport

	// The checkpoints got from providers at the current height.
	currentCheckpoints []*configs.CheckpointExport

	// The number of effective providers should exceed the minimum required.
	minimalCheckpoint int

	// timeout for request checkpoint, uint second
	timeout time.Duration

	sync.RWMutex
}

func NewRuntimeState(
	denyListPath string,
	providers []provider.CheckpointProvider,
	lastCheckpoint *configs.CheckpointExport,
	minimalCheckpoint int,
	fetchTimeout time.Duration,
) *State {
	return &State{
		denyListPath:       denyListPath,
		providers:          providers,
		lastCheckpoint:     lastCheckpoint,
		currentCheckpoints: make([]*configs.CheckpointExport, len(providers)),
		minimalCheckpoint:  minimalCheckpoint,
		timeout:            fetchTimeout,
	}
}

func (s *State) CurrentHeight() uint {
	ck := s.CurrentFirstCheckpoint()
	if ck == nil {
		return 0
	}

	h, err := strconv.ParseUint(ck.Checkpoint.Height, 10, 64)
	if err != nil {
		logs.Error.Printf("ParseUint(checkpoint.Height) failed", "error:", err)
	}
	return uint(h)
}

func (s *State) UpdateCheckpoints(height uint, hash string) error {
	s.Lock()
	defer s.Unlock()
	constant.ApiState = constant.StatusSync

	// Get checkpoints from the providers.
	checkpoints := provider.GetCheckpoints(s.providers, height, hash, s.timeout)
	if len(checkpoints) < int(s.minimalCheckpoint) {
		return errors.New("not enough checkpoints fetched")
	}

	_, _, inconsistent := provider.CheckpointsInconsist(checkpoints)
	if inconsistent {
		constant.ApiState = constant.StatusVerify

		logs.Warn.Printf("Checkpoints retrieved from providers are inconsistent for height %q, hash %q, start the verification and regeneration process...", height, hash)

		// Aggregate checkpoints by commitment
		aggs := make(map[string]*configs.CheckpointExport)
		for _, ck := range checkpoints {
			if _, exist := aggs[ck.Checkpoint.Commitment]; exist {
				continue
			}
			aggs[ck.Checkpoint.Commitment] = ck
		}

		type succCommit struct {
			commitment  string
			transferLen int
		}

		succVerify := make([]succCommit, 0, len(aggs))
		var wg sync.WaitGroup
		for commit, ck := range aggs {
			wg.Add(1)
			go func(commit string, ck *checkpoint.Checkpoint) {
				defer wg.Done()

				stateProof, err := committee.New(context.TODO(), ck.Name, ck.URL).LatestStateProof()
				if err != nil || stateProof.Error != nil {
					logs.Warn.Printf("Failed to get latest state proof from the committee indexer: name=%s, url=%s, err:=%v", ck.Name, ck.URL, err)
					return
				}

				// Verify ordTransfers via Bitcoin
				curHeight, _ := strconv.ParseInt(ck.Height, 10, 64)

				var ordTransfers []getter.OrdTransfer
				for _, tran := range stateProof.Result.OrdTransfers {
					contentBytes, _ := base64.StdEncoding.DecodeString(tran.Content)
					ordTransfers = append(ordTransfers, getter.OrdTransfer{
						ID:            tran.ID,
						InscriptionID: tran.InscriptionID,
						OldSatpoint:   tran.OldSatpoint,
						NewSatpoint:   tran.NewSatpoint,
						NewPkscript:   tran.NewPkscript,
						NewWallet:     tran.NewWallet,
						SentAsFee:     tran.SentAsFee,
						Content:       contentBytes,
						ContentType:   tran.ContentType,
					})
				}

				ok, err := transfer.VerifyOrdTransfer(ordTransfers, uint(curHeight))
				if err != nil || !ok {
					return
				}

				// Generate current checkpoint
				preCheckpoint := s.LastCheckpoint()
				prePointByte, err := base64.StdEncoding.DecodeString(preCheckpoint.Checkpoint.Commitment)
				if err != nil {
					return
				}
				prePoint := &verkle.Point{}
				err = prePoint.SetBytes(prePointByte)
				if err != nil {
					return
				}

				node, err := apis.GeneratePostRoot(prePoint, height, stateProof)
				if err != nil {
					return
				}
				if node == nil {
					return
				}

				postBytes := node.Commit().Bytes()
				curentCommit := base64.StdEncoding.EncodeToString(postBytes[:])

				if curentCommit != commit {
					return
				}

				succVerify = append(succVerify, succCommit{
					commitment:  commit,
					transferLen: len(ordTransfers),
				})

			}(commit, ck.Checkpoint)
		}
		wg.Wait()

		if len(succVerify) == 0 {
			return errors.New("all checkpoints verify failed")
		}

		maxTransfer := succVerify[0].transferLen
		champion := 0
		seemRight := []string{}
		for i := 1; i < len(succVerify); i++ {
			seemRight = append(seemRight, succVerify[i].commitment)
			if succVerify[i].transferLen > maxTransfer {
				maxTransfer = succVerify[i].transferLen
				champion = i
			}
		}
		trustCommitment := succVerify[champion].commitment

		s.lastCheckpoint, s.currentCheckpoints = s.CurrentFirstCheckpoint(), []*configs.CheckpointExport{aggs[trustCommitment]}
		constant.ApiState = constant.StateActive

		// Deny untrusted providers.
		for _, ck := range checkpoints {
			if !slices.Contains(seemRight, ck.Checkpoint.Commitment) {
				provider.DenyCheckpoint(s.denyListPath, aggs[trustCommitment], ck)
			}
		}
	} else {
		s.lastCheckpoint, s.currentCheckpoints = s.CurrentFirstCheckpoint(), checkpoints
		constant.ApiState = constant.StateActive
	}

	c := s.CurrentFirstCheckpoint().Checkpoint.Commitment
	if inconsistent {
		logs.Info.Printf(fmt.Sprintf("Checkpoints fetched from providers have been verified, the commitment: %s, current height %d, hash %s", c, height, hash))
	} else {
		logs.Info.Printf(fmt.Sprintf("Checkpoints fetched from providers are consistent, the commitment: %s, current height %d, hash %s", c, height, hash))
	}

	return nil
}

func (s *State) LastCheckpoint() *configs.CheckpointExport {
	return s.lastCheckpoint
}

func (s *State) CurrentCheckpoints() []*configs.CheckpointExport {
	return s.currentCheckpoints
}

func (s *State) CurrentFirstCheckpoint() *configs.CheckpointExport {
	return s.currentCheckpoints[0]
}
