package wallet

const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// See: https://github.com/satoshilabs/slips/blob/master/slip-0044.md#registered-coin-types
const (
	BTC        = ZeroQuote + 0
	BTCTestnet = ZeroQuote + 1
	LTC        = ZeroQuote + 2
	DOGE       = ZeroQuote + 3
	DASH       = ZeroQuote + 5
	Purpose    = ZeroQuote + 44
	ETH        = ZeroQuote + 60
	ETC        = ZeroQuote + 61
	BCH        = ZeroQuote + 145
	QTUM       = ZeroQuote + 2301

	// BTC-based tokens.

	USDT = BTCToken + 1

	// ETH-based token.

	IOST = ETHToken + 1
	USDC = ETHToken + 2
)

const (
	AesSalt      = "Nubit DA Chain"
	SHA1Checksum = "Nubit DA Chain"
)

const MasterSeedLen = 16

const Bip39SeedLen = 64

const AccountSeedLen = 32

const (
	AccountTypeUndefined = iota // undefined
	AccountTypeSEP0005          // SEP0005 derived from the BIP39 seed
	AccountTypeRandom           // based on a randomly generated key
	AccountTypeWatching         // has a public key only
	AccountTypeAddressBook
	AccountTypePrivateKey // private key is generated
)
