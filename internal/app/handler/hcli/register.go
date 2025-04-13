package hcli

import (
	"io"
	"os"

	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/noop"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/process"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/sample"
)

var OutPutHandlerWriter io.Writer = os.Stdout

var (
	Handlers = append(commonHandlers, applicationHandlers...)

	commonHandlers = []Handler{
		noop.NewHandler(OutPutHandlerWriter),
	}

	applicationHandlers = []Handler{
		process.NewHandler(OutPutHandlerWriter),
		sample.NewHandler(OutPutHandlerWriter),
	}
)
