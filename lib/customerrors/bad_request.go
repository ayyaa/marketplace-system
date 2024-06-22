package customerror

import (
	"fmt"

	"github.com/pkg/errors"
)

type badRequest struct {
	TraceableError
}

type BadRequestError interface {
	error
	IsBadRequestError() bool
}

func (e *badRequest) IsBadRequestError() bool { return true }

func NewBadRequestError(msg string) (err error) {
	err = errors.New(msg)
	return &badRequest{TraceableError{err}}
}

func NewBadRequestErrorf(format string, a ...interface{}) (err error) {
	err = errors.New(fmt.Sprintf(format, a...))
	return &badRequest{TraceableError{err}}
}
