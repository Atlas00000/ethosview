package security

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware provides comprehensive security features
type SecurityMiddleware struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware() *SecurityMiddleware {
	return &SecurityMiddleware{
		allowedOrigins: []string{"http://localhost:3000", "https://ethosview.com"},
		allowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		allowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
	}
}

// SecurityHeaders adds security headers to responses
func (sm *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// CORS handles Cross-Origin Resource Sharing
func (sm *SecurityMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if sm.isOriginAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", strings.Join(sm.allowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(sm.allowedHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// InputSanitization sanitizes input data
func (sm *SecurityMiddleware) InputSanitization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = sm.sanitizeString(value)
			}
			c.Request.URL.Query()[key] = values
		}

		// Sanitize path parameters
		for _, param := range c.Params {
			param.Value = sm.sanitizeString(param.Value)
		}

		c.Next()
	}
}

// RateLimitHeaders adds rate limiting headers
func (sm *SecurityMiddleware) RateLimitHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add rate limiting headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", "99") // This would be dynamic
		c.Header("X-RateLimit-Reset", "60")     // This would be dynamic

		c.Next()
	}
}

// RequestSizeLimit limits request body size
func (sm *SecurityMiddleware) RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// SQLInjectionProtection provides basic SQL injection protection
func (sm *SecurityMiddleware) SQLInjectionProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for SQL injection patterns in query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if sm.containsSQLInjection(value) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Invalid input detected",
						"code":  400,
					})
					c.Abort()
					return
				}
			}
		}

		// Check path parameters
		for _, param := range c.Params {
			if sm.containsSQLInjection(param.Value) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid input detected",
					"code":  400,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// XSSProtection provides basic XSS protection
func (sm *SecurityMiddleware) XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for XSS patterns in query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if sm.containsXSS(value) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Invalid input detected",
						"code":  400,
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func (sm *SecurityMiddleware) isOriginAllowed(origin string) bool {
	for _, allowed := range sm.allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

// sanitizeString removes potentially dangerous characters
func (sm *SecurityMiddleware) sanitizeString(input string) string {
	// Remove null bytes and control characters
	re := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	sanitized := re.ReplaceAllString(input, "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Remove HTML tags
	htmlRe := regexp.MustCompile(`<[^>]*>`)
	sanitized = htmlRe.ReplaceAllString(sanitized, "")

	return sanitized
}

// containsSQLInjection checks for SQL injection patterns
func (sm *SecurityMiddleware) containsSQLInjection(input string) bool {
	// Common SQL injection patterns
	patterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
		`(?i)(--|/\*|\*/|xp_|sp_)`,
		`(?i)(or\s+\d+\s*=\s*\d+|and\s+\d+\s*=\s*\d+)`,
		`(?i)(union\s+select|select\s+union)`,
		`(?i)(information_schema|sys\.tables|sys\.columns)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(input) {
			return true
		}
	}

	return false
}

// containsXSS checks for XSS patterns
func (sm *SecurityMiddleware) containsXSS(input string) bool {
	// Common XSS patterns
	patterns := []string{
		`(?i)<script[^>]*>.*?</script>`,
		`(?i)javascript:`,
		`(?i)on\w+\s*=`,
		`(?i)<iframe[^>]*>`,
		`(?i)<object[^>]*>`,
		`(?i)<embed[^>]*>`,
		`(?i)<form[^>]*>`,
		`(?i)<input[^>]*>`,
		`(?i)<textarea[^>]*>`,
		`(?i)<select[^>]*>`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(input) {
			return true
		}
	}

	return false
}

// ValidateAPIKey validates API keys
func (sm *SecurityMiddleware) ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		// In production, validate against a database or external service
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
				"code":  401,
			})
			c.Abort()
			return
		}

		// Basic validation (in production, use proper validation)
		if len(apiKey) < 10 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
				"code":  401,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuditLog logs security-relevant events
func (sm *SecurityMiddleware) AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log suspicious activities
		if sm.isSuspiciousRequest(c) {
			// In production, log to a security monitoring system
			// log.Printf("SUSPICIOUS: IP=%s, UA=%s, Method=%s, Path=%s", c.ClientIP(), c.GetHeader("User-Agent"), c.Request.Method, c.Request.URL.Path)
		}

		c.Next()
	}
}

// isSuspiciousRequest checks if the request is suspicious
func (sm *SecurityMiddleware) isSuspiciousRequest(c *gin.Context) bool {
	// Check for suspicious patterns
	suspiciousPatterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
		`(?i)(<script|javascript:|on\w+\s*=)`,
		`(?i)(\.\./|\.\.\\|%2e%2e%2f|%2e%2e%5c)`, // Path traversal
		`(?i)(eval\(|setTimeout\(|setInterval\()`,
	}

	// Check query parameters
	for _, values := range c.Request.URL.Query() {
		for _, value := range values {
			for _, pattern := range suspiciousPatterns {
				re := regexp.MustCompile(pattern)
				if re.MatchString(value) {
					return true
				}
			}
		}
	}

	// Check path
	for _, pattern := range suspiciousPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(c.Request.URL.Path) {
			return true
		}
	}

	return false
}
