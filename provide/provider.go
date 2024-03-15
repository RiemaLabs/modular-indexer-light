package provide

import (
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/types"
)

func GetProviders(config *types.Config) []types.CheckPointProvider {
	var Provider []types.CheckPointProvider
	if config == nil || config.CommitteeIndexer == nil || (config.CommitteeIndexer.S3 == nil && config.CommitteeIndexer.Da == nil) {
		return Provider
	}
	if config.CommitteeIndexer.S3 != nil {
		for _, s3 := range config.CommitteeIndexer.S3 {
			Provider = append(Provider, &ProviderS3{
				Name:   constant.ProvideS3Name,
				Config: s3,
			})
		}
	}
	if config.CommitteeIndexer.Da != nil {
		for _, da := range config.CommitteeIndexer.Da {
			Provider = append(Provider, &ProviderDa{
				Name:   constant.ProvideDaName,
				Config: da,
			})
		}
	}
	return Provider
}
