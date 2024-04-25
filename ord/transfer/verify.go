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

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/sirupsen/logrus"

	"github.com/RiemaLabs/modular-indexer-committee/ord/getter"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
	client "github.com/RiemaLabs/modular-indexer-light/ord/getter"
)

// TODO: High.
// Retrieve OrdTransfer directly from the Bitcoin block using the mapping between oldSatPoint and newSatPoint,
// bypassing the need for OrdTransfers verification.

func CheckOrdTransfer(o getter.OrdTransfer) (bool, string) {
	if len(o.InscriptionID) < 66 || len(strings.Split(o.InscriptionID, "i")) != 2 {
		return false, fmt.Sprintf("invalid inscription_id: %s", o.InscriptionID)
	}

	if len(o.NewSatpoint) < 68 || len(strings.Split(o.NewSatpoint, ":")) != 3 {
		return false, fmt.Sprintf("invalid new_satpoint: %s", o.NewSatpoint)
	}

	if len(o.OldSatpoint) > 0 {
		if len(o.OldSatpoint) < 68 || len(strings.Split(o.OldSatpoint, ":")) != 3 {
			return false, fmt.Sprintf("invalid old_satpoint: %s", o.OldSatpoint)
		}
	}
	if len(o.NewPkscript) <= 0 || len(o.NewWallet) <= 0 {
		return false, "invalid new_pkscript or new_wallet"
	}
	return true, ""
}

func OffsetSat(o getter.OrdTransfer) uint64 {
	offset, _ := strconv.ParseInt(strings.Split(o.InscriptionID, "i")[1], 10, 32)
	return uint64(offset)
}

func VerifyOrdTransfer(transfers TransferByInscription, blockHeight uint) (bool, error) {
	if len(transfers) == 0 {
		return false, errors.New("enpty tranfer data")
	}

	sort.Sort(transfers)
	chainClient, _ := client.NewBitcoinOrdGetter(configs.C.Verification.BitcoinRPC)

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
	if err != nil {
		return false, err
	}
	blockBody, err := chainClient.GetBlockDetail(hash)
	if err != nil {
		return false, err
	}

	for _, tx := range blockBody.Tx {
		if trans, exist := batch[tx.Txid]; exist {
			is, err := VerifyEnvelop(chainClient, trans, tx)
			if err != nil {
				logrus.Warnf("envelopVerify failed txid: %s, err %v", tx.Txid, err)
			}
			if !is {
				return is, err
			}
		}

	}
	return true, nil
}

func VerifyEnvelop(chainClient *client.BitcoinOrdGetter, transfers []getter.OrdTransfer, tx btcjson.TxRawResult) (bool, error) {
	sort.Sort(TransferByInscription(transfers))

	buf, err := hex.DecodeString(tx.Hex)
	if err != nil {
		return false, err
	}
	msgTx := new(wire.MsgTx)
	if err := msgTx.Deserialize(bytes.NewReader(buf)); err != nil {
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
				if err != nil {
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
		if err != nil {
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
					// TODO: Low.
					// The parser library is missing the functionality to parse the pointer.
					// So we default to setting the first inscription at the outpoint.
					Offset: offset,
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
		offset := OffsetSat(transfers[p1])

		tmpTr := transfers[p1]
		tmpNewl := new_location[p2]

		if OffsetSat(tmpTr) == new_location[p2].Flotsam.Offset {
			// Verify pkscript
			pkOBj, err := txscript.ParsePkScript(tmpNewl.TxOut.PkScript)
			if err != nil {
				return false, err
			}
			addr, _ := pkOBj.Address(&chaincfg.MainNetParams)

			//TODO: Low. Verify content.
			if string(tmpTr.NewPkscript) != hex.EncodeToString(tmpNewl.TxOut.PkScript) || string(tmpTr.NewWallet) != addr.String() ||
				tmpTr.ContentType != hex.EncodeToString(tmpNewl.Body.Inscription.ContentType) {
				return false, nil
			}
			// TODO: Low. Verify newSatPoint.
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
