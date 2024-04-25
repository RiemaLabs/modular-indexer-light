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

	"github.com/spf13/cobra"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/apis"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/ord/getter"
	"github.com/RiemaLabs/modular-indexer-light/provider"
	"github.com/RiemaLabs/modular-indexer-light/runtime"
)

var (
	version = "latest"
	gitHash = "unknown"
)

const (
	DefaultConfigFile   = "config.json"
	DefaultDenyListFile = "blacklist.jsonlines"
	DefaultPrivateFile  = "private"
)

type App struct {
	ConfigPath, DenyListPath, PrivatePath string
	EnableTest, EnableDAReport            bool
}

func (a *App) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Nubit Light Indexer",
		Short: "Activates the Nubit Light Indexer with optional services.",
		Long: `Light Indexer is an essential component of the Nubit Modular Indexer architecture.
It enables typical users to verify Bitcoin meta-protocols without requiring substantial computing resources.
This command offers multiple flags to tailor the indexer's functionality according to the user's needs.`,
		Run: func(cmd *cobra.Command, args []string) {
			if a.EnableTest {
				logs.Info.Println("Test mode enabled")
			}

			if a.EnableDAReport {
				logs.Info.Println("DA report enabled")
			} else {
				logs.Info.Println("DA report disabled")
			}

			a.Run()
		},
		Version: fmt.Sprintf("modular-indexer-light %v (%v)", version, gitHash),
	}
	cmd.Flags().StringVarP(&a.ConfigPath, "config", "c", DefaultConfigFile, "path to config file")
	cmd.Flags().StringVar(&a.DenyListPath, "deny", DefaultDenyListFile, "path to deny list file")
	cmd.Flags().StringVar(&a.PrivatePath, "private", DefaultPrivateFile, "path to private file")
	cmd.Flags().BoolVarP(&a.EnableTest, "test", "t", false, "Enable this flag to hijack the block height to test the service")
	cmd.Flags().BoolVarP(&a.EnableDAReport, "report", "", true, "Enable this flag to upload verified checkpoint to DA")
	return cmd
}

func (a *App) Run() {
	if err := configs.Init(a.ConfigPath, a.DenyListPath); err != nil {
		logs.Error.Fatalln("Config failed to initialize:", err)
	}

	if a.EnableTest {
		logs.SetDebug()
	}

	getCheckpointsTimeout := time.Minute * 2
	retry := 3
	cfg := configs.C
	reportCfg := &cfg.Report
	verifyCfg := &cfg.Verification

	if a.EnableDAReport {
		if err := configs.C.Report.LoadPrivate(a.PrivatePath); err != nil {
			logs.Error.Fatalln("Failed to read private key:", err)
		}

		if !checkpoint.IsValidNamespaceID(reportCfg.NamespaceID) {
			logs.Info.Println("Invalid Namespace ID found in configurations. Initializing a new namespace.")
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
				logs.Error.Panicf("failed to create namespace: %v", err)
			}
			reportCfg.NamespaceID = nid
			bytes, err := json.MarshalIndent(configs.C, "", "  ")
			if err != nil {
				logs.Error.Panicf("failed to save namespace ID to local file: %v", err)
			}
			err = os.WriteFile("./config.json", bytes, 0644)
			if err != nil {
				logs.Error.Panicf("failed to save namespace ID to local file: %v", err)
			}
			logs.Info.Println("Namespace created successfully, Namespace ID: %s!", nid)
		}
	}

	logs.Info.Printf("Syncing the latest state from committee indexers, please wait.")

	// Create Bitcoin getter.
	bitcoinGetter, err := getter.NewBitcoinOrdGetter(verifyCfg.BitcoinRPC)
	if err != nil {
		logs.Error.Panicf("failed to initialize Bitcoin Getter. Error: %v", err)
	}

	currentBlockHeight, err := bitcoinGetter.GetLatestBlockHeight()
	if err != nil {
		logs.Error.Panicf("failed to get latest block height. Error: %v", err)
	}

	lastBlockHeight := currentBlockHeight - 1
	lastBlockHash, err := bitcoinGetter.GetBlockHash(lastBlockHeight)
	if err != nil {
		logs.Error.Panicf("failed to get block hash at height %d. Error: %v", lastBlockHeight, err)
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
		logs.Error.Panicf("the number of checkpoint providers is below the required minimum: %d", verifyCfg.MinimalCheckpoint)
	}

	// Get last checkpoint.
	// TODO: High. Historical verification.
	checkpoints := provider.GetCheckpoints(providers, lastBlockHeight, lastBlockHash, getCheckpointsTimeout)

	if len(checkpoints) == 0 {
		logs.Error.Panicf("failed to get checkpoints at height %d", lastBlockHeight)
	}

	_, _, inconsistent := provider.CheckpointsInconsist(checkpoints)
	if inconsistent {
		logs.Error.Panicf("inconsistent checkpoints detected at height %d during initialization."+
			"a version of the modular indexer with historical verification capabilities will be released soon", lastBlockHeight)
	}

	lastCheckpoint := checkpoints[0]

	// Create runtime state
	df := runtime.NewRuntimeState(a.DenyListPath, providers, lastCheckpoint, verifyCfg.MinimalCheckpoint, getCheckpointsTimeout)

	logs.Info.Printf("Succeed to sync the latest state!")

	syncCommitteeIndexers(a, df, bitcoinGetter)
}

