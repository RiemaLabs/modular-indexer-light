package utils

import (
	"crypto/ecdsa"

	sdk "github.com/RiemaLabs/nubit-da-sdk/utils"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func EcdsaToPrivateStr(ecd *ecdsa.PrivateKey) string {
	return sdk.EcdsaToPrivateStr(ecd)

}

func KeyToBtcAddress(key *hdkeychain.ExtendedKey) string {
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return ""
	}
	_, pub := btcec.PrivKeyFromBytes(sdk.PrivateStrToByte(EcdsaToPrivateStr(privateKey.ToECDSA())))
	taproot, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pub)), &chaincfg.TestNet3Params)
	if err != nil {
		return ""
	}
	return taproot.EncodeAddress()
}

func PrivateStrToBtcAddress(private string) string {
	_, pub := btcec.PrivKeyFromBytes(sdk.PrivateStrToByte(private))
	taproot, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pub)), &chaincfg.TestNet3Params)
	if err != nil {
		return ""
	}
	return taproot.EncodeAddress()
}
