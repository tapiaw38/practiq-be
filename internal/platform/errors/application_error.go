package errors

import (
	"fmt"
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ApplicationError interface {
	Error() string
	StatusCode() int
	Log(ctx interface{})
}

type applicationError struct {
	Code            string `json:"code"`
	Message         string `json:"message"`
	OriginalMessage string `json:"-"`
	statusCode      int
}

func NewApplicationError(details mappings.ErrorDetails, originalErr error) ApplicationError {
	originalMsg := ""
	if originalErr != nil {
		originalMsg = originalErr.Error()
	}
	return &applicationError{
		Code:            details.Code,
		Message:         details.Message,
		OriginalMessage: originalMsg,
		statusCode:      details.StatusCode,
	}
}

func (e *applicationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.OriginalMessage)
}

func (e *applicationError) StatusCode() int {
	return e.statusCode
}

func (e *applicationError) Log(ctx interface{}) {
	fmt.Printf("ERROR [%d] %s: %s | original: %s\n", e.statusCode, e.Code, e.Message, e.OriginalMessage)
}

func NewInternalError(err error) ApplicationError {
	return NewApplicationError(mappings.InternalServerError, err)
}

func NewBadRequestError(msg string) ApplicationError {
	return &applicationError{
		Code:       "common:bad-request",
		Message:    msg,
		statusCode: http.StatusBadRequest,
	}
}

func NewNotFoundError(msg string) ApplicationError {
	return &applicationError{
		Code:       "common:not-found",
		Message:    msg,
		statusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError() ApplicationError {
	return NewApplicationError(mappings.UnauthorizedError, nil)
}

func NewForbiddenError() ApplicationError {
	return NewApplicationError(mappings.ForbiddenError, nil)
}
