package ordi

import (
	"strconv"
	"strings"

	"github.com/RiemaLabs/modular-indexer-committee/ord/getter"
	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/wire"
)

type ByNewSatpoint []getter.OrdTransfer

func (a ByNewSatpoint) Len() int           { return len(a) }
func (a ByNewSatpoint) Less(i, j int) bool { return a[i].NewSatpoint < a[j].NewSatpoint }
func (a ByNewSatpoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type NewLocation struct {
	SentToCoinbase bool
	TxOut          *wire.TxOut
	NewSatPoint    string
	Flotsam
}

type Flotsam struct {
	InsID  InscriptionID
	Offset uint64
	Body   *parser.TransactionInscription
}

type ByOffset []Flotsam

func (a ByOffset) Len() int           { return len(a) }
func (a ByOffset) Less(i, j int) bool { return a[i].Offset < a[j].Offset }
func (a ByOffset) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type InscriptionID struct {
	TxID  string
	Index int
}

func NewInscriptionID(raw string) InscriptionID {
	txID, indexRaw, _ := strings.Cut(raw, "i")
	index, _ := strconv.ParseInt(indexRaw, 10, 32)
	return InscriptionID{
		TxID:  txID,
		Index: int(index),
	}
}

func FromRawSatpoint(rawSatpoint string) (txID string, index, offset uint64) {
	raws := strings.SplitN(rawSatpoint, ":", 3)
	txID = raws[0]
	if len(raws) >= 2 {
		index, _ = strconv.ParseUint(raws[1], 10, 64)
	}
	if len(raws) >= 3 {
		offset, _ = strconv.ParseUint(raws[2], 10, 64)
	}
	return
}
