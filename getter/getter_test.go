package getter

import (
	"testing"

	"github.com/RiemaLabs/indexer-light/clients/http"
	"github.com/RiemaLabs/indexer-light/config"
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
				Endpoint: config.Config.BitCoinRpc.Host,
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
				Endpoint: config.Config.BitCoinRpc.Host,
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
