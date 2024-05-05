package services

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ethereum/go-verkle"
	"github.com/gin-gonic/gin"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/internal/clients/committee"
	"github.com/RiemaLabs/modular-indexer-light/internal/states"
)

func CheckState() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch s := states.Status(states.S.State.Load()); s {
		case states.StatusActive:
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, s.String())
		}
	}
}

func GetCurrentBalanceOfWallet(c *gin.Context, ck *checkpoint.Checkpoint) {
	tick := c.DefaultQuery("tick", "")
	wallet := c.DefaultQuery("wallet", "")

	cl, err := committee.New(ck.URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	balance, err := cl.CurrentBalanceOfWallet(context.Background(), tick, wallet)
	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}

	commitmentBytes, _ := base64.StdEncoding.DecodeString(ck.Commitment)
	var point verkle.Point
	_ = point.SetBytes(commitmentBytes)

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

	cl, err := committee.New(ck.URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	balance, err := cl.CurrentBalanceOfPkscript(context.Background(), tick, pkscript)
	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}

	commitmentBytes, _ := base64.StdEncoding.DecodeString(ck.Commitment)
	point := verkle.Point{}
	_ = point.SetBytes(commitmentBytes)

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
