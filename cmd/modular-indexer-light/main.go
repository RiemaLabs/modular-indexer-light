package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/internal/clients/ord/getter"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/provider"
	"github.com/RiemaLabs/modular-indexer-light/internal/runtime"
	"github.com/RiemaLabs/modular-indexer-light/internal/services"
)

var (
	version = "latest"
	gitHash = "unknown"
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
		Version: fmt.Sprintf("%v (%v)", version, gitHash),
	}
	cmd.Flags().StringVarP(&a.ConfigPath, "config", "c", "config.json", "path to config file")
	cmd.Flags().StringVar(&a.DenyListPath, "deny", "deny.jsonlines", "path to deny list file")
	cmd.Flags().StringVar(&a.PrivatePath, "private", "private", "path to private file")
	cmd.Flags().BoolVarP(&a.EnableTest, "test", "t", false, "Enable this flag to hijack the block height to test the service")
	cmd.Flags().BoolVarP(&a.EnableDAReport, "report", "", true, "Enable this flag to upload verified checkpoint to DA")
	return cmd
}

func (a *App) Run() {
	if err := configs.Init(a.ConfigPath, a.DenyListPath); err != nil {
		logs.Error.Fatalln("Config failed to initialize:", err)
	}

	a.initDaReport()
	logs.Info.Println("Syncing the latest state from committee indexers, please wait...")

	bitcoinGetter, err := getter.New(configs.C.Verification.BitcoinRPC)
	if err != nil {
		logs.Error.Fatalf("Failed to initialize Bitcoin Getter: %v", err)
	}
	currentBlockHeight, err := bitcoinGetter.GetLatestBlockHeight()
	if err != nil {
		logs.Error.Fatalf("Failed to get latest block height: %v", err)
	}
	lastBlockHeight := currentBlockHeight - 1
	lastBlockHash, err := bitcoinGetter.GetBlockHash(lastBlockHeight)
	if err != nil {
		logs.Error.Fatalf("Failed to get block hash: height=%d, err=%v", lastBlockHeight, err)
	}

	var providers []provider.CheckpointProvider
	for _, sourceS3 := range configs.C.CommitteeIndexers.S3 {
		prov := provider.NewProviderS3(&sourceS3, configs.C.Verification.MetaProtocol, 3)
		providers = append(providers, prov)
	}
	for _, sourceDA := range configs.C.CommitteeIndexers.DA {
		prov := provider.NewProviderDA(&sourceDA, configs.C.Verification.MetaProtocol, 3)
		providers = append(providers, prov)
	}
	actual := len(providers)
	expected := configs.C.Verification.MinimalCheckpoint
	if actual < expected {
		logs.Error.Fatalf("Insufficient checkpoint providers: actual=%d, expected=%d", actual, expected)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	checkpoints, err := provider.GetCheckpoints(ctx, providers, lastBlockHeight, lastBlockHash)
	if err != nil {
		logs.Error.Fatalf("Failed to get checkpoints: height=%d, hash=%s, err=%v", lastBlockHeight, lastBlockHash, err)
	}

	// TODO(MeihaoZuyu): Historical verification.
	if _, _, inconsistent := provider.CheckpointsInconsist(checkpoints); inconsistent {
		logs.Error.Fatalf("inconsistent checkpoints detected at height %q, historical verification is not supported but will be released soon :'(", lastBlockHeight)
	}
	logs.Info.Println("Latest state successfully synced!")

	df := runtime.NewRuntimeState(
		a.DenyListPath,
		providers,
		checkpoints[0],
		configs.C.Verification.MinimalCheckpoint,
		2*time.Minute,
	)

	go services.StartService(df, a.EnableTest, configs.C.ListenAddr)
	runSyncForever(a, df, bitcoinGetter)
}

func (a *App) initDaReport() {
	if !a.EnableDAReport {
		return
	}

	if err := configs.C.Report.LoadPrivate(a.PrivatePath); err != nil {
		logs.Error.Fatalln("Failed to read private key:", err)
	}

	if !checkpoint.IsValidNamespaceID(configs.C.Report.NamespaceID) {
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
		nid, err := checkpoint.CreateNamespace(
			configs.C.Report.PrivateKey,
			configs.C.Report.GasCoupon,
			namespaceName,
			configs.C.Report.Network,
		)
		if err != nil {
			logs.Error.Fatalf("failed to create namespace: %v", err)
		}
		configs.C.Report.NamespaceID = nid
		data, err := json.MarshalIndent(configs.C, "", "  ")
		if err != nil {
			logs.Error.Fatalln("marshal configurations error:", err)
		}
		err = os.WriteFile(a.ConfigPath, data, 0644)
		if err != nil {
			logs.Error.Fatalf("Failed to save namespace ID to configuration file: %v", err)
		}
		logs.Info.Printf("Namespace created successfully: %s", nid)
	}
}

func runSyncForever(app *App, df *runtime.State, bitcoinGetter *getter.Client) {
	for {
		time.Sleep(10 * time.Second)

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
				curCheckpoint := df.CurrentFirstCheckpoint().Checkpoint
				newCheckpoint := checkpoint.Checkpoint{
					Commitment:   curCheckpoint.Commitment,
					Hash:         curCheckpoint.Hash,
					Height:       curCheckpoint.Height,
					MetaProtocol: verifyCfg.MetaProtocol,
					Name:         reportCfg.Name,
					Version:      version,
				}

				time.Sleep(time.Duration(rand.Intn(40)+1) * time.Second)

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
		logs.Error.Fatalf("failed to execute: %v", err)
	}
}
