package middleware

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware creates caching middleware
func CacheMiddleware(redis *redis.Client, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip caching for non-GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Create cache key from request
		cacheKey := generateCacheKey(c.Request.URL.String())

		// Try to get from cache
		ctx := context.Background()
		cachedData, err := redis.Get(ctx, cacheKey).Result()
		if err == nil {
			// Return cached response
			var response map[string]interface{}
			if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
				c.JSON(http.StatusOK, response)
				c.Abort()
				return
			}
		}

		// Create a custom response writer to capture the response
		responseWriter := &responseWriter{ResponseWriter: c.Writer, body: []byte{}}
		c.Writer = responseWriter

		// Continue with the request
		c.Next()

		// Cache the response if it's successful
		if c.Writer.Status() == http.StatusOK && len(responseWriter.body) > 0 {
			// Set cache with TTL
			redis.Set(ctx, cacheKey, string(responseWriter.body), ttl)
		}
	}
}

// responseWriter captures the response body
type responseWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

// generateCacheKey creates a unique cache key from URL
func generateCacheKey(url string) string {
	hash := md5.Sum([]byte(url))
	return "cache:" + hex.EncodeToString(hash[:])
}

// InvalidateCache invalidates cache for a specific pattern
func InvalidateCache(redis *redis.Client, pattern string) error {
	ctx := context.Background()
	keys, err := redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		_, err = redis.Del(ctx, keys...).Result()
	}
	return err
}
