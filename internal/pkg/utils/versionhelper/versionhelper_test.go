package versionhelper

import (
	"reflect"
	"runtime/debug"
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
		name string
		args args

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
				t.Errorf("getVersionLogic() gotReturnedData.Version= %v, want.Version %v", gotReturnData.Version, tt.wantReturnData.Version)
			}
			if !reflect.DeepEqual(gotReturnData.Environment, tt.wantReturnData.Environment) {
				t.Errorf("getVersionLogic() gotReturnedData.Environment= %v, want.Environment %v", gotReturnData.Environment, tt.wantReturnData.Environment)
			}
		})
	}
}

func TestGetVersionLogic(t *testing.T) {
	type args struct {
		buildInfo   *debug.BuildInfo
		version     string
		buildDate   string
		commitHash  string
		environment string
	}

	tests := []struct {
		name string
		args args

		wantReturnData VersionStruct
	}{
		{
			name: "Ok",
			args: args{
				buildInfo: func() *debug.BuildInfo {
					return &debug.BuildInfo{
						Main: debug.Module{
							Version: "v1.2.3",
						},
						Settings: []debug.BuildSetting{
							{Key: "vcs.revision", Value: "abcdef123456"},
							{Key: "vcs.time", Value: "2025-10-10T10:00:00Z"},
						},
					}
				}(),
				version:     "",
				environment: "",
			},
			wantReturnData: VersionStruct{
				Version:     "v1.2.3",
				Environment: "default",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReturnData := getVersionLogic(tt.args.buildInfo, tt.args.version, tt.args.buildDate, tt.args.commitHash, tt.args.environment)
			if !reflect.DeepEqual(gotReturnData.Version, tt.wantReturnData.Version) {
				t.Errorf("getVersionLogic() gotReturnedData.Version= %v, want.Version %v", gotReturnData.Version, tt.wantReturnData.Version)
			}
			if !reflect.DeepEqual(gotReturnData.Environment, tt.wantReturnData.Environment) {
				t.Errorf("getVersionLogic() gotReturnedData.Environment= %v, want.Environment %v", gotReturnData.Environment, tt.wantReturnData.Environment)
			}
		})
	}
}
