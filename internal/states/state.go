package states

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-committee/ord/getter"
	"github.com/ethereum/go-verkle"

	"github.com/RiemaLabs/modular-indexer-light/internal/checkpoints"
	"github.com/RiemaLabs/modular-indexer-light/internal/clients/committee"
	"github.com/RiemaLabs/modular-indexer-light/internal/clients/ordi"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
)

// TODO: Medium. Uniform the error report.

type Status int

const (
	StatusActive Status = iota + 1
	StatusSync
	StatusVerify
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "ready"
	case StatusSync:
		return "syncing"
	case StatusVerify:
		return "verifying"
	default:
		return ""
	}
}

type State struct {
	State atomic.Int64

	denyListPath string

	providers []checkpoints.CheckpointProvider

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
	providers []checkpoints.CheckpointProvider,
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
	providers []checkpoints.CheckpointProvider,
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

	s.State.Store(int64(StatusSync))

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	cps, err := checkpoints.GetCheckpoints(ctx, s.providers, height, hash)
	if err != nil {
		return err
	}
	if l := len(cps); l < s.minimalCheckpoint {
		return fmt.Errorf("not enough checkpoints fetched: expected=%d, actual=%d", s.minimalCheckpoint, l)
	}

	if checkpoints.Inconsistent(cps) {
		logs.Warn.Printf("Inconsistent checkpoints at: height=%d, hash=%s", height, hash)

		aggregates := make(map[string]*configs.CheckpointExport)
		for _, ck := range cps {
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
			go func(checkpointCommit string, ck *checkpoint.Checkpoint) {
				defer wg.Done()

				committeeCl, err := committee.New(ck.URL)
				if err != nil {
					logs.Error.Printf(
						"Failed to create committee indexer client: commit=%s, name=%s, url=%s, err=%v",
						checkpointCommit,
						ck.Name,
						ck.URL,
						err,
					)
					return
				}
				stateProof, err := committeeCl.LatestStateProof(context.Background())
				if err != nil {
					logs.Error.Printf(
						"Failed to get latest state proof from the committee indexer: commit=%s, name=%s, url=%s, err=%v",
						checkpointCommit,
						ck.Name,
						ck.URL,
						err,
					)
					return
				}
				if errMsg := stateProof.Error; errMsg != nil {
					logs.Error.Printf(
						"Non-nil error message from the latest state proof: commit=%s, name=%s, url=%s, errMsg=%s",
						checkpointCommit,
						ck.Name,
						ck.URL,
						*errMsg,
					)
					return
				}

				var ordTransfers []getter.OrdTransfer
				for _, t := range stateProof.Result.OrdTransfers {
					contentBytes, err := base64.StdEncoding.DecodeString(t.Content)
					if err != nil {
						logs.Error.Printf(
							"Invalid Ordinals transfer content: commit=%s, name=%s, url=%s, err=%v",
							checkpointCommit,
							ck.Name,
							ck.URL,
							err,
						)
						return
					}
					ordTransfers = append(ordTransfers, getter.OrdTransfer{
						ID:            t.ID,
						InscriptionID: t.InscriptionID,
						OldSatpoint:   t.OldSatpoint,
						NewSatpoint:   t.NewSatpoint,
						NewPkscript:   t.NewPkscript,
						NewWallet:     t.NewWallet,
						SentAsFee:     t.SentAsFee,
						Content:       contentBytes,
						ContentType:   t.ContentType,
					})
				}

				curHeight, _ := strconv.ParseInt(ck.Height, 10, 64)
				if err := ordi.VerifyOrdTransfer(ordTransfers, uint(curHeight)); err != nil {
					logs.Error.Printf("Ordinals transfers verification error: err=%v", err)
					return
				}

				preCheckpoint := s.lastCheckpoint
				prePointByte, err := base64.StdEncoding.DecodeString(preCheckpoint.Checkpoint.Commitment)
				if err != nil {
					return
				}
				prePoint := new(verkle.Point)
				if err := prePoint.SetBytes(prePointByte); err != nil {
					return
				}

				node, err := apis.GeneratePostRoot(prePoint, height, stateProof)
				if err != nil {
					logs.Error.Printf("generate post root error: %v", err)
					return
				}
				if node == nil {
					return
				}

				postBytes := node.Commit().Bytes()
				calCommit := base64.StdEncoding.EncodeToString(postBytes[:])
				if calCommit != checkpointCommit {
					logs.Warn.Printf(
						"inconsistent commits: calCommit=%s, checkpointCommit=%s",
						calCommit,
						checkpointCommit,
					)
					return
				}

				succCommits <- succCommit{
					commitment:  checkpointCommit,
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
			return errors.New("all cps verify failed")
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
		s.State.Store(int64(StatusActive))

		for _, ck := range cps {
			if !slices.Contains(seemRight, ck.Checkpoint.Commitment) && s.denyListPath != "" {
				checkpoints.Deny(s.denyListPath, aggregates[trustCommitment], ck)
			}
		}

		c := s.currentCheckpoints[0].Checkpoint.Commitment
		logs.Info.Printf("Checkpoints fetched from providers have been verified, the commitment: %s, current height %d, hash %s", c, height, hash)
		return nil
	}

	s.lastCheckpoint, s.currentCheckpoints = s.currentCheckpoints[0], cps
	s.State.Store(int64(StatusActive))
	c := s.currentCheckpoints[0].Checkpoint.Commitment
	logs.Info.Printf("Checkpoints fetched from providers are all consistent: commitment=%s, height=%d, hash=%s", c, height, hash)

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
