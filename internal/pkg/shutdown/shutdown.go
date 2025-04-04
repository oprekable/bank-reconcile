package shutdown

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

const (
	ErrTermSig = "termination signal caught"
)

type SignalTrap chan os.Signal

func TermSignalTrap() SignalTrap {
	trap := SignalTrap(make(chan os.Signal, 1))
	signal.Notify(trap, syscall.SIGINT, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSEGV)

	return trap
}

func (t SignalTrap) Wait(ctx context.Context, f func()) error {
	select {
	case <-t:
		{
			f()
			return errors.New(ErrTermSig)
		}
	case <-ctx.Done():
		{
			f()
			return ctx.Err()
		}
	}
}
