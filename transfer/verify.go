package transfer

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var defaultURL = "https://frosty-serene-emerald.btc.quiknode.pro/402f5ac57de95e38c0a33d1a5e6f6c2f66709262/"

type OrdTransfer struct {
	ID            uint
	InscriptionID string
	OldSatpoint   string
	NewPkscript   string
	NewWallet     string
	SentAsFee     bool
	Content       []byte
	ContentType   string

	// BRC-20 Special
	TransferInscribeDone     bool
	TransferInscribeWallet   string
	TransferInscribePkscript string
	TransferTransferDone     bool
}

func (o OrdTransfer) Offset() uint64 {
	offset, _ := strconv.ParseInt(strings.Split(o.InscriptionID, "i")[1], 10, 32)
	return uint64(offset)
}

// The verifiableOrdTransfer.
// For the verification.
type VerifiableOrdTransfer struct {
	ordTransfer  OrdTransfer
	satPointPath []string
}

func Verify(transfers TransferByInscription, blockHeight uint) (bool, error) {
	if len(transfers) == 0 {
		return false, errors.New("enpty tranfer data")
	}

	sort.Sort(transfers)

	batch := make(map[string]TransferByInscription)
	// find a batch of inscriptions with same txid
	f, n := 0, 1
	for n < len(transfers) {
		first := strings.Split(transfers[f].ordTransfer.InscriptionID, "i")[0]
		cur := strings.Split(transfers[n].ordTransfer.InscriptionID, "i")[0]
		if first == cur {
			n++
			continue
		}
		batch[first] = transfers[f:n]
		f = n
		n++
	}
	batch[strings.Split(transfers[f].ordTransfer.InscriptionID, "i")[0]] = transfers[f:n]

	client := NewHttpGetter(defaultURL, "", "")
	hash, err := client.GetBlockHash(blockHeight)
	if nil != err {
		return false, err
	}
	blockBody, err := client.GetBlock1(hash)
	if nil != err {
		return false, err
	}
	txids := blockBody.Tx

	for txid, trs := range batch {
		if slices.Contains(txids, txid) {
			// new inscription verify
			rawtx, err := client.GetRawTransaction(txid)
			if nil != err {
				return false, err
			}
			if is, err := NewTransferVerify(trs, *rawtx); !is {
				return is, err
			}
		}
		if is, err := oldTransferVerify(trs); !is {
			return is, err
		}
	}
	return true, nil
}

func NewTransferVerify(transfers []VerifiableOrdTransfer, tx btcjson.TxRawResult) (bool, error) {
	buf, err := hex.DecodeString(tx.Hex)
	if nil != err {
		return false, err
	}
	msgTx := new(wire.MsgTx)
	if err := msgTx.Deserialize(bytes.NewReader(buf)); nil != err {
		return false, err
	}

	inscriptions := parser.ParseInscriptionsFromTransaction(msgTx)
	if len(transfers) > len(inscriptions) {
		return false, nil
	}
	id_counter := 0
	allIns := make([]Flotsam, 0, len(inscriptions))
	total_input_value := uint64(0)
	for index, tx_in := range msgTx.TxIn {
		offset := total_input_value
		pOut, err := DefaultBitcoinClient.GetOutput(tx_in.PreviousOutPoint.Hash.String(), int(tx_in.PreviousOutPoint.Index))
		if nil != err {
			return false, err
		}
		total_input_value += uint64(pOut.Value * math.Pow10(8))
		for _, ii := range inscriptions {
			if index == int(ii.TxInIndex) {
				allIns = append(allIns, Flotsam{
					InsID: InscriptionID{
						tx.Txid,
						id_counter,
					},
					Offset: offset, // TODO parser lib missing pointer, So we default inscription first sat at outpoint
					Body:   ii,
				})
				id_counter++
			}
		}
	}
	sort.Sort(ArrayFloatsam(allIns))

	new_location := make([]NewLocation, 0)
	output_value := uint64(0)
	for _, out := range msgTx.TxOut {
		end := output_value + uint64(out.Value)
		for i, flot := range allIns {
			if flot.Offset >= uint64(end) {
				allIns = allIns[i:]
				break
			}
			new_location = append(new_location, NewLocation{
				SentToCoinbasse: false,
				TxOut:           *out,
				Flotsam:         flot,
			})
		}
		output_value = end
	}

	p1, p2 := 0, 0
	for p1 < len(transfers) && p2 < len(new_location) {
		offset := transfers[p1].ordTransfer.Offset()

		tmpTr := transfers[p1].ordTransfer
		tmpNewl := new_location[p2]

		if tmpTr.Offset() == new_location[p2].Flotsam.Offset {
			// pkscript verify ||  wallet  verify
			pkOBj, err := txscript.ParsePkScript(tmpNewl.TxOut.PkScript)
			if nil != err {
				return false, err
			}
			addr, _ := pkOBj.Address(&chaincfg.MainNetParams)

			if tmpTr.NewPkscript != string(tmpNewl.TxOut.PkScript) || tmpTr.NewWallet != addr.String() ||
				!bytes.Equal(tmpTr.Content, tmpNewl.Flotsam.Body.Inscription.ContentBody) ||
				!bytes.Equal([]byte(tmpTr.ContentType), tmpNewl.Body.Inscription.ContentType) {
				return false, nil
			}
			p1++
			p2++
		} else if offset > inscriptions[p2].TxInOffset {
			p2++
		} else {
			p1++
		}
	}
	if p1 < len(transfers) {
		return false, nil
	}

	return true, nil
}

