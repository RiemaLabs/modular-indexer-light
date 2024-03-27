package apis

import (
	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/RiemaLabs/modular-indexer-light/constant"
)

type Status interface {
	Status() constant.ApiStatus
}

type Brc20VerifiableLightCheckpointsResponse struct {
	Checkpoints []*config.CheckpointExport `json:"checkpoints"`
}

type Brc20VerifiableLightLastCheckpointResponse struct {
	Checkpoint *config.CheckpointExport `json:"checkpoint"`
}

type Brc20VerifiableLightStateResponse struct {
	State string `json:"state"`
}
