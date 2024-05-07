package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/states"
)

const DefaultAddr = ":8080"

func StartService(enableDebug bool, addr string) {
	if !enableDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"POST", "GET"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)
	r.GET("/v1/brc20_verifiable/light/state", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			State fmt.Stringer `json:"state"`
		}{
			State: states.Status(states.S.Status.Load()),
		})
	})
	g := r.Group("v1")
	{
		g.Use(CheckState)
		g.GET("/brc20_verifiable/light/block_height", func(c *gin.Context) { c.String(http.StatusOK, strconv.Itoa(int(states.S.CurrentHeight()))) })
		g.GET("/brc20_verifiable/light/current_balance_of_wallet", HandleGetCurrentBalanceOfWallet)
		g.GET("/brc20_verifiable/light/current_balance_of_pkscript", HandleGetCurrentBalanceOfPkscript)
		g.GET("/brc20_verifiable/light/checkpoints", func(c *gin.Context) { c.JSON(http.StatusOK, states.S.CurrentCheckpoints()) })
		g.GET("/brc20_verifiable/light/last_checkpoint", func(c *gin.Context) { c.JSON(http.StatusOK, states.S.LastCheckpoint()) })
	}

	if addr == "" {
		addr = DefaultAddr
	}
	if err := r.Run(addr); !errors.Is(err, http.ErrServerClosed) {
		logs.Error.Fatalln("Server exit with error:", err)
	}
}
