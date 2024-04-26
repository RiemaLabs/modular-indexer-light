package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"

	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
)

type S3 struct {
	Config       *configs.SourceS3
	MetaProtocol string
	Retry        int
}

func NewProviderS3(sourceS3 *configs.SourceS3, metaProtocol string, retry int) *S3 {
	return &S3{
		Config:       sourceS3,
		MetaProtocol: metaProtocol,
		Retry:        retry,
	}
}

func (p *S3) Get(ctx context.Context, height uint, hash string) (*configs.CheckpointExport, error) {
	var ck *checkpoint.Checkpoint
	var err error
	for i := 0; i < p.Retry; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ck, err = DownloadCheckpointByS3(p.Config.Region, p.Config.Bucket, p.Config.Name, p.MetaProtocol, strconv.Itoa(int(height)), hash, 100*time.Second)
			if err != nil {
				logs.Error.Println("Download S3 checkpoint error:", err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &configs.CheckpointExport{Checkpoint: ck, SourceS3: p.Config}, nil
}
