package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
	"github.com/RiemaLabs/modular-indexer-light/config"
)

func TestProviderS3_GetCheckpoint(t *testing.T) {
	config.InitConfig()
	type fields struct {
		Config       *config.SourceS3
		MetaProtocol string
		Retry        int
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
			name: "comnon test",
			fields: fields{
				Config: &config.SourceS3{
					Region: "us-west-2",
					Bucket: "nubit-modular-indexer-brc-20",
					Name:   "nubit-official-02",
				},
				MetaProtocol: "brc-20",
				Retry:        1,
			},
			args: args{
				ctx:    context.TODO(),
				height: 836223,
				hash:   "000000000000000000020443905dc3ba40f729313523fba9a01d0d0c44fdd693",
			},
			want: &config.CheckpointExport{
				Checkpoint: &checkpoint.Checkpoint{Commitment: "QVzPh4JnwHGuPz9qRWm+q0z2HED4SvxY1uS0prvW+ZE=", Hash: "000000000000000000020443905dc3ba40f729313523fba9a01d0d0c44fdd693", Height: "836223", MetaProtocol: "brc-20", Name: "nubit-official-02", URL: "https://committee.modular.nubit.org", Version: "v0.1.0-rc.0"},
				SourceS3: &config.SourceS3{
					Region: "us-west-2",
					Bucket: "nubit-modular-indexer-brc-20",
					Name:   "nubit-official-02",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProviderS3{
				Config:       tt.fields.Config,
				MetaProtocol: tt.fields.MetaProtocol,
				Retry:        tt.fields.Retry,
			}
			got, err := p.GetCheckpoint(tt.args.ctx, tt.args.height, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderS3.GetCheckpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderS3.GetCheckpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
