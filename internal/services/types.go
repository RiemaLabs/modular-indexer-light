package services

import "github.com/RiemaLabs/modular-indexer-light/internal/constant"

type Status interface {
	Status() constant.ApiStatus
}

type Brc20VerifiableLightStateResponse struct {
	State string `json:"state"`
}
