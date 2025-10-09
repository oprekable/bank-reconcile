package appcontext

import (
	"context"
	"embed"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/server"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/pkg/shutdown"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
	"golang.org/x/sync/errgroup"
)

type IsProfilerActive bool
type AppContext struct {
	ctx          context.Context
	ctxCancel    context.CancelFunc
	eg           *errgroup.Group
	embedFS      *embed.FS
	repositories *repository.Repositories
	services     *service.Services
	components   *component.Components
	servers      *server.Server
}

var _ IAppContext = (*AppContext)(nil)

// NewAppContext initiate AppContext object
func NewAppContext(
	ctx context.Context,
	embedFS *embed.FS,
	repository *repository.Repositories,
	services *service.Services,
	components *component.Components,
	servers *server.Server,
) (*AppContext, func()) {
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)

	return &AppContext{
		ctx:          ctx,
		ctxCancel:    cancel,
		eg:           eg,
		embedFS:      embedFS,
		repositories: repository,
		services:     services,
		components:   components,
		servers:      servers,
	}, cancel
}

func (a *AppContext) GetCtx() context.Context {
	if a.components != nil && a.components.Logger != nil {
		return a.components.Logger.GetCtx()
	}

	return a.ctx
}

func (a *AppContext) GetComponents() *component.Components {
	return a.components
}

func (a *AppContext) Start() error {
	atexit.Add(a.Shutdown)
	a.eg.Go(func() error {
		log.Msg(a.GetCtx(), "[application] start")

		if a.components.Config != nil && a.components.Config.IsProfilerActive {
			a.components.Profiler.StartProfiler()
		}

		return shutdown.TermSignalTrap().Wait(a.ctx, func() {
			atexit.AtExit()
		})
	})

	if a.servers != nil {
		a.servers.Run(a.eg)
	}

	err := a.eg.Wait()
	if context.Cause(a.ctx).Error() == "done" {
		err = nil
	}

	log.Err(a.GetCtx(), "[application] exit", err)
	return err
}

func (a *AppContext) Shutdown() {
	a.components.Profiler.StopProfiler()
	log.Msg(a.GetCtx(), "[application] shutdown")
}
