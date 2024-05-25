package v1

import (
	"net/http"
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
	default:
		return buildApiError(http.StatusInternalServerError, "Internal Server Error")
	}
}
