package provide

import (
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/types"
)

func NewDA(config *types.SourceDa) *ProviderDa {
	return &ProviderDa{
		Name:   constant.ProvideDaName,
		Config: config,
	}
}
func (p *ProviderDa) GetCheckpoint(height uint) *types.CheckPointObject {
	// TODO::
	return nil
}
