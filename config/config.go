package config

import (
	"encoding/json"
	"os"

	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/types"
)

var Config *types.Config

func init() {
	Config = &types.Config{}
	file, err := os.ReadFile(constant.ConfigFileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		return
	}
}

func GetCommitteeIndexerApi(config *types.Config) []string {
	var url []string
	if config == nil || config.CommitteeIndexer == nil {
		return url
	}
	if config.CommitteeIndexer.Da != nil {
		for i, _ := range config.CommitteeIndexer.Da {
			url = append(url, config.CommitteeIndexer.Da[i].ApiUrl)
		}
	}
	if config.CommitteeIndexer.S3 != nil {
		for i, _ := range config.CommitteeIndexer.S3 {
			url = append(url, config.CommitteeIndexer.S3[i].ApiUrl)
		}
	}
	return url
}
