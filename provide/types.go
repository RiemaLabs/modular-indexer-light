package provide

import (
	"context"

	"github.com/RiemaLabs/indexer-light/indexer"
	"github.com/RiemaLabs/indexer-light/types"
	sdk "github.com/RiemaLabs/nubit-da-sdk"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ProviderS3 struct {
	Name        string          `json:"name"`
	Config      *types.SourceS3 `json:"config"`
	awsS3Client *s3.Client
}

type ProviderDa struct {
	ctx       context.Context
	Name      string          `json:"name"`
	Config    *types.SourceDa `json:"config"`
	Client    *sdk.NubitSDK
	Committee *indexer.CommitteeIndexer
}
