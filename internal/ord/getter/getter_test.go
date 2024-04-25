package getter

import (
	"reflect"
	"testing"

	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"

	"github.com/RiemaLabs/modular-indexer-light/internal/clients/http"
	"github.com/RiemaLabs/modular-indexer-light/internal/configs"
)

func TestBitcoinOrdGetter_GetLatestBlockHeight(t *testing.T) {
	type fields struct {
		client   *http.Client
		Endpoint string
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				client:   http.NewClient(),
				Endpoint: configs.C.Verification.BitcoinRPC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BitcoinOrdGetter{
				client:   tt.fields.client,
				Endpoint: tt.fields.Endpoint,
			}
			got, err := r.GetLatestBlockHeight()
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinOrdGetter.GetLatestBlockHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got <= 0 {
				t.Errorf("BitcoinOrdGetter.GetLatestBlockHeight() = %v", got)
			}
		})
	}
}

func TestBitcoinOrdGetter_GetBlockHash(t *testing.T) {
	configs.InitConfig()
	type fields struct {
		client   *http.Client
		Endpoint string
	}
	type args struct {
		blockHeight uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				client:   http.NewClient(),
				Endpoint: configs.C.Verification.BitcoinRPC,
			},
			args: args{
				blockHeight: 835161,
			},
			want: "000000000000000000021a731d2106dda997d6eaf6228252c7abdc259c1fca5e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BitcoinOrdGetter{
				client:   tt.fields.client,
				Endpoint: tt.fields.Endpoint,
			}
			got, err := r.GetBlockHash(tt.args.blockHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinOrdGetter.GetBlockHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BitcoinOrdGetter.GetBlockHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitcoinOrdGetter_GetRawTransaction(t *testing.T) {
	configs.InitConfig()
	type fields struct {
		client   *http.Client
		Endpoint string
	}
	type args struct {
		txID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *btcjson.TxRawResult
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				client:   http.NewClient(),
				Endpoint: configs.C.Verification.BitcoinRPC,
			},
			args: args{
				txID: "26a08b3ac578f1fe01bde9d0268353121f22461fcb48dc3144f1dd5210d0f8ad",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BitcoinOrdGetter{
				client:   tt.fields.client,
				Endpoint: tt.fields.Endpoint,
			}
			got, err := r.GetRawTransaction(tt.args.txID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinOrdGetter.GetRawTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Txid != tt.args.txID {
				t.Errorf("BitcoinOrdGetter.GetRawTransaction()  txid =  %v, want %v", got.Txid, tt.args.txID)
			}
		})
	}
}

func TestBitcoinOrdGetter_GetOutput(t *testing.T) {
	configs.InitConfig()
	type fields struct {
		client   *http.Client
		Endpoint string
	}
	type args struct {
		txID  string
		index int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *btcjson.Vout
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				client:   http.NewClient(),
				Endpoint: configs.C.Verification.BitcoinRPC,
			},
			args: args{
				txID:  "a071d2a7abb989bd47d186b6b4bfe74b9673d5529dbfcf8f76229720f6b867c4",
				index: 0,
			},
			want: &btcjson.Vout{
				Value: 8.29761253,
				N:     0,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm: "0 ad0bb347149ac1889c8b92140daa7ad06a14a3b2",
					// Desc: "addr(bc1q459mx3c5ntqc38ytjg2qm2n66p4pfgajlpzvwu)#w4ky3fht",
					Hex:     "0014ad0bb347149ac1889c8b92140daa7ad06a14a3b2",
					Address: "bc1q459mx3c5ntqc38ytjg2qm2n66p4pfgajlpzvwu",
					Type:    "witness_v0_keyhash",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BitcoinOrdGetter{
				client:   tt.fields.client,
				Endpoint: tt.fields.Endpoint,
			}
			got, err := r.GetOutput(tt.args.txID, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinOrdGetter.GetOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitcoinOrdGetter.GetOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitcoinOrdGetter_GetBlock(t *testing.T) {
	configs.InitConfig()
	type fields struct {
		client   *http.Client
		Endpoint string
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *btcjson.GetBlockVerboseResult
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				client:   http.NewClient(),
				Endpoint: configs.C.Verification.BitcoinRPC,
			},
			args: args{
				hash: "0000000000000000000454a3a654c88ab5ad9824ca8506c1f7f65cc0ea193503",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BitcoinOrdGetter{
				client:   tt.fields.client,
				Endpoint: tt.fields.Endpoint,
			}
			got, err := r.GetBlock(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinOrdGetter.GetBlock1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Hash != tt.args.hash {
				t.Errorf("BitcoinOrdGetter.GetBlock1() hash =  %v, want %v", got.Hash, tt.args.hash)
			}
		})
	}
}

func TestBitcoinOrdGetter_GetBlockDetail(t *testing.T) {
	configs.InitConfig()
	type fields struct {
		client   *http.Client
		Endpoint string
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *btcjson.GetBlockVerboseTxResult
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				client:   http.NewClient(),
				Endpoint: configs.C.Verification.BitcoinRPC,
			},
			args: args{
				hash: "0000000000000000000454a3a654c88ab5ad9824ca8506c1f7f65cc0ea193503",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BitcoinOrdGetter{
				client:   tt.fields.client,
				Endpoint: tt.fields.Endpoint,
			}
			got, err := r.GetBlockDetail(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinOrdGetter.GetBlock2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Hash != tt.args.hash {
				t.Errorf("BitcoinOrdGetter.GetBlock2() hash =  %v, want %v", got.Hash, tt.args.hash)
			}
		})
	}
}

func TestBitcoinOrdGetter_GetAllInscriptions(t *testing.T) {
	configs.InitConfig()
	type fields struct {
		client   *http.Client
		Endpoint string
	}
	type args struct {
		txID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*parser.TransactionInscription
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				client:   http.NewClient(),
				Endpoint: configs.C.Verification.BitcoinRPC,
			},
			args: args{
				txID: "9db3938b6ae166668e35e6f219a5c3a6146b613eed2f088644ce1fe829309b55",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BitcoinOrdGetter{
				client:   tt.fields.client,
				Endpoint: tt.fields.Endpoint,
			}
			got, err := r.GetAllInscriptions(tt.args.txID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinOrdGetter.GetAllInscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Error("BitcoinOrdGetter.GetAllInscriptions() failed")
			}
		})
	}
}
