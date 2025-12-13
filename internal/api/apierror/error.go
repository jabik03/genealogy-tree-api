package apierror

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int    // (400, 404, 500...)
	Message    string // Сообщение для пользователя
	Err        error  // Оригинальная ошибка (для логирования)
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// Конструкторы для популярных ошибок

func BadRequest(message string, err error) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Err:        err,
	}
}

func NotFound(message string, err error) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Err:        err,
	}
}

func InternalError(message string, err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Err:        err,
	}
}
