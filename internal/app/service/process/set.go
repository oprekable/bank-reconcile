package process

import (
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"

	"github.com/google/wire"
)

func ProviderSvc(
	comp *component.Components,
	repo *repository.Repositories,
) *Svc {
	return NewSvc(comp, repo)
}

var Set = wire.NewSet(
	ProviderSvc,
	wire.Bind(new(Service), new(*Svc)),
)
