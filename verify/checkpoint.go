package verify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/RiemaLabs/indexer-committee/apis"
	"github.com/RiemaLabs/indexer-committee/checkpoint"
	ordgetter "github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-committee/ord/stateless"
	"github.com/RiemaLabs/indexer-light/config"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/indexer"
	"github.com/RiemaLabs/indexer-light/log"
	"github.com/RiemaLabs/indexer-light/provide"
	"github.com/RiemaLabs/indexer-light/types"
	"github.com/ethereum/go-verkle"
)

// VerifyCheckpoint Obtain and verify whether the checkpoints of m committee members are consistent.
func VerifyCheckpoint(getter ordgetter.OrdGetter, config *types.Config) error {
	constant.ApiState = constant.ApiStateLoading
	ctx := context.Background()
	height, err := getter.GetLatestBlockHeight()
	if err != nil {
		return err
	}
	hash, err := getter.GetBlockHash(height)
	if err != nil {
		return err
	}
	config.StartHeight = int(height)
	config.StartBlockHash = hash
	committeeIndexer := provide.GetCommitteeIndexers(config)
	var Checkpoints []*types.CheckPointObject
	if len(committeeIndexer) == 0 {
		panic("Invalid CommitteeIndexer provide")
	}
	for i, _ := range committeeIndexer {
		go func(p types.CheckPointProvider, height uint) {
			Checkpoints = append(Checkpoints, p.GetCheckpoint(ctx, height, hash))
		}(committeeIndexer[i], height)
	}
	for len(Checkpoints) < config.MinimalCheckPoint {
		continue
	}

	diffMap := make(map[string]*types.CheckPointObject)
	for i, _ := range Checkpoints {
		if _, ok := diffMap[Checkpoints[i].CheckPoint.Commitment]; ok {
			err = DefiniteState.Update(getter, Checkpoints[i].CheckPoint)
			if err != nil {
				return err
			}
			continue
		}
		diffMap[Checkpoints[i].CheckPoint.Commitment] = Checkpoints[i]
	}

	// The checkpoints of m committee members are inconsistent
	if len(diffMap) > 0 {
		err = verifyCheckpoint(getter, config, diffMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func verifyCheckpoint(getter ordgetter.OrdGetter, config *types.Config, diffCheckpoint map[string]*types.CheckPointObject) error {
	constant.ApiState = constant.ApiStateSync
	ctx := context.Background()
	if len(diffCheckpoint) == 0 {
		return nil
	}

	height, err := getter.GetLatestBlockHeight()
	if err != nil {
		return nil
	}

	preHeight := height - 1
	preHash, err := getter.GetBlockHash(preHeight)
	if err != nil {
		return err
	}

	// Request the two Committee Indexers respectively: the set of state changes from the state corresponding to the parent block to the current block state
	// The state change set that transitions from the state  corresponding to the parent block to the current block state
	var diffState = make(map[string]*apis.Brc20VerifiableLatestStateProofResponse)
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
			preCp := provide.NewS3(cpo.Source.SourceS3).GetCheckpoint(ctx, preHeight, preHash)
			diffPreCheckpoint[key] = preCp
		case cpo.Source.SourceDa != nil:
			preCp := provide.NewDA(cpo.Source.SourceDa).GetCheckpoint(ctx, preHeight, preHash)
			diffPreCheckpoint[key] = preCp
		}
	}

	//transfers, err := getter.GetOrdTransfers(height)
	//if err != nil {
	//	return nil
	//}

	// Definitely something wrong checkpoint
	var wrongCheckpoint []*types.CheckPointObject

	// Verify whether the proof is legal
	for key, state := range diffState {
		preCheckpoint := diffPreCheckpoint[key].CheckPoint
		prePointByte, err := base64.StdEncoding.DecodeString(preCheckpoint.Commitment)
		if err != nil {
			return nil
		}
		prePoint := &verkle.Point{}
		err = prePoint.SetBytes(prePointByte)
		if err != nil {
			return nil
		}

		preProofByte, err := base64.StdEncoding.DecodeString(state.Proof)
		if err != nil {
			return nil
		}
		preVProof := &verkle.VerkleProof{}
		err = preVProof.UnmarshalJSON(preProofByte)
		if err != nil {
			return err
		}

		preProof, err := verkle.DeserializeProof(preVProof, nil)
		if err != nil {
			return err
		}

		//preStatePartial  == preCheckpoint
		preStatePartial, err := verkle.PreStateTreeFromProof(preProof, prePoint) // preStatePartial===VerkleNode
		if err != nil {
			return nil
		}

		// preCheckpoint   preProof  preStatePartial
		err = verkle.VerifyVerkleProofWithPreState(preProof, preStatePartial)
		if err != nil {
			return nil
		}

		preState := &stateless.Header{
			Root:   preStatePartial,
			Height: height,
			Hash:   "",
			KV:     make(stateless.KeyValueMap),
		}

		// Light clients computes the partial postState from the partial preState.
		// Then verifies: partial postState->partial preState is consistent with the stateDiff in the proofOfStateTrans.
		// calculate State
		//stateless.Exec(preState, transfers)
		var ordTransfers []ordgetter.OrdTransfer
		if len(state.OrdTrans) > 0 {
			for _, tran := range state.OrdTrans {
				decodeString, err := base64.StdEncoding.DecodeString(tran.Content)
				if err != nil {
					return err
				}
				ordTransfers = append(ordTransfers, ordgetter.OrdTransfer{
					ID:            tran.ID,
					InscriptionID: tran.InscriptionID,
					OldSatpoint:   tran.NewSatpoint,
					NewSatpoint:   tran.NewSatpoint,
					//NewPkScript:   tran.NewPkScript,
					NewWallet:   tran.NewWallet,
					SentAsFee:   tran.SentAsFee,
					Content:     decodeString,
					ContentType: tran.ContentType,
				})
			}
		}
		stateless.Exec(preState, ordTransfers)
		calculatebytes := preState.Root.Commit().Bytes()
		decodeString, err := base64.StdEncoding.DecodeString(diffCheckpoint[key].CheckPoint.Commitment)
		if err != nil {
			return nil
		}
		if !bytes.Equal(calculatebytes[:], decodeString) {
			// Not equal, indicating that there is a problem with the current state.Proof
			wrongCheckpoint = append(wrongCheckpoint, diffCheckpoint[key])
		}
	}
	return eliminateBadCommittee(config, wrongCheckpoint)
}

func VerifyProof(preProof *verkle.Proof, point *verkle.Point) error {
	// preStatePartial  == preCheckpoint
	preStatePartial, err := verkle.PreStateTreeFromProof(preProof, point)
	if err != nil {
		return err
	}
	// preCheckpoint   preProof  preStatePartial
	err = verkle.VerifyVerkleProofWithPreState(preProof, preStatePartial)
	if err != nil {
		return err
	}
	return nil
}

// Eliminate bad committee members
func eliminateBadCommittee(cfg *types.Config, wrongCheckpoint []*types.CheckPointObject) error {
	if len(wrongCheckpoint) > 0 {
		for _, object := range wrongCheckpoint {
			if object.Source != nil {
				switch true {
				case object.Source.SourceS3 != nil && cfg.CommitteeIndexer.S3 != nil:
					var ns3 []*types.SourceS3
					for _, s3 := range cfg.CommitteeIndexer.S3 {
						if s3.IndexerName == object.Source.SourceS3.IndexerName && s3.ApiUrl == object.Source.SourceS3.Url {
							log.Error("Verify", "RemoveS3", object.Source.SourceS3.IndexerName, "msg", "This Committee Indexer generated an untrusted checkpoint", "errCheckpoint", checkpointStr(object.CheckPoint))
							continue
						}
						ns3 = append(ns3, s3)
					}
					cfg.CommitteeIndexer.S3 = ns3
				case object.Source.SourceDa != nil && cfg.CommitteeIndexer.Da != nil:
					var nda []*types.SourceDa
					for _, da := range cfg.CommitteeIndexer.Da {
						if da.IndexerName == object.Source.SourceS3.IndexerName && da.ApiUrl == object.Source.SourceS3.Url {
							log.Error("Verify", "RemoveDA", object.Source.SourceS3.IndexerName, "msg", "This Committee Indexer generated an untrusted checkpoint", "errCheckpoint", checkpointStr(object.CheckPoint))
							continue
						}
						nda = append(nda, da)
					}
					cfg.CommitteeIndexer.Da = nda
				}
			}
		}
		err := config.UpdateConfig(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkpointStr(point *checkpoint.Checkpoint) string {
	if point == nil {
		return ""
	}
	marshal, err := json.Marshal(point)
	if err != nil {
		return ""
	}
	return string(marshal)
}
