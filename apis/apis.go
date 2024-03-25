package apis

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/transfer"
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

func Start() {
	r := gin.Default()

	r.Use(gin.Recovery(), CheckState(), gin.Logger())

	r.GET(constant.LightBlockHigh, setupReverseProxy(constant.BlockHigh))

	r.GET(constant.LightState, func(context *gin.Context) {
		context.JSON(http.StatusOK, Brc20VerifiableLightStateResponse{
			State: constant.ApiState,
		})
	})

	r.GET(constant.LightBalance, GetcurrentBalanceOfWallet)

	r.GET(constant.LightCheckpoint, func(context *gin.Context) {
		//TODO::
		context.JSON(http.StatusOK, Brc20VerifiableLightLastCheckpointResponse{})
	})

	r.GET(constant.LightLastCheckpoint, func(context *gin.Context) {
		//TODO::
		context.JSON(http.StatusOK, Brc20VerifiableLightLastCheckpointResponse{})
	})

	r.POST(constant.LightTransfer, func(ctx *gin.Context) {
		req := Brc20VerifiableLightTransferVerifyRequest{}
		if err := ctx.BindJSON(&req); nil != err {
			ctx.JSON(http.StatusBadRequest, Brc20VerifiableLightTransferVerifyResponse{false, errors.New("unauthed parameter")})
			return
		}

		if is, msg := req.Check(); !is {
			ctx.JSON(http.StatusOK, Brc20VerifiableLightTransferVerifyResponse{false, errors.New(msg)})
			return
		}

		is, err := transfer.Verify(req.Transfers, req.BlockHeight)
		ctx.JSON(http.StatusOK, Brc20VerifiableLightTransferVerifyResponse{is, err})
	})

	fmt.Println("Starting Gin HTTP reverse proxy server on :8081...")
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
