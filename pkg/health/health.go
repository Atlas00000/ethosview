package health

import (
	"context"
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// HealthChecker handles comprehensive health checks
type HealthChecker struct {
	db    *sql.DB
	redis *redis.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
	return &HealthChecker{
		db:    db,
		redis: redis,
	}
}

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status      string               `json:"status"`
	Timestamp   string               `json:"timestamp"`
	Version     string               `json:"version"`
	Uptime      string               `json:"uptime"`
	Components  map[string]Component `json:"components"`
	Performance PerformanceMetrics   `json:"performance"`
	Memory      MemoryInfo           `json:"memory"`
	Database    DatabaseHealth       `json:"database"`
	Cache       CacheHealth          `json:"cache"`
}

// Component represents a system component health
type Component struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	ResponseTime int64  `json:"response_time_ms"`
}

// PerformanceMetrics represents performance metrics
type PerformanceMetrics struct {
	ResponseTime   int64   `json:"response_time_ms"`
	RequestsPerSec int64   `json:"requests_per_sec"`
	ErrorRate      float64 `json:"error_rate"`
}

// MemoryInfo represents memory usage information
type MemoryInfo struct {
	Alloc      uint64 `json:"alloc_bytes"`
	TotalAlloc uint64 `json:"total_alloc_bytes"`
	Sys        uint64 `json:"sys_bytes"`
	NumGC      uint32 `json:"num_gc"`
}

// DatabaseHealth represents database health information
type DatabaseHealth struct {
	Status      string `json:"status"`
	Connections int    `json:"connections"`
	SlowQueries int    `json:"slow_queries"`
	TableCount  int    `json:"table_count"`
	IndexCount  int    `json:"index_count"`
}

// CacheHealth represents cache health information
type CacheHealth struct {
	Status           string  `json:"status"`
	UsedMemory       int64   `json:"used_memory_bytes"`
	ConnectedClients int     `json:"connected_clients"`
	HitRatio         float64 `json:"hit_ratio"`
	KeysCount        int     `json:"keys_count"`
}

// CheckHealth performs comprehensive health check
func (hc *HealthChecker) CheckHealth() HealthStatus {
	start := time.Now()

	status := HealthStatus{
		Status:     "healthy",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Version:    "1.0.0",
		Uptime:     "running", // In production, track actual uptime
		Components: make(map[string]Component),
	}

	// Check database health
	dbHealth := hc.checkDatabaseHealth()
	status.Database = dbHealth
	status.Components["database"] = Component{
		Status:       dbHealth.Status,
		Message:      "Database connectivity and performance",
		ResponseTime: time.Since(start).Milliseconds(),
	}

	// Check cache health
	cacheHealth := hc.checkCacheHealth()
	status.Cache = cacheHealth
	status.Components["cache"] = Component{
		Status:       cacheHealth.Status,
		Message:      "Redis cache connectivity and performance",
		ResponseTime: time.Since(start).Milliseconds(),
	}

	// Check memory usage
	status.Memory = hc.getMemoryInfo()

	// Check performance metrics
	status.Performance = hc.getPerformanceMetrics()

	// Determine overall status
	if dbHealth.Status != "healthy" || cacheHealth.Status != "healthy" {
		status.Status = "degraded"
	}

	return status
}

