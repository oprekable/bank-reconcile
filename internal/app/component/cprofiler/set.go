package cprofiler

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewProfiler,
)
