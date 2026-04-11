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

const (
	MUTEX_PPROF_FILE     = "./mutex.pprof"
	CPU_PPROF_FILE       = "./cpu.pprof"
	MEM_PPROF_FILE       = "./mem.pprof"
	BLOCK_PPROF_FILE     = "./block.pprof"
	TRACE_PPROF_FILE     = "./trace.pprof"
	GOROUTINE_PPROF_FILE = "./goroutine.pprof"
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
				logger:   logger,
				profiler: make(map[string]interface{ Stop() }),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProfiler(tt.args.logger); !reflect.DeepEqual(got.profiler, tt.want.profiler) || !reflect.DeepEqual(got.logger, tt.want.logger) {
				t.Errorf("NewProfiler() = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

func TestProfilerStartProfiler(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		logger *clogger.Logger
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
				logger: logger,
			},
			wantLog: `[profiler] starting profiler`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProfiler(tt.fields.logger)

			p.StartProfiler()

			checkPprofFiles(t, []string{
				CPU_PPROF_FILE, MEM_PPROF_FILE, MUTEX_PPROF_FILE, BLOCK_PPROF_FILE,
				TRACE_PPROF_FILE, GOROUTINE_PPROF_FILE,
			})

			cleanupPprofFiles(t, []string{
				CPU_PPROF_FILE, MEM_PPROF_FILE, MUTEX_PPROF_FILE, BLOCK_PPROF_FILE,
				TRACE_PPROF_FILE, GOROUTINE_PPROF_FILE,
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
		logger *clogger.Logger
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
				logger: logger,
			},
			wantLog: `[profiler] stop profiler`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProfiler(tt.fields.logger)
			p.StartProfiler()
			bf.Reset()
			p.StopProfiler()

			checkPprofFiles(t, []string{
				CPU_PPROF_FILE, MEM_PPROF_FILE, MUTEX_PPROF_FILE, BLOCK_PPROF_FILE,
				TRACE_PPROF_FILE, GOROUTINE_PPROF_FILE,
			})

			cleanupPprofFiles(t, []string{
				CPU_PPROF_FILE, MEM_PPROF_FILE, MUTEX_PPROF_FILE, BLOCK_PPROF_FILE,
				TRACE_PPROF_FILE, GOROUTINE_PPROF_FILE,
			})

			got := bf.String()
			if !strings.Contains(got, tt.wantLog) {
				t.Errorf("StopProfiler() log = %v, wantlog %v", got, tt.wantLog)
			}

			bf.Reset()
		})
	}
}
