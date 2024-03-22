package transfer

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/RiemaLabs/indexer-light/config"
	"github.com/RiemaLabs/indexer-light/getter"
	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/sirupsen/logrus"
)

type OrdTransfer struct {
	ID            uint
	InscriptionID string
	OldSatpoint   string
	NewSatpoint   string
	NewPkscript   string
	NewWallet     string
	SentAsFee     bool
	Content       string
	ContentType   string
}

func (o OrdTransfer) Offset() uint64 {
	offset, _ := strconv.ParseInt(strings.Split(o.InscriptionID, "i")[1], 10, 32)
	return uint64(offset)
}

func Verify(transfers TransferByInscription, blockHeight uint) (bool, error) {
	if len(transfers) == 0 {
		return false, errors.New("enpty tranfer data")
	}

	sort.Sort(transfers)
	chainClient, _ := getter.NewGetter(config.Config)

	batch := make(map[string]TransferByInscription)
	// find a batch of inscriptions in same txid
	f, n := 0, 1
	for n < len(transfers) {
		first := strings.Split(transfers[f].NewSatpoint, ":")[0]
		cur := strings.Split(transfers[n].NewSatpoint, ":")[0]
		if first == cur {
			n++
			continue
		}
		batch[first] = transfers[f:n]
		f = n
		n++
	}
	batch[strings.Split(transfers[f].NewSatpoint, ":")[0]] = transfers[f:n]

	hash, err := chainClient.GetBlockHash(blockHeight)
	if nil != err {
		return false, err
	}
	blockBody, err := chainClient.GetBlock2(hash)
	if nil != err {
		return false, err
	}

	for _, tx := range blockBody.Tx {
		if trans, exist := batch[tx.Txid]; exist {
			is, err := envelopVerify(chainClient, trans, tx)
			if nil != err {
				logrus.Warnf("envelopVerify failed txid: %s, err %v", tx.Txid, err)
			}
			if !is {
				return is, err
			}
		}

	}
	return true, nil
}

func envelopVerify(chainClient *getter.BitcoinOrdGetter, transfers []OrdTransfer, tx btcjson.TxRawResult) (bool, error) {
	sort.Sort(TransferByInscription(transfers))

	buf, err := hex.DecodeString(tx.Hex)
	if nil != err {
		return false, err
	}
	msgTx := new(wire.MsgTx)
	if err := msgTx.Deserialize(bytes.NewReader(buf)); nil != err {
		return false, err
	}

	inscriptions := parser.ParseInscriptionsFromTransaction(msgTx)
	id_counter := 0
	allIns := make([]Flotsam, 0, len(inscriptions))
	total_input_value := uint64(0)
	for index, tx_in := range msgTx.TxIn {
		// find oldSatpoin for privious output
		for _, obj := range transfers {
			if obj.OldSatpoint != "" && strings.Join(strings.Split(obj.OldSatpoint, ":")[:2], ":") == tx_in.PreviousOutPoint.String() {
				arr := strings.Split(obj.OldSatpoint, ":")
				satOff, _ := strconv.ParseInt(arr[2], 10, 64)
				offset := total_input_value + uint64(satOff)
				// find old inscription content && content type
				beforeIns, err := chainClient.GetAllInscriptions(arr[0])
				if nil != err {
					return false, err
				}
				body, exist := beforeIns[obj.InscriptionID]
				if !exist {
					return false, fmt.Errorf("old inscription not found: %s", obj.InscriptionID)
				}
				allIns = append(allIns, Flotsam{
					InsID:  InsFromStr(obj.InscriptionID),
					Offset: offset,
					Body:   body,
				})
			}
		}

		// parse new Inscripitons
		offset := total_input_value
		pOut, err := chainClient.GetOutput(tx_in.PreviousOutPoint.Hash.String(), int(tx_in.PreviousOutPoint.Index))
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
	for vout, out := range msgTx.TxOut {
		end := output_value + uint64(out.Value)
		for _, flot := range allIns {
			if flot.Offset >= uint64(end) {
				break
			}
			new_location = append(new_location, NewLocation{
				SentToCoinbasse: false,
				TxOut:           *out,
				Flotsam:         flot,
				NewSatpoint:     fmt.Sprintf("%s:%d:%d", flot.InsID.TxID, vout, flot.Offset),
			})
			allIns = allIns[1:]
		}
		output_value = end
	}

	p1, p2 := 0, 0
	for p1 < len(transfers) && p2 < len(new_location) {
		offset := transfers[p1].Offset()

		tmpTr := transfers[p1]
		tmpNewl := new_location[p2]

		if tmpTr.Offset() == new_location[p2].Flotsam.Offset {
			// pkscript verify ||  wallet  verify
			pkOBj, err := txscript.ParsePkScript(tmpNewl.TxOut.PkScript)
			if nil != err {
				return false, err
			}
			addr, _ := pkOBj.Address(&chaincfg.MainNetParams)

			// fmt.Println(tmpTr.NewPkscript)
			// fmt.Println(hex.EncodeToString(tmpNewl.TxOut.PkScript))

			// fmt.Println(tmpTr.NewWallet)
			// fmt.Println(addr.String())

			// fmt.Println([]byte(tmpTr.Content))
			// fmt.Println(tmpNewl.Flotsam.Body.Inscription.ContentBody)

			// fmt.Println(tmpTr.ContentType)
			// fmt.Println(hex.EncodeToString(tmpNewl.Flotsam.Body.Inscription.ContentType))

			//TODO verify content,why ?
			if tmpTr.NewPkscript != hex.EncodeToString(tmpNewl.TxOut.PkScript) || tmpTr.NewWallet != addr.String() ||
				tmpTr.ContentType != hex.EncodeToString(tmpNewl.Body.Inscription.ContentType) {
				return false, nil
			}
			// TODO verify newSatPoint
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
