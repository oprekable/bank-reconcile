package process

import (
	"github.com/google/wire"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func ProviderSvc(
	comp *component.Components,
	repo *repository.Repositories,
	parserRegistry *banks.ParserRegistry, // <-- Dependensi baru ditambahkan
) *Svc {
	return NewSvc(comp, repo, parserRegistry) // <-- Diteruskan ke konstruktor
}

var Set = wire.NewSet(
	ProviderSvc,
	wire.Bind(new(ServiceGenerator), new(*Svc)),
)
