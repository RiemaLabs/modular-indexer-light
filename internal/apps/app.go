package apps

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

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/spf13/cobra"

	"github.com/RiemaLabs/modular-indexer-light/internal/checkpoints"
	"github.com/RiemaLabs/modular-indexer-light/internal/clients/btcutl"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/services"
	"github.com/RiemaLabs/modular-indexer-light/internal/states"
)

type App struct {
	version, gitHash string

	ConfigPath, DenyListPath, PrivatePath string
	EnableTest, EnableDAReport            bool
}

func NewApp(version, gitHash string) *App {
	return &App{version: version, gitHash: gitHash}
}

func (a *App) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modular-indexer-light",
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

			if err := configs.Init(a.ConfigPath, a.DenyListPath); err != nil {
				logs.Error.Fatalln("Config failed to initialize:", err)
			}

			a.Run()
		},
		Version: fmt.Sprintf("%v (%v)", a.version, a.gitHash),
	}
	cmd.Flags().StringVarP(&a.ConfigPath, "config", "c", "config.json", "path to config file")
	cmd.Flags().StringVar(&a.DenyListPath, "deny", "deny.jsonlines", "path to deny list file")
	cmd.Flags().StringVar(&a.PrivatePath, "private", "private", "path to private file")
	cmd.Flags().BoolVarP(&a.EnableTest, "test", "t", false, "Enable this flag to hijack the block height to test the service")
	cmd.Flags().BoolVarP(&a.EnableDAReport, "report", "", true, "Enable this flag to upload verified checkpoint to DA")
	return cmd
}

func (a *App) Run() {
	a.initDaReport()
	btcutl.Init(configs.C.Verification.BitcoinRPC)

	logs.Info.Println("Syncing the latest state from committee indexers, please wait...")

	currentBlockHeight, err := btcutl.BTC.GetLatestBlockHeight(context.Background())
	if err != nil {
		logs.Error.Fatalf("Failed to get latest block height: %v", err)
	}
	lastBlockHeight := currentBlockHeight - 1
	lastBlockHash, err := btcutl.BTC.GetBlockHash(context.Background(), lastBlockHeight)
	if err != nil {
		logs.Error.Fatalf("Failed to get last block hash: height=%d, err=%v", lastBlockHeight, err)
	}

	var providers []checkpoints.CheckpointProvider
	if raw := configs.C.CommitteeIndexers.Raw; a.EnableTest && len(raw) > 0 {
		for _, sourceRaw := range raw {
			providers = append(providers, sourceRaw)
		}
	} else {
		for _, sourceS3 := range configs.C.CommitteeIndexers.S3 {
			providers = append(providers, checkpoints.NewProviderS3(&sourceS3, configs.C.Verification.MetaProtocol))
		}
		for _, sourceDA := range configs.C.CommitteeIndexers.DA {
			providers = append(providers, checkpoints.NewProviderDA(&sourceDA, configs.C.Verification.MetaProtocol))
		}
	}
	actual := len(providers)
	expected := configs.C.Verification.MinimalCheckpoint
	if actual < expected {
		logs.Error.Fatalf("Insufficient checkpoint providers: actual=%d, expected=%d", actual, expected)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cps, err := checkpoints.GetCheckpoints(ctx, providers, lastBlockHeight, lastBlockHash)
	if err != nil {
		logs.Error.Fatalf("Failed to get checkpoints: height=%d, hash=%s, err=%v", lastBlockHeight, lastBlockHash, err)
	}

	// TODO: Historical verification.
	if inconsistent := checkpoints.Inconsistent(cps); inconsistent {
		logs.Error.Fatalf("inconsistent checkpoints detected at height %q, historical verification is not supported but will be released soon :'(", lastBlockHeight)
	}
	logs.Info.Println("Latest state successfully synced!")

	states.Init(
		a.DenyListPath,
		providers,
		cps[0],
		configs.C.Verification.MinimalCheckpoint,
		2*time.Minute,
	)

	go services.StartService(a.EnableTest, configs.C.ListenAddr)
	a.runSyncForever()
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

func (a *App) runSyncForever() {
	for {
		time.Sleep(10 * time.Second)
		logs.Info.Println("Syncing latest state...")

		currentHeight, err := btcutl.BTC.GetLatestBlockHeight(context.Background())
		if err != nil {
			logs.Error.Printf("Failed to get latest block height: %v", err)
			continue
		}
		currentHash, err := btcutl.BTC.GetBlockHash(context.Background(), currentHeight)
		if err != nil {
			logs.Error.Printf("Failed to get block hash: %v", err)
			continue
		}

		if first := states.S.CurrentFirstCheckpoint(); first == nil ||
			first.Checkpoint.Height != strconv.Itoa(int(currentHeight)) ||
			first.Checkpoint.Hash != currentHash {
			// Checkpoints are not the latest, start syncing.

			if err := states.S.UpdateCheckpoints(currentHeight, currentHash); err != nil {
				logs.Error.Printf("Failed to update checkpoints: %v", err)
				continue
			}

			if a.EnableDAReport {
				cp := states.S.CurrentFirstCheckpoint().Checkpoint
				newCp := checkpoint.Checkpoint{
					Commitment:   cp.Commitment,
					Hash:         cp.Hash,
					Height:       cp.Height,
					MetaProtocol: configs.C.Verification.MetaProtocol,
					Name:         configs.C.Report.Name,
					Version:      a.version,
				}

				time.Sleep(time.Duration(rand.Intn(40)+1) * time.Second)

				if err := checkpoint.UploadCheckpointByDA(
					&newCp,
					configs.C.Report.PrivateKey,
					configs.C.Report.GasCoupon,
					configs.C.Report.NamespaceID,
					configs.C.Report.Network,
					configs.C.Report.Timeout.Duration,
				); err != nil {
					logs.Error.Printf("Unable to upload the checkpoint via DA: %v", err)
				} else {
					logs.Info.Printf("Checkpoint successfully uploaded via DA at height: %s", newCp.Height)
				}
			}
		}

		logs.Info.Printf("Listening for new Bitcoin block: height=%d", states.S.CurrentHeight())
	}
}
