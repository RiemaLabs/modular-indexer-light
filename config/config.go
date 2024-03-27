package config

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RiemaLabs/modular-indexer-light/log"
)

var BlacklistFile = "./blacklist.jsonlines"

//go:embed config.json
var configBody []byte

var GlobalConfig *Config

func InitConfig() {
	GlobalConfig = &Config{}
	err := json.Unmarshal(configBody, &GlobalConfig)
	if err != nil {
		return
	}

	blacks := LoadBlacklist()

	for _, b := range blacks {
		if b.SourceDA != nil {
			for i, source := range GlobalConfig.CommitteeIndexers.DA {
				if source.Name == b.SourceDA.Name && source.NamespaceID == b.SourceDA.NamespaceID && source.Network == b.SourceDA.Network {
					GlobalConfig.CommitteeIndexers.DA[0], GlobalConfig.CommitteeIndexers.DA[i] = GlobalConfig.CommitteeIndexers.DA[i], GlobalConfig.CommitteeIndexers.DA[0]
					GlobalConfig.CommitteeIndexers.DA = GlobalConfig.CommitteeIndexers.DA[1:]
					break
				}
			}
		}

		if b.SourceS3 != nil {
			for i, source := range GlobalConfig.CommitteeIndexers.S3 {
				if source.Name == b.SourceS3.Name && source.Region == b.SourceS3.Region && source.Bucket == b.SourceS3.Bucket {
					GlobalConfig.CommitteeIndexers.S3[0], GlobalConfig.CommitteeIndexers.S3[i] = GlobalConfig.CommitteeIndexers.S3[i], GlobalConfig.CommitteeIndexers.S3[0]
					GlobalConfig.CommitteeIndexers.S3 = GlobalConfig.CommitteeIndexers.S3[1:]
					break
				}
			}
		}

	}
}

func AppendBlacklist(in *Blacklist) {
	f, err := os.OpenFile(BlacklistFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Error(fmt.Sprintf("failed to open the blacklist file, error: %v", err))
	}
	defer f.Close()

	data, err := json.Marshal(in)
	if err != nil {
		log.Error(fmt.Sprintf("failed to open the blacklist file, error: %v", err))
	}
	f.WriteString(string(data) + "\n")
}

func LoadBlacklist() []*Blacklist {

	f, err := os.Open(BlacklistFile)
	if err != nil {
		return []*Blacklist{}
	}
	defer f.Close()

	r := bufio.NewReader(f)
	body := []*Blacklist{}
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF || err != nil {
			break
		}
		line = strings.TrimRight(line, "\n")
		tmp := Blacklist{}
		if err := json.Unmarshal([]byte(line), &tmp); err != nil {
			return body
		}
		body = append(body, &tmp)
	}
	return body
}
