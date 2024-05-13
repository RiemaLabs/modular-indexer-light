package btcutl

import (
	"testing"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/httputl"
)

func testClient(t *testing.T) *Client {
	cl, err := New("https://bitcoin-mainnet-archive.allthatnode.com")
	if err != nil {
		t.Fatal(err)
	}
	return cl
}

func TestBitcoinOrdGetter_GetLatestBlockHeight(t *testing.T) {
	h, err := testClient(t).GetLatestBlockHeight(httputl.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if h < 840000 {
		t.Fatal(h)
	}
}

func TestBitcoinOrdGetter_GetBlockHash(t *testing.T) {
	h, err := testClient(t).GetBlockHash(httputl.TODO(), 835161)
	if err != nil {
		t.Fatal(err)
	}
	if h != "000000000000000000021a731d2106dda997d6eaf6228252c7abdc259c1fca5e" {
		t.Fatal(h)
	}
}

func TestBitcoinOrdGetter_GetRawTransaction(t *testing.T) {
	const txID = "26a08b3ac578f1fe01bde9d0268353121f22461fcb48dc3144f1dd5210d0f8ad"
	rawTx, err := testClient(t).GetRawTransaction(httputl.TODO(), txID)
	if err != nil {
		t.Fatal(err)
	}
	if rawTx.Txid != txID {
		t.Fatal(rawTx)
	}
}

func TestBitcoinOrdGetter_GetOutput(t *testing.T) {
	out, err := testClient(t).GetOutput(httputl.TODO(), "a071d2a7abb989bd47d186b6b4bfe74b9673d5529dbfcf8f76229720f6b867c4", 0)
	if err != nil {
		t.Fatal(err)
	}
	if out.Value != 8.29761253 || out.ScriptPubKey.Asm != "0 ad0bb347149ac1889c8b92140daa7ad06a14a3b2" {
		t.Fatal(out)
	}
}

func TestBitcoinOrdGetter_GetBlock(t *testing.T) {
	const hash = "0000000000000000000454a3a654c88ab5ad9824ca8506c1f7f65cc0ea193503"
	b, err := testClient(t).GetBlock(httputl.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	if b.Hash != hash {
		t.Fatal(b)
	}
	if len(b.RawTx) != 0 {
		t.Fatal(b)
	}
}

func TestBitcoinOrdGetter_GetBlockDetail(t *testing.T) {
	const hash = "0000000000000000000454a3a654c88ab5ad9824ca8506c1f7f65cc0ea193503"
	b, err := testClient(t).GetBlock(httputl.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	if b.Hash != hash {
		t.Fatal(b)
	}
	if len(b.RawTx) == 0 {
		t.Fatal(b)
	}
}

func TestBitcoinOrdGetter_GetAllInscriptions(t *testing.T) {
	ins, err := testClient(t).GetAllInscriptions(httputl.TODO(), "9db3938b6ae166668e35e6f219a5c3a6146b613eed2f088644ce1fe829309b55")
	if err != nil {
		t.Fatal(err)
	}
	if len(ins) == 0 {
		t.Fatal(ins)
	}
}
