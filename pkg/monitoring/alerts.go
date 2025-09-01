package monitoring

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// AlertManager handles performance monitoring and alerting
type AlertManager struct {
	db          *sql.DB
	redis       *redis.Client
	mu          sync.RWMutex
	alerts      []Alert
	thresholds  Thresholds
	isMonitoring bool
}

// Alert represents a system alert
type Alert struct {
	ID          string    `json:"id"`
	Type        AlertType `json:"type"`
	Severity    Severity  `json:"severity"`
	Message     string    `json:"message"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Timestamp   time.Time `json:"timestamp"`
	Resolved    bool      `json:"resolved"`
	ResolvedAt  time.Time `json:"resolved_at,omitempty"`
}

// AlertType defines different types of alerts
type AlertType string

const (
	DatabaseResponseTime AlertType = "database_response_time"
	DatabaseConnections  AlertType = "database_connections"
	CacheHitRate        AlertType = "cache_hit_rate"
	MemoryUsage         AlertType = "memory_usage"
	ErrorRate           AlertType = "error_rate"
	RequestRate         AlertType = "request_rate"
	DiskSpace           AlertType = "disk_space"
	QueryPerformance    AlertType = "query_performance"
)

// Severity defines alert severity levels
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Thresholds defines monitoring thresholds
type Thresholds struct {
	DatabaseResponseTimeMs   float64 `json:"database_response_time_ms"`
	DatabaseMaxConnections   float64 `json:"database_max_connections"`
	CacheMinHitRate         float64 `json:"cache_min_hit_rate"`
	MaxMemoryUsagePercent   float64 `json:"max_memory_usage_percent"`
	MaxErrorRatePercent     float64 `json:"max_error_rate_percent"`
	MaxRequestsPerSecond    float64 `json:"max_requests_per_second"`
	MaxQueryTimeMs          float64 `json:"max_query_time_ms"`
	MinDiskSpacePercent     float64 `json:"min_disk_space_percent"`
}

// MonitoringData represents current system metrics
type MonitoringData struct {
	DatabaseMetrics DatabaseMetrics `json:"database"`
	CacheMetrics    CacheMetrics    `json:"cache"`
	SystemMetrics   SystemMetrics   `json:"system"`
	AppMetrics      AppMetrics      `json:"app"`
	Timestamp       time.Time       `json:"timestamp"`
}

type DatabaseMetrics struct {
	ResponseTimeMs    float64 `json:"response_time_ms"`
	ActiveConnections int     `json:"active_connections"`
	SlowQueries       int     `json:"slow_queries"`
	LocksCount        int     `json:"locks_count"`
}

type CacheMetrics struct {
	HitRate         float64 `json:"hit_rate"`
	UsedMemoryMB    float64 `json:"used_memory_mb"`
	ConnectedClients int     `json:"connected_clients"`
	KeysCount       int     `json:"keys_count"`
}

type SystemMetrics struct {
	MemoryUsagePercent float64 `json:"memory_usage_percent"`
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	DiskUsagePercent   float64 `json:"disk_usage_percent"`
	Goroutines         int     `json:"goroutines"`
}

type AppMetrics struct {
	RequestsPerSecond float64 `json:"requests_per_second"`
	ErrorRatePercent  float64 `json:"error_rate_percent"`
	AvgResponseTimeMs float64 `json:"avg_response_time_ms"`
	ActiveUsers       int     `json:"active_users"`
}

// NewAlertManager creates a new alert manager
func NewAlertManager(db *sql.DB, redis *redis.Client) *AlertManager {
	return &AlertManager{
		db:    db,
		redis: redis,
		alerts: make([]Alert, 0),
		thresholds: Thresholds{
			DatabaseResponseTimeMs:   500,   // 500ms
			DatabaseMaxConnections:   80,    // 80 connections
			CacheMinHitRate:         80.0,  // 80%
			MaxMemoryUsagePercent:   85.0,  // 85%
			MaxErrorRatePercent:     5.0,   // 5%
			MaxRequestsPerSecond:    1000,  // 1000 req/s
			MaxQueryTimeMs:          200,   // 200ms
			MinDiskSpacePercent:     15.0,  // 15% free space
		},
	}
}

// StartMonitoring starts the monitoring and alerting system
func (am *AlertManager) StartMonitoring(interval time.Duration) {
	am.mu.Lock()
	am.isMonitoring = true
	am.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Println("🔍 Starting performance monitoring and alerting...")

		for range ticker.C {
			if !am.isMonitoring {
				break
			}

			if err := am.checkMetrics(); err != nil {
				log.Printf("Error checking metrics: %v", err)
			}

			am.cleanupResolvedAlerts()
		}
	}()
}

// StopMonitoring stops the monitoring system
func (am *AlertManager) StopMonitoring() {
	am.mu.Lock()
	am.isMonitoring = false
	am.mu.Unlock()
}

// checkMetrics collects and analyzes metrics for alerting
func (am *AlertManager) checkMetrics() error {
	data, err := am.collectMetrics()
	if err != nil {
		return err
	}

	// Check each metric against thresholds
	am.checkDatabaseMetrics(data.DatabaseMetrics)
	am.checkCacheMetrics(data.CacheMetrics)
	am.checkSystemMetrics(data.SystemMetrics)
	am.checkAppMetrics(data.AppMetrics)

	// Store metrics for historical analysis
	am.storeMetrics(data)

	return nil
}

// collectMetrics gathers current system metrics
func (am *AlertManager) collectMetrics() (*MonitoringData, error) {
	data := &MonitoringData{
		Timestamp: time.Now().UTC(),
	}

	// Collect database metrics
	start := time.Now()
	err := am.db.Ping()
	data.DatabaseMetrics.ResponseTimeMs = float64(time.Since(start).Nanoseconds()) / 1e6

	if err == nil {
		am.collectDatabaseMetrics(&data.DatabaseMetrics)
	}

	// Collect cache metrics
	if am.redis != nil {
		am.collectCacheMetrics(&data.CacheMetrics)
	}

	// Collect system metrics
	am.collectSystemMetrics(&data.SystemMetrics)

	// Collect app metrics (simplified)
	am.collectAppMetrics(&data.AppMetrics)

	return data, nil
}

func (am *AlertManager) collectDatabaseMetrics(metrics *DatabaseMetrics) {
	// Active connections
	var connections int
	am.db.QueryRow("SELECT count(*) FROM pg_stat_activity WHERE state = 'active'").Scan(&connections)
	metrics.ActiveConnections = connections

	// Slow queries
	var slowQueries int
	am.db.QueryRow(`
		SELECT COUNT(*) 
		FROM pg_stat_activity 
		WHERE state = 'active' 
		AND query_start < NOW() - INTERVAL '1 second'
	`).Scan(&slowQueries)
	metrics.SlowQueries = slowQueries

	// Locks count
	var locks int
	am.db.QueryRow("SELECT COUNT(*) FROM pg_locks").Scan(&locks)
	metrics.LocksCount = locks
}

func (am *AlertManager) collectCacheMetrics(metrics *CacheMetrics) {
	ctx := context.Background()
	
	// Cache hit rate
	info, err := am.redis.Info(ctx, "stats").Result()
	if err == nil {
		metrics.HitRate = am.parseCacheHitRate(info)
	}

	// Memory usage
	memInfo, err := am.redis.Info(ctx, "memory").Result()
	if err == nil {
		metrics.UsedMemoryMB = am.parseMemoryUsage(memInfo)
	}

	// Connected clients
	clients, err := am.redis.Info(ctx, "clients").Result()
	if err == nil {
		metrics.ConnectedClients = am.parseConnectedClients(clients)
	}

	// Keys count
	keys, err := am.redis.Keys(ctx, "*").Result()
	if err == nil {
		metrics.KeysCount = len(keys)
	}
}

func (am *AlertManager) collectSystemMetrics(metrics *SystemMetrics) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Memory usage (simplified calculation)
	metrics.MemoryUsagePercent = float64(m.Alloc) / float64(m.Sys) * 100
	metrics.Goroutines = runtime.NumGoroutine()

	// CPU and disk usage would require additional libraries in production
	metrics.CPUUsagePercent = 15.0  // Placeholder
	metrics.DiskUsagePercent = 45.0 // Placeholder
}

func (am *AlertManager) collectAppMetrics(metrics *AppMetrics) {
	// In production, these would be collected from actual application metrics
	metrics.RequestsPerSecond = 150.0 // Placeholder
	metrics.ErrorRatePercent = 1.2    // Placeholder
	metrics.AvgResponseTimeMs = 85.0  // Placeholder
	metrics.ActiveUsers = 45          // Placeholder
}

// Check methods for different metric types

func (am *AlertManager) checkDatabaseMetrics(metrics DatabaseMetrics) {
	// Check response time
	if metrics.ResponseTimeMs > am.thresholds.DatabaseResponseTimeMs {
		am.createAlert(DatabaseResponseTime, SeverityCritical,
			fmt.Sprintf("Database response time is %0.2fms (threshold: %0.2fms)",
				metrics.ResponseTimeMs, am.thresholds.DatabaseResponseTimeMs),
			metrics.ResponseTimeMs, am.thresholds.DatabaseResponseTimeMs)
	}

	// Check connections
	if float64(metrics.ActiveConnections) > am.thresholds.DatabaseMaxConnections {
		am.createAlert(DatabaseConnections, SeverityWarning,
			fmt.Sprintf("High database connections: %d (threshold: %0.0f)",
				metrics.ActiveConnections, am.thresholds.DatabaseMaxConnections),
			float64(metrics.ActiveConnections), am.thresholds.DatabaseMaxConnections)
	}

	// Check slow queries
	if metrics.SlowQueries > 5 {
		am.createAlert(QueryPerformance, SeverityWarning,
			fmt.Sprintf("High number of slow queries: %d", metrics.SlowQueries),
			float64(metrics.SlowQueries), 5.0)
	}
}

func (am *AlertManager) checkCacheMetrics(metrics CacheMetrics) {
	// Check hit rate
	if metrics.HitRate < am.thresholds.CacheMinHitRate {
		am.createAlert(CacheHitRate, SeverityWarning,
			fmt.Sprintf("Low cache hit rate: %0.2f%% (threshold: %0.2f%%)",
				metrics.HitRate, am.thresholds.CacheMinHitRate),
			metrics.HitRate, am.thresholds.CacheMinHitRate)
	}
}

func (am *AlertManager) checkSystemMetrics(metrics SystemMetrics) {
	// Check memory usage
	if metrics.MemoryUsagePercent > am.thresholds.MaxMemoryUsagePercent {
		am.createAlert(MemoryUsage, SeverityCritical,
			fmt.Sprintf("High memory usage: %0.2f%% (threshold: %0.2f%%)",
				metrics.MemoryUsagePercent, am.thresholds.MaxMemoryUsagePercent),
			metrics.MemoryUsagePercent, am.thresholds.MaxMemoryUsagePercent)
	}

	// Check goroutines
	if metrics.Goroutines > 1000 {
		am.createAlert(MemoryUsage, SeverityWarning,
			fmt.Sprintf("High number of goroutines: %d", metrics.Goroutines),
			float64(metrics.Goroutines), 1000.0)
	}
}

func (am *AlertManager) checkAppMetrics(metrics AppMetrics) {
	// Check error rate
	if metrics.ErrorRatePercent > am.thresholds.MaxErrorRatePercent {
		am.createAlert(ErrorRate, SeverityCritical,
			fmt.Sprintf("High error rate: %0.2f%% (threshold: %0.2f%%)",
				metrics.ErrorRatePercent, am.thresholds.MaxErrorRatePercent),
			metrics.ErrorRatePercent, am.thresholds.MaxErrorRatePercent)
	}

	// Check request rate
	if metrics.RequestsPerSecond > am.thresholds.MaxRequestsPerSecond {
		am.createAlert(RequestRate, SeverityWarning,
			fmt.Sprintf("High request rate: %0.2f req/s (threshold: %0.2f req/s)",
				metrics.RequestsPerSecond, am.thresholds.MaxRequestsPerSecond),
			metrics.RequestsPerSecond, am.thresholds.MaxRequestsPerSecond)
	}
}

// createAlert creates a new alert if it doesn't already exist
func (am *AlertManager) createAlert(alertType AlertType, severity Severity, message string, value, threshold float64) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Check if similar alert already exists and is not resolved
	for _, alert := range am.alerts {
		if alert.Type == alertType && !alert.Resolved {
			return // Alert already exists
		}
	}

	alert := Alert{
		ID:        fmt.Sprintf("%s_%d", alertType, time.Now().Unix()),
		Type:      alertType,
		Severity:  severity,
		Message:   message,
		Value:     value,
		Threshold: threshold,
		Timestamp: time.Now().UTC(),
		Resolved:  false,
	}

	am.alerts = append(am.alerts, alert)

	// Log alert
	log.Printf("🚨 ALERT [%s]: %s", severity, message)

	// Store alert in Redis for external consumption
	am.storeAlert(alert)
}

// GetActiveAlerts returns all active (unresolved) alerts
func (am *AlertManager) GetActiveAlerts() []Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var activeAlerts []Alert
	for _, alert := range am.alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// GetAllAlerts returns all alerts (resolved and unresolved)
func (am *AlertManager) GetAllAlerts() []Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alertsCopy := make([]Alert, len(am.alerts))
	copy(alertsCopy, am.alerts)
	return alertsCopy
}

// ResolveAlert marks an alert as resolved
func (am *AlertManager) ResolveAlert(alertID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	for i, alert := range am.alerts {
		if alert.ID == alertID && !alert.Resolved {
			am.alerts[i].Resolved = true
			am.alerts[i].ResolvedAt = time.Now().UTC()
			log.Printf("✅ Alert resolved: %s", alert.Message)
			return nil
		}
	}

	return fmt.Errorf("alert not found or already resolved")
}

// Helper methods

func (am *AlertManager) storeMetrics(data *MonitoringData) {
	if am.redis == nil {
		return
	}

	ctx := context.Background()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Store current metrics
	am.redis.Set(ctx, "monitoring:current", jsonData, time.Hour)

	// Store historical metrics (keep last 24 hours)
	timestamp := data.Timestamp.Format("20060102_1504")
	key := fmt.Sprintf("monitoring:history:%s", timestamp)
	am.redis.Set(ctx, key, jsonData, 24*time.Hour)
}

func (am *AlertManager) storeAlert(alert Alert) {
	if am.redis == nil {
		return
	}

	ctx := context.Background()
	jsonData, err := json.Marshal(alert)
	if err != nil {
		return
	}

	// Store in alerts list
	am.redis.LPush(ctx, "alerts:active", jsonData)
	am.redis.LTrim(ctx, "alerts:active", 0, 99) // Keep last 100 alerts
}

func (am *AlertManager) cleanupResolvedAlerts() {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Remove resolved alerts older than 1 hour
	var activeAlerts []Alert
	cutoff := time.Now().Add(-time.Hour)

	for _, alert := range am.alerts {
		if !alert.Resolved || alert.ResolvedAt.After(cutoff) {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	am.alerts = activeAlerts
}

// Parsing helper methods (simplified implementations)

func (am *AlertManager) parseCacheHitRate(info string) float64 {
	// Simplified implementation - in production use proper parsing
	return 85.0 // Placeholder
}

func (am *AlertManager) parseMemoryUsage(info string) float64 {
	// Simplified implementation - in production use proper parsing
	return 125.5 // Placeholder MB
}

func (am *AlertManager) parseConnectedClients(info string) int {
	// Simplified implementation - in production use proper parsing
	return 15 // Placeholder
}
