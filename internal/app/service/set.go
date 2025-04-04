package service

import (
	"github.com/oprekable/bank-reconcile/internal/app/service/process"
	"github.com/oprekable/bank-reconcile/internal/app/service/sample"

	"github.com/google/wire"
)

func NewServices(
	svcSample sample.Service,
	svcProcess process.Service,
) *Services {
	return &Services{
		SvcSample:  svcSample,
		SvcProcess: svcProcess,
	}
}

var Set = wire.NewSet(
	sample.Set,
	process.Set,
	NewServices,
)
