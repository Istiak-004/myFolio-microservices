package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Standard error codes
const (
	ErrCodeValidation     = "validation_error"
	ErrCodeUnauthorized   = "unauthorized"
	ErrCodeForbidden      = "forbidden"
	ErrCodeNotFound       = "not_found"
	ErrCodeInternal       = "internal_error"
	ErrCodeBadRequest     = "bad_request"
	ErrCodeConflict       = "conflict"
	ErrCodeRateLimited    = "rate_limited"
	ErrCodeNotImplemented = "not_implemented"
)

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")
	ErrInternal       = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrConflict       = errors.New("conflict")
	ErrRateLimited    = errors.New("rate limited")
	ErrNotImplemented = errors.New("not implemented")
)

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

// NewError creates a new error with code
func NewError(code string, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// ErrorToStatusCode maps errors to HTTP status codes
func ErrorToStatusCode(err error) int {
	switch err {
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	case ErrRateLimited:
		return http.StatusTooManyRequests
	case ErrNotImplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}

// HandleError is a centralized error handler middleware
func HandleError(c *gin.Context) {
	c.Next()

	// Check if there are any errors to handle
	errs := c.Errors
	if len(errs) == 0 {
		return
	}

	// Get the last error
	err := errs.Last().Err

	// Handle different error types
	switch e := err.(type) {
	case *ErrorResponse:
		SendError(c, ErrorToStatusCode(err), e.Code, e.Message)
	case validator.ValidationErrors:
		handleValidationError(c, e)
	default:
		SendError(c, ErrorToStatusCode(err), ErrCodeInternal, "Internal Server Error")
	}
}

func handleValidationError(c *gin.Context, errs validator.ValidationErrors) {
	details := make(map[string]string)
	for _, err := range errs {
		details[err.Field()] = err.Tag()
	}

	SendError(c, http.StatusBadRequest, ErrCodeValidation, "Validation failed")
}
