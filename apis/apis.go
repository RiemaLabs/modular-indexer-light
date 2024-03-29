package apis

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/runtime"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartService(df *runtime.RuntimeState, enableDebug bool) {

	if !enableDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// TODO: Medium. Add the TRUSTED_PROXIES to our config
	// trustedProxies := os.Getenv("TRUSTED_PROXIES")
	// if trustedProxies != "" {
	//     r.SetTrustedProxies([]string{trustedProxies})
	// }

	r.Use(gin.Recovery(), CheckState(), gin.Logger(), cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// r.LoadHTMLGlob("index/*")
	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })

	r.StaticFile("/", "./build/index.html")
	r.StaticFile("/logo192.png", "./build/logo192.png")
	r.StaticFile("/manifest.json", "./build/manifest.json")
	r.StaticFile("/favicon.ico", "./build/favicon.ico")
	r.Static("/static", "./build/static")

	r.GET(constant.LightBlockHeight, func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("%d", df.CurrentHeight())))
	})

	r.GET(constant.LightState, func(c *gin.Context) {
		c.JSON(http.StatusOK, Brc20VerifiableLightStateResponse{
			State: constant.ApiState.String(),
		})
	})

	r.GET(constant.LightCurrentBalanceOfWallet, func(c *gin.Context) {
		ck := df.CurrentFirstCheckpoint().Checkpoint

		GetCurrentBalanceOfWallet(c, ck)
	})

	r.GET(constant.LightCurrentBalanceOfPkscript, func(c *gin.Context) {
		ck := df.CurrentFirstCheckpoint().Checkpoint

		GetCurrentBalanceOfPkscript(c, ck)
	})

	r.GET(constant.LightCurrentCheckpoints, func(c *gin.Context) {
		cur := df.CurrentCheckpoints()
		c.JSON(http.StatusOK, cur)
	})

	r.GET(constant.LightLastCheckpoint, func(c *gin.Context) {
		lt := df.LastCheckpoint()
		c.JSON(http.StatusOK, lt)
	})

	// TODO: Medium. Allow user to setup port.
	r.Run(":8080")
}
