package apis

import (
	"github.com/RiemaLabs/modular-indexer-light/transfer"
)

type Brc20VerifiableLightGetCurrentBalanceOfWalletRequest struct {
	Tick     string `json:"tick"`
	Pkscript string `json:"pkscript"`
}

type Brc20VerifiableLightGetCurrentBalanceOfWalletResponse struct {
	Result      string `json:"result"`
	BlockHeight int    `json:"blockHeight"`
}

type Brc20VerifiableLightCheckpointsResponse struct {
	Checkpoints []struct {
		CheckpointHash   string `json:"checkpointHash"`
		SubmissionMethod string `json:"submissionMethod"`
		IndexerS3URL     string `json:"indexerS3URL"`
		IndexerAddress   string `json:"indexerAddress"`
		TransactionID    string `json:"transactionID"`
		Checkpoint       struct {
		} `json:"checkpoint"`
	} `json:"checkpoints"`
}

type Brc20VerifiableLightLastCheckpointResponse struct {
	CheckpointHash   string `json:"checkpointHash"`
	SubmissionMethod string `json:"submissionMethod"`
	IndexerS3URL     string `json:"indexerS3URL"`
	IndexerAddress   string `json:"indexerAddress"`
	TransactionID    string `json:"transactionID"`
	Checkpoint       struct {
	} `json:"checkpoint"`
}

type Brc20VerifiableLightStateResponse struct {
	State string `json:"state"`
}

type Brc20VerifiableLightTransferVerifyRequest struct {
	BlockHeight uint                   `json:"block_height"`
	Transfers   []transfer.OrdTransfer `json:"transfers"`
}

func (o Brc20VerifiableLightTransferVerifyRequest) Check() (bool, string) {
	if o.BlockHeight <= 0 {
		return false, "invalid block_height"
	}

	for _, tr := range o.Transfers {
		if is, msg := tr.Check(); !is {
			return is, msg
		}
	}

	return true, ""
}

type Brc20VerifiableLightTransferVerifyResponse struct {
	Result bool  `json:"result"`
	Error  error `json:"error"`
}
