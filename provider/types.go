package provider

import (
	"context"

	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
)

type CheckpointProvider interface {
	GetCheckpoint(ctx context.Context, height uint, hash string) (*configs.CheckpointExport, error)
}
