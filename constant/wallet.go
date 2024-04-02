package constant

const DefaultPassword string = "light-indexer"

// mnemonic language
const (
	English            = "english"
	ChineseSimplified  = "chinese_simplified"
	ChineseTraditional = "chinese_traditional"
)

// zero is deafult of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// wallet type from bip44
const (
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md#registered-coin-types
	BTC        = ZeroQuote + 0
	BTCTestnet = ZeroQuote + 1
	LTC        = ZeroQuote + 2
	DOGE       = ZeroQuote + 3
	DASH       = ZeroQuote + 5
	ETH        = ZeroQuote + 60
	ETC        = ZeroQuote + 61
	BCH        = ZeroQuote + 145
	QTUM       = ZeroQuote + 2301

	// btc token
	USDT = BTCToken + 1

	// eth token
	IOST = ETHToken + 1
	USDC = ETHToken + 2
	// Purpose
	Purpose = ZeroQuote + 44
)

const (
	AesSalt      = "Nubit DA Chain"
	SHA1Checksum = "Nubit DA Chain"
	//Purpose      uint32 = 44
)

const MasterSeedLen = 16

const Bip39SeedLen = 64

const AccountSeedLen = 32

const AccountTypeUndefined = 0   // Account type undefined
const AccountTypeSEP0005 = 1     // Account is SEP0005 derived from the BIP39 seed
const AccountTypeRandom = 2      // Account is based on a randomly generated key
const AccountTypeWatching = 3    // Account has a public key only
const AccountTypeAddressBook = 4 // Account is an ad
const AccountTypePrivateKey = 5  // Account The private key is generated
