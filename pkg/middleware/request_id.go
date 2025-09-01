package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already set in headers
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a new request ID
			requestID = generateRequestID()
		}

		// Set request ID in context
		c.Set("request_id", requestID)

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}
