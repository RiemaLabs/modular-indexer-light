package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/RiemaLabs/modular-indexer-light/log"
)

type Client struct {
	Client *http.Client
}

func NewClient() *Client {
	return &Client{
		Client: &http.Client{Timeout: time.Minute},
	}
}

func (c *Client) Post(ctx context.Context, url string, data interface{}, headers map[string]string) (target []byte, err error) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Error("http_client", "Marshal.data.err", err, "data", data)
		return nil, err
	}
	log.Debug("http_client", "url", url, "data", string(body))
	req, err := http.NewRequestWithContext(ctx, "POST", url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		log.Error("http_client", "NewRequestWithContext.err", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, header := range headers {
		req.Header.Set(key, header)
	}

	do, err := c.Client.Do(req)
	if err != nil {
		log.Error("http_client", "Client.Do", err)
		return nil, err
	}

	defer do.Body.Close()
	return io.ReadAll(do.Body)
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	log.Debug("http_client", "method", "GET", "url", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error("http_client", "NewRequestWithContext.err", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	do, err := c.Client.Do(req)
	if err != nil {
		log.Error("http_client", "Client.Do", err, "req", req)
		return nil, err
	}
	defer do.Body.Close()
	return io.ReadAll(do.Body)
}
