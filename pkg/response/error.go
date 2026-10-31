package response

import "net/http"

type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

func NewNotFoundError(message string) *AppError {
	return &AppError{StatusCode: 404, Message: message, Code: http.StatusText(404)}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{StatusCode: 400, Message: message, Code: http.StatusText(400)}
}

func NewError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message, Code: http.StatusText(statusCode)}
}
