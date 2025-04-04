package hcli

import (
	"os"

	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/noop"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/process"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/sample"
)

var (
	Handlers = append(commonHandlers, applicationHandlers...)

	commonHandlers = []Handler{
		noop.NewHandler(os.Stdout),
	}

	applicationHandlers = []Handler{
		process.NewHandler(os.Stdout),
		sample.NewHandler(os.Stdout),
	}
)
