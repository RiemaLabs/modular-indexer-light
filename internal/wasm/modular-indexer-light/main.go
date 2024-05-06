//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/RiemaLabs/modular-indexer-light/internal/apps"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/states"
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

func SetConfig(_ js.Value, args []js.Value) interface{} {
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

func Initialize(js.Value, []js.Value) interface{} {
	go app.Run()
	return nil
}

func GetBlockHeight(js.Value, []js.Value) interface{} {
	return Promise.New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		if states.S == nil || states.Status(states.S.State.Load()) != states.StatusActive {
			reject.Invoke(Error.New("light indexer not active yet"))
		} else {
			resolve.Invoke(states.S.CurrentHeight())
		}
		return nil
	}))
}

func main() {
	js.Global().Set("lightSetConfig", js.FuncOf(SetConfig))
	js.Global().Set("lightInitialize", js.FuncOf(Initialize))
	js.Global().Set("lightGetBlockHeight", js.FuncOf(GetBlockHeight))
	select {}
}
