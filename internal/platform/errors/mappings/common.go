package mappings

import "net/http"

var (
	RequestBodyParsingError = ErrorDetails{
		Code:       "common:request:body-parsing-error",
		StatusCode: http.StatusBadRequest,
		Message:    "invalid request body",
	}

	InternalServerError = ErrorDetails{
		Code:       "common:internal-server-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}

	UnauthorizedError = ErrorDetails{
		Code:       "common:unauthorized",
		StatusCode: http.StatusUnauthorized,
		Message:    "unauthorized",
	}

	ForbiddenError = ErrorDetails{
		Code:       "common:forbidden",
		StatusCode: http.StatusForbidden,
		Message:    "forbidden",
	}

	NotFoundError = ErrorDetails{
		Code:       "common:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "resource not found",
	}
)
