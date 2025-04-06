package core

// App ..
type App struct {
	Secret           string `default:"-"     mapstructure:"secret"`
	IsShowLog        bool   `default:"true"  mapstructure:"is_show_log"`
	IsDebug          bool   `default:"false" mapstructure:"is_debug"`
	IsProfilerActive bool   `default:"false" mapstructure:"is_profiler_active"`
}
