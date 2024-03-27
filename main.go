package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RiemaLabs/modular-indexer-light/apis"
	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/RiemaLabs/modular-indexer-light/log"
	"github.com/RiemaLabs/modular-indexer-light/ord/getter"
	"github.com/RiemaLabs/modular-indexer-light/provider"
	"github.com/RiemaLabs/modular-indexer-light/runtime"
)

func Execution(arguments *RuntimeArguments) {
	getCheckpointsTimeout := time.Minute * 2
	retry := 3
	cfg := config.GlobalConfig

	// Create Bitcoin getter.
	bitcoinGetter, err := getter.NewBitcoinOrdGetter(cfg.BitcoinRPC)
	if err != nil {
		log.Panicf(fmt.Errorf("failed to initiate Bitcoin Getter, error: %v", err))
	}

	currentBlockHeight, err := bitcoinGetter.GetLatestBlockHeight()
	if err != nil {
		log.Panicf(fmt.Errorf("failed to GetLatestBlockHeight, error: %v", err))
	}

	lastBlockHeight := currentBlockHeight - 1
	lastBlockHash, err := bitcoinGetter.GetBlockHash(lastBlockHeight)
	if err != nil {
		log.Panicf(fmt.Errorf("failed to GetBlockHash at height %d, error: %v", lastBlockHeight, err))
	}

	// Create checkpoint providers
	providers := make([]provider.CheckpointProvider, 0)

	for _, sourceS3 := range cfg.CommitteeIndexers.S3 {
		prov := provider.NewProviderS3(&sourceS3, cfg.MetaProtocol, retry)
		providers = append(providers, prov)
	}

	for _, sourceDA := range cfg.CommitteeIndexers.DA {
		prov := provider.NewProviderDA(&sourceDA, cfg.MetaProtocol, retry)
		providers = append(providers, prov)
	}

	if len(providers) < cfg.MinimalCheckpoint {
		log.Panicf(fmt.Errorf("the number of checkpoint providers is below the minimum required amount: %d", cfg.MinimalCheckpoint))
	}

	// Get last checkpoint.
	// TODO: High. Historical verification.
	checkpoints := provider.GetCheckpoints(providers, lastBlockHeight, lastBlockHash, getCheckpointsTimeout)

	if len(checkpoints) == 0 {
		log.Panicf(fmt.Errorf("failed to GetCheckpoints at height %d", lastBlockHeight))
	}

	_, _, inconsistent := provider.CheckpointsInconsist(checkpoints)
	if inconsistent {
		log.Panicf(fmt.Errorf("inconsistent checkpoints detected at height %d during initialization. "+
			"a version of the modular indexer with historical verification capabilities will be released soon", lastBlockHeight))
	}

	lastCheckpoint := checkpoints[0]

	// Create runtime state
	df := runtime.NewRuntimeState(providers, lastCheckpoint, cfg.MinimalCheckpoint, getCheckpointsTimeout)

	syncCommitteeIndexers(arguments, df, bitcoinGetter)
}

func syncCommitteeIndexers(arguments *RuntimeArguments, df *runtime.RuntimeState, bitcoinGetter *getter.BitcoinOrdGetter) {
	log.Info("Providing API service at: 8080")
	go apis.StartService(df, arguments.EnableTest)

	sleepInterval := time.Minute * 1
	for {
		currentHeight, err := bitcoinGetter.GetLatestBlockHeight()
		if err != nil {
			log.Panicf(fmt.Errorf("failed to GetLatestBlockHeight in syncCommitteeIndexers, error: %v", err))
		}
		hash, err := bitcoinGetter.GetBlockHash(currentHeight)
		if err != nil {
			log.Panicf(fmt.Errorf("failed to GetBlockHash in syncCommitteeIndexers, error: %v", err))
		}

		var unsynced = false

		firstCheckpoint := df.CurrentFirstCheckpoint()
		if firstCheckpoint == nil {
			unsynced = true
		} else if strconv.Itoa(int(currentHeight)) != firstCheckpoint.Checkpoint.Height || hash != firstCheckpoint.Checkpoint.Hash {
			unsynced = true
		}

		if unsynced {
			err = df.UpdateCheckpoints(currentHeight, hash)
			if err != nil {
				log.Panicf(fmt.Errorf("failed to UpdateCheckpoints in syncCommitteeIndexers, error: %v", err))
			}
		}

		log.Info(fmt.Sprintf("Listening for new Bitcoin block, current height: %d", df.CurrentHeight()))
		time.Sleep(sleepInterval)
	}
}

// TODO: Medium. Uniform the expression of Bitcoin block height and hash.
func main() {
	arguments := NewRuntimeArguments()
	rootCmd := arguments.MakeCmd()

	config.InitConfig()
	log.SetLevel(log.LevelError)
	if arguments.EnableTest {
		log.SetLevel(log.LevelDebug)
	}
	log.SetVerion(config.GlobalConfig.Version, time.Now().Format("20060102"))

	if err := rootCmd.Execute(); err != nil {
		log.Panicf(fmt.Errorf("failed to execute: %v", err))
	}
}
