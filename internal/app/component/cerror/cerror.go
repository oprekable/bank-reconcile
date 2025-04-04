package cerror

import "github.com/oprekable/bank-reconcile/internal/app/err/core"

type Error struct {
	errors []error
}

type ErType []core.ErrorType

func NewError(erType ErType) *Error {
	var e []error
	for k := range erType {
		e = append(e, erType[k].Error())
	}

	return &Error{
		errors: e,
	}
}

func (e *Error) GetErrors() []error {
	return e.errors
}
