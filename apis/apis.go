package apis

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"

	"github.com/RiemaLabs/indexer-light/config"
	"github.com/RiemaLabs/indexer-light/constant"
	"github.com/RiemaLabs/indexer-light/indexer"
	"github.com/RiemaLabs/indexer-light/verify"
	"github.com/ethereum/go-verkle"
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

func Start() {
	r := gin.Default()

	r.POST(constant.LightBlockHigh, setupReverseProxy(constant.BlockHigh))

	r.POST(constant.LightState, func(context *gin.Context) {
		context.JSON(http.StatusOK, Brc20VerifiableLightStateResponse{
			State: constant.ApiState,
		})
	})

	r.POST(constant.LightBalance, func(c *gin.Context) {
		balancer := NewRoundRobinBalancer(config.GetCommitteeIndexerApi(config.Config))
		target := balancer.Next()
		if target == "" {
			c.String(http.StatusForbidden, "No available backends")
			return
		}
		req := &Brc20VerifiableLightGetCurrentBalanceOfWalletRequest{}
		err := c.BindJSON(req)
		if err != nil {
			c.String(http.StatusForbidden, "Parameter error")
			return
		}
		balance, err := indexer.NewClient(c, target, target).GetBalance(req.Tick, req.Pkscript)
		if err != nil {
			return
		}
		if balance != nil {
			//TODO:: balance.Proof to verkle.Proof
			preProof := &verkle.Proof{}
			prePointByte, err := base64.StdEncoding.DecodeString(verify.DefiniteState.PreCheckpoint.Commitment)
			if err != nil {
				c.JSON(http.StatusOK, "")
				return
			}
			prePoint := &verkle.Point{}
			err = prePoint.SetBytes(prePointByte)
			if err != nil {
				c.JSON(http.StatusOK, "")
				return
			}
			err = verify.VerifyProof(preProof, prePoint)
			if err != nil {
				c.JSON(http.StatusOK, "")
				return
			}
			h, _ := strconv.Atoi(verify.DefiniteState.PostCheckpoint.Height)
			c.JSON(http.StatusOK, Brc20VerifiableLightGetCurrentBalanceOfWalletResponse{
				Result:      balance.Result,
				BlockHeight: h,
			})
		}

	})

	r.POST(constant.LightCheckpoint, func(context *gin.Context) {
		//TODO::
		context.JSON(http.StatusOK, Brc20VerifiableLightLastCheckpointResponse{})
	})

	r.POST(constant.LightLastCheckpoint, func(context *gin.Context) {
		//TODO::
		context.JSON(http.StatusOK, Brc20VerifiableLightLastCheckpointResponse{})
	})

	fmt.Println("Starting Gin HTTP reverse proxy server on :8081...")
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}
