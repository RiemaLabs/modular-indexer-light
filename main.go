package main

import (
	"strings"
	"time"

	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/apis"
	"github.com/RiemaLabs/indexer-light/cmd"
	"github.com/RiemaLabs/indexer-light/config"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/log"
	"github.com/RiemaLabs/indexer-light/verify"

	getter2 "github.com/RiemaLabs/indexer-light/getter"
)

func main() {
	log.SetLevel(log.LevelVerbose)
	log.SetVerion("v0.0.1", time.Now().Format("20060102"))
	rootCmd := cmd.NewCmd()
	if err := rootCmd.Init().Execute(); err != nil {
		log.Error("main", "Failed to parse the arguments", err)
	}

	go apis.Start()
	rpcGetter, err := getter2.NewGetter(config.Config)
	if err != nil {
		return
	}
	fetchHeight(rpcGetter)
}

// TODO:: Not completed
func fetchHeight(getter getter.OrdGetter) {
	for {
		latestHeight, err := getter.GetLatestBlockHeight()
		if err != nil {
			log.Error("fetchHeight", "GetLatestBlockHeight", err)
			return
		}
		log.Debug("fetchHeight", "latestHeight", latestHeight, "Config.StartHeight", config.Config.StartHeight)
		hash, err := getter.GetBlockHash(latestHeight)
		if err != nil {
			log.Error("fetchHeight", "GetBlockHash", err)
			return
		}
		log.Debug("fetchHeight", "latestHash", hash, "Config.StartBlockHash", config.Config.StartBlockHash)
		if latestHeight > uint(config.Config.StartHeight) && !strings.EqualFold(hash, config.Config.StartBlockHash) {
			log.Debug("fetchHeight", "msg", "sync...")
			constant.ApiState = constant.ApiStateInit
			err = verify.VerifyCheckpoint(getter, config.Config)
			if err != nil {
				log.Error("fetchHeight", "VerifyCheckpoint", err)
				continue
			}
			constant.ApiState = constant.ApiStateActive
		}
		time.Sleep(10 * time.Minute)
	}
}
