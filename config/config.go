package config

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/types"
)

//go:embed config.json
var configBody []byte

var Config *types.Config

func init() {
	Config = &types.Config{}
	err := json.Unmarshal(configBody, &Config)
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

func UpdateConfig(config *types.Config) error {
	file, err := os.OpenFile(constant.ConfigFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = file.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}
