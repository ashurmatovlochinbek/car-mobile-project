package httpErrors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	BadRequest             = errors.New("bad request")
	NotFound               = errors.New("not Found")
	Unauthorized           = errors.New("unauthorized")
	RequestTimeoutError    = errors.New("request Timeout")
	ExistsPhoneNumberError = errors.New("user with given phone number already exists")
)

type RestErr interface {
	Status() int
	Error() string
}

type RestError struct {
	ErrStatus int    `json:"status,omitempty"`
	ErrError  string `json:"error,omitempty"`
}

func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s ", e.ErrStatus, e.ErrError)
}

func (e RestError) Status() int {
	return e.ErrStatus
}

func NewRestError(status int, err string) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
	}
}

func ParseErrors(err error) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, fmt.Sprintf("%s", NotFound.Error()))
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, fmt.Sprintf("%s", RequestTimeoutError.Error()))
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewRestError(http.StatusBadRequest, fmt.Sprintf("%s", BadRequest.Error()))
	case strings.Contains(err.Error(), "UUID"):
		return NewRestError(http.StatusBadRequest, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), "token"):
		return NewRestError(http.StatusUnauthorized, fmt.Sprintf("%s", Unauthorized.Error()))
	default:
		return NewRestError(http.StatusInternalServerError, err.Error())
	}
}

func parseSqlErrors(err error) RestErr {
	if strings.Contains(err.Error(), "23505") {
		return NewRestError(http.StatusBadRequest, fmt.Sprintf("%s", ExistsPhoneNumberError.Error()))
	}
	return NewRestError(http.StatusBadRequest, fmt.Sprintf("%s", BadRequest.Error()))
}
