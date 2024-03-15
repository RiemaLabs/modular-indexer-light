package main

import (
	"strings"

	"github.com/RiemaLabs/indexer-committee/ord"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/apis"
	"github.com/RiemaLabs/indexer-light/checkpoint"
	"github.com/RiemaLabs/indexer-light/config"
)

func main() {
	go apis.ApiSrv()
	queue, err := fetchHeight(nil)
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
			checkpoint.ExamineCheckpoint(getter, config.Config)
		}
	}
}
