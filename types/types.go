package types

import (
	"context"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
)

type (
	Config struct {
		*CommitteeIndexer `json:"committeeIndexer"`
		*BitCoinRpc       `json:"bitCoinRpc"`
		// minimal entry
		MinimalCheckPoint int    `json:"minimalCheckPoint"`
		StartHeight       int    `json:"startHeight"`
		StartBlockHash    string `json:"startBlockHash"`
	}

	BitCoinRpc struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
	}

	CommitteeIndexer struct {
		S3      []*SourceS3 `json:"s3"`
		Da      []*SourceDa `json:"da"`
		TimeOut int         `json:"timeOut"`
	}

	// data source

	SourceS3 struct {
		Bucket       string `json:"bucket"`
		AccessKey    string `json:"accessKey"`
		ApiUrl       string `json:"apiUrl"`
		IndexerName  string `json:"indexerName"`
		MetaProtocol string `json:"metaProtocol"`
		Region       string `json:"region"`
	}

	SourceDa struct {
		NamespaceID  string `json:"namespaceID"`
		IndexerName  string `json:"indexerName"`
		MetaProtocol string `json:"metaProtocol"`
		ApiUrl       string `json:"apiUrl"`
		Rpc          string `json:"rpc"`
	}

	Source struct {
		*SourceS3 `json:"*SourceS3,omitempty"`
		*SourceDa `json:"*SourceDa,omitempty"`
	}

	CheckPointObject struct {
		CheckPoint *checkpoint.Checkpoint `json:"checkPoint,omitempty"`
		Name       string                 `json:"name,omitempty"`
		Type       string                 `json:"type,omitempty"`
		Source     *Source                `json:"source,omitempty"`
	}
)

type CheckPointProvider interface {
	GetCheckpoint(ctx context.Context, height uint, hash string) *CheckPointObject
}
