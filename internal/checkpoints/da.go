package checkpoints

import (
	"context"
	"errors"

	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
)

type DA struct{}

func NewProviderDA(sourceDA *configs.SourceDA, metaProtocol string) *DA {
	_ = sourceDA
	_ = metaProtocol
	return new(DA)
}

func (p *DA) Get(context.Context, uint, string) (*configs.CheckpointExport, error) {
	return nil, errors.New("DA closed for revamps, please stay tuned")
}
