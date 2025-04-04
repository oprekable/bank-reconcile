package server

import "golang.org/x/sync/errgroup"

//go:generate mockery --name "IServer" --output "./_mock" --outpkg "_mock"
type IServer interface {
	Name() string
	Start(eg *errgroup.Group)
	Shutdown()
}
