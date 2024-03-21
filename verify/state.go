package verify

import (
	"strconv"
	"sync"

	"github.com/RiemaLabs/indexer-committee/checkpoint"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
	"github.com/RiemaLabs/indexer-light/config"
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

func (d *DefiniteCheckpoint) Update(getter getter.OrdGetter, post *checkpoint.Checkpoint) error {
	d.mux.TryLock()
	defer d.mux.Unlock()
	if d.PreCheckpoint == nil {
		//TODO:: check post checkpoint is string `S3` or `DA` ?
		var committee types.CheckPointProvider
		switch post.Name {
		case "S3":
			for _, s3 := range d.cfg.CommitteeIndexer.S3 {
				if s3.IndexerName == post.Name {
					committee = provide.NewS3(s3)
				}
			}
		case "DA":
			for _, da := range d.cfg.CommitteeIndexer.Da {
				if da.IndexerName == post.Name {
					committee = provide.NewDA(da)
				}
			}
		}
		h, err := strconv.Atoi(post.Height)
		if err != nil {
			return err
		}
		preH := uint(h - 1)
		hash, err := getter.GetBlockHash(preH)
		if err != nil {
			return err
		}
		ckObj := committee.GetCheckpoint(nil, preH, hash)
		d.SetPre(ckObj.CheckPoint)
	}

	d.SetPre(d.PostCheckpoint)
	d.SetPost(post)

	return nil
}

func (d *DefiniteCheckpoint) SetPre(pre *checkpoint.Checkpoint) {
	d.mux.TryLock()
	defer d.mux.Unlock()
	d.PreCheckpoint = pre
}

func (d *DefiniteCheckpoint) SetPost(post *checkpoint.Checkpoint) {
	d.mux.TryLock()
	defer d.mux.Unlock()
	d.PostCheckpoint = post
}
