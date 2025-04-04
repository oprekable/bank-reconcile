package clogger

import (
	"context"
	"io"
	"os"

	"github.com/google/wire"
)

type IsShowLog bool
type LogShowWriter io.Writer
type LogNoShowWriter io.Writer

func ProviderLogger(ctx context.Context, isShowLog IsShowLog, logShowWriter LogShowWriter, logNoShowWriter LogNoShowWriter) *Logger {
	var logWriter io.Writer
	if isShowLog {
		logWriter = logShowWriter
	} else {
		logWriter = logNoShowWriter
	}

	return NewLogger(
		ctx,
		logWriter,
	)
}

var Set = wire.NewSet(
	wire.InterfaceValue(new(LogShowWriter), os.Stdout),
	wire.InterfaceValue(new(LogNoShowWriter), io.Discard),
	ProviderLogger,
)
