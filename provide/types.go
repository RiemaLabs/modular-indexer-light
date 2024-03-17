package provide

import (
	"github.com/RiemaLabs/indexer-light/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ProviderS3 struct {
	Name        string          `json:"name"`
	Config      *types.SourceS3 `json:"config"`
	awsS3Client *s3.Client
}

type ProviderDa struct {
	Name   string          `json:"name"`
	Config *types.SourceDa `json:"config"`
}
