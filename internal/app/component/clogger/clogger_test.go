package clogger

import (
	"bytes"
	"context"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.chromium.org/luci/common/clock"

	"github.com/rs/zerolog"
	"go.chromium.org/luci/common/clock/testclock"
)

func TestLoggerGetCtx(t *testing.T) {
	type fields struct {
		log zerolog.Logger
	}

	tests := []struct {
		fields fields
		want   context.Context
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				log: zerolog.Logger{},
			},
			want: context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				log: tt.fields.log,
				ctx: context.Background(),
			}

			got := l.GetCtx()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoggerGetLogger(t *testing.T) {
	type fields struct {
		log zerolog.Logger
	}

	tests := []struct {
		fields fields
		want   zerolog.Logger
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				log: zerolog.Logger{},
			},
			want: zerolog.Logger{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				log: tt.fields.log,
				ctx: context.Background(),
			}

			if got := l.GetLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	ctx := context.Background()
	timeCtx, _ := testclock.UseTime(ctx, time.Unix(1742017753, 0))

	type args struct {
		ctxFn func(context.Context) context.Context
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Ok",
			args: args{
				ctxFn: func(c context.Context) context.Context {
					timeCtx, _ := testclock.UseTime(c, time.Unix(1742017753, 0))
					return timeCtx
				},
			},
			want: `"\x1b[90m2025-03-15T00:00:00+07:00\x1b[0m | INFO  | *** foo **** | uptime:"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logWriter := &bytes.Buffer{}
			l := NewLogger(tt.args.ctxFn(ctx), logWriter)
			zerolog.TimeFieldFormat = time.DateOnly
			zerolog.TimestampFunc = func() time.Time {
				return clock.Get(timeCtx).Now()
			}
			l.log.Info().Msg("foo")
			got := strings.TrimRight(logWriter.String(), "\n")
			got = strconv.Quote(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}

			logWriter.Reset()
		})
	}
}
