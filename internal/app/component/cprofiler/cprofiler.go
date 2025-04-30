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
	ctx      context.Context
	logger   *clogger.Logger
	profiler map[string]interface{ Stop() }
}

func NewProfiler(logger *clogger.Logger) *Profiler {
	return &Profiler{
		ctx:      logger.GetLogger().With().Str("component", "Profiler").Ctx(context.Background()).Logger().WithContext(logger.GetCtx()),
		logger:   logger,
		profiler: make(map[string]interface{ Stop() }),
	}
}

func (p *Profiler) StartProfiler() {
	log.Msg(p.ctx, "[profiler] starting profiler")
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
	for k := range p.profiler {
		p.profiler[k].Stop()
		log.Msg(p.ctx, fmt.Sprintf("[profiler] stop profiler - %s", k))
	}
}
