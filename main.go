package main

import (
	"strings"

	"github.com/RiemaLabs/indexer-committee/ord"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/apis"
	"github.com/RiemaLabs/indexer-light/checkpoint"
	"github.com/RiemaLabs/indexer-light/config"

	getter2 "github.com/RiemaLabs/indexer-light/getter"
)

func main() {
	go apis.ApiSrv()
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
			checkpoint.VerifyCheckpoint(getter, config.Config)
		}
	}
}
