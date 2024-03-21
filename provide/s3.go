package provide

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/RiemaLabs/indexer-committee/checkpoint"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/log"
	"github.com/RiemaLabs/indexer-light/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3(config *types.SourceS3) *ProviderS3 {
	cfg, err := awscfg.LoadDefaultConfig(context.Background(), awscfg.WithRegion(config.Region), awscfg.WithEC2IMDSEndpoint(config.Url))
	if err != nil {
		log.Error("ProviderS3", "s3.LoadDefaultConfig", err)
		return nil
	}
	awsS3Client := s3.NewFromConfig(cfg)
	return &ProviderS3{
		Name:        constant.ProvideS3Name,
		Config:      config,
		awsS3Client: awsS3Client,
	}
}

func (p *ProviderS3) GetCheckpoint(ctx context.Context, height uint, hash string) *types.CheckPointObject {
	select {
	case <-ctx.Done():
		log.Error("ProviderS3", "context", ctx.Err())
		return nil
	default:

	}
	downloader := manager.NewDownloader(p.awsS3Client)
	objectKey := fmt.Sprintf("checkpoint-%s-%s-%d-%s.json",
		p.Config.IndexerName, constant.DefaultMetaProtocol, height, hash)
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(p.Config.Bucket),
		Key:    aws.String(objectKey),
	}
	log.Debug("ProviderS3", "objectKey", objectKey, "Bucket", p.Config.Bucket)
	object, err := downloader.S3.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		log.Error("ProviderS3", "s3.GetObject", err)
		return nil
	}
	defer object.Body.Close()
	var (
		data = []byte{}
		ck   = &checkpoint.Checkpoint{}
	)
	read, err := object.Body.Read(data)
	if err != nil {
		log.Error("ProviderS3", "s3.Read", err)
		return nil
	}
	err = json.Unmarshal(data[:read], &ck)
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

func getMACAddress() (string, error) {
	// all interfaces info
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// the first MAC addr of non-vertical interface
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			// filter virtual and loop interfaces
			// println(iface.HardwareAddr.String())
			return iface.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("no active non-loopback network interface found")
}