func syncCommitteeIndexers(app *App, df *runtime.State, bitcoinGetter *getter.BitcoinOrdGetter) {
	cfg := configs.C
	reportCfg := &cfg.Report
	verifyCfg := &cfg.Verification

	go apis.StartService(df, app.EnableTest, cfg.ListenAddr)

	sleepInterval := time.Second * 10
	for {
		time.Sleep(sleepInterval)
		currentHeight, err := bitcoinGetter.GetLatestBlockHeight()
		if err != nil {
			logs.Error.Printf("failed to GetLatestBlockHeight in syncCommitteeIndexers: %v", err)
			continue
		}
		hash, err := bitcoinGetter.GetBlockHash(currentHeight)
		if err != nil {
			logs.Error.Printf("failed to get block hash in syncCommitteeIndexers: %v", err)
			continue
		}

		notSynced := false

		firstCheckpoint := df.CurrentFirstCheckpoint()
		if firstCheckpoint == nil {
			notSynced = true
		} else if strconv.Itoa(int(currentHeight)) != firstCheckpoint.Checkpoint.Height || hash != firstCheckpoint.Checkpoint.Hash {
			notSynced = true
		}

		if notSynced {

			err = df.UpdateCheckpoints(currentHeight, hash)
			if err != nil {
				logs.Error.Printf("failed to UpdateCheckpoints in syncCommitteeIndexers: %v", err)
				continue
			}

			if app.EnableDAReport {
				// upload verified checkpoint to DA
				curCheckpoint := df.CurrentFirstCheckpoint().Checkpoint
				newCheckpoint := checkpoint.Checkpoint{
					Commitment:   curCheckpoint.Commitment,
					Hash:         curCheckpoint.Hash,
					Height:       curCheckpoint.Height,
					MetaProtocol: verifyCfg.MetaProtocol,
					Name:         reportCfg.Name,
					Version:      version,
				}

				timeGap := time.Duration(rand.Intn(40)+1) * time.Second

				time.Sleep(timeGap)

				err := checkpoint.UploadCheckpointByDA(&newCheckpoint, reportCfg.PrivateKey, reportCfg.GasCoupon, reportCfg.NamespaceID, reportCfg.Network, time.Duration(reportCfg.Timeout)*time.Millisecond)
				if err != nil {
					logs.Error.Printf(fmt.Sprintf("Unable to upload the checkpoint via DA: %v", err))
				} else {
					logs.Info.Printf(fmt.Sprintf("Checkpoint successfully uploaded via DA at height: %s", newCheckpoint.Height))
				}
			}
		}

		logs.Info.Printf(fmt.Sprintf("Listening for new Bitcoin block, current height: %d", df.CurrentHeight()))
	}
}

// TODO: Medium. Uniform the expression of Bitcoin block height and hash.
func main() {
	if err := new(App).Command().Execute(); err != nil {
		logs.Error.Panicf("failed to execute: %v", err)
	}
}
