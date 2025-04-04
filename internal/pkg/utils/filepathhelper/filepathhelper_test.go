package filepathhelper

import (
	"os"
	"testing"
)

type MockErrorSystemCalls struct{}

func (MockErrorSystemCalls) Getwd() (string, error) {
	return "", os.ErrPermission
}

func (MockErrorSystemCalls) Executable() (string, error) {
	return "/foo", nil
}

func (MockErrorSystemCalls) FilepathDir(_ string) string {
	return "/"
}

func TestGetWorkDir(t *testing.T) {
	tests := []struct {
		name        string
		systemCalls ISystemCalls
		want        string
	}{
		{
			name:        "Success get work dir",
			systemCalls: SystemCalls{},
			want: func() string {
				wDirPath, _ := os.Getwd()
				return wDirPath
			}(),
		},
		{
			name:        "Error get work dir",
			systemCalls: MockErrorSystemCalls{},
			want:        "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWorkDir(tt.systemCalls); got != tt.want {
				t.Errorf("GetWorkDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemCallsGetwd(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			want: func() string {
				r, _ := os.Getwd()
				return r
			}(),
			wantErr: func() bool {
				_, e := os.Getwd()
				return e != nil
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sy := SystemCalls{}
			got, err := sy.Getwd()

			if (err != nil) != tt.wantErr {
				t.Errorf("Getwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Getwd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemCallsExecutable(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			want: func() string {
				r, _ := os.Executable()
				return r
			}(),
			wantErr: func() bool {
				_, e := os.Executable()
				return e != nil
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sy := SystemCalls{}
			got, err := sy.Executable()
			if (err != nil) != tt.wantErr {
				t.Errorf("Executable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Executable() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemCallsFilepathDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Ok",
			args: args{
				path: "/foo/bar",
			},
			want: "/foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sy := SystemCalls{}
			if got := sy.FilepathDir(tt.args.path); got != tt.want {
				t.Errorf("FilepathDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
