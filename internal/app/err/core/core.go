package core

import (
	"errors"
	"strconv"
)

type ErrorType int

// General Error
const (
	CErrInternal ErrorType = 100000 + iota
	CErrDBConn
)

func (e ErrorType) Error() error {
	return errors.New(strconv.Itoa(int(e)))
}
