package customerror

import (
	"fmt"

	"github.com/pkg/errors"
)

type internal struct {
	TraceableError
}

type InternalError interface {
	error
	IsInternalError() bool
}

func (e *internal) IsInternalError() bool { return true }

func NewInternalError(msg string) (err error) {
	err = errors.New(msg)
	return &internal{TraceableError{err}}
}

func NewInternalErrorf(format string, a ...interface{}) (err error) {
	err = errors.New(fmt.Sprintf(format, a...))
	return &internal{TraceableError{err}}
}
