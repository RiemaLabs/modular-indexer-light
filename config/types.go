package config

import (
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
)

// TODO: Medium. Re-organize type structs.

type (
	Config struct {
		CommitteeIndexers CommitteeIndexers `json:"committeeIndexers"`

		Verification Verification `json:"verification"`

		Report Report `json:"report"`
	}

	CommitteeIndexers struct {
		S3 []SourceS3 `json:"s3"`
		DA []SourceDA `json:"da"`
	}

	Verification struct {
		BitcoinRPC        string `json:"bitcoinRPC"`
		MinimalCheckpoint int    `json:"minimalCheckpoint"`
		MetaProtocol      string `json:"metaProtocol"`
	}

	Report struct {
		Name        string `json:"name"`
		Network     string `json:"network"`
		NamespaceID string `json:"namespaceID"`
		GasCoupon   string `json:"gasCoupon"`
		PrivateKey  string `json:"privateKey"`
		Timeout     int    `json:"timeout"`
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
