package apis

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