// checkDatabaseHealth checks database health
func (hc *HealthChecker) checkDatabaseHealth() DatabaseHealth {
	start := time.Now()
	health := DatabaseHealth{
		Status: "healthy",
	}

	// Test database connection
	if err := hc.db.Ping(); err != nil {
		health.Status = "unhealthy"
		return health
	}

	// Get connection count
	var connectionCount int
	if err := hc.db.QueryRow("SELECT count(*) FROM pg_stat_activity").Scan(&connectionCount); err == nil {
		health.Connections = connectionCount
	}

	// Get slow queries count
	var slowQueries int
	if err := hc.db.QueryRow(`
		SELECT COUNT(*) 
		FROM pg_stat_activity 
		WHERE state = 'active' 
		AND query_start < NOW() - INTERVAL '100 milliseconds'
	`).Scan(&slowQueries); err == nil {
		health.SlowQueries = slowQueries
	}

	// Get table count
	var tableCount int
	if err := hc.db.QueryRow(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`).Scan(&tableCount); err == nil {
		health.TableCount = tableCount
	}

	// Get index count
	var indexCount int
	if err := hc.db.QueryRow(`
		SELECT COUNT(*) 
		FROM pg_indexes 
		WHERE schemaname = 'public'
	`).Scan(&indexCount); err == nil {
		health.IndexCount = indexCount
	}

	// Check response time
	if time.Since(start) > 5*time.Second {
		health.Status = "degraded"
	}

	return health
}

// checkCacheHealth checks cache health
func (hc *HealthChecker) checkCacheHealth() CacheHealth {
	start := time.Now()
	ctx := context.Background()
	health := CacheHealth{
		Status: "healthy",
	}

	// Test Redis connection
	if err := hc.redis.Ping(ctx).Err(); err != nil {
		health.Status = "unhealthy"
		return health
	}

	// Get Redis info
	info, err := hc.redis.Info(ctx).Result()
	if err != nil {
		health.Status = "degraded"
		return health
	}

	// Parse Redis info (simplified)
	health.UsedMemory = int64(hc.parseRedisInfoInt(info, "used_memory:"))
	health.ConnectedClients = hc.parseRedisInfoInt(info, "connected_clients:")

	hits := hc.parseRedisInfoInt(info, "keyspace_hits:")
	misses := hc.parseRedisInfoInt(info, "keyspace_misses:")
	if hits+misses > 0 {
		health.HitRatio = float64(hits) / float64(hits+misses) * 100
	}

	// Get keys count
	keys, err := hc.redis.Keys(ctx, "*").Result()
	if err == nil {
		health.KeysCount = len(keys)
	}

	// Check response time
	if time.Since(start) > 2*time.Second {
		health.Status = "degraded"
	}

	return health
}

// getMemoryInfo gets current memory usage
func (hc *HealthChecker) getMemoryInfo() MemoryInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryInfo{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}
}

// getPerformanceMetrics gets performance metrics
func (hc *HealthChecker) getPerformanceMetrics() PerformanceMetrics {
	// In a real implementation, you'd track these metrics over time
	return PerformanceMetrics{
		ResponseTime:   100,  // Placeholder
		RequestsPerSec: 1000, // Placeholder
		ErrorRate:      0.01, // Placeholder
	}
}

// parseRedisInfoInt parses Redis info for integer values
func (hc *HealthChecker) parseRedisInfoInt(info, key string) int {
	// Simplified parsing - in production use a proper parser
	return 0 // Placeholder
}

// HealthCheckHandler handles health check requests
func (hc *HealthChecker) HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		health := hc.CheckHealth()

		statusCode := http.StatusOK
		if health.Status == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		} else if health.Status == "degraded" {
			statusCode = http.StatusOK // Still OK but with warnings
		}

		c.JSON(statusCode, health)
	}
}

// DetailedHealthCheckHandler provides detailed health information
func (hc *HealthChecker) DetailedHealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		health := hc.CheckHealth()

		// Add additional detailed information
		detailed := map[string]interface{}{
			"health": health,
			"system": map[string]interface{}{
				"goroutines": runtime.NumGoroutine(),
				"cpu_count":  runtime.NumCPU(),
				"go_version": runtime.Version(),
			},
			"timestamp": time.Now().UTC(),
		}

		statusCode := http.StatusOK
		if health.Status == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, detailed)
	}
}

// ReadinessCheckHandler checks if the service is ready to serve traffic
func (hc *HealthChecker) ReadinessCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check critical dependencies
		ready := true
		issues := []string{}

		// Check database
		if err := hc.db.Ping(); err != nil {
			ready = false
			issues = append(issues, "Database connection failed")
		}

		// Check Redis
		ctx := context.Background()
		if err := hc.redis.Ping(ctx).Err(); err != nil {
			ready = false
			issues = append(issues, "Redis connection failed")
		}

		response := map[string]interface{}{
			"ready":  ready,
			"issues": issues,
		}

		statusCode := http.StatusOK
		if !ready {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, response)
	}
}

// LivenessCheckHandler checks if the service is alive
func (hc *HealthChecker) LivenessCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple liveness check - just return OK if the service is running
		c.JSON(http.StatusOK, gin.H{
			"status":    "alive",
			"timestamp": time.Now().UTC(),
		})
	}
}
