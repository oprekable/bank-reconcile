package versionhelper

import "testing"

func TestGetVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok - SNAPSHOT",
			args: args{
				version: "",
			},
			want: "SNAPSHOT",
		},
		{
			name: "ok - non SNAPSHOT",
			args: args{
				version: "v1.0.0",
			},
			want: "v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersion(tt.args.version); got != tt.want {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
