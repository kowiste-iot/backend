package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeBadRequest     ErrorType = "BAD_REQUEST"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeInternal       ErrorType = "INTERNAL"
	ErrorTypeValidation     ErrorType = "VALIDATION"
	ErrorTypeAlreadyExists  ErrorType = "ALREADY_EXISTS"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
)

type AppError struct {
	Type      ErrorType `json:"type"`
	Message   string    `json:"message"`
	Code      int       `json:"code"`
	Operation string    `json:"operation,omitempty"`
	Err       error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     err,
	}
}

func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    http.StatusNotFound,
		Err:     err,
	}
}

func NewInternal(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Code:    http.StatusUnauthorized,
		Err:     err,
	}
}

func NewValidation(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    http.StatusUnprocessableEntity,
		Err:     err,
	}
}

func NewAlreadyExists(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeAlreadyExists,
		Message: message,
		Code:    http.StatusConflict,
		Err:     err,
	}
}

func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Code:    http.StatusForbidden,
		Err:     err,
	}
}