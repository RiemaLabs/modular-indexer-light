package wallet

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

type Wallet struct {
	desc                string
	masterSeed          []byte
	masterSeedLen       int
	bip39Seed           []byte
	bip39SeedLen        int
	sep0005AccountCount uint16
	accounts            []Account
	assets              []*Asset
}

type Account struct {
	Key                   *hdkeychain.ExtendedKey
	ChainType             uint32
	wallet                *Wallet
	active                bool
	desc                  string
	accountType           uint16
	publicKey             string
	privateKey            []byte
	sep0005DerivationPath string
	memoText              string
	memoId                uint64
	memoIdSet             bool
}

type Asset struct {
	wallet  *Wallet
	active  bool
	desc    string
	issuer  string
	assetId string
}
