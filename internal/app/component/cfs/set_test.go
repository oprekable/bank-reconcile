package cfs

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestProviderCFs(t *testing.T) {
	type args struct {
		fsType afero.Fs
	}
	tests := []struct {
		args args
		want *Fs
		name string
	}{
		{
			name: "Ok",
			args: args{fsType: afero.NewMemMapFs()},
			want: &Fs{
				LocalStorageFs: afero.NewMemMapFs(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProviderCFs(tt.args.fsType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderCFs() = %v, want %v", got, tt.want)
			}
		})
	}
}
