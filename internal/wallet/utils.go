package wallet

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/pbkdf2"
)

// Create AES key from given password string
func deriveAesKey(password *string) (key []byte) {
	return pbkdf2.Key([]byte(*password), []byte(AesSalt), 4096, 32, sha1.New)
}

// AES encrypts data with given key.
// Data will be padded with random data to match AES block size.
// Encrypted data is returned in newly allocated slice. Input data will not be modified.
func aesEncrypt(data, key []byte) []byte {
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	l := len(data)
	blockLen := block.BlockSize()

	blocks := l / blockLen

	padding := l % blockLen

	if padding != 0 {
		blocks += 1
	}

	buf := make([]byte, blockLen*blocks)

	if padding != 0 {
		// fill padding bytes with random numbers
		_, err = rand.Read(buf[blockLen*(blocks-1)+padding:])
		if err != nil {
			panic(err)
		}
	}

	copy(buf, data)

	itr := buf

	for i := 0; i < blocks; i++ {
		block.Encrypt(itr, itr)
		itr = itr[blockLen:]
	}

	return buf
}

// AES decypts data with given key.
// Input data slice is overwritten by decrypted data.
func aesDecrypt(data, key []byte) {
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	l := len(data)
	blockLen := block.BlockSize()
	blocks := l / blockLen

	if l%blockLen != 0 {
		return
	}

	itr := data

	for i := 0; i < blocks; i++ {
		block.Decrypt(itr, itr)
		itr = itr[blockLen:]
	}
}

// Adds a SHA1 checksum to data and encrypts it with given key.
func encryptWithCheckSum(data, key []byte) []byte {
	// build checksum using a hash
	mac := hmac.New(sha1.New, []byte(SHA1Checksum))
	_, err := mac.Write(data)
	if err != nil {
		panic(err)
	}
	sum := mac.Sum(nil)
	dataChk := make([]byte, len(data))
	copy(dataChk, data)
	dataChk = append(dataChk, sum...)
	enc := aesEncrypt(dataChk, key)
	return enc
}

// Decrypts data with given key and verifies SHA1 checksum.
// If verification fails nil is returned. Otherwise a newly
// allocated slice is returned containing the decrypted data.
// The input data is not modified.
func decryptWithCheckSum(encData []byte, resLen int, key []byte) []byte {
	buf := make([]byte, len(encData))
	copy(buf, encData)

	aesDecrypt(buf, key)

	if len(buf) < resLen {
		return nil
	}

	// check checksum using a hash
	mac := hmac.New(sha1.New, []byte(SHA1Checksum))
	_, err := mac.Write(buf[:resLen])
	if err != nil {
		panic(err)
	}
	sum := mac.Sum(nil)

	if !hmac.Equal(sum, buf[resLen:resLen+mac.Size()]) {
		return nil
	}

	return buf[:resLen]
}

// CheckDescription checks for valid wallet, account or asset description string.
// If given description is not valid returned error contains details about failed check.
func CheckDescription(s string) error {
	if len(s) > 2000 {
		return errors.New("exceeds max length (2000 characters)")
	}

	return nil
}

// CheckMemoText checks for valid transaction memo text.
// If given memo text is not valid returned error contains details about failed check.
func CheckMemoText(s string) error {
	if len(s) > 28 {
		return errors.New("exceeds max length (28 characters)")
	}

	return nil
}

func newKeyBySeed(seed []byte, path []uint32) (*hdkeychain.ExtendedKey, error) {
	var child *hdkeychain.ExtendedKey
	param := &chaincfg.TestNet3Params
	child, err := hdkeychain.NewMaster(seed, param)
	if err != nil {
		return nil, err
	}

	return newKeyByMasterKey(child, path)
}

func newKeyByMasterKey(master *hdkeychain.ExtendedKey, path []uint32) (*hdkeychain.ExtendedKey, error) {
	var err error
	child := master
	for _, p := range path {
		child, err = child.Derive(p)
		//child, err = child.Child(p)
		if err != nil {
			return nil, err
		}
	}
	return child, nil
}

// Path set to options
// example: m/44'/0'/0'/0/0
// example: m/Purpose'/CoinType'/Account'/Change/AddressIndex
func Path(path string) []uint32 {
	path = strings.TrimPrefix(path, "m/")
	paths := strings.Split(path, "/")
	if len(paths) != 5 {
		return nil
	}
	purpose := PathNumber(paths[0])
	coinType := PathNumber(paths[1])
	account := PathNumber(paths[2])
	change := PathNumber(paths[3])
	addressIndex := PathNumber(paths[4])
	return []uint32{
		purpose,
		coinType,
		account,
		change,
		addressIndex,
	}
}

// PathNumber 44' => 0x80000000 + 44
func PathNumber(str string) uint32 {
	num64, _ := strconv.ParseInt(strings.TrimSuffix(str, "'"), 10, 64)
	num := uint32(num64)
	if strings.HasSuffix(str, "'") {
		num += ZeroQuote
	}
	return num
}
