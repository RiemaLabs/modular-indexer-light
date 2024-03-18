package transfer

import (
	"net/http"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
)

func TestHttpGetter_GetRawTransaction(t *testing.T) {
	type fields struct {
		URL      string
		Username string
		Password string
		client   *http.Client
	}
	type args struct {
		txID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *btcutil.Tx
		wantErr bool
	}{
		{
			name: "rawtransaction test",
			fields: fields{
				URL:      "https://frosty-serene-emerald.btc.quiknode.pro/402f5ac57de95e38c0a33d1a5e6f6c2f66709262/",
				Username: "",
				Password: "",
				client:   &http.Client{},
			},
			args: args{
				txID: "26a08b3ac578f1fe01bde9d0268353121f22461fcb48dc3144f1dd5210d0f8ad",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpGetter{
				URL:      tt.fields.URL,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				client:   tt.fields.client,
			}
			got, err := r.GetRawTransaction(tt.args.txID)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpGetter.GetRawTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Txid != tt.args.txID {
				t.Errorf("HttpGetter.GetRawTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpGetter_GetBlock(t *testing.T) {
	type fields struct {
		URL      string
		Username string
		Password string
		client   *http.Client
	}
	type args struct {
		hash    string
		verbose int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *btcjson.GetBlockVerboseResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "getblock test",
			fields: fields{
				URL:      "https://frosty-serene-emerald.btc.quiknode.pro/402f5ac57de95e38c0a33d1a5e6f6c2f66709262/",
				Username: "",
				Password: "",
				client:   &http.Client{},
			},
			args: args{
				hash:    "0000000000000000000454a3a654c88ab5ad9824ca8506c1f7f65cc0ea193503",
				verbose: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpGetter{
				URL:      tt.fields.URL,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				client:   tt.fields.client,
			}
			got, err := r.GetBlock(tt.args.hash, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpGetter.GetBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Hash != tt.args.hash {
				t.Errorf("HttpGetter.GetBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
