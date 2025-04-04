package log

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.chromium.org/luci/common/clock/testclock"
)

func TestAddErr(t *testing.T) {
	timeCtx, _ := testclock.UseTime(context.Background(), time.Unix(1742017753, 0))

	type args struct {
		ctx context.Context
		err error
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Err nil",
			args: args{
				ctx: context.WithValue(timeCtx, StartTime, time.Unix(1742017751, 0)),
				err: nil,
			},
			want: []string{`{"level":"info","uptime":"2s","message":"Test"}`},
		},
		{
			name: "Err not nil",
			args: args{
				ctx: context.WithValue(timeCtx, StartTime, time.Unix(1742017751, 0)),
				err: errors.New("test error"),
			},
			want: []string{
				`{"level":"info","1742017753000000000":{"file":"`,
				`","line":`,
				`,"function":"`,
				`"},"1742017753000000000":"test error","uptime":"2s","message":"Test"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			loggerCtx := zerolog.New(&bf).Hook(UptimeHook{}).WithContext(tt.args.ctx)
			AddErr(loggerCtx, tt.args.err)
			Msg(loggerCtx, "Test")
			got := bf.String()
			for _, want := range tt.want {
				if !strings.Contains(strings.TrimRight(got, "\n"), want) {
					t.Errorf("AddErr() output = %v, want %v", got, want)
				}
			}

			bf.Reset()
		})
	}
}

func TestErr(t *testing.T) {
	timeCtx, _ := testclock.UseTime(context.Background(), time.Unix(1742017752, 0))
	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Err nil",
			args: args{
				ctx: context.WithValue(timeCtx, StartTime, time.Unix(1742017751, 0)),
				err: nil,
			},
			want: []string{`{"level":"info","uptime":"1s","message":"Test"}`},
		},
		{
			name: "Err not nil",
			args: args{
				ctx: context.WithValue(timeCtx, StartTime, time.Unix(1742017751, 0)),
				err: errors.New("test error"),
			},
			want: []string{`{"level":"error","error":"test error","caller":"`, `","uptime":"1s","message":"Test"}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			loggerCtx := zerolog.New(&bf).Hook(UptimeHook{}).WithContext(tt.args.ctx)
			Err(loggerCtx, "Test", tt.args.err)
			got := bf.String()
			for _, want := range tt.want {
				if !strings.Contains(strings.TrimRight(got, "\n"), want) {
					t.Errorf("Err() output = %v, want %v", got, want)
				}
			}

			bf.Reset()
		})
	}
}

func TestMsg(t *testing.T) {
	timeCtx, _ := testclock.UseTime(context.Background(), time.Unix(1742017752, 0))
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Ok",
			args: args{
				ctx: context.WithValue(timeCtx, StartTime, time.Unix(1742017751, 0)),
			},
			want: `{"level":"info","uptime":"1s","message":"Test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			loggerCtx := zerolog.New(&bf).Hook(UptimeHook{}).WithContext(tt.args.ctx)
			Msg(loggerCtx, "Test")

			if got := bf.String(); strings.TrimRight(got, "\n") != tt.want {
				t.Errorf("Msg() output = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

func TestUptimeHookRun(t *testing.T) {
	timeCtx, _ := testclock.UseTime(context.Background(), time.Unix(1742017752, 0))
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty StartTime",
			args: args{
				ctx: timeCtx,
			},
			want: `{"level":"info","uptime":"","message":"Test"}`,
		},
		{
			name: "Non Empty StartTime",
			args: args{
				ctx: context.WithValue(timeCtx, StartTime, time.Unix(1742017751, 0)),
			},
			want: `{"level":"info","uptime":"1s","message":"Test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			loggerCtx := zerolog.New(&bf).Hook(UptimeHook{}).WithContext(tt.args.ctx)
			zerolog.Ctx(loggerCtx).Info().Ctx(loggerCtx).Msg("Test")

			if got := bf.String(); strings.TrimRight(got, "\n") != tt.want {
				t.Errorf("UptimeHook() output = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}
