package ordi

import (
	"testing"

	"github.com/RiemaLabs/modular-indexer-committee/ord/getter"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/btcutl"
)

func TestVerify(t *testing.T) {
	btcutl.Init("https://bitcoin-mainnet-archive.allthatnode.com")
	transfers := []getter.OrdTransfer{
		{
			InscriptionID: "fb0d434af0bebb1808b6454614020306a5dcd49209ae463eaa58643848d344dfi0",
			OldSatpoint:   "",
			NewSatpoint:   "fb0d434af0bebb1808b6454614020306a5dcd49209ae463eaa58643848d344df:0:0",
			NewPkscript:   "0014acea3e647df1bcc0559308ea776eb8d45ce327b0",
			NewWallet:     "bc1q4n4ruera7x7vq4vnpr48wm4c63wwxfast6vume",
			Content:       []byte(`{"p":"brc-20","op":"transfer","amt":"2000000","tick":"DRCR"}`),
			ContentType:   "746578742f706c61696e",
		},
		{
			InscriptionID: "9db3938b6ae166668e35e6f219a5c3a6146b613eed2f088644ce1fe829309b55i0",
			OldSatpoint:   "9db3938b6ae166668e35e6f219a5c3a6146b613eed2f088644ce1fe829309b55:0:0",
			NewSatpoint:   "ea2986b56db47bfe36a057db0e3b4668cad89fcb1f237335b378a8b31b2ee22e:0:0",
			NewPkscript:   "5120f9d29c2c8ce283ad4751d63847baa86587da2b04e620d56beda494e1a794f397",
			NewWallet:     "bc1pl8ffctyvu2p663636cuy0w4gvkra52cyucsd26ld5j2wrfu57wtse23er9",
			Content:       []byte(`{"p":"brc-20","op":"mint","amt":"100","tick":"HUHU"}`),
			ContentType:   "746578742f706c61696e3b636861727365743d7574662d38",
		},
	}
	if err := VerifyOrdTransfer(transfers, 835477); err != nil {
		t.Fatal(err)
	}
}
