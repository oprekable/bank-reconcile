package cconfig

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"
	"time"

	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/oprekable/bank-reconcile/internal/app/config/core"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"

	"github.com/spf13/afero"
)

//go:embed all:embeds
var testEmbedFS embed.FS

func TestConfigGetAppName(t *testing.T) {
	type fields struct {
		Data         *config.Data
		timeLocation *time.Location
		appName      AppName
		workDirPath  WorkDirPath
		timeZone     TimeZone
		timeOffset   TimeOffset
	}

	tests := []struct {
		name   string
		want   AppName
		fields fields
	}{
		{
			name: "Ok",
			fields: fields{
				Data:         nil,
				timeLocation: nil,
				appName:      "foo",
				workDirPath:  "",
				timeZone:     "",
				timeOffset:   0,
			},
			want: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Data:         tt.fields.Data,
				timeLocation: tt.fields.timeLocation,
				appName:      tt.fields.appName,
				workDirPath:  tt.fields.workDirPath,
				timeZone:     tt.fields.timeZone,
				timeOffset:   tt.fields.timeOffset,
			}

			if got := c.GetAppName(); got != tt.want {
				t.Errorf("GetAppName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigGetTimeLocation(t *testing.T) {
	type fields struct {
		Data         *config.Data
		timeLocation *time.Location
		appName      AppName
		workDirPath  WorkDirPath
		timeZone     TimeZone
		timeOffset   TimeOffset
	}

	tests := []struct {
		want   *time.Location
		name   string
		fields fields
	}{
		{
			name: "Ok",
			fields: fields{
				Data: nil,
				timeLocation: func() *time.Location {
					_, loc, _, _ := initTimeZone("Asia/Jakarta")
					return loc
				}(),
				appName:     "",
				workDirPath: "",
				timeZone:    "Asia/Jakarta",
				timeOffset:  0,
			},
			want: func() *time.Location {
				_, loc, _, _ := initTimeZone("Asia/Jakarta")
				return loc
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Data:         tt.fields.Data,
				timeLocation: tt.fields.timeLocation,
				appName:      tt.fields.appName,
				workDirPath:  tt.fields.workDirPath,
				timeZone:     tt.fields.timeZone,
				timeOffset:   tt.fields.timeOffset,
			}

			if got := c.GetTimeLocation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTimeLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigGetTimeOffset(t *testing.T) {
	type fields struct {
		Data         *config.Data
		timeLocation *time.Location
		appName      AppName
		workDirPath  WorkDirPath
		timeZone     TimeZone
		timeOffset   TimeOffset
	}

	tests := []struct {
		name   string
		fields fields
		want   TimeOffset
	}{
		{
			name: "Ok",
			fields: fields{
				Data:         nil,
				timeLocation: nil,
				appName:      "",
				workDirPath:  "",
				timeZone:     "",
				timeOffset: func() TimeOffset {
					_, _, to, _ := initTimeZone("Asia/Jakarta")
					return TimeOffset(to)
				}(),
			},
			want: func() TimeOffset {
				_, _, to, _ := initTimeZone("Asia/Jakarta")
				return TimeOffset(to)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Data:         tt.fields.Data,
				timeLocation: tt.fields.timeLocation,
				appName:      tt.fields.appName,
				workDirPath:  tt.fields.workDirPath,
				timeZone:     tt.fields.timeZone,
				timeOffset:   tt.fields.timeOffset,
			}

			got := c.GetTimeOffset()
			if got != tt.want {
				t.Errorf("GetTimeOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigGetTimeZone(t *testing.T) {
	type fields struct {
		Data         *config.Data
		timeLocation *time.Location
		appName      AppName
		workDirPath  WorkDirPath
		timeZone     TimeZone
		timeOffset   TimeOffset
	}

	tests := []struct {
		name   string
		want   TimeZone
		fields fields
	}{
		{
			name: "Ok",
			fields: fields{
				Data:         nil,
				timeLocation: nil,
				appName:      "",
				workDirPath:  "",
				timeZone: func() TimeZone {
					tz, _, _, _ := initTimeZone("Asia/Jakarta")
					return TimeZone(tz)
				}(),
				timeOffset: 0,
			},
			want: func() TimeZone {
				tz, _, _, _ := initTimeZone("Asia/Jakarta")
				return TimeZone(tz)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Data:         tt.fields.Data,
				timeLocation: tt.fields.timeLocation,
				appName:      tt.fields.appName,
				workDirPath:  tt.fields.workDirPath,
				timeZone:     tt.fields.timeZone,
				timeOffset:   tt.fields.timeOffset,
			}

			if got := c.GetTimeZone(); got != tt.want {
				t.Errorf("GetTimeZone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigGetWorkDirPath(t *testing.T) {
	type fields struct {
		Data         *config.Data
		timeLocation *time.Location
		appName      AppName
		workDirPath  WorkDirPath
		timeZone     TimeZone
		timeOffset   TimeOffset
	}

	tests := []struct {
		name   string
		want   WorkDirPath
		fields fields
	}{
		{
			name: "Ok",
			fields: fields{
				Data:         nil,
				timeLocation: nil,
				appName:      "",
				workDirPath:  "/foo",
				timeZone:     "",
				timeOffset:   0,
			},
			want: "/foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Data:         tt.fields.Data,
				timeLocation: tt.fields.timeLocation,
				appName:      tt.fields.appName,
				workDirPath:  tt.fields.workDirPath,
				timeZone:     tt.fields.timeZone,
				timeOffset:   tt.fields.timeOffset,
			}

			if got := c.GetWorkDirPath(); got != tt.want {
				t.Errorf("GetWorkDirPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		ctx         context.Context
		afs         afero.Fs
		embedFS     *embed.FS
		appName     AppName
		tzArgs      TimeZone
		configPaths ConfigPaths
	}

	tests := []struct {
		want    *Config
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				ctx:     context.Background(),
				embedFS: &testEmbedFS,
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					workDirPath := initWorkDirPath()
					appToml, _ := f.Create(string(workDirPath) + "/app.toml")
					_, _ = appToml.Write([]byte(
						`[app]
secret = "36189d92-a31e-412c-ab50-ebd94ab52698"
`,
					))

					_ = appToml.Close()

					reconToml, _ := f.Create(string(workDirPath) + "/params/reconciliation.toml")
					_, _ = reconToml.Write([]byte(
						`[reconciliation]
system_trx_path = "/tmp/system"
bank_trx_path= "/tmp/bank"
report_trx_path = "/tmp/report"
total_data = 1000
percentage_match = 100
list_bank = [ "bca", "danamon", "bri", "mandiri" ]
is_delete_current_sample_directory = true
`,
					))

					_ = reconToml.Close()

					envFile, _ := f.Create(string(workDirPath) + "/params/.env")
					_, _ = reconToml.Write([]byte(
						`APP_SECRET="secret-from-memory"
`,
					))

					_ = envFile.Close()

					return f
				}(),
				appName:     "foo",
				tzArgs:      "",
				configPaths: nil,
			},
			want: &Config{
				Data: &config.Data{
					App: core.App{
						Secret:    "796bb93e-6269-4aa1-96e4-a8b54aad40aa",
						IsShowLog: "true",
						IsDebug:   false,
					},
					Sqlite: core.Sqlite{
						Write: core.SqliteParameters{
							DBPath:      ":memory:",
							Cache:       "shared",
							JournalMode: "WAL",
							IsEnabled:   false,
						},
						Read: core.SqliteParameters{
							DBPath:      ":memory:",
							Cache:       "shared",
							JournalMode: "WAL",
							IsEnabled:   false,
						},
						IsEnabled:   false,
						IsDoLogging: false,
					},
					Reconciliation: reconciliation.Reconciliation{
						TotalData:                      1000,
						PercentageMatch:                100,
						NumberWorker:                   10,
						IsDeleteCurrentSampleDirectory: true,
						SystemTRXPath:                  "/tmp/system",
						BankTRXPath:                    "/tmp/bank",
						ReportTRXPath:                  "/tmp/report",
						ListBank: []string{
							"bca", "danamon", "bri", "mandiri",
						},
					},
				},
				timeLocation: func() *time.Location {
					_, loc, _, _ := initTimeZone("Asia/Jakarta")
					return loc
				}(),
				appName:     "foo",
				workDirPath: initWorkDirPath(),
				timeZone:    "Asia/Jakarta",
				timeOffset:  TimeOffset(25200),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.ctx, tt.args.embedFS, tt.args.afs, tt.args.configPaths, tt.args.appName, tt.args.tzArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}

			t.Cleanup(func() {
				os.Clearenv()
			})
		})
	}
}

func TestFromFS(t *testing.T) {
	type args struct {
		conf     interface{}
		embedFS  fs.ReadFileFS
		patterns []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				conf: &config.Data{},
				embedFS: func() fs.ReadFileFS {
					return fstest.MapFS{
						"embeds/params/app.toml": &fstest.MapFile{
							Data: []byte(
								`[app]
secret = "36189d92-a31e-412c-ab50-ebd94ab52697"
`,
							),
							Mode:    fs.ModePerm,
							ModTime: time.Time{},
							Sys:     nil,
						},
					}
				}(),
				patterns: []string{
					"embeds/params/*.toml",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fromFS(tt.args.embedFS, tt.args.patterns, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromFS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromFiles(t *testing.T) {
	type args struct {
		conf     interface{}
		afs      afero.Fs
		patterns []string
	}

	tests := []struct {
		name                                      string
		wantConfigDataAppSecret                   string
		wantConfigDataReconciliationSystemTrxPath string
		args                                      args
		wantErr                                   bool
	}{
		{
			name: "Ok",
			args: args{
				patterns: []string{
					"./*.toml",
					"./params/*.toml",
				},
				conf: &config.Data{},
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					appToml, _ := f.Create("./app.toml")
					_, _ = appToml.Write([]byte(
						`[app]
secret = "36189d92-a31e-412c-ab50-ebd94ab52698"
`,
					))

					_ = appToml.Close()

					reconToml, _ := f.Create("./params/reconciliation.toml")
					_, _ = reconToml.Write([]byte(
						`[reconciliation]
system_trx_path = "/tmp/system"
bank_trx_path= "/tmp/bank"
report_trx_path = "/tmp/report"
total_data = 1000
percentage_match = 100
list_bank = [ "bca", "danamon", "bri", "mandiri" ]
is_delete_current_sample_directory = true
`,
					))

					_ = reconToml.Close()

					return f
				}(),
			},
			wantConfigDataAppSecret:                   "36189d92-a31e-412c-ab50-ebd94ab52698",
			wantConfigDataReconciliationSystemTrxPath: "/tmp/system",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fromFiles(tt.args.patterns, tt.args.conf, tt.args.afs)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromFiles() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotConfigDataAppSecret := tt.args.conf.(*config.Data).App.Secret
			if !reflect.DeepEqual(gotConfigDataAppSecret, tt.wantConfigDataAppSecret) {
				t.Errorf("initTimeZone() gotConfigDataAppSecret = %v, want %v", gotConfigDataAppSecret, tt.wantConfigDataAppSecret)
			}

			gotConfigDataReconciliationSystemTrxPath := tt.args.conf.(*config.Data).Reconciliation.SystemTRXPath
			if !reflect.DeepEqual(gotConfigDataReconciliationSystemTrxPath, tt.wantConfigDataReconciliationSystemTrxPath) {
				t.Errorf("initTimeZone() gotConfigDataAppSecret = %v, want %v", gotConfigDataReconciliationSystemTrxPath, tt.wantConfigDataReconciliationSystemTrxPath)
			}
		})
	}
}

func TestInitTimeZone(t *testing.T) {
	type args struct {
		tzArgs TimeZone
	}

	tests := []struct {
		name       string
		args       args
		wantTz     string
		wantLoc    string
		wantOffset int
		wantErr    bool
	}{
		{
			name: "Ok",
			args: args{
				tzArgs: TimeZone("Asia/Jakarta"),
			},
			wantTz:     "Asia/Jakarta",
			wantLoc:    "Asia/Jakarta",
			wantOffset: 25200,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTz, gotLoc, gotOffset, err := initTimeZone(tt.args.tzArgs)

			if (err != nil) != tt.wantErr {
				t.Errorf("initTimeZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotTz != tt.wantTz {
				t.Errorf("initTimeZone() gotTz = %v, want %v", gotTz, tt.wantTz)
			}

			if !reflect.DeepEqual(gotLoc.String(), tt.wantLoc) {
				t.Errorf("initTimeZone() gotLoc = %v, want %v", gotLoc.String(), tt.wantLoc)
			}

			if gotOffset != tt.wantOffset {
				t.Errorf("initTimeZone() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

func TestInitWorkDirPath(t *testing.T) {
	tests := []struct {
		name string
		want WorkDirPath
	}{
		{
			name: "Ok",
			want: func() WorkDirPath {
				ex, _ := os.Executable()
				return WorkDirPath(filepath.Dir(ex))
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := initWorkDirPath()
			if got != tt.want {
				t.Errorf("initWorkDirPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
