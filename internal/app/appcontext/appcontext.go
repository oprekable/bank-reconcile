package appcontext

import (
	"context"
	"embed"
	"fmt"
	"os"

	profile "github.com/bygui86/multi-profile/v2"

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
	profiler     map[string]interface{ Stop() }
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

func (a *AppContext) startProfiler() {
	if a.components.Config.IsProfilerActive {
		log.Msg(a.GetCtx(), "[profiler] starting profiler")
		a.profiler = make(map[string]interface{ Stop() })
		dir, _ := os.Getwd()

		a.profiler["CPUProfile"] = profile.CPUProfile(
			&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
		).Start()

		a.profiler["BlockProfile"] = profile.BlockProfile(
			&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
		).Start()

		a.profiler["GoroutineProfile"] = profile.GoroutineProfile(
			&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
		).Start()

		a.profiler["MutexProfile"] = profile.MutexProfile(
			&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true}).Start()

		a.profiler["MemProfile"] = profile.MemProfile(
			&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
		).Start()

		a.profiler["TraceProfile"] = profile.TraceProfile(
			&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
		).Start()
	}
}

func (a *AppContext) stopProfiler() {
	for k := range a.profiler {
		a.profiler[k].Stop()
		log.Msg(a.GetCtx(), fmt.Sprintf("[profiler] stop profiler - %s", k))
	}
}

func (a *AppContext) Start() {
	atexit.Add(a.Shutdown)
	a.eg.Go(func() error {
		log.Msg(a.GetCtx(), "[application] start")
		a.startProfiler()
		return shutdown.TermSignalTrap().Wait(a.ctx, func() {
			defer func() {
				if r := recover(); r != nil {
					errRecovery := fmt.Errorf("recovered from panic: %s", r)
					log.AddErr(context.Background(), errRecovery)
					return
				}
			}()

			atexit.AtExit()

			select {
			case <-a.ctx.Done():
				{
					if context.Cause(a.ctx).Error() == "done" {
						os.Exit(0)
					} else {
						os.Exit(1)
					}
				}
			default:
			}
		})
	})

	if a.servers != nil {
		a.servers.Run(a.eg)
	}

	log.Err(a.GetCtx(), "[application] shutdown", a.eg.Wait())
}

func (a *AppContext) Shutdown() {
	a.stopProfiler()
	log.Msg(a.GetCtx(), "[application] shutdown")
}
