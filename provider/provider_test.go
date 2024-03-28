package provider

import (
	"reflect"
	"testing"
	"time"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"
)

func TestDownloadCheckpointByDA(t *testing.T) {
	type args struct {
		namespaceID   string
		network       string
		name          string
		metaProtocol  string
		height        string
		hash          string
		runtimeOffset int
		timeout       time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *checkpoint.Checkpoint
		want1   int
		wantErr bool
	}{
		{
			name: "common test",
			args: args{
				namespaceID:   "0x00000003",
				network:       "Pre-Alpha Testnet",
				name:          "nubit-official-00",
				metaProtocol:  "brc-20",
				height:        "836614",
				hash:          "0000000000000000000303c500b5359801231aa44ef53e6f4d8017aa2e97aeb0",
				runtimeOffset: 709,
				timeout:       10 * time.Second,
			},
			want: &checkpoint.Checkpoint{Commitment: "KiHR43Oqvcl5Bbunbc69/ObmqakjfPiyT9v8jJ9DDBU=", Hash: "0000000000000000000303c500b5359801231aa44ef53e6f4d8017aa2e97aeb0", Height: "836614", MetaProtocol: "brc-20", Name: "nubit-official-00", URL: "https://committee.modular.nubit.org", Version: "v0.1.0-rc.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := DownloadCheckpointByDA(tt.args.namespaceID, tt.args.network, tt.args.name, tt.args.metaProtocol, tt.args.height, tt.args.hash, tt.args.runtimeOffset, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadCheckpointByDA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DownloadCheckpointByDA() got = %v, want %v", got, tt.want)
			}
			// if got1 != tt.want1 {
			// 	t.Errorf("DownloadCheckpointByDA() got1 = %v, want %v", got1, tt.want1)
			// }
		})
	}
}
