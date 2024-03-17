package main

import (
	"strings"

	"github.com/RiemaLabs/indexer-committee/ord"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/apis"
	"github.com/RiemaLabs/indexer-light/config"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/verify"

	getter2 "github.com/RiemaLabs/indexer-light/getter"
)

func main() {
	go apis.Start()
	rpcGetter, err2 := getter2.NewGetter(config.Config)
	if err2 != nil {
		return
	}
	queue, err := fetchHeight(rpcGetter)
	if err != nil {
		return
	}
	queue.Println()
}

// TODO:: Not completed
func fetchHeight(getter getter.OrdGetter) (*ord.StateQueue, error) {
	for {
		latestHeight, err := getter.GetLatestBlockHeight()
		if err != nil {
			return nil, err
		}
		hash, err := getter.GetBlockHash(latestHeight)
		if err != nil {
			return nil, err
		}
		if latestHeight > uint(config.Config.StartHeight) && !strings.EqualFold(hash, config.Config.StartBlockHash) {
			constant.ApiState = constant.ApiStateInit
			err := verify.VerifyCheckpoint(getter, config.Config)
			if err != nil {
				return nil, err
			}
			constant.ApiState = constant.ApiStateActive
		}
	}
}
