package cerror

import (
	"github.com/oprekable/bank-reconcile/internal/app/err/core"

	"github.com/google/wire"
)

func ProvideErType(errType []core.ErrorType) ErType {
	return errType
}

var Set = wire.NewSet(
	ProvideErType,
	NewError,
)
