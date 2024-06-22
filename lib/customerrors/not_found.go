package customerror

import (
	"fmt"

	"github.com/pkg/errors"
)

type notFound struct {
	TraceableError
}

type NotFoundError interface {
	error
	IsNotFoundError() bool
}

func (e *notFound) IsNotFoundError() bool { return true }

func NewNotFoundError(msg string) (err error) {
	err = errors.New(msg)
	return &notFound{TraceableError{err}}
}

func NewNotFoundErrorf(format string, a ...interface{}) (err error) {
	err = errors.New(fmt.Sprintf(format, a...))
	return &notFound{TraceableError{err}}
}
