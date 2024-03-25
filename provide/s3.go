package provide

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/clients/http"
	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/log"
	"github.com/RiemaLabs/modular-indexer-light/types"
)

func NewS3(config *types.SourceS3) *ProviderS3 {
	return &ProviderS3{
		Name:   constant.ProvideS3Name,
		Config: config,
	}
}

func (p *ProviderS3) GetCheckpoint(ctx context.Context, height uint, hash string) *types.CheckPointObject {
	objectKey := fmt.Sprintf("test/checkpoint-%s-%s-%d-%s.json",
		p.Config.IndexerName, constant.DefaultMetaProtocol, height, hash)

	endpoint, err := url.JoinPath(p.Config.ApiUrl, objectKey)
	if err != nil {
		log.Error("ProviderS3", "s3.JoinPath", err)
		return nil
	}
	//https://nubit-modular-indexer.s3.us-west-2.amazonaws.com/test/checkpoint-world-brc-20-835660-000000000000000000033e6eaf21a28cafffd89892c62b468ca73b91684690a6.json
	//https://nubit-modular-indexer.s3.us-west-2.amazonaws.com/test/checkpoint-world-brc-20-780212-000000000000000000046f13e9532206eb5432baea2eb502fb2ec4bba803b434.json
	log.Debug("ProviderS3", "objectKey", objectKey, "endpoint", endpoint)
	client := http.NewClient()
	get, err := client.Get(ctx, endpoint)
	if err != nil {
		log.Error("ProviderS3", "s3.Get", err)
		return nil
	}
	var (
		ck = &checkpoint.Checkpoint{}
	)
	err = json.Unmarshal(get, &ck)
	if err != nil {
		log.Error("ProviderS3", "s3.Unmarshal", err)
		return nil
	}
	return &types.CheckPointObject{
		CheckPoint: ck,
		Name:       p.Name,
		Type:       "",
		Source: &types.Source{
			SourceS3: p.Config,
		},
	}
}

//func (p *ProviderS3) GetCheckpoint(ctx context.Context, height uint, hash string) *types.CheckPointObject {
//
//	cfg, err := awscfg.LoadDefaultConfig(ctx,
//		awscfg.WithRegion(p.Config.Region),
//		awscfg.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
//			if service == s3.ServiceID && p.Config.Url != "" {
//				return aws.Endpoint{
//					URL:           p.Config.Url,
//					SigningRegion: p.Config.Region,
//				}, nil
//			}
//			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
//		})),
//		awscfg.WithS3DisableMultiRegionAccessPoints(false),
//	)
//	if err != nil {
//		log.Error("ProviderS3", "s3.LoadDefaultConfig", err)
//		return nil
//	}
//	awsS3Client := s3.NewFromConfig(cfg)
//	downloader := manager.NewDownloader(awsS3Client)
//	objectKey := fmt.Sprintf("test/checkpoint-%s-%s-%d-%s.json",
//		p.Config.IndexerName, constant.DefaultMetaProtocol, height, hash)
//
//	objectKey = "test/checkpoint-world-brc-20-780212-000000000000000000046f13e9532206eb5432baea2eb502fb2ec4bba803b434.json"
//	getObjectInput := &s3.GetObjectInput{
//		Bucket: aws.String(p.Config.Bucket),
//		Key:    aws.String(objectKey),
//	}
//	log.Debug("ProviderS3", "objectKey", objectKey, "Bucket", p.Config.Bucket)
//	object, err := downloader.S3.GetObject(ctx, getObjectInput)
//	if err != nil {
//		log.Error("ProviderS3", "s3.GetObject", err)
//		return nil
//	}
//	defer object.Body.Close()
//	var (
//		data = []byte{}
//		ck   = &checkpoint.Checkpoint{}
//	)
//	read, err := object.Body.Read(data)
//	if err != nil {
//		log.Error("ProviderS3", "s3.Read", err)
//		return nil
//	}
//	err = json.Unmarshal(data[:read], &ck)
//	if err != nil {
//		log.Error("ProviderS3", "s3.Unmarshal", err)
//		return nil
//	}
//	return &types.CheckPointObject{
//		CheckPoint: ck,
//		Name:       p.Name,
//		Type:       "",
//		Source: &types.Source{
//			SourceS3: p.Config,
//		},
//	}
//}
