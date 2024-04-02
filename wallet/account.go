package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/utils"
)

func (a *Account) init(wallet *Wallet) {
	a.wallet = wallet
	a.ChainType = constant.BTC
	a.accountType = constant.AccountTypeUndefined
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
	path := fmt.Sprintf("m/%d'/%d'/%d'/%d/%d", constant.Purpose, constant.BTC, constant.Zero, constant.Zero, uint32(w.sep0005AccountCount))
	key, err := newKeyBySeed(seed, Path(path))
	if err != nil {
		return err
	}
	a.Key = key
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return err
	}
	a.privateKey = crypto.FromECDSA(privateKey.ToECDSA())
	a.accountType = constant.AccountTypeSEP0005
	a.sep0005DerivationPath = path
	a.generatePubAddr()
	w.sep0005AccountCount++
	a.active = true
	return nil
}

func (a *Account) generatePubAddr() {
	switch a.ChainType {
	case constant.BTC:
		a.publicKey = utils.KeyToBtcAddress(a.Key)
	case constant.ETH:
		a.publicKey = utils.KeyTo0xAddress(a.Key)
	}
}

// Type returns account type.
func (a *Account) Type() uint16 {
	return a.accountType
}

// IsOwnAccount checks if current account is an own account, i.e. of type generated, random or watching.
func (a *Account) IsOwnAccount() bool {
	switch a.accountType {
	case constant.AccountTypeSEP0005:
		return true
	case constant.AccountTypeRandom:
		return true
	case constant.AccountTypeWatching:
		return true
	}

	return false
}

// IsAddressBookAccount checks if current account is an address book account..
func (a *Account) IsAddressBookAccount() bool {
	if a.accountType == constant.AccountTypeAddressBook {
		return true
	}

	return false
}

// HasPrivateKey checks true if current account holds a private key.
func (a *Account) HasPrivateKey() bool {
	switch a.accountType {
	case constant.AccountTypeSEP0005:
		return true
	case constant.AccountTypeRandom:
		return true
	}

	return false
}

// Description returns description of account. Empty string is returned if no description is defined.
func (a *Account) Description() string {
	return a.desc
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

// MemoText returns the optional memo text of account. Empty string is returned if no memo text is defined.
func (a *Account) MemoText() string {
	return a.memoText
}

// SetMemoText sets memo text on account. If given memo text string is not valid a descriptive error is returned.
func (a *Account) SetMemoText(memo string) error {
	err := CheckMemoText(memo)
	if err != nil {
		return err
	}

	a.memoText = memo

	return nil
}

// MemoId returns memo id of account. If no memo id is defined for current account, the boolean return value is false.
func (a *Account) MemoId() (bool, uint64) {
	if a.memoIdSet {
		return true, a.memoId
	}

	return false, 0
}

// SetMemoId sets memo id on account.
func (a *Account) SetMemoId(memo uint64) {
	a.memoId = memo
	a.memoIdSet = true
}

// ClearMemoId clears memo id from account.
func (a *Account) ClearMemoId() {
	a.memoId = 0
	a.memoIdSet = false
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
	case constant.AccountTypeSEP0005:
		toECDSA, err := crypto.ToECDSA(a.privateKey)
		if err != nil {
			return nil
		}
		return toECDSA
	case constant.AccountTypeRandom:
	}
	return nil
}
