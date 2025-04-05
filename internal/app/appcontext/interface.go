package appcontext

import (
	"context"

	"github.com/oprekable/bank-reconcile/internal/app/component"
)

//go:generate mockery --name "IAppContext" --output "./_mock" --outpkg "_mock"
type IAppContext interface {
	GetCtx() context.Context
	GetComponents() *component.Components
	Start()
	Shutdown()
}
