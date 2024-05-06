package main

import (
	"github.com/RiemaLabs/modular-indexer-light/internal/apps"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
)

var (
	version = "latest"
	gitHash = "unknown"
)

// TODO: Medium. Uniform the expression of Bitcoin block height and hash.
func main() {
	if err := apps.NewApp(version, gitHash).Command().Execute(); err != nil {
		logs.Error.Fatalf("failed to execute: %v", err)
	}
}
