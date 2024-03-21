package constant

const (
	// committee api

	StateDiff        = "/v1/brc20_verifiable/latest_state_proof"
	BlockHigh        = "/v1/brc20_verifiable/block_height"
	Balance          = "/v1/brc20_verifiable/get_current_balance_of_wallet"
	BalanceOfPkscrip = "/v1/brc20_verifiable/current_balance_of_pkscript"
)

const (
	// Light Indexer

	LightBalance        = "/v1/brc20_verifiable/light/get_current_balance_of_wallet"
	LightCheckpoint     = "/v1/brc20_verifiable/light/checkpoints"
	LightLastCheckpoint = "/v1/brc20_verifiable/light/last_checkpoint"
	LightBlockHigh      = "/v1/brc20_verifiable/light/block_height"
	LightState          = "/v1/brc20_verifiable/light/state"
)
