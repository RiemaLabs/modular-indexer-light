package transfer

import (
	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/wire"
)

type NewLocation struct {
	SentToCoinbasse bool
	TxOut           wire.TxOut
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
