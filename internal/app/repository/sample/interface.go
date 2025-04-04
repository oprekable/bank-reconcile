package sample

import (
	"context"
	"time"
)

//go:generate mockery --name "Repository" --output "./_mock" --outpkg "_mock"
type Repository interface {
	Pre(
		ctx context.Context,
		listBank []string,
		startDate time.Time,
		toDate time.Time,
		limitTrxData int64,
		matchPercentage int,
	) (err error)

	GetTrx(ctx context.Context) (returnData []TrxData, err error)
	Post(ctx context.Context) (err error)
	Close() (err error)
}
