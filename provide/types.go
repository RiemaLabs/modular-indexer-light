package provide

import "github.com/RiemaLabs/indexer-light/types"

type ProviderS3 struct {
	Name   string          `json:"name"`
	Config *types.SourceS3 `json:"config"`
}

type ProviderDa struct {
	Name   string          `json:"name"`
	Config *types.SourceDa `json:"config"`
}
