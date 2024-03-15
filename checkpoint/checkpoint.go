package checkpoint

import (
	"bytes"
	"encoding/base64"

	"github.com/RiemaLabs/indexer-committee/checkpoint"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/provide"
	"github.com/RiemaLabs/indexer-light/types"
)

// VerifyCheckpoint Obtain and verify whether the checkpoints of m committee members are consistent.
func VerifyCheckpoint(getter getter.OrdGetter, config *types.Config) {
	provides := provide.GetProviders(config)
	var Checkpoints []*checkpoint.Checkpoint
	if len(provides) > 0 {
		for i, _ := range provides {
			go func(p types.CheckPointProvider) {
				Checkpoints = append(Checkpoints, p.GetCheckpoint())
			}(provides[i])
		}
	}
	for len(Checkpoints) < config.MinimalCheckPoint {
		continue
	}
	baseCheckpoint := Checkpoints[0]
	baseStateByte, err := base64.StdEncoding.DecodeString(baseCheckpoint.Commitment)
	if err != nil {
		return
	}
	diffCheckpoint := []*checkpoint.Checkpoint{}
	for _, c := range Checkpoints {
		stateByte, err := base64.StdEncoding.DecodeString(c.Commitment)
		if err != nil {
			return
		}
		// TODO:: Use bytes to check whether the checkpoint is consistent？
		if !bytes.Equal(baseStateByte, stateByte) {
			diffCheckpoint = append(diffCheckpoint, c)
		}
	}

	if len(diffCheckpoint) == len(Checkpoints)-1 {
		diffCheckpoint = []*checkpoint.Checkpoint{baseCheckpoint}
	}

	// TODO:: ....
}

func GetCheckpoint(config *types.Config) *checkpoint.Checkpoint {
	if config == nil || config.CommitteeIndexer == nil {
		return nil
	}
	switch true {
	case config.CommitteeIndexer.S3 != nil:
		return getCheckpointByS3(config.CommitteeIndexer.S3)
	case config.CommitteeIndexer.Da != nil:
		return getCheckpointByDa(config.CommitteeIndexer.Da)
	}
	return nil
}

func getCheckpointByDa(cfg []*types.SourceDa) *checkpoint.Checkpoint {

	return nil
}

func getCheckpointByS3(cfg []*types.SourceS3) *checkpoint.Checkpoint {

	return nil
}
