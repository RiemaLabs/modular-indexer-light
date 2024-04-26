package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/httputl"
)

const DefaultVersion = "2.0"

type Request struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
	ReqID   string `json:"id"`
}

func NewRequest(ctx context.Context, method string, params any) *Request {
	return &Request{
		Version: DefaultVersion,
		Method:  method,
		Params:  params,
		ReqID:   httputl.RequestID(ctx),
	}
}

type Response[T, E any] struct {
	Version string `json:"jsonrpc"`
	Result  T      `json:"result"`
	Error   E      `json:"error"`
	ReqID   string `json:"id"`
}

type HasErr interface {
	Err() error
}

// Client is the JSON-RPC client.
//
// Why don't we just use Go's standard `net/rpc/jsonrpc`? Because it's not backed by an *http.Transport, there's just a
// poor man's TCP connection, it's not efficient and robust.
type Client interface {
	Call(ctx context.Context, method string, params, out interface{}) error
}

type client struct {
	nodeURL string
	cl      *http.Client
}

func New(rawURL string) (Client, error) {
	if _, err := url.Parse(rawURL); err != nil {
		return nil, fmt.Errorf("invalid node URL: %v", err)
	}
	return &client{nodeURL: rawURL, cl: httputl.Client}, nil
}

func (c *client) Call(ctx context.Context, method string, params, out any) error {
	reqID := httputl.RequestID(ctx)
	in := NewRequest(ctx, method, params)

	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf(
			"marshal error: nodeURL=%s, method=%s, params=%v, reqID=%s, err=%v",
			c.nodeURL,
			method,
			params,
			reqID,
			err,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.nodeURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf(
			"invalid request: nodeURL=%s, method=%s, params=%v, reqID=%s, err=%v",
			c.nodeURL,
			method,
			params,
			reqID,
			err,
		)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", in.ReqID) // for tracing

	rsp, err := c.cl.Do(req)
	if err != nil {
		return fmt.Errorf(
			"HTTP transport error: nodeURL=%s, method=%s, params=%v, reqID=%s, err=%v",
			c.nodeURL,
			method,
			params,
			reqID,
			err,
		)
	}
	defer func() { _ = rsp.Body.Close() }()

	rspBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf(
			"read response body error: nodeURL=%s, method=%s, params=%v, reqID=%s, err=%v",
			c.nodeURL,
			method,
			params,
			reqID,
			err,
		)
	}

	if err := json.Unmarshal(rspBody, out); err != nil {
		return fmt.Errorf(
			"unmarshal error: nodeURL=%s, method=%s, params=%v, rspBody=%s, reqID=%s, err=%v",
			c.nodeURL,
			method,
			params,
			string(rspBody),
			reqID,
			err,
		)
	}

	if hasErr, ok := out.(HasErr); ok {
		if err := hasErr.Err(); err != nil {
			return fmt.Errorf(
				"response error: nodeURL=%s, method=%s, params=%v, rspBody=%s, reqID=%s, err=%v",
				c.nodeURL,
				method,
				params,
				string(rspBody),
				reqID,
				err,
			)
		}
	}

	return nil
}
