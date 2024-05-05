package httputl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is a simple wrapper with timeout, as core client of other services.
var Client = &http.Client{Timeout: time.Minute}

func GetJSON(ctx context.Context, u *url.URL, out any) error {
	rawURL := u.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("invalid request: rawURL=%s, err=%v", rawURL, err)
	}
	if reqID := RequestID(ctx); reqID != "" {
		req.Header.Set("X-Request-Id", reqID)
	}

	rsp, err := Client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP transport error: rawURL=%s, err=%v", rawURL, err)
	}
	defer func() { _ = rsp.Body.Close() }()

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("response body read error: rawURL=%s, err=%v", rawURL, err)
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("unmarshal error: rawURL=%s, err=%v", rawURL, err)
	}

	return nil
}
