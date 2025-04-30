package cprofiler

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/stretchr/testify/assert"
)

// checkPprofFile checks if input pprof files exist
func checkPprofFiles(t *testing.T, pprofFilesPath []string) {
	for _, pprof := range pprofFilesPath {
		info, err := os.Stat(pprof)
		assert.Nil(t, err)
		assert.False(t, os.IsNotExist(err))
		assert.False(t, info.IsDir())
	}
}

// cleanupPprofFiles deletes all specified pprof files
func cleanupPprofFiles(t *testing.T, pprofFilesPath []string) {
	for _, pprof := range pprofFilesPath {
		err := os.Remove(pprof)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestNewProfiler(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		logger *clogger.Logger
	}

	var logger = clogger.NewLogger(
		context.Background(),
		&bf,
	)

	tests := []struct {
		args args
		want *Profiler
		name string
	}{
		{
			name: "Ok",
			args: args{
				logger: logger,
			},
			want: &Profiler{
				ctx:      logger.GetLogger().With().Str("component", "Profiler").Ctx(context.Background()).Logger().WithContext(logger.GetCtx()),
				logger:   logger,
				profiler: make(map[string]interface{ Stop() }),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProfiler(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProfiler() = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

func TestProfilerStartProfiler(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		ctx      context.Context
		logger   *clogger.Logger
		profiler map[string]interface{ Stop() }
	}

	var logger = clogger.NewLogger(
		context.Background(),
		&bf,
	)

	tests := []struct {
		name    string
		fields  fields
		wantLog string
	}{
		{
			name: "Ok",
			fields: fields{
				ctx:      logger.GetLogger().With().Str("component", "Profiler").Ctx(context.Background()).Logger().WithContext(logger.GetCtx()),
				logger:   logger,
				profiler: make(map[string]interface{ Stop() }),
			},
			wantLog: `[profiler] starting profiler`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Profiler{
				ctx:      tt.fields.ctx,
				logger:   tt.fields.logger,
				profiler: tt.fields.profiler,
			}

			p.StartProfiler()

			checkPprofFiles(t, []string{
				"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
				"./trace.pprof", "./goroutine.pprof",
			})

			cleanupPprofFiles(t, []string{
				"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
				"./trace.pprof", "./goroutine.pprof",
			})

			got := bf.String()
			if !strings.Contains(got, tt.wantLog) {
				t.Errorf("StartProfiler() log = %v, wantlog %v", got, tt.wantLog)
			}

			bf.Reset()
		})
	}
}

func TestProfilerStopProfiler(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		ctx      context.Context
		logger   *clogger.Logger
		profiler map[string]interface{ Stop() }
	}

	var logger = clogger.NewLogger(
		context.Background(),
		&bf,
	)

	tests := []struct {
		name    string
		fields  fields
		wantLog string
	}{
		{
			name: "Ok",
			fields: fields{
				ctx:      logger.GetLogger().With().Str("component", "Profiler").Ctx(context.Background()).Logger().WithContext(logger.GetCtx()),
				logger:   logger,
				profiler: make(map[string]interface{ Stop() }),
			},
			wantLog: `[profiler] stop profiler`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Profiler{
				ctx:      tt.fields.ctx,
				logger:   tt.fields.logger,
				profiler: tt.fields.profiler,
			}

			p.StartProfiler()
			bf.Reset()
			p.StopProfiler()

			checkPprofFiles(t, []string{
				"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
				"./trace.pprof", "./goroutine.pprof",
			})

			cleanupPprofFiles(t, []string{
				"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
				"./trace.pprof", "./goroutine.pprof",
			})

			got := bf.String()
			if !strings.Contains(got, tt.wantLog) {
				t.Errorf("StopProfiler() log = %v, wantlog %v", got, tt.wantLog)
			}

			bf.Reset()
		})
	}
}
