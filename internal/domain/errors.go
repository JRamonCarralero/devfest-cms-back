package domain

import (
	"fmt"
	"runtime"
)

type ErrorType string

const (
	TypeNotFound        ErrorType = "NOT_FOUND"
	TypeBadRequest      ErrorType = "BAD_REQUEST"
	TypeInternal        ErrorType = "INTERNAL"
	TypeAlreadyExists   ErrorType = "ALREADY_EXISTS"
	TypeUnauthenticated ErrorType = "UNAUTHENTICATED"
	TypeUnauthorized    ErrorType = "UNAUTHORIZED"
)

type AppError struct {
	Type     ErrorType `json:"type"`
	Message  string    `json:"message"`
	Code     string    `json:"code,omitempty"`
	TraceID  string    `json:"trace_id"`
	Location string    `json:"-"`
	Err      error     `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
}

// Constructor
func NewAppError(errType ErrorType, message string, err error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)

	return &AppError{
		Type:     errType,
		Message:  message,
		Location: location,
		Err:      err,
	}
}
