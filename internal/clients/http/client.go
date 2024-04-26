package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Client *http.Client
}

func New() *Client {
	return &Client{Client: &http.Client{Timeout: time.Minute}}
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP GET request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP transport error: %v", err)
	}
	defer func() { _ = rsp.Body.Close() }()
	return io.ReadAll(rsp.Body)
}