func oldTransferVerify(transfers []VerifiableOrdTransfer) (bool, error) {
	txid := strings.Split(transfers[0].ordTransfer.InscriptionID, "i")[0]
	rawTx, err := DefaultBitcoinClient.GetRawTransaction(txid)
	if nil != err {
		return false, err
	}

	buf, err := hex.DecodeString(rawTx.Hex)
	if nil != err {
		return false, err
	}
	msgTx := new(wire.MsgTx)
	if err := msgTx.Deserialize(bytes.NewReader(buf)); nil != err {
		return false, err
	}

	inscriptions := parser.ParseInscriptionsFromTransaction(msgTx)
	if len(transfers) > len(inscriptions) {
		return false, nil
	}

	id_counter := 0
	allIns := make([]Flotsam, 0, len(inscriptions))
	total_input_value := uint64(0)
	for index, tx_in := range msgTx.TxIn {
		offset := total_input_value
		pOut, err := DefaultBitcoinClient.GetOutput(tx_in.PreviousOutPoint.Hash.String(), int(tx_in.PreviousOutPoint.Index))
		if nil != err {
			return false, err
		}
		total_input_value += uint64(pOut.Value * math.Pow10(8))
		for _, ii := range inscriptions {
			if index == int(ii.TxInIndex) {
				allIns = append(allIns, Flotsam{
					InsID: InscriptionID{
						txid,
						id_counter,
					},
					Offset: offset, // TODO parser lib missing pointer, So we default inscription first sat at outpoint
					Body:   ii,
				})
				id_counter++
			}
		}
	}
	sort.Sort(ArrayFloatsam(allIns))
	new_location := make([]NewLocation, 0)
	output_value := uint64(0)
	for _, out := range msgTx.TxOut {
		end := output_value + uint64(out.Value)
		for i, flot := range allIns {
			if flot.Offset >= uint64(end) {
				allIns = allIns[i:]
				break
			}
			new_location = append(new_location, NewLocation{
				SentToCoinbasse: false,
				TxOut:           *out,
				Flotsam:         flot,
			})
		}
		output_value = end
	}

	p1, p2 := 0, 0
	for p1 < len(transfers) && p2 < len(new_location) {
		offset := transfers[p1].ordTransfer.Offset()

		tmpTr := transfers[p1].ordTransfer
		tmpNewl := new_location[p2]

		if tmpTr.Offset() == new_location[p2].Flotsam.Offset {
			// pkscript verify ||  wallet  verify
			pkOBj, err := txscript.ParsePkScript(tmpNewl.TxOut.PkScript)
			if nil != err {
				return false, err
			}
			addr, _ := pkOBj.Address(&chaincfg.MainNetParams)

			if tmpTr.NewPkscript != string(tmpNewl.TxOut.PkScript) || tmpTr.NewWallet != addr.String() ||
				!bytes.Equal(tmpTr.Content, tmpNewl.Flotsam.Body.Inscription.ContentBody) ||
				!bytes.Equal([]byte(tmpTr.ContentType), tmpNewl.Body.Inscription.ContentType) {
				return false, nil
			}
			p1++
			p2++
		} else if offset > inscriptions[p2].TxInOffset {
			p2++
		} else {
			p1++
		}
	}
	if p1 < len(transfers) {
		return false, nil
	}
	// TODO verify satpoint path

	return true, nil
}
