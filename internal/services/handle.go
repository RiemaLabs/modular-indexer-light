package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/RiemaLabs/modular-indexer-committee/apis"
	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/ethereum/go-verkle"
	"github.com/gin-gonic/gin"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/committee"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/states"
)

const errMsgBalanceNotFound = "proof of absence"

func CheckState(c *gin.Context) {
	switch s := states.Status(states.S.Status.Load()); s {
	case states.StatusVerified:
		c.Next()
	default:
		c.AbortWithStatusJSON(http.StatusForbidden, s.String())
	}
}

func HandleGetCurrentBalanceOfWallet(c *gin.Context) {
	balance, err := GetCurrentBalanceOfWallet(
		states.S.CurrentFirstCheckpoint().Checkpoint,
		c.DefaultQuery("tick", ""),
		c.DefaultQuery("wallet", ""),
	)
	if err != nil {
		msg := err.Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.Brc20VerifiableCurrentBalanceOfWalletResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func HandleGetCurrentBalanceOfPkscript(c *gin.Context) {
	balance, err := GetCurrentBalanceOfPkscript(
		states.S.CurrentFirstCheckpoint().Checkpoint,
		c.DefaultQuery("tick", ""),
		c.DefaultQuery("pkscript", ""),
	)
	if err != nil {
		msg := err.Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func GetCurrentBalanceOfWallet(ck *checkpoint.Checkpoint, tick, wallet string) (*apis.Brc20VerifiableCurrentBalanceOfWalletResponse, error) {
	cl, err := committee.New(ck.URL)
	if err != nil {
		logs.Error.Printf("Create committee client failed: ck=%+v, tick=%s, wallet=%s, err=%v", ck, tick, wallet, err)
		return nil, err
	}

	balance, err := cl.CurrentBalanceOfWallet(context.Background(), tick, wallet)
	if err != nil {
		logs.Error.Printf("Get balance of wallet error: ck=%+v, tick=%s, wallet=%s, err=%v", ck, tick, wallet, err)
		return nil, err
	}

	commitmentBytes, _ := base64.StdEncoding.DecodeString(ck.Commitment)
	var point verkle.Point
	_ = point.SetBytes(commitmentBytes)

	ok, err := apis.VerifyCurrentBalanceOfWallet(&point, tick, wallet, balance)
	if err != nil {
		if strings.HasPrefix(err.Error(), errMsgBalanceNotFound) {
			return balance, nil
		}
		logs.Error.Printf("Verify balance of wallet error: ck=%+v, tick=%s, wallet=%s, balance=%+v, err=%v", ck, tick, wallet, balance, err)
		return nil, err
	}

	if !ok {
		logs.Error.Printf("Verify balance of wallet not OK: ck=%+v, tick=%s, wallet=%s, balance=%+v, err=%v", ck, tick, wallet, balance, err)
		return nil, fmt.Errorf("verify balance of wallet not OK")
	}

	return balance, nil
}

func GetCurrentBalanceOfPkscript(ck *checkpoint.Checkpoint, tick, pkscript string) (*apis.Brc20VerifiableCurrentBalanceOfPkscriptResponse, error) {
	cl, err := committee.New(ck.URL)
	if err != nil {
		logs.Error.Printf("Create committee client failed: ck=%+v, tick=%s, pkscript=%s, err=%v", ck, tick, pkscript, err)
		return nil, err
	}

	balance, err := cl.CurrentBalanceOfPkscript(context.Background(), tick, pkscript)
	if err != nil {
		logs.Error.Printf("Get balance of PkScript error: ck=%+v, tick=%s, pkscript=%s, err=%v", ck, tick, pkscript, err)
		return nil, err
	}

	commitmentBytes, _ := base64.StdEncoding.DecodeString(ck.Commitment)
	var point verkle.Point
	_ = point.SetBytes(commitmentBytes)

	ok, err := apis.VerifyCurrentBalanceOfPkscript(&point, tick, pkscript, balance)
	if err != nil {
		if strings.HasPrefix(err.Error(), errMsgBalanceNotFound) {
			return balance, nil
		}
		logs.Error.Printf("Verify balance of PkScript error: ck=%+v, tick=%s, pkscript=%s, balance=%+v, err=%v", ck, tick, pkscript, balance, err)
		return nil, err
	}

	if !ok {
		logs.Error.Printf("Verify balance of PkScript not OK: ck=%+v, tick=%s, pkscript=%s, balance=%+v, err=%v", ck, tick, pkscript, balance, err)
		return nil, fmt.Errorf("verify balance of PkScript not OK")
	}

	return balance, nil
}
