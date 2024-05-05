package wallet

import (
	"crypto/ecdsa"
	"fmt"

	sdk "github.com/RiemaLabs/nubit-da-sdk/utils"
	"github.com/btcsuite/btcd/btcec/v2"

	"github.com/RiemaLabs/modular-indexer-light/internal/utils"
)

func (a *Account) init(wallet *Wallet) {
	a.wallet = wallet
	a.ChainType = BTC
	a.accountType = AccountTypeUndefined
	a.desc = ""
	a.publicKey = ""
	a.privateKey = nil
	a.sep0005DerivationPath = ""
	a.memoText = ""
	a.memoId = 0
	a.memoIdSet = false
}

func (a *Account) GenerateAddress(seed []byte) error {
	return a.generateBip32Address(seed)
}

func (a *Account) generateBip32Address(seed []byte) error {
	w := a.wallet
	// example: m/Purpose'/CoinType'/Account'/Change/AddressIndex
	path := fmt.Sprintf("m/%d'/%d'/%d'/%d/%d", Purpose, BTC, Zero, Zero, uint32(w.sep0005AccountCount))
	key, err := newKeyBySeed(seed, Path(path))
	if err != nil {
		return err
	}
	a.Key = key
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return err
	}

	a.privateKey = sdk.FromECDSA(privateKey.ToECDSA())
	a.accountType = AccountTypeSEP0005
	a.sep0005DerivationPath = path
	a.generatePubAddr()
	w.sep0005AccountCount++
	a.active = true
	return nil
}

func (a *Account) generatePubAddr() {
	switch a.ChainType {
	case BTC:
		a.publicKey = utils.KeyToBtcAddress(a.Key)
	}
}

// Type returns account type.
func (a *Account) Type() uint16 {
	return a.accountType
}

// IsOwnAccount checks if current account is an own account, i.e. of type generated, random or watching.
func (a *Account) IsOwnAccount() bool {
	switch a.accountType {
	case AccountTypeSEP0005:
		return true
	case AccountTypeRandom:
		return true
	case AccountTypeWatching:
		return true
	}

	return false
}

// HasPrivateKey checks true if current account holds a private key.
func (a *Account) HasPrivateKey() bool {
	switch a.accountType {
	case AccountTypeSEP0005:
		return true
	case AccountTypeRandom:
		return true
	}

	return false
}

// SetDescription sets description on account. If give description string is not valid a descriptive error is returned.
func (a *Account) SetDescription(desc string) error {
	err := CheckDescription(desc)
	if err != nil {
		return err
	}

	a.desc = desc

	return nil
}

// PublicKey returns public key of account.
func (a *Account) PublicKey() string {
	if !a.active {
		panic("account not active")
	}
	return a.publicKey
}

// PrivateKey returns private key of account.
// Emptry string is returned if wallet password is not valid or current account does not hold a private key.
func (a *Account) PrivateKey(walletPassword *string) *ecdsa.PrivateKey {
	wkey := a.wallet.checkPassword(walletPassword)
	if wkey == nil {
		return nil
	}
	switch a.accountType {
	case AccountTypeSEP0005:

		pri, _ := btcec.PrivKeyFromBytes(a.privateKey)
		prv := pri.ToECDSA()
		return prv
	case AccountTypeRandom:
	}
	return nil
}
