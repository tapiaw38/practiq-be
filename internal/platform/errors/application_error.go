package errors

import (
	"context"
	"encoding/json"

	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ApplicationError interface {
	InternalCode() string
	StatusCode() int
	Message() string
	OriginalMessage() string
	AddExtraFields(fields map[string]interface{}) ApplicationError
	IsBasedOn(mappings.ErrorDetails) bool
	Log(context.Context)
	json.Marshaler
	error
}

type applicationError struct {
	errorDetails    mappings.ErrorDetails
	originalMessage string
	extraFields     map[string]interface{}
}

func (r *applicationError) InternalCode() string {
	return r.errorDetails.InternalCode
}

func (r *applicationError) StatusCode() int {
	return r.errorDetails.StatusCode
}

func (r *applicationError) Message() string {
	return r.errorDetails.Message
}

func (r *applicationError) OriginalMessage() string {
	return r.originalMessage
}

func (r *applicationError) Error() string {
	return r.errorDetails.Message
}

func (r *applicationError) AddExtraFields(fields map[string]interface{}) ApplicationError {
	if r.extraFields == nil {
		r.extraFields = make(map[string]interface{})
	}
	for k, v := range fields {
		r.extraFields[k] = v
	}
	return r
}

func (r *applicationError) IsBasedOn(errorDetails mappings.ErrorDetails) bool {
	return r.errorDetails == errorDetails
}

func (r *applicationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    r.InternalCode(),
		Message: r.Message(),
	})
}

func NewApplicationError(
	code mappings.ErrorDetails,
	originalError error,
) ApplicationError {
	var originalMessage string

	if originalError != nil {
		originalMessage = originalError.Error()
	}

	return &applicationError{code, originalMessage, make(map[string]interface{})}
}

func NewInternalError(err error) ApplicationError {
	return NewApplicationError(mappings.InternalServerError, err)
}

func NewBadRequestError(msg string) ApplicationError {
	return NewApplicationError(mappings.ErrorDetails{
		InternalCode: "common:bad-request",
		StatusCode:   400,
		Message:      msg,
	}, nil)
}

func NewNotFoundError(msg string) ApplicationError {
	return NewApplicationError(mappings.ErrorDetails{
		InternalCode: "common:not-found",
		StatusCode:   404,
		Message:      msg,
	}, nil)
}

func NewUnauthorizedError() ApplicationError {
	return NewApplicationError(mappings.UnauthorizedError, nil)
}

func NewForbiddenError() ApplicationError {
	return NewApplicationError(mappings.ForbiddenError, nil)
}
