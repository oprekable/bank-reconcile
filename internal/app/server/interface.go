package server

import "golang.org/x/sync/errgroup"

type IServer interface {
	Name() string
	Start(eg *errgroup.Group)
	Shutdown()
}
