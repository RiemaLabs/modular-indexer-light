package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/crypto"
)

func EcdsaToPrivateStr(ecd *ecdsa.PrivateKey) string {
	PirKeyByte := crypto.FromECDSA(ecd)
	return hex.EncodeToString(PirKeyByte)
}

func KeyTo0xAddress(key *hdkeychain.ExtendedKey) string {
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return ""
	}
	publicKey := privateKey.ToECDSA().Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA).String()
}

func KeyToBtcAddress(key *hdkeychain.ExtendedKey) string {
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return ""
	}
	_, pub := btcec.PrivKeyFromBytes(PrivateStrToByte(EcdsaToPrivateStr(privateKey.ToECDSA())))
	taproot, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pub)), &chaincfg.TestNet3Params)
	if err != nil {
		return ""
	}
	return taproot.EncodeAddress()
}

func PrivateStrToBtcAddress(private string) string {
	_, pub := btcec.PrivKeyFromBytes(PrivateStrToByte(private))
	taproot, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pub)), &chaincfg.TestNet3Params)
	if err != nil {
		return ""
	}
	return taproot.EncodeAddress()
}

func PrivateStrToEcdsa(private string) *ecdsa.PrivateKey {
	toECDSA, err := crypto.HexToECDSA(private)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return nil
	}
	return toECDSA
}

func PrivateStrToByte(private string) []byte {
	ecd := PrivateStrToEcdsa(private)
	ecd.Public()
	return crypto.FromECDSA(ecd)
}
