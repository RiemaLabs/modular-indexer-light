package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/config"
)

func TestProviderDA_GetCheckpoint(t *testing.T) {
	type fields struct {
		Config               *config.SourceDA
		MetaProtocol         string
		Retry                int
		LastCheckpointOffset int
	}
	type args struct {
		ctx    context.Context
		height uint
		hash   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *config.CheckpointExport
		wantErr bool
	}{
		{
			name: "common test",
			fields: fields{
				Config: &config.SourceDA{
					Network:     "Pre-Alpha Testnet",
					NamespaceID: "0x00000003",
					Name:        "nubit-official-00",
				},
				MetaProtocol:         "brc-20",
				Retry:                1,
				LastCheckpointOffset: 709,
			},
			args: args{
				ctx:    context.TODO(),
				height: 836614,
				hash:   "0000000000000000000303c500b5359801231aa44ef53e6f4d8017aa2e97aeb0",
			},
			want: &config.CheckpointExport{
				Checkpoint: &checkpoint.Checkpoint{Commitment: "KiHR43Oqvcl5Bbunbc69/ObmqakjfPiyT9v8jJ9DDBU=", Hash: "0000000000000000000303c500b5359801231aa44ef53e6f4d8017aa2e97aeb0", Height: "836614", MetaProtocol: "brc-20", Name: "nubit-official-00", URL: "https://committee.modular.nubit.org", Version: "v0.1.0-rc.0"},
				SourceDA: &config.SourceDA{
					Network:     "Pre-Alpha Testnet",
					NamespaceID: "0x00000003",
					Name:        "nubit-official-00",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProviderDA{
				Config:               tt.fields.Config,
				MetaProtocol:         tt.fields.MetaProtocol,
				Retry:                tt.fields.Retry,
				LastCheckpointOffset: tt.fields.LastCheckpointOffset,
			}
			got, err := p.GetCheckpoint(tt.args.ctx, tt.args.height, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderDA.GetCheckpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderDA.GetCheckpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
