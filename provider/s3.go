package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/RiemaLabs/modular-indexer-light/log"
)

type ProviderS3 struct {
	Config       *config.SourceS3 `json:"config"`
	MetaProtocol string
	Retry        int
}

func NewProviderS3(sourceS3 *config.SourceS3, metaProtocol string, retry int) *ProviderS3 {
	return &ProviderS3{
		Config:       sourceS3,
		MetaProtocol: metaProtocol,
		Retry:        retry,
	}
}

func (p *ProviderS3) GetCheckpoint(ctx context.Context, height uint, hash string) (*config.CheckpointExport, error) {
	var ck *checkpoint.Checkpoint
	var err error
OuterLoop:
	for i := 0; i < int(p.Retry); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ck, err = DownloadCheckpointByS3(p.Config.Region, p.Config.Bucket, p.Config.Name, p.MetaProtocol, strconv.Itoa(int(height)), hash, 100*time.Second)
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			break OuterLoop
		}
	}
	if err != nil {
		return nil, err
	}

	return &config.CheckpointExport{
		Checkpoint: ck,
		SourceS3:   p.Config,
	}, nil
}
