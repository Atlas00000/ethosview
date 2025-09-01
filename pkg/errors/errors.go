package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Error types for different scenarios
var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrDatabaseError     = errors.New("database error")
	ErrValidationError   = errors.New("validation error")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// AppError represents an application error
type AppError struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Error     string      `json:"error"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// NewAppError creates a new application error
func NewAppError(code int, message, err string, details interface{}) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Error:     err,
		Details:   details,
		Timestamp: getCurrentTimestamp(),
	}
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	appErr := &AppError{
		Code:      statusCode,
		Message:   message,
		Error:     err.Error(),
		Timestamp: getCurrentTimestamp(),
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		appErr.RequestID = requestID.(string)
	}

	c.JSON(statusCode, appErr)
}

// HandleDatabaseError handles database-specific errors
func HandleDatabaseError(c *gin.Context, err error, operation string) {
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("%s not found", operation), ErrNotFound)
		return
	}

	// Check for constraint violations
	if strings.Contains(err.Error(), "duplicate key") {
		ErrorResponse(c, http.StatusConflict, "Resource already exists", ErrValidationError)
		return
	}

	if strings.Contains(err.Error(), "foreign key") {
		ErrorResponse(c, http.StatusBadRequest, "Referenced resource does not exist", ErrValidationError)
		return
	}

	// Default database error
	ErrorResponse(c, http.StatusInternalServerError, "Database operation failed", ErrDatabaseError)
}

// HandleValidationError handles validation errors
func HandleValidationError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, "Validation failed", err)
}

// HandleAuthError handles authentication errors
func HandleAuthError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusUnauthorized, "Authentication failed", err)
}

// HandleRateLimitError handles rate limiting errors
func HandleRateLimitError(c *gin.Context, limit int) {
	appErr := &AppError{
		Code:      http.StatusTooManyRequests,
		Message:   "Rate limit exceeded",
		Error:     ErrRateLimitExceeded.Error(),
		Details:   map[string]interface{}{"limit": limit, "reset": "in 1 minute"},
		Timestamp: getCurrentTimestamp(),
	}

	if requestID, exists := c.Get("request_id"); exists {
		appErr.RequestID = requestID.(string)
	}

	c.JSON(http.StatusTooManyRequests, appErr)
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, data interface{}) {
	response := gin.H{
		"success":   true,
		"data":      data,
		"timestamp": getCurrentTimestamp(),
	}

	if requestID, exists := c.Get("request_id"); exists {
		response["request_id"] = requestID.(string)
	}

	c.JSON(http.StatusOK, response)
}

// getCurrentTimestamp returns the current timestamp in ISO format
func getCurrentTimestamp() string {
	return getCurrentTime().Format("2006-01-02T15:04:05Z07:00")
}

// getCurrentTime is a function that can be overridden for testing
var getCurrentTime = func() time.Time {
	return time.Now()
}
