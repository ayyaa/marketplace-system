package customerror

import (
	"net/http"

	"github.com/pkg/errors"
)

type TraceableError struct {
	err error
}

func (e TraceableError) Error() string { return e.err.Error() }
func (e TraceableError) Cause() error  { return e.err }

func New(msg string) error {
	return errors.New(msg)
}

func GetStatusCode(err error) int {
	if _, ok := err.(NotFoundError); ok {
		return http.StatusNotFound
	} else if _, ok := err.(BadRequestError); ok {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError

}
