package apis

import (
	"github.com/RiemaLabs/modular-indexer-light/constant"
)

type Status interface {
	Status() constant.ApiStatus
}

type Brc20VerifiableLightStateResponse struct {
	State string `json:"state"`
}
