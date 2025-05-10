package models

import "net/http"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

func ValidationError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

func UnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, message)
}

func NotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message)
}

func InternalServerError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}
