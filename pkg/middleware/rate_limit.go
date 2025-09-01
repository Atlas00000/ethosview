package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"ethosview-backend/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter handles rate limiting using Redis
type RateLimiter struct {
	redis *redis.Client
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redis *redis.Client) *RateLimiter {
	return &RateLimiter{redis: redis}
}

// RateLimitMiddleware creates rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = "unknown"
		}

		// Create rate limiter key
		key := "rate_limit:" + clientIP

		// Check current count
		ctx := context.Background()
		count, err := rl.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		// Check if limit exceeded
		if count >= requestsPerMinute {
			errors.HandleRateLimitError(c, requestsPerMinute)
			c.Abort()
			return
		}

		// Increment counter
		pipe := rl.redis.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit update failed"})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(requestsPerMinute-count-1))
		c.Header("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))

		c.Next()
	}
}

// UserRateLimitMiddleware creates rate limiting middleware for authenticated users
func (rl *RateLimiter) UserRateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Create rate limiter key for user
		key := "user_rate_limit:" + strconv.Itoa(userID.(int))

		// Check current count
		ctx := context.Background()
		count, err := rl.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		// Check if limit exceeded
		if count >= requestsPerMinute {
			errors.HandleRateLimitError(c, requestsPerMinute)
			c.Abort()
			return
		}

		// Increment counter
		pipe := rl.redis.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit update failed"})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(requestsPerMinute-count-1))
		c.Header("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))

		c.Next()
	}
}
