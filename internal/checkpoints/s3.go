package checkpoints

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
)

type S3 struct {
	Config       *configs.SourceS3
	MetaProtocol string
}

func NewProviderS3(sourceS3 *configs.SourceS3, metaProtocol string) *S3 {
	return &S3{
		Config:       sourceS3,
		MetaProtocol: metaProtocol,
	}
}

func (p *S3) Get(ctx context.Context, height uint, hash string) (*configs.CheckpointExport, error) {
	var (
		ck  *checkpoint.Checkpoint
		err error
	)
	for i := 0; i < DefaultRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if ck, err = p.doDownload(height, hash); err != nil {
				logs.Error.Println("Download S3 checkpoint error:", err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &configs.CheckpointExport{Checkpoint: ck, SourceS3: p.Config}, nil
}

func (p *S3) doDownload(height uint, hash string) (*checkpoint.Checkpoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	u := &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.s3.%s.amazonaws.com", p.Config.Bucket, p.Config.Region),
		Path:   fmt.Sprintf("checkpoint-%s-%s-%d-%s.json", p.Config.Name, p.MetaProtocol, height, hash),
	}
	obj := u.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, obj, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid S3 request: obj=%s, err=%v", obj, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport S3 error: obj=%s, err=%v", obj, err)
	}
	defer func() { _ = resp.Body.Close() }()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read S3 response error: obj=%s, err=%v", obj, err)
	}

	var c checkpoint.Checkpoint
	if err := json.Unmarshal(bytes, &c); err != nil {
		type s3Resp struct {
			Code, Message string
		}
		var s3resp s3Resp
		if err := xml.Unmarshal(bytes, &s3resp); err == nil {
			return nil, fmt.Errorf("download from S3 error: obj=%s, code=%s, msg=%s", obj, s3resp.Code, s3resp.Message)
		}
		return nil, fmt.Errorf("unmarshal checkpoint error: obj=%s, body=%q, err=%v", obj, string(bytes), err)
	}

	return &c, nil
}
