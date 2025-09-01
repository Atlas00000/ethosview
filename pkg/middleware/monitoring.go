package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Metrics stores basic request metrics
type Metrics struct {
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	AverageResponseTime time.Duration
}

var globalMetrics = &Metrics{}

// GetMetrics returns current metrics
func GetMetrics() *Metrics {
	return globalMetrics
}

// MonitoringMiddleware creates monitoring middleware
func MonitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		duration := time.Since(start)

		// Update metrics
		globalMetrics.TotalRequests++
		if c.Writer.Status() < 400 {
			globalMetrics.SuccessfulRequests++
		} else {
			globalMetrics.FailedRequests++
		}

		// Update average response time
		if globalMetrics.TotalRequests > 0 {
			totalTime := globalMetrics.AverageResponseTime * time.Duration(globalMetrics.TotalRequests-1)
			globalMetrics.AverageResponseTime = (totalTime + duration) / time.Duration(globalMetrics.TotalRequests)
		}

		// Log request details
		logrus.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("Request processed")
	}
}

// HealthCheckMiddleware adds health check endpoint
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"metrics": gin.H{
				"total_requests":        globalMetrics.TotalRequests,
				"successful_requests":   globalMetrics.SuccessfulRequests,
				"failed_requests":       globalMetrics.FailedRequests,
				"average_response_time": globalMetrics.AverageResponseTime.String(),
			},
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
}
