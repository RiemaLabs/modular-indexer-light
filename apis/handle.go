package apis

import (
	"encoding/base64"
	"net/http"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/clients/committee"
	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/ethereum/go-verkle"
	"github.com/gin-gonic/gin"
)

func CheckState() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch constant.ApiState {
		case constant.StateActive:
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, constant.ApiState.String())
		}
	}
}

func GetCurrentBalanceOfWallet(c *gin.Context, ck *checkpoint.Checkpoint) {
	tick := c.DefaultQuery("tick", "")
	wallet := c.DefaultQuery("wallet", "")

	balance, err := committee.NewCommitteeIndexerClient(c, ck.Name, ck.URL).CurrentBalanceOfWallet(tick, wallet)
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

	balance, err := committee.NewCommitteeIndexerClient(c, ck.Name, ck.URL).CurrentBalanceOfPkscript(tick, pkscript)
	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusOK, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}

	pbytes, _ := base64.StdEncoding.DecodeString(ck.Commitment)
	point := &verkle.Point{}
	point.SetBytes(pbytes)

	ok, err := apis.VerifyCurrentBalanceOfPkscript(point, tick, pkscript, balance)
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
