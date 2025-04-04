package csvhelper

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestDeleteDirectory(t *testing.T) {
	type args struct {
		ctx      context.Context
		fs       afero.Fs
		filePath string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Directory does not exist",
			args: args{
				ctx: context.Background(),
				fs:  afero.NewMemMapFs(),
			},
			wantErr: false,
		},
		{
			name: "Directory  exist, success",
			args: args{
				ctx: context.Background(),
				fs: func() afero.Fs {
					fs := afero.NewMemMapFs()
					_ = fs.Mkdir("test", os.ModeDir)
					return fs
				}(),
			},
			wantErr: false,
		},
		{
			name: "Directory  exist, error",
			args: args{
				ctx: context.Background(),
				fs: func() afero.Fs {
					fs := afero.NewMemMapFs()
					_ = fs.Mkdir("test", os.ModeDir)
					return afero.NewReadOnlyFs(fs)
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteDirectory(tt.args.ctx, tt.args.fs, tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("DeleteDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStructToCSVFile(t *testing.T) {
	type T struct {
		Name string
	}

	type args struct {
		ctx               context.Context
		fs                afero.Fs
		structData        interface{}
		filePath          string
		isDeleteDirectory bool
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "isDeleteDirectory & structData have value - Success",
			args: args{
				ctx:      context.Background(),
				fs:       afero.NewMemMapFs(),
				filePath: "/test/test.csv",
				structData: []T{
					{
						Name: "test",
					},
				},
				isDeleteDirectory: true,
			},
			wantErr: false,
		},
		{
			name: "isDeleteDirectory = false & structData have no value - Success",
			args: args{
				ctx:               context.Background(),
				fs:                afero.NewMemMapFs(),
				filePath:          "/test/test.csv",
				structData:        []T{},
				isDeleteDirectory: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StructToCSVFile(tt.args.ctx, tt.args.fs, tt.args.filePath, tt.args.structData, tt.args.isDeleteDirectory); (err != nil) != tt.wantErr {
				t.Errorf("StructToCSVFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
