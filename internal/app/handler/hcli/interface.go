package hcli

import (
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
)

type Handler interface {
	SetComponents(c *component.Components)
	SetServices(s *service.Services)
	SetRepositories(r *repository.Repositories)
	Exec() error
	Name() string
}
