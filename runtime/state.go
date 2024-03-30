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
	"github.com/RiemaLabs/modular-indexer-light/clients/committee"
	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/log"
	"github.com/RiemaLabs/modular-indexer-light/ord/transfer"
	"github.com/RiemaLabs/modular-indexer-light/provider"
	"github.com/ethereum/go-verkle"
)

// TODO: Medium. Uniform the error report.

type RuntimeState struct {
	providers []provider.CheckpointProvider

	// The consistent check point at the current height - 1.
	lastCheckpoint *config.CheckpointExport

	// The checkpoints got from providers at the current height.
	currentCheckpoints []*config.CheckpointExport

	// The number of effective providers should exceed the minimum required.
	minimalCheckpoint int

	// timeout for request checkpoint, uint second
	timeout time.Duration

	sync.RWMutex
}

func NewRuntimeState(providers []provider.CheckpointProvider, lastCheckpoint *config.CheckpointExport, minimalCheckpoint int, fetchTimeout time.Duration) *RuntimeState {
	return &RuntimeState{
		providers:          providers,
		lastCheckpoint:     lastCheckpoint,
		currentCheckpoints: make([]*config.CheckpointExport, len(providers)),
		minimalCheckpoint:  minimalCheckpoint,
		timeout:            fetchTimeout,
	}
}

func (s *RuntimeState) CurrentHeight() uint {
	ck := s.CurrentFirstCheckpoint()

	h, err := strconv.ParseUint(ck.Checkpoint.Height, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(h)
}

func (s *RuntimeState) UpdateCheckpoints(height uint, hash string) error {
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

		log.Warn("Checkpoints retrieved from providers are inconsistent for height %d, hash %s. Beginning verification and regeneration process...", height, hash)

		// Aggregate checkpoints by commitment
		aggs := make(map[string]*config.CheckpointExport)
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

				// Get stateProof
				stateProof, err := committee.NewCommitteeIndexerClient(context.TODO(), ck.Name, ck.URL).LatestStateProof()
				if err != nil || stateProof.Error != nil {
					log.Warn("failed to getLatestStateProof from the committee indexer. Its name: %s, url: %s, error: %v", ck.Name, ck.URL, err)
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
						//OldSatpoint:   tran.NewSatpoint,
						NewSatpoint: tran.NewSatpoint,
						NewPkscript: tran.NewPkscript,
						NewWallet:   tran.NewWallet,
						SentAsFee:   tran.SentAsFee,
						Content:     contentBytes,
						ContentType: tran.ContentType,
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

		s.currentCheckpoints = []*config.CheckpointExport{aggs[trustCommitment]}
		constant.ApiState = constant.StateActive

		// black untrust provider
		for _, ck := range checkpoints {
			if !slices.Contains(seemRight, ck.Checkpoint.Commitment) {
				provider.RecordBlacklist(aggs[trustCommitment], ck)
			}
		}
	}

	s.currentCheckpoints = checkpoints
	constant.ApiState = constant.StateActive

	c := s.CurrentFirstCheckpoint().Checkpoint.Commitment
	if inconsistent {
		log.Info(fmt.Sprintf("checkpoints fetched from providers have been verified, the commitment: %s, current height %d, hash %s", c, height, hash))
	} else {
		log.Info(fmt.Sprintf("checkpoints fetched from providers are consistent, the commitment: %s, current height %d, hash %s", c, height, hash))
	}

	return nil
}

func (s *RuntimeState) LastCheckpoint() *config.CheckpointExport {
	return s.lastCheckpoint
}

func (s *RuntimeState) CurrentCheckpoints() []*config.CheckpointExport {
	return s.currentCheckpoints
}

func (s *RuntimeState) CurrentFirstCheckpoint() *config.CheckpointExport {
	return s.currentCheckpoints[0]
}
