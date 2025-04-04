package log

import (
	"context"
	"runtime"
	"strconv"
	"time"

	"go.chromium.org/luci/common/clock"

	"github.com/rs/zerolog"
)

type KeyType string

const (
	// StartTime ...
	StartTime KeyType = "start_time"
)

type UptimeHook struct{}

func (u UptimeHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	startTime := ctx.Value(StartTime)
	uptime := ""

	if startTime != nil {
		uptime = clock.Since(ctx, startTime.(time.Time)).String()
	}

	e.Str("uptime", uptime)
}

type Caller struct {
	File     string
	Function string
	Line     int
}

func New() Caller {
	pc, f, l, _ := runtime.Caller(4)
	return Caller{
		File:     f,
		Line:     l,
		Function: runtime.FuncForPC(pc).Name(),
	}
}

func (lc Caller) MarshalZerologObject(e *zerolog.Event) {
	e.Str("file", lc.File).
		Int("line", lc.Line).
		Str("function", lc.Function)
}

func AddErr(ctx context.Context, er error) {
	if er == nil {
		return
	}

	nowUnixNano := clock.Get(ctx).Now().UnixNano()
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Object(strconv.FormatInt(nowUnixNano, 10), New()).
			AnErr(strconv.FormatInt(nowUnixNano, 10), er)
	})
}

func Msg(ctx context.Context, msg string) {
	zerolog.Ctx(ctx).
		Info().
		Ctx(ctx).
		Msg(msg)
}

func Err(ctx context.Context, msg string, er error) {
	if er == nil {
		zerolog.Ctx(ctx).
			Info().
			Ctx(ctx).
			Msg(msg)
	} else {
		zerolog.Ctx(ctx).
			Err(er).
			Ctx(ctx).
			Caller(2).
			Msg(msg)
	}
}
