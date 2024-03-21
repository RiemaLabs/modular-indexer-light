package constant

const (
	ApiStateActive = "active"

	ApiStateSync = "sync"

	ApiStateError = "error"

	ApiStateInit = "init"

	ApiStateLoading = "loading"
)

const (
	ProvideS3Name       = "S3"
	ProvideDaName       = "Da"
	DefaultMetaProtocol = "brc-20"
)

var (
	ConfigFileName = "./config.json"
	ApiState       = ApiStateActive
)