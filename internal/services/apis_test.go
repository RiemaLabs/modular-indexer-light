package services

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Medium. Uniform the default service port.
var localHost = "http://127.0.0.1:8080"

func TestBlockHeight(t *testing.T) {
	client := &http.Client{}
	url, _ := url.JoinPath(localHost, "/brc20_verifiable/light/block_height")
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLightCurrentBalanceOfWallet(t *testing.T) {
	client := &http.Client{}
	urlString, _ := url.JoinPath(localHost, "/brc20_verifiable/light/current_balance_of_wallet")

	p := url.Values{}
	p.Add("tick", "btcs")
	p.Add("wallet", "bc1qqqpx5690calxc5q2x83mhyftk6zmtvlprvdujz")

	req, _ := http.NewRequest("GET", urlString+"?"+p.Encode(), nil)
	resp, err := client.Do(req)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestLightCurrentBalanceOfPkscript(t *testing.T) {
	client := &http.Client{}
	urlString, _ := url.JoinPath(localHost, "/brc20_verifiable/light/current_balance_of_pkscript")

	p := url.Values{}
	p.Add("tick", "btcs")
	p.Add("wallet", "bc1qqqpx5690calxc5q2x83mhyftk6zmtvlprvdujz")

	req, _ := http.NewRequest("GET", urlString+"?"+p.Encode(), nil)
	resp, err := client.Do(req)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestLightLastCheckpoint(t *testing.T) {
	client := &http.Client{}
	urlString, _ := url.JoinPath(localHost, "/brc20_verifiable/light/last_checkpoint")

	req, _ := http.NewRequest("GET", urlString, nil)
	resp, err := client.Do(req)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestLightCurrentCheckpoints(t *testing.T) {
	client := &http.Client{}
	urlString, _ := url.JoinPath(localHost, "/brc20_verifiable/light/checkpoints")

	req, _ := http.NewRequest("GET", urlString, nil)
	resp, err := client.Do(req)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
