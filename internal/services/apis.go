package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/RiemaLabs/modular-indexer-light/internal/constant"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/runtime"
)

const DefaultAddr = ":8080"

func StartService(enableDebug bool, addr string) {
	if !enableDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger(), cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/v1/brc20_verifiable/light/state", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			State string `json:"state"`
		}{
			State: constant.ApiState.String(),
		})
	})
	serv := r.Group("v1")
	{
		serv.Use(CheckState())
		serv.GET("/brc20_verifiable/light/block_height", func(c *gin.Context) {
			c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("%d", df.CurrentHeight())))
		})

		serv.GET("/brc20_verifiable/light/current_balance_of_wallet", func(c *gin.Context) {
			ck := runtime.S.CurrentFirstCheckpoint().Checkpoint

			GetCurrentBalanceOfWallet(c, ck)
		})

		serv.GET("/brc20_verifiable/light/current_balance_of_pkscript", func(c *gin.Context) {
			ck := df.CurrentFirstCheckpoint().Checkpoint

			GetCurrentBalanceOfPkscript(c, ck)
		})

		serv.GET("/brc20_verifiable/light/checkpoints", func(c *gin.Context) {
			cur := df.CurrentCheckpoints()
			c.JSON(http.StatusOK, cur)
		})

		serv.GET("/brc20_verifiable/light/last_checkpoint", func(c *gin.Context) {
			lt := df.LastCheckpoint()
			c.JSON(http.StatusOK, lt)
		})
	}

	if addr == "" {
		addr = DefaultAddr
	}
	if err := r.Run(addr); !errors.Is(err, http.ErrServerClosed) {
		logs.Error.Fatal("Server exit with error: ", err)
	}
}
