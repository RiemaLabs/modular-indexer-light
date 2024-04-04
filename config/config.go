package config

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/log"
	"github.com/RiemaLabs/modular-indexer-light/utils"
	"github.com/RiemaLabs/modular-indexer-light/wallet"
)

var BlacklistFile = "./blacklist.jsonlines"
var ConfigFile = "./config.json"
var PrivateFile = "./private"

var GlobalConfig *Config
var Version string

func InitConfig() {
	configFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		log.Panicf(fmt.Errorf("failed to read config file: %v", err))
	}

	err = json.Unmarshal(configFile, &GlobalConfig)
	if err != nil {
		log.Panicf(fmt.Errorf("failed to parse config file: %v", err))
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

func ReadPrivate() string {
	_, err := os.Stat(PrivateFile)
	if err == nil {
		// read private key
		data, err := os.ReadFile(PrivateFile)
		if err == nil {
			log.Info("Read private key from", "file", PrivateFile)
			return string(data)
		}

	}

	log.Info("Failed to read the private key from the local directory. Generate a new one")
	var pwd = constant.DefaultPassword

	wall := wallet.NewWallet(&pwd)
	if !wall.GenerateBip39Seed(&pwd, &pwd) {
		log.Panicf(errors.New("failed to generate seeds"))
	}
	account := wall.GenerateAccount(&pwd)

	pri := utils.EcdsaToPrivateStr(account.PrivateKey(&pwd))

	file, err := os.OpenFile(PrivateFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panicf(err)
	}
	defer file.Close()

	_, err = file.WriteString(pri)
	if err != nil {
		log.Panicf(errors.New("failed to write private key"))
	}
	log.Info("Store your private file carefully and don't share it!")
	return pri
}
