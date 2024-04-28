package services

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ethereum/go-verkle"
	"github.com/gin-gonic/gin"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/internal/clients/committee"
	"github.com/RiemaLabs/modular-indexer-light/internal/constant"
	"github.com/RiemaLabs/modular-indexer-light/internal/states"
)

func CheckState() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch s := constant.ApiStatus(states.S.State.Load()); s {
		case constant.StateActive:
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, s.String())
		}
	}
}

func GetCurrentBalanceOfWallet(c *gin.Context, ck *checkpoint.Checkpoint) {
	tick := c.DefaultQuery("tick", "")
	wallet := c.DefaultQuery("wallet", "")

	balance, err := committee.New(c, ck.Name, ck.URL).CurrentBalanceOfWallet(tick, wallet)
	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}

	pbytes, _ := base64.StdEncoding.DecodeString(ck.Commitment)
	var point verkle.Point
	point.SetBytes(pbytes)

	ok, err := apis.VerifyCurrentBalanceOfWallet(&point, tick, wallet, balance)
	if err != nil {
		msg := err.Error()
		if strings.HasPrefix(msg, "proof of absence") {
			msg = "Balance does not exist"
		}
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}

	if !ok {
		msg := "Failed to verify the result obtained from the committee indexer."
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func GetCurrentBalanceOfPkscript(c *gin.Context, ck *checkpoint.Checkpoint) {
	tick := c.DefaultQuery("tick", "")
	pkscript := c.DefaultQuery("pkscript", "")

	balance, err := committee.New(c, ck.Name, ck.URL).CurrentBalanceOfPkscript(tick, pkscript)
	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}

	pbytes, _ := base64.StdEncoding.DecodeString(ck.Commitment)
	point := verkle.Point{}
	point.SetBytes(pbytes)

	ok, err := apis.VerifyCurrentBalanceOfPkscript(&point, tick, pkscript, balance)
	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}

	if !ok {
		msg := "Failed to verify the result obtained from the committee indexer."
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusOK, balance)
}
