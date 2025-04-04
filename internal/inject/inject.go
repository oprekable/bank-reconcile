//go:build wireinject
// +build wireinject

package inject

import (
	"embed"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/err/core"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/server"
	"github.com/oprekable/bank-reconcile/internal/app/service"

	"context"

	"github.com/google/wire"
)

func WireApp(
	ctx context.Context,
	embedFS *embed.FS,
	appName cconfig.AppName,
	tz cconfig.TimeZone,
	errType []core.ErrorType,
	isShowLog clogger.IsShowLog,
	dBPath csqlite.DBPath,
) (*appcontext.AppContext, func(), error) {
	wire.Build(
		component.Set,
		repository.Set,
		service.Set,
		server.Set,
		appcontext.NewAppContext,
	)

	return new(appcontext.AppContext), nil, nil
}
