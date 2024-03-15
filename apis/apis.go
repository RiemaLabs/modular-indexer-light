package apis

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/RiemaLabs/indexer-light/config"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/gin-gonic/gin"
)

type RoundRobinBalancer struct {
	sync.Mutex
	index    int
	backends []string
}

func (r *RoundRobinBalancer) Next() string {
	r.Lock()
	defer r.Unlock()

	if len(r.backends) == 0 {
		return ""
	}

	next := r.backends[r.index]
	r.index = (r.index + 1) % len(r.backends)
	return next
}

func NewRoundRobinBalancer(backends []string) *RoundRobinBalancer {
	return &RoundRobinBalancer{index: -1, backends: backends}
}

func setupReverseProxy(proxyPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch constant.ApiState {
		case constant.ApiStateSync:
			c.JSON(http.StatusForbidden, "API is syncing")
			return
		case constant.ApiStateLoading:
			c.JSON(http.StatusForbidden, "API is loading")
			return
		case constant.ApiStateActive:

		}

		c.Request.URL.Path = proxyPath
		balancer := NewRoundRobinBalancer(config.GetCommitteeIndexerApi(config.Config))
		target := balancer.Next()
		if target == "" {
			c.String(http.StatusForbidden, "No available backends")
			return
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			c.String(http.StatusForbidden, fmt.Sprintf("Failed to parse backend URL: %s", err))
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func ApiSrv() {
	r := gin.Default()

	r.POST("/brc20_verifiable_light_get_current_balance_of_wallet", setupReverseProxy("brc20_verifiable_get_current_balance_of_wallet"))

	r.POST("/brc20_verifiable_light_block_height", setupReverseProxy("/brc20_verifiable_block_height"))

	r.POST("/brc20_verifiable_light_state", func(context *gin.Context) {
		context.JSON(http.StatusOK, Brc20VerifiableLightStateResponse{
			State: constant.ApiState,
		})
	})

	r.POST("/brc20_verifiable_light_last_checkpoint", func(context *gin.Context) {
		//TODO::
		context.JSON(http.StatusOK, Brc20VerifiableLightLastCheckpointResponse{})
	})

	r.POST("/brc20_verifiable_light_last_checkpoint", func(context *gin.Context) {
		//TODO::
		context.JSON(http.StatusOK, Brc20VerifiableLightLastCheckpointResponse{})
	})

	fmt.Println("Starting Gin HTTP reverse proxy server on :8081...")
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}
