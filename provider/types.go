package provider

import (
	"context"

	"github.com/RiemaLabs/modular-indexer-light/config"
)

type CheckpointProvider interface {
	GetCheckpoint(ctx context.Context, height uint, hash string) (*config.CheckpointExport, error)
}
