package transfer

import (
	"strconv"
	"strings"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/wire"
)

type NewLocation struct {
	SentToCoinbase bool
	TxOut          wire.TxOut
	NewSatPoint    string
	Flotsam
}

type Flotsam struct {
	InsID  InscriptionID
	Offset uint64
	Body   *parser.TransactionInscription
}

type ByOffset []Flotsam

func (a ByOffset) Len() int {
	return len(a)
}

func (a ByOffset) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByOffset) Less(i, j int) bool {
	return a[i].Offset < a[j].Offset
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
