package verify

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RiemaLabs/indexer-committee/checkpoint"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/config"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/provide"
	"github.com/RiemaLabs/indexer-light/types"
)

type DefiniteCheckpoint struct {
	mux            sync.RWMutex
	cfg            *types.Config
	PreCheckpoint  *checkpoint.Checkpoint
	PostCheckpoint *checkpoint.Checkpoint
}

func init() {
	if DefiniteState == nil {
		DefiniteState = NewState()
	}
}

// DefiniteState ...
var DefiniteState *DefiniteCheckpoint

func NewState() *DefiniteCheckpoint {
	return &DefiniteCheckpoint{mux: sync.RWMutex{}, cfg: config.Config}
}

func (d *DefiniteCheckpoint) Update(getter getter.OrdGetter, post *types.CheckPointObject) error {
	ctx := context.Background()
	ctx, CancelFunc := context.WithTimeout(ctx, time.Duration(config.Config.CommitteeIndexer.TimeOut)*time.Second)
	defer CancelFunc()
	if d.PreCheckpoint == nil {
		var committee types.CheckPointProvider
		switch strings.ToLower(post.Name) {
		case strings.ToLower(constant.ProvideS3Name):
			for _, s3 := range d.cfg.CommitteeIndexer.S3 {
				if s3.IndexerName == post.CheckPoint.Name {
					committee = provide.NewS3(s3)
				}
			}
		case strings.ToLower(constant.ProvideDaName):
			for _, da := range d.cfg.CommitteeIndexer.Da {
				if da.IndexerName == post.CheckPoint.Name {
					committee = provide.NewDA(da)
				}
			}
		}
		h, err := strconv.Atoi(post.CheckPoint.Height)
		if err != nil {
			return err
		}
		preH := uint(h - 1)
		hash, err := getter.GetBlockHash(preH)
		if err != nil {
			return err
		}
		if committee == nil {
			return errors.New("provide is nil")
		}
		ckObj := committee.GetCheckpoint(ctx, preH, hash)
		if ckObj != nil {
			d.SetPre(ckObj.CheckPoint)
		}
	}
	if d.PostCheckpoint != nil {
		d.SetPre(d.PostCheckpoint)
	}
	if post != nil {
		d.SetPost(post.CheckPoint)
	}
	return nil
}

func (d *DefiniteCheckpoint) SetPre(pre *checkpoint.Checkpoint) {
	d.mux.TryLock()
	d.PreCheckpoint = pre
	d.mux.Unlock()
}

func (d *DefiniteCheckpoint) SetPost(post *checkpoint.Checkpoint) {
	d.mux.TryLock()
	d.PostCheckpoint = post
	d.mux.Unlock()
}
