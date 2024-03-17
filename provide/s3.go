package provide

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/RiemaLabs/indexer-committee/checkpoint"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3(config *types.SourceS3) *ProviderS3 {
	// the SDK uses its default credential chain to find AWS credentials. This default credential chain looks for credentials in the following order:aws.Configconfig.LoadDefaultConfig
	// creds := credentials.NewStaticCredentialsProvider(your_access_key, your_secret_key, "")
	cfg, err := awscfg.LoadDefaultConfig(context.Background(), awscfg.WithRegion(config.Region))
	if err != nil {
		log.Fatal(err)
	}
	awsS3Client := s3.NewFromConfig(cfg)

	return &ProviderS3{
		Name:        constant.ProvideS3Name,
		Config:      config,
		awsS3Client: awsS3Client,
	}
}

func (p *ProviderS3) GetCheckpoint(height uint, hash string) *types.CheckPointObject {
	downloader := manager.NewDownloader(p.awsS3Client)
	objectKey := fmt.Sprintf("checkpoint-%s-%s-%d-%s.json",
		p.Config.IndexerName, constant.DefaultMetaProtocol, height, hash)
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(p.Config.Bucket),
		Key:    aws.String(objectKey),
	}
	object, err := downloader.S3.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		return nil
	}
	defer object.Body.Close()
	var (
		data = []byte{}
		ck   = &checkpoint.Checkpoint{}
	)
	read, err := object.Body.Read(data)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data[:read], &ck)
	if err != nil {
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
