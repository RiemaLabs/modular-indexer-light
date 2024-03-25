package provide

import (
	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/types"
)

func GetCommitteeIndexers(config *types.Config) []types.CheckPointProvider {
	var Provider []types.CheckPointProvider
	if config == nil || config.CommitteeIndexer == nil || (config.CommitteeIndexer.S3 == nil && config.CommitteeIndexer.Da == nil) {
		return Provider
	}
	if config.CommitteeIndexer.S3 != nil && len(config.CommitteeIndexer.S3) > 0 {
		for _, s3 := range config.CommitteeIndexer.S3 {
			Provider = append(Provider, &ProviderS3{
				Name:   constant.ProvideS3Name,
				Config: s3,
			})
		}
	}
	if config.CommitteeIndexer.Da != nil && len(config.CommitteeIndexer.Da) > 0 {
		for _, da := range config.CommitteeIndexer.Da {
			Provider = append(Provider, &ProviderDa{
				Name:   constant.ProvideDaName,
				Config: da,
			})
		}
	}
	return Provider
}
