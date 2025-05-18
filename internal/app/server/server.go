package server

import (
	"reflect"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"golang.org/x/sync/errgroup"
)

type CliServer IServer
type Server struct {
	Cli CliServer
}

func NewServer(
	cli CliServer,
) *Server {
	return &Server{
		Cli: cli,
	}
}

func (s *Server) Run(eg *errgroup.Group) {
	v := reflect.ValueOf(*s)
	for i := 0; i < v.NumField(); i++ {
		if s, ok := v.Field(i).Interface().(IServer); ok {
			atexit.Add(s.Shutdown)
			s.Start(eg)
		}
	}
}
