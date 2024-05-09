//go:build js && wasm

package main

import (
	"syscall/js"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"

	"github.com/RiemaLabs/modular-indexer-light/internal/apps"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/services"
	"github.com/RiemaLabs/modular-indexer-light/internal/states"
	"github.com/RiemaLabs/modular-indexer-light/internal/utils"
)

var (
	version = "latest"
	gitHash = "unknown"
)

var (
	Error   = js.Global().Get("Error")
	Promise = js.Global().Get("Promise")
)

var app = apps.NewApp(version, gitHash)

func SetConfig(_ js.Value, args []js.Value) any {
	configs.C = new(configs.Config)

	verifyCfg := &configs.C.Verification
	verifyInput := args[0].Get("verification")

	verifyCfg.BitcoinRPC = verifyInput.Get("bitcoinRPC").String()
	verifyCfg.MinimalCheckpoint = verifyInput.Get("minimalCheckpoint").Int()
	verifyCfg.MetaProtocol = verifyInput.Get("metaProtocol").String()

	committeeCfg := &configs.C.CommitteeIndexers
	committeeInput := args[0].Get("committeeIndexers")

	sourceS3Input := committeeInput.Get("s3")
	for i := 0; i < sourceS3Input.Length(); i++ {
		committeeCfg.S3 = append(committeeCfg.S3, configs.SourceS3{
			Region: sourceS3Input.Index(i).Get("region").String(),
			Bucket: sourceS3Input.Index(i).Get("bucket").String(),
			Name:   sourceS3Input.Index(i).Get("name").String(),
		})
	}

	sourceDaInput := committeeInput.Get("da")
	for i := 0; i < sourceDaInput.Length(); i++ {
		committeeCfg.DA = append(committeeCfg.DA, configs.SourceDA{
			Network:     sourceDaInput.Index(i).Get("network").String(),
			NamespaceID: sourceDaInput.Index(i).Get("namespaceID").String(),
			Name:        sourceDaInput.Index(i).Get("name").String(),
		})
	}

	return nil
}

func Initialize(js.Value, []js.Value) any {
	go app.Run()
	return nil
}

func Warmup(js.Value, []js.Value) any {
	logs.Info.Println("Warming up...")
	go func() {
		for {
			if isVerifying() {
				logs.Info.Println("Still verifying, waiting...")
				time.Sleep(3 * time.Second)
				continue
			}
			_, _ = services.GetCurrentBalanceOfWallet(
				states.S.CurrentFirstCheckpoint().Checkpoint,
				"ordi",
				"bc1qhuv3dhpnm0wktasd3v0kt6e4aqfqsd0uhfdu7d",
			)
			logs.Info.Println("Warmed up successfully")
			return
		}
	}()
	return nil
}

func Status(js.Value, []js.Value) any {
	return Promise.New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve := args[0]
		var status states.Status
		if s := states.S; s != nil {
			status = states.Status(s.Status.Load())
		} else {
			status = states.StatusVerifying
		}
		resolve.Invoke(status.String())
		return nil
	}))
}

func isVerifying() bool {
	return states.S == nil || states.Status(states.S.Status.Load()) != states.StatusVerified
}

func GetBlockHeight(js.Value, []js.Value) any {
	return Promise.New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]
		if isVerifying() {
			reject.Invoke(Error.New("light indexer still verifying"))
		} else {
			resolve.Invoke(states.S.CurrentHeight())
		}
		return nil
	}))
}

func GetCurrentBalanceOfPkScript(_ js.Value, args []js.Value) any {
	tick := args[0].String()
	pkscript := args[1].String()
	return Promise.New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		if isVerifying() {
			reject.Invoke(Error.New("light Indexer still verifying"))
			return nil
		}

		go func() {
			balance, err := services.GetCurrentBalanceOfPkscript(states.S.CurrentFirstCheckpoint().Checkpoint, tick, pkscript)
			if err != nil {
				reject.Invoke(Error.New(err.Error()))
				return
			}
			resolve.Invoke(utils.Raw[utils.RawMap](balance))
		}()

		return nil
	}))
}

func GetCurrentBalanceOfWallet(_ js.Value, args []js.Value) any {
	tick := args[0].String()
	wallet := args[1].String()
	return Promise.New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		if isVerifying() {
			reject.Invoke(Error.New("light Indexer still verifying"))
			return nil
		}

		go func() {
			balance, err := services.GetCurrentBalanceOfWallet(states.S.CurrentFirstCheckpoint().Checkpoint, tick, wallet)
			if err != nil {
				reject.Invoke(Error.New(err.Error()))
				return
			}
			resolve.Invoke(utils.Raw[utils.RawMap](balance))
		}()

		return nil
	}))
}

func GetCurrentCheckpoints(js.Value, []js.Value) any {
	return Promise.New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]
		if isVerifying() {
			reject.Invoke(Error.New("light Indexer still verifying"))
		} else {
			var ret []*checkpoint.Checkpoint
			for _, export := range states.S.CurrentCheckpoints() {
				ret = append(ret, export.Checkpoint)
			}
			resolve.Invoke(utils.Raw[utils.RawSlice](ret))
		}
		return nil
	}))
}

func GetLastCheckpoint(js.Value, []js.Value) any {
	return Promise.New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]
		if isVerifying() {
			reject.Invoke(Error.New("light Indexer still verifying"))
			return nil
		}
		if last := states.S.LastCheckpoint(); last != nil {
			resolve.Invoke(utils.Raw[utils.RawMap](last.Checkpoint))
		} else {
			reject.Invoke(Error.New("empty last checkpoint"))
		}
		return nil
	}))
}

func main() {
	js.Global().Set("lightSetConfig", js.FuncOf(SetConfig))
	js.Global().Set("lightInitialize", js.FuncOf(Initialize))
	js.Global().Set("lightWarmup", js.FuncOf(Warmup))
	js.Global().Set("lightStatus", js.FuncOf(Status))
	js.Global().Set("lightGetBlockHeight", js.FuncOf(GetBlockHeight))
	js.Global().Set("lightGetBalanceOfPkScript", js.FuncOf(GetCurrentBalanceOfPkScript))
	js.Global().Set("lightGetBalanceOfWallet", js.FuncOf(GetCurrentBalanceOfWallet))
	js.Global().Set("lightGetCurrentCheckpoints", js.FuncOf(GetCurrentCheckpoints))
	js.Global().Set("lightGetLastCheckpoint", js.FuncOf(GetLastCheckpoint))
	select {}
}
