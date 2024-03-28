package config

import (
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
)

// TODO: Medium. Re-organize type structs.

type (
	Config struct {
		CommitteeIndexers CommitteeIndexers `json:"committeeIndexers"`

		BitcoinRPC        string `json:"bitcoinRPC"`
		MinimalCheckpoint int    `json:"minimalCheckpoint"`
		MetaProtocol      string `json:"metaProtocol"`
		Version           string `json:"version"`
	}

	CommitteeIndexers struct {
		S3 []SourceS3 `json:"s3"`
		DA []SourceDA `json:"da"`
	}

	// data source

	SourceS3 struct {
		Region string `json:"region"`
		Bucket string `json:"bucket"`
		Name   string `json:"name"`
	}

	SourceDA struct {
		Network     string `json:"network"`
		NamespaceID string `json:"namespaceID"`
		Name        string `json:"name"`
	}

	CheckpointExport struct {
		Checkpoint *checkpoint.Checkpoint `json:"checkPoint"`
		SourceS3   *SourceS3              `json:"sourceS3,omitempty"`
		SourceDA   *SourceDA              `json:"sourceDa,omitempty"`
	}

	Evidence struct {
		Height            uint   `json:"height"`
		Hash              string `json:"hash"`
		CorrectCommitment string `json:"correctCommitment"`
		FraudCommitment   string `json:"fraudCommitment"`
	}

	Blacklist struct {
		Evidence *Evidence `json:"evidence"`
		SourceS3 *SourceS3 `json:"sourceS3"`
		SourceDA *SourceDA `json:"sourceDa"`
	}
)
