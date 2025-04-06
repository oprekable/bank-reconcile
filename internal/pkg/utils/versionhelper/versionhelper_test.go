package versionhelper

import (
	"reflect"
	"testing"
)

func TestGetVersion(t *testing.T) {
	type args struct {
		version     string
		buildDate   string
		commitHash  string
		environment string
	}

	tests := []struct {
		name           string
		args           args
		wantReturnData VersionStruct
	}{
		{
			name: "Ok",
			args: args{
				version:     "",
				environment: "",
			},
			wantReturnData: VersionStruct{
				Version:     "(devel)",
				Environment: "default",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReturnData := GetVersion(tt.args.version, tt.args.buildDate, tt.args.commitHash, tt.args.environment)
			if !reflect.DeepEqual(gotReturnData.Version, tt.wantReturnData.Version) {
				t.Errorf("GetVersion() gotReturnedData.Version= %v, want.Version %v", gotReturnData.Version, tt.wantReturnData.Version)
			}
			if !reflect.DeepEqual(gotReturnData.Environment, tt.wantReturnData.Environment) {
				t.Errorf("GetVersion() gotReturnedData.Environment= %v, want.Environment %v", gotReturnData.Environment, tt.wantReturnData.Environment)
			}
		})
	}
}
