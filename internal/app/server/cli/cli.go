package cli

import (
	"context"
	"errors"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
	"golang.org/x/sync/errgroup"
)

const name = "cli"

type Cli struct {
	ctx      context.Context
	comp     *component.Components
	svc      *service.Services
	repo     *repository.Repositories
	handlers []hcli.Handler
}

func NewCli(
	comp *component.Components,
	svc *service.Services,
	repo *repository.Repositories,
	handlers []hcli.Handler,
) (*Cli, error) {
	returnData := &Cli{
		ctx:      comp.Logger.GetCtx(),
		comp:     comp,
		svc:      svc,
		repo:     repo,
		handlers: handlers,
	}

	for k := range handlers {
		handlers[k].SetComponents(comp)
		handlers[k].SetServices(svc)
		handlers[k].SetRepositories(repo)
	}

	return returnData, nil
}

func (c *Cli) Name() string {
	return name
}

func (c *Cli) Start(eg *errgroup.Group) {
	eg.Go(func() (err error) {
		ctx := c.ctx

		for k := range c.handlers {
			if c.handlers[k].Name() == c.comp.Config.Reconciliation.Action {
				err = c.handlers[k].Exec()
				if err != nil {
					log.Err(ctx, "error", err)
				} else {
					err = errors.New("done")
				}

				break
			}
		}

		return
	})
}

func (c *Cli) Shutdown() {
	log.Msg(c.ctx, "["+name+"] shutdown")
}
