package wallet

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strings"

	bip39 "github.com/tyler-smith/go-bip39"

	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/log"
	"github.com/RiemaLabs/modular-indexer-light/utils"
	sdk "github.com/RiemaLabs/nubit-da-sdk/utils"
)

// NewWallet creates a new empty wallet, encrypted with given password.
// Each new wallet has an associated encrypted 256 bit entropy, which is the source for the mnemonic words list,
// i.e. the mnemonic word list is defined when a new wallet is created.
func NewWallet(password *string) *Wallet {
	wallet := new(Wallet)
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		panic(err)
	}
	key := deriveAesKey(password)
	wallet.encryptMasterSeed(entropy, key)
	return wallet
}

// ImportBinary creates a new wallet from an exported binary serialization of the wallet content.
// This method can be used to restore a wallet from a permanent storage location.
// This method panics if the build-in self test fails (see method SelfTest()).
func ImportBinary(buf []byte) (w *Wallet, err error) {
	w = new(Wallet)
	err = w.readFromBufferCompressed(buf)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// ImportBase64 creates a new wallet from an exported ascii (base 64) serialization of the wallet content.
// This method can be used to restore a wallet from a permanent storage location.
// This method panics if the build-in self test fails (see method SelfTest()).
func ImportBase64(data string) (w *Wallet, err error) {
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, errors.New("base64 decoding failed: " + err.Error())
	}
	w = new(Wallet)
	err = w.readFromBufferCompressed(buf)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// Bip39Mnemonic returns mnemonic word list (24 words) associated with the current wallet.
// After creating a new wallet this word list should be presented to the user.
func (w *Wallet) Bip39Mnemonic(walletPassword *string) (words []string) {
	key := deriveAesKey(walletPassword)
	words = w.getBip39Mnemonic(key)
	return
}

func (w *Wallet) GetSeed(walletPassword *string) []byte {
	if w.bip39Seed == nil {
		return nil
	}
	wkey := w.checkPassword(walletPassword)
	if wkey == nil {
		return nil
	}
	seed := w.decryptBip39Seed(wkey)
	if seed == nil {
		return nil
	}
	return seed
}

// GenerateAccount generates a new account according to SEP-0005. The wallet password
// is required to decrypt the BIP39 seed.
// Before this method can be used, method GenerateBip39Seed() must have been called before once.
// nil is returned if GenerateBip39Seed() was not called before or the wallet password is not valid.
func (w *Wallet) GenerateAccount(walletPassword *string) *Account {
	a := w.newAccount()
	err := a.GenerateAddress(w.GetSeed(walletPassword))
	if err != nil {
		log.Error("Wallet", "func", "GenerateAccount", "GenerateAddress.err", err)
		return nil
	}
	return a
}

// GenerateBip39Seed generates the seed used for key derivation (generated accounts). This
// method mus be called before the first call to GenerateAccount().
// It uses the mnemonic word list,
// which is internally derived from the master seed (same as returned by Bip39Mnemonic()),
// and combines it with the given mnemonic password.
// The wallet password is required for decrypting and master seed and encrypting the generated
// key derivation seed (BIP39 seed).
func (w *Wallet) GenerateBip39Seed(walletPassword *string, mnemonicPassword *string) bool {
	if w.bip39Seed != nil {
		return false
	}
	key := deriveAesKey(walletPassword)
	words := w.getBip39Mnemonic(key)
	if words == nil {
		// may happen if wallet password is not correct
		return false
	}
	return w.generateBip39Seed(key, words, mnemonicPassword)
}

// Returns word list derived from stored master seed.
// key is used to decrypt the master seed.
func (w *Wallet) getBip39Mnemonic(key []byte) (words []string) {
	seed := w.decryptMasterSeed(key)
	if seed == nil {
		return
	}
	wordString, err := bip39.NewMnemonic(seed)
	if err != nil {
		panic(err)
		return
	}
	words = strings.Split(wordString, " ")
	return
}

func (w *Wallet) generateBip39Seed(key []byte, words []string, mnemonicPassword *string) bool {
	var seed, prevSeed []byte
	var err error
	mnemonic := strings.Join(words, " ")
	// paranoia mode: calculate seed 5 times and check for identical results
	// to reduce risk of generating an unreproducable seed caused by faulty hardware
	for i := 0; i < 5; i++ {
		prevSeed = seed
		seed, err = bip39.NewSeedWithErrorChecking(mnemonic, *mnemonicPassword)
		if err != nil {
			return false
		}

		if prevSeed != nil && bytes.Compare(prevSeed, seed) != 0 {
			panic("calculation error")
		}
	}
	if len(seed) != constant.Bip39SeedLen {
		panic("Unexpected length of seed")
	}
	w.bip39SeedLen = constant.Bip39SeedLen
	w.encryptBip39Seed(seed, key)
	return true
}

func (w *Wallet) newAccount() *Account {
	for i, _ := range w.accounts {
		a := &w.accounts[i]
		if !a.active {
			a.init(w)
			return a
		}
	}
	w.accounts = append(w.accounts, Account{})
	a := &w.accounts[len(w.accounts)-1]
	a.init(w)
	return a
}

func (w *Wallet) encryptBip39Seed(seed, key []byte) {
	w.bip39Seed = encryptWithCheckSum(seed, key)
}

func (w *Wallet) decryptBip39Seed(key []byte) []byte {
	return decryptWithCheckSum(w.bip39Seed, constant.Bip39SeedLen, key)
}

func (w *Wallet) encryptMasterSeed(seed, key []byte) {
	w.masterSeed = encryptWithCheckSum(seed, key)
}

func (w *Wallet) decryptMasterSeed(key []byte) []byte {
	return decryptWithCheckSum(w.masterSeed, constant.MasterSeedLen, key)
}

// Checks if given wallet password is valid.
// Returns derived AES key on success else nil
func (w *Wallet) checkPassword(walletPassword *string) []byte {
	key := deriveAesKey(walletPassword)
	seed := w.decryptMasterSeed(key)
	if seed != nil {
		return key
	}
	return nil
}

// CheckPassword checks if given wallet password is valid.
func (w *Wallet) CheckPassword(walletPassword *string) bool {
	key := w.checkPassword(walletPassword)
	if key != nil {
		return true
	}
	return false
}

// ExportBinary creates a binary serialization of the wallet content, e.g. for permanent storage of the wallet on disk.
func (w *Wallet) ExportBinary() []byte {
	return w.writeToBufferCompressed()
}

// ExportBase64 creates an ascii (base 64) serialization of the wallet content, e.g. for permanent storage of the wallet on disk.
func (w *Wallet) ExportBase64() string {
	buf := w.writeToBufferCompressed()

	if buf == nil {
		return ""
	} else {
		return base64.StdEncoding.EncodeToString(buf)
	}
}

// Clears all accounts.
func (w *Wallet) clearAccounts() {
	for _, a := range w.accounts {
		a.active = false
		a.privateKey = nil
		a.publicKey = ""
		a.sep0005DerivationPath = ""
	}
}

// Description returns the optional wallet description.
func (w *Wallet) Description() string {
	return w.desc
}

// SetDescription sets wallet description. Error is returned if given string does not pass the description check.
func (w *Wallet) SetDescription(desc string) error {
	err := CheckDescription(desc)
	if err != nil {
		return err
	}

	w.desc = desc

	return nil
}

func (w *Wallet) AddPirKeyAccount(PrivateKey string, walletPassword *string) *Account {
	key := w.checkPassword(walletPassword)
	if key == nil {
		return nil
	}
	a := w.newAccount()
	a.privateKey = sdk.PrivateStrToByte(PrivateKey)
	a.accountType = constant.AccountTypePrivateKey
	a.publicKey = utils.PrivateStrToBtcAddress(PrivateKey)
	a.active = true
	a.ChainType = constant.BTC
	return a
}

// AddRandomAccount adds a new account with given private key (seed) and returns a new Account object.
// The private key is stored encrypted.
// Application implementors should make the user aware that this type of account cannot be
// recovered with the mnemonic word list and password.
// nil is returend if the wallet password is invalid or an invald seed string was provided.
func (w *Wallet) AddRandomAccount(seed string, walletPassword *string) *Account {
	key := w.checkPassword(walletPassword)
	if key == nil {
		return nil
	}
	encSeed := []byte(seed)
	//w.encryptBip39Seed(encSeed, key)
	a := w.newAccount()
	a.GenerateAddress(encSeed)
	return a
}

// AddWatchingAccount adds a watching account and return a new Account object for it.
// Watching accounts just store the public account key.
// Watching accounts are treated as "own" accounts - in contrast to address book accounts.
// nil is returned if the given public key string is not valid.
func (w *Wallet) AddWatchingAccount(pubkey string) *Account {
	a := w.FindAccountByPublicKey(pubkey)
	if a != nil {
		return nil
	}
	a = w.newAccount()
	a.accountType = constant.AccountTypeWatching
	a.publicKey = pubkey
	a.active = true
	return a
}

// AddAddressBookAccount adds an address book account and return a new Account object for it.
// Address book accounts just store the public account key.
// Address book accounts are treated as "foreign" accounts - in contrast to watching accounts.
// nil is returned if the given public key string is not valid.
func (w *Wallet) AddAddressBookAccount(pubkey string) *Account {
	a := w.FindAccountByPublicKey(pubkey)

	if a != nil {
		return nil
	}

	a = w.newAccount()

	a.accountType = constant.AccountTypeAddressBook
	a.publicKey = pubkey

	a.active = true

	return a
}

// DeleteAccount deletes given account. false is returned if given account does not belong to current wallet object.
func (w *Wallet) DeleteAccount(acc *Account) bool {
	if acc.wallet == w {
		nAccount := []Account{}
		for i, _ := range w.Accounts() {
			if w.accounts[i].publicKey != acc.publicKey {
				nAccount = append(nAccount, w.accounts[i])
			}
		}
		w.accounts = nAccount
		return true
	}

	return false
}

// FindAccountByPublicKey returns account object for given public account key.
// If not matching account is found, nil is returned.
func (w *Wallet) FindAccountByPublicKey(pubkey string) *Account {

	for i, _ := range w.accounts {
		if w.accounts[i].active && w.accounts[i].publicKey == pubkey {
			return &w.accounts[i]
		}
	}

	return nil
}

// FindAccountByDescription returns first account matching given description string.
// Matching is performed case insensitive on sub string level..
// If not matching account is found, nil is returned.
func (w *Wallet) FindAccountByDescription(desc string) *Account {
	desc = strings.ToLower(desc)

	for i, _ := range w.accounts {
		if w.accounts[i].active {
			if strings.Contains(strings.ToLower(w.accounts[i].desc), desc) {
				return &w.accounts[i]
			}
		}
	}

	return nil
}

// Accounts returns a slice containing all "own" accounts if current wallet, i.e.
// all but address book accounts.
func (w *Wallet) Accounts() []*Account {
	accounts := make([]*Account, 0, len(w.accounts))

	for i, _ := range w.accounts {
		if w.accounts[i].active && w.accounts[i].IsOwnAccount() {
			accounts = append(accounts, &w.accounts[i])
		}
	}

	return accounts
}

// AddressBook returns a slice containing all address book accounts.
func (w *Wallet) AddressBook() []*Account {
	accounts := make([]*Account, 0, len(w.accounts))

	for i, _ := range w.accounts {
		if w.accounts[i].active {
			accounts = append(accounts, &w.accounts[i])
		}
	}

	return accounts
}
