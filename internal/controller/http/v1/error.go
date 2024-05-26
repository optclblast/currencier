package v1

import (
	"errors"
	"net/http"
)

var (
	ErrorDateRequired = errors.New("date required")
	ErrorDateInvalid  = errors.New("date invalid")
	ErrorValRequired  = errors.New("val required")
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func buildApiError(code int, message string) apiError {
	return apiError{
		Code:    code,
		Message: message,
	}
}

func mapError(err error) apiError {
	switch {
	case errors.Is(err, ErrorDateInvalid):
		return buildApiError(http.StatusBadRequest, "Date Invalid")
	case errors.Is(err, ErrorDateRequired):
		return buildApiError(http.StatusBadRequest, "Date Required")
	case errors.Is(err, ErrorValRequired):
		return buildApiError(http.StatusBadRequest, "Val Required")
	default:
		return buildApiError(http.StatusInternalServerError, "Internal Server Error")
	}
}
