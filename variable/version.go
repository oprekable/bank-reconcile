package variable

import (
	"fmt"
	"runtime"
)

// this variable value is set on makefile

// GitCommit The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Environment the environment app is running at the moment
var Environment string

// TimeZone default build timezone
var TimeZone string

// Version The main version number that is being run at the moment.
var Version string

// BuildDate the timestamp image is built
var BuildDate = ""

// GoVersion version of GO used
var GoVersion = runtime.Version()

// OsArch which architecture is used when compiled
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
