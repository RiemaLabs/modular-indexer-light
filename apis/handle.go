package apis

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/indexer"
	"github.com/RiemaLabs/modular-indexer-light/verify"
	"github.com/ethereum/go-verkle"
	"github.com/gin-gonic/gin"
)

func CheckState() gin.HandlerFunc {
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
		c.Next()
	}
}

func GetcurrentBalanceOfWallet(c *gin.Context) {
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

		preProofByte, err := base64.StdEncoding.DecodeString(*balance.Proof)
		if err != nil {
			return
		}
		preVProof := &verkle.VerkleProof{}
		err = preVProof.UnmarshalJSON(preProofByte)
		if err != nil {
			return
		}

		preProof, err := verkle.DeserializeProof(preVProof, nil)
		if err != nil {
			return
		}

		err = verify.VerifyProof(preProof, prePoint)
		if err != nil {
			c.JSON(http.StatusOK, "")
			return
		}
		h, _ := strconv.Atoi(verify.DefiniteState.PostCheckpoint.Height)
		c.JSON(http.StatusOK, Brc20VerifiableLightGetCurrentBalanceOfWalletResponse{
			Result:      balance.Result.AvailableBalance,
			BlockHeight: h,
		})
	}
}
