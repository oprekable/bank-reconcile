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

func TestProviderLogger(t *testing.T) {
	var bfIsShowLogTrue bytes.Buffer
	var bfIsShowLogFalse bytes.Buffer
	timeCtx, _ := testclock.UseTime(context.Background(), time.Unix(1742017753, 0))

	type args struct {
		ctx             context.Context
		logShowWriter   LogShowWriter
		logNoShowWriter LogNoShowWriter
		isShowLog       IsShowLog
	}

	tests := []struct {
		name               string
		args               args
		wantIsShowLogTrue  string
		wantIsShowLogFalse string
	}{
		{
			name: "Ok - isShowLog false",
			args: args{
				ctx:             timeCtx,
				isShowLog:       false,
				logShowWriter:   &bfIsShowLogTrue,
				logNoShowWriter: &bfIsShowLogFalse,
			},
			wantIsShowLogTrue:  `""`,
			wantIsShowLogFalse: `"\x1b[90m2025-03-15T00:00:00+07:00\x1b[0m | INFO  | *** foo **** | uptime:"`,
		},
		{
			name: "Ok - isShowLog true",
			args: args{
				ctx:             timeCtx,
				isShowLog:       true,
				logShowWriter:   &bfIsShowLogTrue,
				logNoShowWriter: &bfIsShowLogFalse,
			},
			wantIsShowLogTrue:  `"\x1b[90m2025-03-15T00:00:00+07:00\x1b[0m | INFO  | *** foo **** | uptime:"`,
			wantIsShowLogFalse: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := ProviderLogger(tt.args.ctx, tt.args.isShowLog, tt.args.logShowWriter, tt.args.logNoShowWriter)
			zerolog.TimeFieldFormat = time.DateOnly
			zerolog.TimestampFunc = func() time.Time {
				return clock.Get(timeCtx).Now()
			}

			l.log.Info().Msg("foo")

			gotIsShowLogTrue := strconv.Quote(strings.TrimRight(bfIsShowLogTrue.String(), "\n"))
			gotIsShowLogFalse := strconv.Quote(strings.TrimRight(bfIsShowLogFalse.String(), "\n"))

			if !reflect.DeepEqual(gotIsShowLogTrue, tt.wantIsShowLogTrue) {
				t.Errorf("ProviderLogger() gotIsShowLogTrue = %v, wantIsShowLogTrue %v", gotIsShowLogTrue, tt.wantIsShowLogTrue)
			}

			if !reflect.DeepEqual(gotIsShowLogFalse, tt.wantIsShowLogFalse) {
				t.Errorf("ProviderLogger() gotIsShowLogFalse = %v, wantIsShowLogTrue %v", gotIsShowLogFalse, tt.wantIsShowLogFalse)
			}

			bfIsShowLogFalse.Reset()
			bfIsShowLogTrue.Reset()
		})
	}
}
