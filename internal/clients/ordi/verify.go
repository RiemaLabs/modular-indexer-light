package ordi

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/sirupsen/logrus"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/btcutl"
)

// TODO: High.
// Retrieve OrdTransfer directly from the Bitcoin block using the mapping between oldSatPoint and newSatPoint,
// bypassing the need for OrdTransfers verification.

func VerifyOrdTransfer(transfers ByNewSatpoint, blockHeight uint) error {
	if len(transfers) == 0 {
		return errors.New("empty transfer data")
	}

	transfersByID := make(map[string]ByNewSatpoint)
	for _, t := range transfers {
		txID, _, _ := strings.Cut(t.NewSatpoint, ":")
		transfersByID[txID] = append(transfersByID[txID], t)
	}
	for _, ts := range transfersByID {
		sort.Sort(ts)
	}

	hash, err := btcutl.BTC.GetBlockHash(context.Background(), blockHeight)
	if err != nil {
		return err
	}
	blockBody, err := btcutl.BTC.GetBlockDetail(context.Background(), hash)
	if err != nil {
		return err
	}

	for _, tx := range blockBody.Tx {
		trans, found := transfersByID[tx.Txid]
		if !found {
			continue
		}
		if err := VerifyEnvelop(trans, tx); err != nil {
			logrus.Warnf("Envelop verify failed: txid=%s, err=%v", tx.Txid, err)
			return err
		}
	}

	return nil
}

func VerifyEnvelop(transfers ByNewSatpoint, txRaw btcjson.TxRawResult) error {
	txRawBytes, err := hex.DecodeString(txRaw.Hex)
	if err != nil {
		return err
	}

	tx := new(wire.MsgTx)
	if err := tx.Deserialize(bytes.NewReader(txRawBytes)); err != nil {
		return err
	}

	inscriptions := parser.ParseInscriptionsFromTransaction(tx)
	idCnt := 0
	var allInscriptions ByOffset
	curOff := uint64(0)
	for index, txIn := range tx.TxIn {
		for _, transfer := range transfers {
			if transfer.OldSatpoint == "" {
				continue
			}
			if !strings.HasPrefix(transfer.OldSatpoint, txIn.PreviousOutPoint.String()) {
				continue
			}
			txID, _, offset := FromRawSatpoint(transfer.OldSatpoint)
			beforeIns, err := btcutl.BTC.GetAllInscriptions(context.Background(), txID)
			if err != nil {
				return err
			}
			body, found := beforeIns[transfer.InscriptionID]
			if !found {
				return fmt.Errorf("old inscription not found: %s", transfer.InscriptionID)
			}
			allInscriptions = append(allInscriptions, Flotsam{
				InsID:  NewInscriptionID(transfer.InscriptionID),
				Offset: curOff + offset,
				Body:   body,
			})
		}

		output, err := btcutl.BTC.GetOutput(
			context.Background(),
			txIn.PreviousOutPoint.Hash.String(),
			int(txIn.PreviousOutPoint.Index),
		)
		if err != nil {
			return err
		}
		for _, inscription := range inscriptions {
			if index != int(inscription.TxInIndex) {
				continue
			}
			allInscriptions = append(allInscriptions, Flotsam{
				InsID: InscriptionID{txRaw.Txid, idCnt},
				// TODO: Low.
				// The parser library is missing the functionality to parse the pointer.
				// So we default to setting the first inscription at the outpoint.
				Offset: curOff,
				Body:   inscription,
			})
			idCnt++
		}
		curOff += uint64(output.Value * math.Pow10(8))
	}
	sort.Sort(allInscriptions)

	var newLocation []NewLocation
	curOff = 0
	for idx, out := range tx.TxOut {
		end := curOff + uint64(out.Value)
		for _, flot := range allInscriptions {
			if flot.Offset >= end {
				break
			}
			newLocation = append(newLocation, NewLocation{
				TxOut:       out,
				Flotsam:     flot,
				NewSatPoint: fmt.Sprintf("%s:%d:%d", flot.InsID.TxID, idx, flot.Offset),
			})
		}
		curOff = end
	}

	p1, p2 := 0, 0
	for p1 < len(transfers) && p2 < len(newLocation) {
		actual := transfers[p1]
		expected := newLocation[p2]

		offset := uint64(NewInscriptionID(actual.InscriptionID).Index)
		if offset == expected.Flotsam.Offset {
			actualPkScript := string(actual.NewPkscript)
			expectedPkScript := hex.EncodeToString(expected.TxOut.PkScript)
			if actualPkScript != expectedPkScript {
				return fmt.Errorf("unmatched new PkScript: actual=%s, expected=%s", actualPkScript, expectedPkScript)
			}

			actualNewWallet := string(actual.NewWallet)
			pkscript, err := txscript.ParsePkScript(expected.TxOut.PkScript)
			if err != nil {
				return err
			}
			expectedNewAddr, _ := pkscript.Address(&chaincfg.MainNetParams)
			expectedNewWallet := expectedNewAddr.String()
			if actualNewWallet != expectedNewWallet {
				return fmt.Errorf("unmatched new wallet: actual=%s, expected=%s", actualNewWallet, expectedNewWallet)
			}

			actualContentType := actual.ContentType
			expectedContentType := hex.EncodeToString(expected.Body.Inscription.ContentType)
			if actualContentType != expectedContentType {
				return fmt.Errorf("unmatched content type: actual=%s, expected=%s", actualContentType, expectedContentType)
			}

			// TODO: Low. Verify content.

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
		return fmt.Errorf("invalid transfer: %+v", transfers[p1])
	}

	return nil
}
