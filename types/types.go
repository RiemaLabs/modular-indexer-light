package types

import "github.com/RiemaLabs/indexer-committee/checkpoint"

type (
	Config struct {
		*CommitteeIndexer    `json:"committeeIndexer"`
		*CommitteeIndexerApi `json:"committeeIndexerApi"`
		// minimal entry
		MinimalCheckPoint int    `json:"minimalCheckPoint"`
		StartHeight       int    `json:"startHeight"`
		StartBlockHash    string `json:"startBlockHash"`
	}

	CommitteeIndexerApi struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	CommitteeIndexer struct {
		S3 []*SourceS3 `json:"s3"`
		Da []*SourceDa `json:"da"`
	}

	// data source

	SourceS3 struct {
		Bucket       string `json:"bucket"`
		AccessKey    string `json:"accessKey"`
		Url          string `json:"url"`
		IndexerName  string `json:"indexerName"`
		MetaProtocol string `json:"metaProtocol"`
		ApiUrl       string `json:"apiUrl"`
	}

	SourceDa struct {
		Address       string `json:"address"`
		TransactionID string `json:"transactionID"`
		IndexerName   string `json:"indexerName"`
		MetaProtocol  string `json:"metaProtocol"`
		ApiUrl        string `json:"apiUrl"`
	}
)

type CheckPointProvider interface {
	GetCheckpoint() *checkpoint.Checkpoint
}
