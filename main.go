package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
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
	reportCfg := &cfg.Report
	verifyCfg := &cfg.Verification

	if arguments.EnableDAReport {
		reportCfg.PrivateKey = config.ReadPrivate()

		if !checkpoint.IsValidNamespaceID(reportCfg.NamespaceID) {
			log.Info("Invalid Namespace ID found in config.json. Initializing a new namespace.")
			scanner := bufio.NewScanner(os.Stdin)
			namespaceName := ""
			for {
				fmt.Print("Please enter your desired namespace name:")
				if scanner.Scan() {
					namespaceName = scanner.Text()
					if strings.TrimSpace(namespaceName) == "" {
						fmt.Print("Namespace name required!")
					} else {
						break
					}
				}
			}

			nid, err := checkpoint.CreateNamespace(reportCfg.PrivateKey, reportCfg.GasCoupon, namespaceName, reportCfg.Network)
			if err != nil {
				log.Panicf(fmt.Errorf("failed to create namespace: %v", err))
			}
			reportCfg.NamespaceID = nid
			bytes, err := json.Marshal(config.GlobalConfig)
			if err != nil {
				log.Panicf(fmt.Errorf("failed to save namespace ID to local file: %v", err))
			}
			err = os.WriteFile("./config.json", bytes, 0644)
			if err != nil {
				log.Panicf(fmt.Errorf("failed to save namespace ID to local file: %v", err))
			}
			log.Info(fmt.Sprintf("Namespace created successfully, Namespace ID: %s!", nid))
		}
	}

	log.Info("Syncing the latest state from committee indexers, please wait.")

	// Create Bitcoin getter.
	bitcoinGetter, err := getter.NewBitcoinOrdGetter(verifyCfg.BitcoinRPC)
	if err != nil {
		log.Panicf(fmt.Errorf("failed to initialize Bitcoin Getter. Error: %v", err))
	}

	currentBlockHeight, err := bitcoinGetter.GetLatestBlockHeight()
	if err != nil {
		log.Panicf(fmt.Errorf("failed to get latest block height. Error: %v", err))
	}

	lastBlockHeight := currentBlockHeight - 1
	lastBlockHash, err := bitcoinGetter.GetBlockHash(lastBlockHeight)
	if err != nil {
		log.Panicf(fmt.Errorf("failed to get block hash at height %d. Error: %v", lastBlockHeight, err))
	}

	// Create checkpoint providers
	providers := make([]provider.CheckpointProvider, 0)

	for _, sourceS3 := range cfg.CommitteeIndexers.S3 {
		prov := provider.NewProviderS3(&sourceS3, verifyCfg.MetaProtocol, retry)
		providers = append(providers, prov)
	}

	for _, sourceDA := range cfg.CommitteeIndexers.DA {
		prov := provider.NewProviderDA(&sourceDA, verifyCfg.MetaProtocol, retry)
		providers = append(providers, prov)
	}

	if len(providers) < verifyCfg.MinimalCheckpoint {
		log.Panicf(fmt.Errorf("the number of checkpoint providers is below the required minimum: %d", verifyCfg.MinimalCheckpoint))
	}

	// Get last checkpoint.
	// TODO: High. Historical verification.
	checkpoints := provider.GetCheckpoints(providers, lastBlockHeight, lastBlockHash, getCheckpointsTimeout)

	if len(checkpoints) == 0 {
		log.Panicf(fmt.Errorf("failed to get checkpoints at height %d", lastBlockHeight))
	}

	_, _, inconsistent := provider.CheckpointsInconsist(checkpoints)
	if inconsistent {
		log.Panicf(fmt.Errorf("inconsistent checkpoints detected at height %d during initialization."+
			"a version of the modular indexer with historical verification capabilities will be released soon", lastBlockHeight))
	}

	lastCheckpoint := checkpoints[0]

	// Create runtime state
	df := runtime.NewRuntimeState(providers, lastCheckpoint, verifyCfg.MinimalCheckpoint, getCheckpointsTimeout)

	log.Info("Succeed to sync the latest state!")

	syncCommitteeIndexers(arguments, df, bitcoinGetter)
}

func syncCommitteeIndexers(arguments *RuntimeArguments, df *runtime.RuntimeState, bitcoinGetter *getter.BitcoinOrdGetter) {
	cfg := config.GlobalConfig
	reportCfg := &cfg.Report
	verifyCfg := &cfg.Verification
	log.Info("API service available at: 8080")
	go apis.StartService(df, arguments.EnableTest)

	sleepInterval := time.Second * 10
	for {
		time.Sleep(sleepInterval)
		currentHeight, err := bitcoinGetter.GetLatestBlockHeight()
		if err != nil {
			log.Error("failed to GetLatestBlockHeight in syncCommitteeIndexers", "error", err)
			continue
		}
		hash, err := bitcoinGetter.GetBlockHash(currentHeight)
		if err != nil {
			log.Error("failed to get block hash in syncCommitteeIndexers", "error", err)
			continue
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
				log.Error("failed to UpdateCheckpoints in syncCommitteeIndexers", "error", err)
				continue
			}

			if arguments.EnableDAReport {
				// upload verified checkpoint to DA
				curCheckpoint := df.CurrentFirstCheckpoint().Checkpoint
				newCheckpoint := checkpoint.Checkpoint{
					Commitment:   curCheckpoint.Commitment,
					Hash:         curCheckpoint.Hash,
					Height:       curCheckpoint.Height,
					MetaProtocol: verifyCfg.MetaProtocol,
					Name:         reportCfg.Name,
					URL:          "",
					Version:      config.Version,
				}

				timeGap := time.Duration(rand.Intn(40)+1) * time.Second

				time.Sleep(timeGap)

				err := checkpoint.UploadCheckpointByDA(&newCheckpoint, reportCfg.PrivateKey, reportCfg.GasCoupon, reportCfg.NamespaceID, reportCfg.Network, time.Duration(reportCfg.Timeout)*time.Millisecond)
				if err != nil {
					log.Error(fmt.Sprintf("Unable to upload the checkpoint via DA: %v", err))
				} else {
					log.Info(fmt.Sprintf("Checkpoint successfully uploaded via DA at height: %s", newCheckpoint.Height))
				}
			}
		}

		log.Info(fmt.Sprintf("Listening for new Bitcoin block, current height: %d", df.CurrentHeight()))
	}
}

// TODO: Medium. Uniform the expression of Bitcoin block height and hash.
func main() {
	arguments := NewRuntimeArguments()
	rootCmd := arguments.MakeCmd()

	config.Version = "v0.1.0-rc.3"

	config.InitConfig()
	log.SetLevel(log.LevelError)
	if arguments.EnableTest {
		log.SetLevel(log.LevelDebug)
	}
	log.SetVerion(config.Version, time.Now().Format("20060102"))

	if err := rootCmd.Execute(); err != nil {
		log.Panicf(fmt.Errorf("failed to execute: %v", err))
	}
}
