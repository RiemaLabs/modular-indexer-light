package checkpoint

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/RiemaLabs/indexer-committee/apis"
	"github.com/RiemaLabs/indexer-committee/ord"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/indexer"
	"github.com/RiemaLabs/indexer-light/provide"
	"github.com/RiemaLabs/indexer-light/types"
	"github.com/ethereum/go-verkle"
)

// VerifyCheckpoint Obtain and verify whether the checkpoints of m committee members are consistent.
func VerifyCheckpoint(getter getter.OrdGetter, config *types.Config) {
	height, err := getter.GetLatestBlockHeight()
	if err != nil {
		return
	}
	committeeIndexer := provide.GetCommitteeIndexers(config)
	var Checkpoints []*types.CheckPointObject
	if len(committeeIndexer) > 0 {
		for i, _ := range committeeIndexer {
			go func(p types.CheckPointProvider, height uint) {
				Checkpoints = append(Checkpoints, p.GetCheckpoint(height))
			}(committeeIndexer[i], height)
		}
	}
	for len(Checkpoints) < config.MinimalCheckPoint {
		continue
	}

	diffMap := make(map[string]*types.CheckPointObject)
	for i, _ := range Checkpoints {
		if _, ok := diffMap[Checkpoints[i].CheckPoint.Commitment]; ok {
			continue
		}
		diffMap[Checkpoints[i].CheckPoint.Commitment] = Checkpoints[i]
	}

	// The checkpoints of m committee members are inconsistent
	if len(diffMap) > 0 {
		panic("The checkpoints of m committee members are inconsistent")
	}

}

func GePretCheckpoint(getter getter.OrdGetter, config *types.Config, diffCheckpoint map[string]*types.CheckPointObject) *types.CheckPointObject {
	if len(diffCheckpoint) == 0 {
		return nil
	}
	height, err := getter.GetLatestBlockHeight()
	if err != nil {
		return nil
	}

	// Request the two Committee Indexers respectively: the set of state changes from the state corresponding to the parent block to the current block state

	// The state change set that transitions from the state  corresponding to the parent block to the current block state
	var diffState = make(map[string]*apis.Brc20VerifiableGetCurrentStateProofResponse)
	// Request the Checkpoint of the parent block (under h-1 height)
	var diffPreCheckpoint = make(map[string]*types.CheckPointObject)
	for key, cpo := range diffCheckpoint {
		Committee := indexer.NewClient(context.Background(), cpo.CheckPoint.Name, cpo.CheckPoint.URL)
		state, err := Committee.StateDiff()
		if err != nil {
			continue
		}
		diffState[key] = state
		switch true {
		case cpo.Source.SourceS3 != nil:
			preCp := provide.NewS3(cpo.Source.SourceS3).GetCheckpoint(height - 1)
			diffPreCheckpoint[key] = preCp
		case cpo.Source.SourceDa != nil:
			preCp := provide.NewDA(cpo.Source.SourceDa).GetCheckpoint(height - 1)
			diffPreCheckpoint[key] = preCp
		}
	}

	transfers, err := getter.GetOrdTransfers(height)
	if err != nil {
		return nil
	}

	// 验证证明proof是否合法
	for key, state := range diffState {
		prePointByte, err := base64.StdEncoding.DecodeString(diffPreCheckpoint[key].CheckPoint.Commitment)
		if err != nil {
			return nil
		}
		prePoint := &verkle.Point{}
		err = prePoint.SetBytes(prePointByte)
		if err != nil {
			return nil
		}

		//TODO:: state.Proof to verkle.Proof
		waitProof := &verkle.Proof{}

		preStatePartial, err := verkle.PreStateTreeFromProof(waitProof, prePoint)
		if err != nil {
			return nil
		}

		err = verkle.VerifyVerkleProofWithPreState(waitProof, preStatePartial)
		if err != nil {
			return nil
		}

		tmpState := ord.State{
			Root:   preStatePartial,
			Height: height,
			Hash:   "",
			KV:     make(ord.KeyValueMap),
		}

		// Light clients computes the partial postState from the partial preState.
		// Then verifies: partial postState->partial preState is consistent with the stateDiff in the proofOfStateTrans.
		// calculate State
		calculateState := ord.Exec(tmpState, transfers)
		calculatebytes := calculateState.Root.Commit().Bytes()
		decodeString, err := base64.StdEncoding.DecodeString(diffCheckpoint[key].CheckPoint.Commitment)
		if err != nil {
			return nil
		}
		if !bytes.Equal(calculatebytes[:], decodeString) {
			// Not equal, indicating that there is a problem with the current state.Proof

		}

	}

}

func VerifyProof(getter getter.OrdGetter, preCheckpoint, postCheckpoint *types.CheckPointObject) {

}

func verifyStateDiff(getter getter.OrdGetter, preCheckpoint, postCheckpoint *types.CheckPointObject) {
	preCommittee := indexer.NewClient(context.Background(), preCheckpoint.CheckPoint.Name, preCheckpoint.CheckPoint.URL)
	preDiff, err := preCommittee.StateDiff()
	if err != nil {
		return
	}

	height, err := getter.GetLatestBlockHeight()
	if err != nil {
		return
	}

	transfers, err := getter.GetOrdTransfers(height)
	if err != nil {
		return
	}

	keys := preDiff.Keys
	PreValues := preDiff.PreValues
	PostValues := preDiff.PostValues
	proof := preDiff.Proof

	postCommittee := indexer.NewClient(context.Background(), postCheckpoint.CheckPoint.Name, postCheckpoint.CheckPoint.URL)
	postDiff, err := postCommittee.StateDiff()
	if err != nil {
		return
	}
}
