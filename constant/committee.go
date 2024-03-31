package constant

const (
	// Committee Indexer

	LatestStateProof         = "/v1/brc20_verifiable/latest_state_proof"
	BlockHeight              = "/v1/brc20_verifiable/block_height"
	CurrentBalanceOfWallet   = "/v1/brc20_verifiable/current_balance_of_wallet"
	CurrentBalanceOfPkscript = "/v1/brc20_verifiable/current_balance_of_pkscript"
)

const (
	// Light Indexer

	LightCurrentBalanceOfPkscript = "/brc20_verifiable/light/current_balance_of_pkscript"
	LightCurrentBalanceOfWallet   = "/brc20_verifiable/light/current_balance_of_wallet"
	LightCurrentCheckpoints       = "/brc20_verifiable/light/checkpoints"
	LightLastCheckpoint           = "/brc20_verifiable/light/last_checkpoint"
	LightBlockHeight              = "/brc20_verifiable/light/block_height"
	LightState                    = "/v1/brc20_verifiable/light/state"
)
