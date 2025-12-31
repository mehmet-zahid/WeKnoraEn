package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode defines the error code type
type ErrorCode int

// System error codes
const (
	// Common error codes (1000-1999)
	ErrBadRequest         ErrorCode = 1000
	ErrUnauthorized       ErrorCode = 1001
	ErrForbidden          ErrorCode = 1002
	ErrNotFound           ErrorCode = 1003
	ErrMethodNotAllowed   ErrorCode = 1004
	ErrConflict           ErrorCode = 1005
	ErrTooManyRequests    ErrorCode = 1006
	ErrInternalServer     ErrorCode = 1007
	ErrServiceUnavailable ErrorCode = 1008
	ErrTimeout            ErrorCode = 1009
	ErrValidation         ErrorCode = 1010

	// Tenant related error codes (2000-2099)
	ErrTenantNotFound      ErrorCode = 2000
	ErrTenantAlreadyExists ErrorCode = 2001
	ErrTenantInactive      ErrorCode = 2002
	ErrTenantNameRequired  ErrorCode = 2003
	ErrTenantInvalidStatus ErrorCode = 2004

	// Agent related error codes (2100-2199)
	ErrAgentMissingThinkingModel ErrorCode = 2100
	ErrAgentMissingAllowedTools  ErrorCode = 2101
	ErrAgentInvalidMaxIterations ErrorCode = 2102
	ErrAgentInvalidTemperature   ErrorCode = 2103

	// Add more error codes here
)

// AppError defines the application error structure
type AppError struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	Details  any       `json:"details,omitempty"`
	HTTPCode int       `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("error code: %d, error message: %s", e.Code, e.Message)
}

// WithDetails adds error details
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:     ErrBadRequest,
		Message:  message,
		HTTPCode: http.StatusBadRequest,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:     ErrUnauthorized,
		Message:  message,
		HTTPCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:     ErrForbidden,
		Message:  message,
		HTTPCode: http.StatusForbidden,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:     ErrNotFound,
		Message:  message,
		HTTPCode: http.StatusNotFound,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:     ErrConflict,
		Message:  message,
		HTTPCode: http.StatusConflict,
	}
}

// NewInternalServerError creates an internal server error
func NewInternalServerError(message string) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return &AppError{
		Code:     ErrInternalServer,
		Message:  message,
		HTTPCode: http.StatusInternalServerError,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:     ErrValidation,
		Message:  message,
		HTTPCode: http.StatusBadRequest,
	}
}

// Tenant related errors
func NewTenantNotFoundError() *AppError {
	return &AppError{
		Code:     ErrTenantNotFound,
		Message:  "Tenant does not exist",
		HTTPCode: http.StatusNotFound,
	}
}

// NewTenantAlreadyExistsError creates a tenant already exists error
func NewTenantAlreadyExistsError() *AppError {
	return &AppError{
		Code:     ErrTenantAlreadyExists,
		Message:  "Tenant already exists",
		HTTPCode: http.StatusConflict,
	}
}

// NewTenantInactiveError creates a tenant inactive error
func NewTenantInactiveError() *AppError {
	return &AppError{
		Code:     ErrTenantInactive,
		Message:  "Tenant is inactive",
		HTTPCode: http.StatusForbidden,
	}
}

// Agent related errors
func NewAgentMissingThinkingModelError() *AppError {
	return &AppError{
		Code:     ErrAgentMissingThinkingModel,
		Message:  "Please select a thinking model before enabling Agent mode",
		HTTPCode: http.StatusBadRequest,
	}
}

func NewAgentMissingAllowedToolsError() *AppError {
	return &AppError{
		Code:     ErrAgentMissingAllowedTools,
		Message:  "At least one allowed tool must be selected",
		HTTPCode: http.StatusBadRequest,
	}
}

func NewAgentInvalidMaxIterationsError() *AppError {
	return &AppError{
		Code:     ErrAgentInvalidMaxIterations,
		Message:  "Maximum iterations must be between 1 and 20",
		HTTPCode: http.StatusBadRequest,
	}
}

func NewAgentInvalidTemperatureError() *AppError {
	return &AppError{
		Code:     ErrAgentInvalidTemperature,
		Message:  "Temperature parameter must be between 0 and 2",
		HTTPCode: http.StatusBadRequest,
	}
}

// IsAppError checks if the error is an AppError type
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
