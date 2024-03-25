package getter

import (
	"reflect"
	"testing"

	"github.com/RiemaLabs/modular-indexer-light/clients/http"
	"github.com/RiemaLabs/modular-indexer-light/config"
	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/btcjson"
)

func TestBitcoinOrdGetter_GetRawTransaction(t *testing.T) {
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
				Endpoint: config.Config.BitCoinRpc.Host,
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
				Endpoint: config.Config.BitCoinRpc.Host,
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

func TestBitcoinOrdGetter_GetBlock1(t *testing.T) {
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
				Endpoint: config.Config.BitCoinRpc.Host,
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
			got, err := r.GetBlock1(tt.args.hash)
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

func TestBitcoinOrdGetter_GetBlock2(t *testing.T) {
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
				Endpoint: config.Config.BitCoinRpc.Host,
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
			got, err := r.GetBlock2(tt.args.hash)
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
				Endpoint: config.Config.BitCoinRpc.Host,
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
