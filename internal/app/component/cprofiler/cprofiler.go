package cprofiler

import (
	"context"
	"fmt"
	"os"

	profile "github.com/bygui86/multi-profile/v2"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
)

type Profiler struct {
	ctxFn    func(ctx context.Context) context.Context
	logger   *clogger.Logger
	profiler map[string]interface{ Stop() }
}

func NewProfiler(logger *clogger.Logger) *Profiler {
	returnData := &Profiler{
		logger:   logger,
		profiler: make(map[string]interface{ Stop() }),
	}

	returnData.ctxFn = func(ctx context.Context) context.Context {
		return returnData.
			logger.
			GetLogger().
			With().
			Str("component", "Profiler").
			Ctx(ctx).
			Logger().
			WithContext(returnData.logger.GetCtx())
	}

	return returnData
}

func (p *Profiler) StartProfiler() {
	ctx := p.ctxFn(context.Background())
	log.Msg(ctx, "[profiler] starting profiler")
	p.profiler = make(map[string]interface{ Stop() })
	dir, _ := os.Getwd()

	p.profiler["CPUProfile"] = profile.CPUProfile(
		&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
	).Start()

	p.profiler["BlockProfile"] = profile.BlockProfile(
		&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
	).Start()

	p.profiler["GoroutineProfile"] = profile.GoroutineProfile(
		&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
	).Start()

	p.profiler["MutexProfile"] = profile.MutexProfile(
		&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true}).Start()

	p.profiler["MemProfile"] = profile.MemProfile(
		&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
	).Start()

	p.profiler["TraceProfile"] = profile.TraceProfile(
		&profile.Config{Path: dir, EnableInterruptHook: true, Quiet: true},
	).Start()
}

func (p *Profiler) StopProfiler() {
	ctx := p.ctxFn(context.Background())
	for k := range p.profiler {
		p.profiler[k].Stop()
		log.Msg(ctx, fmt.Sprintf("[profiler] stop profiler - %s", k))
	}
}
