package appcontext

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/server"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/pkg/shutdown"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"

	"github.com/bygui86/multi-profile/v2"

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

func (a *AppContext) Start() {
	atexit.Add(a.Shutdown)
	a.eg.Go(func() error {
		var profiler map[string]interface{ Stop() }
		if a.components.Config.IsProfilerActive {
			profiler = make(map[string]interface{ Stop() })
			dir, _ := os.Getwd()

			profiler["CPUProfile"] = profile.CPUProfile(
				&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
			).Start()

			profiler["BlockProfile"] = profile.BlockProfile(
				&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
			).Start()

			profiler["GoroutineProfile"] = profile.GoroutineProfile(
				&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
			).Start()

			profiler["MutexProfile"] = profile.MutexProfile(
				&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true}).Start()

			profiler["MemProfile"] = profile.MemProfile(
				&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
			).Start()

			profiler["TraceProfile"] = profile.TraceProfile(
				&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
			).Start()
		}

		log.Msg(a.GetCtx(), "[start] application")
		return shutdown.TermSignalTrap().Wait(a.ctx, func() {
			defer func() {
				if r := recover(); r != nil {
					errRecovery := fmt.Errorf("recovered from panic: %s", r)
					log.AddErr(context.Background(), errRecovery)
					return
				}
			}()

			atexit.AtExit()

			if context.Cause(a.ctx).Error() == "done" {
				for k := range profiler {
					if profiler[k] != nil {
						profiler[k].Stop()
					}
				}

				os.Exit(0)
			}
		})
	})

	if a.servers != nil {
		a.servers.Run(a.eg)
	}

	if err := a.eg.Wait(); err != nil {
		log.Err(a.GetCtx(), "[shutdown] application", err)
	}
}

func (a *AppContext) Shutdown() {
	log.Msg(a.GetCtx(), "[shutdown] application")
}
