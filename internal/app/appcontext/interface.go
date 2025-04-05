package appcontext

import (
	"context"

	"github.com/oprekable/bank-reconcile/internal/app/component"
)

type IAppContext interface {
	GetCtx() context.Context
	GetComponents() *component.Components
	Start()
	Shutdown()
}
