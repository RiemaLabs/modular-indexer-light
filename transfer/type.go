package transfer

import (
	"strconv"
	"strings"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/wire"
)

type NewLocation struct {
	SentToCoinbasse bool
	TxOut           wire.TxOut
	NewSatpoint     string
	Flotsam
}

type Flotsam struct {
	InsID  InscriptionID
	Offset uint64
	Body   *parser.TransactionInscription
}

type InscriptionID struct {
	TxID  string
	Index int
}

func InsFromStr(in string) InscriptionID {
	arr := strings.Split(in, "i")
	index, _ := strconv.ParseInt(arr[1], 10, 32)
	return InscriptionID{
		TxID:  arr[0],
		Index: int(index),
	}
}
