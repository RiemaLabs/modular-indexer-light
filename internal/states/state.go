package states

import (
	"context"
	"encoding/base64"
	"errors"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
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
	State atomic.Int64

	denyListPath string

	providers []provider.CheckpointProvider

	// The consistent check point at the current height - 1.
	lastCheckpoint *configs.CheckpointExport

	// The checkpoints got from providers at the current height.
	currentCheckpoints []*configs.CheckpointExport

	// The number of effective providers should exceed the minimum required.
	minimalCheckpoint int

	// timeout for request checkpoint.
	timeout time.Duration

	sync.RWMutex
}

var S *State

func New(
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

func Init(
	denyListPath string,
	providers []provider.CheckpointProvider,
	lastCheckpoint *configs.CheckpointExport,
	minimalCheckpoint int,
	fetchTimeout time.Duration,
) {
	S = New(denyListPath, providers, lastCheckpoint, minimalCheckpoint, fetchTimeout)
}

func (s *State) CurrentHeight() uint {
	ck := s.CurrentFirstCheckpoint()
	if ck == nil {
		return 0
	}
	h, err := strconv.ParseUint(ck.Checkpoint.Height, 10, 64)
	if err != nil {
		logs.Error.Printf("parse checkpoint height failed: %v", err)
	}
	return uint(h)
}

func (s *State) UpdateCheckpoints(height uint, hash string) error {
	s.Lock()
	defer s.Unlock()

	s.State.Store(int64(constant.StatusSync))

	// Get checkpoints from the providers.
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	checkpoints, err := provider.GetCheckpoints(ctx, s.providers, height, hash)
	if err != nil {
		return err
	}
	if len(checkpoints) < s.minimalCheckpoint {
		return errors.New("not enough checkpoints fetched")
	}

	inconsistent := provider.CheckpointsInconsistent(checkpoints)
	if inconsistent {
		logs.Warn.Printf("Inconsistent checkpoints: height=%d, hash=%s, starting verification and reconstruction...", height, hash)

		// Aggregate checkpoints by commitment.
		aggregates := make(map[string]*configs.CheckpointExport)
		for _, ck := range checkpoints {
			if _, exist := aggregates[ck.Checkpoint.Commitment]; exist {
				continue
			}
			aggregates[ck.Checkpoint.Commitment] = ck
		}

		type succCommit struct {
			commitment  string
			transferLen int
		}

		succCommits := make(chan succCommit, len(aggregates))
		var wg sync.WaitGroup
		for commit, ck := range aggregates {
			wg.Add(1)
			go func(commit string, ck *checkpoint.Checkpoint) {
				defer wg.Done()

				stateProof, err := committee.New(context.Background(), ck.Name, ck.URL).LatestStateProof()
				if err != nil {
					logs.Error.Printf(
						"Failed to get latest state proof from the committee indexer: commit=%s, name=%s, url=%s, err=%v",
						commit,
						ck.Name,
						ck.URL,
						err,
					)
					return
				}
				if errMsg := stateProof.Error; errMsg != nil {
					logs.Error.Printf(
						"Latest state proof error from the committee indexer: commit=%s, name=%s, url=%s, err=%s",
						commit,
						ck.Name,
						ck.URL,
						*errMsg,
					)
					return
				}

				// Verify Ordinals transfers via Bitcoin.
				var ordTransfers []getter.OrdTransfer
				for _, tran := range stateProof.Result.OrdTransfers {
					contentBytes, err := base64.StdEncoding.DecodeString(tran.Content)
					if err != nil {
						logs.Error.Printf(
							"Invalid Ordinals transfer content: commit=%s, name=%s, url=%s, err=%v",
							commit,
							ck.Name,
							ck.URL,
							err,
						)
						return
					}
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

				curHeight, _ := strconv.ParseInt(ck.Height, 10, 64)
				ok, err := transfer.VerifyOrdTransfer(ordTransfers, uint(curHeight))
				if err != nil || !ok {
					logs.Error.Printf("Ordinals transfers verification error: err=%v, ok=%v", err, ok)
					return
				}

				// Generate current checkpoint
				preCheckpoint := s.lastCheckpoint
				prePointByte, err := base64.StdEncoding.DecodeString(preCheckpoint.Checkpoint.Commitment)
				if err != nil {
					return
				}
				prePoint := new(verkle.Point)
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
				currentCommit := base64.StdEncoding.EncodeToString(postBytes[:])

				if currentCommit != commit {
					return
				}

				succCommits <- succCommit{
					commitment:  commit,
					transferLen: len(ordTransfers),
				}

			}(commit, ck.Checkpoint)
		}
		wg.Wait()

		close(succCommits)
		var succVerify []succCommit
		for c := range succCommits {
			succVerify = append(succVerify, c)
		}
		if len(succVerify) == 0 {
			return errors.New("all checkpoints verify failed")
		}

		maxTransfer := succVerify[0].transferLen
		champion := 0
		var seemRight []string
		for i := 1; i < len(succVerify); i++ {
			seemRight = append(seemRight, succVerify[i].commitment)
			if succVerify[i].transferLen > maxTransfer {
				maxTransfer = succVerify[i].transferLen
				champion = i
			}
		}
		trustCommitment := succVerify[champion].commitment

		s.lastCheckpoint, s.currentCheckpoints = s.currentCheckpoints[0], []*configs.CheckpointExport{aggregates[trustCommitment]}
		s.State.Store(int64(constant.StateActive))

		// Deny untrusted providers.
		for _, ck := range checkpoints {
			if !slices.Contains(seemRight, ck.Checkpoint.Commitment) && s.denyListPath != "" {
				provider.DenyCheckpoint(s.denyListPath, aggregates[trustCommitment], ck)
			}
		}
	} else {
		s.lastCheckpoint, s.currentCheckpoints = s.currentCheckpoints[0], checkpoints
		s.State.Store(int64(constant.StateActive))
	}

	c := s.currentCheckpoints[0].Checkpoint.Commitment
	if inconsistent {
		logs.Info.Printf("Checkpoints fetched from providers have been verified, the commitment: %s, current height %d, hash %s", c, height, hash)
	} else {
		logs.Info.Printf("Checkpoints fetched from providers are consistent, the commitment: %s, current height %d, hash %s", c, height, hash)
	}

	return nil
}

func (s *State) LastCheckpoint() *configs.CheckpointExport {
	s.RLock()
	defer s.RUnlock()
	return s.lastCheckpoint
}

func (s *State) CurrentCheckpoints() []*configs.CheckpointExport {
	s.RLock()
	defer s.RUnlock()
	return s.currentCheckpoints
}

func (s *State) CurrentFirstCheckpoint() *configs.CheckpointExport {
	s.RLock()
	defer s.RUnlock()
	return s.currentCheckpoints[0]
}
