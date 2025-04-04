package err

import "github.com/oprekable/bank-reconcile/internal/app/err/core"

// RegisteredErrorType Register new errors here!
var RegisteredErrorType = []core.ErrorType{
	core.CErrInternal,
	core.CErrDBConn,
}
