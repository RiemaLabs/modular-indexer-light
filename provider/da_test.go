package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/config"
)

func TestProviderDA_GetCheckpoint(t *testing.T) {
	config.InitConfig()
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
				Config:               &config.GlobalConfig.CommitteeIndexers.DA[0],
				MetaProtocol:         config.GlobalConfig.MetaProtocol,
				Retry:                1,
				LastCheckpointOffset: 685,
			},
			args: args{
				ctx:    context.TODO(),
				height: 836595,
				hash:   "00000000000000000001b813fe10fb6db185d9a8ccf92541951044a3c92a2bbd",
			},
			want: &config.CheckpointExport{
				Checkpoint: &checkpoint.Checkpoint{Commitment: "Ms1MH3NxzKNnaLiIHI/NwhTqhg0DC4cI8qW9Zq4mYzk=", Hash: "00000000000000000001b813fe10fb6db185d9a8ccf92541951044a3c92a2bbd", Height: "836595", MetaProtocol: "brc-20", Name: "nubit-official-00", URL: "https://committee.modular.nubit.org", Version: "v0.1.0-rc.0"},
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
