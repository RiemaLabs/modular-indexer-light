package provide

import (
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/types"
)

func NewS3(config *types.SourceS3) *ProviderS3 {
	return &ProviderS3{
		Name:   constant.ProvideS3Name,
		Config: config,
	}
}

func (p *ProviderS3) GetCheckpoint(height uint) *types.CheckPointObject {
	// TODO::
	return nil
}
