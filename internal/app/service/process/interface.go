package process

import (
	"context"

	"github.com/schollz/progressbar/v3"

	"github.com/spf13/afero"
)

//go:generate mockery --name "Service" --output "./_mock" --outpkg "_mock"
type Service interface {
	GenerateReconciliation(ctx context.Context, fs afero.Fs, bar *progressbar.ProgressBar) (returnSummary ReconciliationSummary, err error)
}
