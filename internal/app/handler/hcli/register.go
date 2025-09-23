package hcli

import (
	"io"
	"os"

	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/noop"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/process"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/sample"
)

var outPutHandlerWriter io.Writer = os.Stdout

var (
	handlers = append(commonHandlers, applicationHandlers...)

	commonHandlers = []Handler{
		noop.NewHandler(outPutHandlerWriter),
	}

	applicationHandlers = []Handler{
		process.NewHandler(outPutHandlerWriter),
		sample.NewHandler(outPutHandlerWriter),
	}
)
