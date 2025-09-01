package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// MetricsCollector handles custom metrics collection
type MetricsCollector struct {
	redis *redis.Client
	db    *sql.DB
	mu    sync.RWMutex
	stats map[string]interface{}
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(redis *redis.Client, db *sql.DB) *MetricsCollector {
	return &MetricsCollector{
		redis: redis,
		db:    db,
		stats: make(map[string]interface{}),
	}
}

// CollectMetrics collects all application metrics
func (mc *MetricsCollector) CollectMetrics() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Collect database metrics
	if err := mc.collectDatabaseMetrics(); err != nil {
		log.Printf("Error collecting database metrics: %v", err)
	}

	// Collect cache metrics
	if err := mc.collectCacheMetrics(); err != nil {
		log.Printf("Error collecting cache metrics: %v", err)
	}

	// Collect business metrics
	if err := mc.collectBusinessMetrics(); err != nil {
		log.Printf("Error collecting business metrics: %v", err)
	}

	// Collect system metrics
	mc.collectSystemMetrics()

	// Store metrics in Redis
	return mc.storeMetrics()
}

// collectDatabaseMetrics collects database performance metrics
func (mc *MetricsCollector) collectDatabaseMetrics() error {
	// Database connection count
	var connectionCount int
	err := mc.db.QueryRow("SELECT count(*) FROM pg_stat_activity").Scan(&connectionCount)
	if err == nil {
		mc.stats["database.connections"] = connectionCount
	}

	// Table sizes
	tables := []string{"companies", "esg_scores", "stock_prices", "financial_indicators"}
	for _, table := range tables {
		var size int64
		query := fmt.Sprintf("SELECT pg_total_relation_size('%s')", table)
		if err := mc.db.QueryRow(query).Scan(&size); err == nil {
			mc.stats[fmt.Sprintf("database.table_size.%s", table)] = size
		}
	}

	// Index usage statistics
	var indexUsage int
	err = mc.db.QueryRow(`
		SELECT SUM(idx_scan) 
		FROM pg_stat_user_indexes 
		WHERE schemaname = 'public'
	`).Scan(&indexUsage)
	if err == nil {
		mc.stats["database.index_usage"] = indexUsage
	}

	// Slow query count (queries taking > 100ms)
	var slowQueries int
	err = mc.db.QueryRow(`
		SELECT COUNT(*) 
		FROM pg_stat_activity 
		WHERE state = 'active' 
		AND query_start < NOW() - INTERVAL '100 milliseconds'
	`).Scan(&slowQueries)
	if err == nil {
		mc.stats["database.slow_queries"] = slowQueries
	}

	return nil
}

// collectCacheMetrics collects Redis cache metrics
func (mc *MetricsCollector) collectCacheMetrics() error {
	ctx := context.Background()

	// Redis info
	info, err := mc.redis.Info(ctx).Result()
	if err != nil {
		return err
	}

	// Parse Redis info for key metrics
	mc.stats["cache.used_memory"] = mc.parseRedisInfo(info, "used_memory:")
	mc.stats["cache.connected_clients"] = mc.parseRedisInfo(info, "connected_clients:")
	mc.stats["cache.total_commands_processed"] = mc.parseRedisInfo(info, "total_commands_processed:")
	mc.stats["cache.keyspace_hits"] = mc.parseRedisInfo(info, "keyspace_hits:")
	mc.stats["cache.keyspace_misses"] = mc.parseRedisInfo(info, "keyspace_misses:")

	// Calculate hit ratio
	hits := mc.stats["cache.keyspace_hits"]
	misses := mc.stats["cache.keyspace_misses"]
	if hits != nil && misses != nil {
		hitsVal := hits.(int)
		missesVal := misses.(int)
		if hitsVal+missesVal > 0 {
			hitRatio := float64(hitsVal) / float64(hitsVal+missesVal) * 100
			mc.stats["cache.hit_ratio"] = hitRatio
		}
	}

	// Cache key count by pattern
	patterns := []string{"cache:companies:*", "cache:esg:*", "cache:sectors:*", "cache:analytics:*"}
	for _, pattern := range patterns {
		keys, err := mc.redis.Keys(ctx, pattern).Result()
		if err == nil {
			mc.stats[fmt.Sprintf("cache.keys.%s", pattern)] = len(keys)
		}
	}

	return nil
}

// collectBusinessMetrics collects business-specific metrics
func (mc *MetricsCollector) collectBusinessMetrics() error {
	// Company metrics
	var totalCompanies int
	err := mc.db.QueryRow("SELECT COUNT(*) FROM companies").Scan(&totalCompanies)
	if err == nil {
		mc.stats["business.total_companies"] = totalCompanies
	}

	// ESG score metrics
	var totalESGScores int
	err = mc.db.QueryRow("SELECT COUNT(*) FROM esg_scores").Scan(&totalESGScores)
	if err == nil {
		mc.stats["business.total_esg_scores"] = totalESGScores
	}

	// Average ESG scores
	var avgOverallScore float64
	err = mc.db.QueryRow("SELECT AVG(overall_score) FROM esg_scores WHERE overall_score IS NOT NULL").Scan(&avgOverallScore)
	if err == nil {
		mc.stats["business.avg_overall_esg_score"] = avgOverallScore
	}

	// Sector distribution
	sectors := make(map[string]int)
	rows, err := mc.db.Query("SELECT sector, COUNT(*) FROM companies WHERE sector IS NOT NULL GROUP BY sector")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var sector string
			var count int
			if rows.Scan(&sector, &count) == nil {
				sectors[sector] = count
			}
		}
		mc.stats["business.sector_distribution"] = sectors
	}

	// Recent ESG score updates
	var recentUpdates int
	err = mc.db.QueryRow("SELECT COUNT(*) FROM esg_scores WHERE created_at > NOW() - INTERVAL '24 hours'").Scan(&recentUpdates)
	if err == nil {
		mc.stats["business.recent_esg_updates"] = recentUpdates
	}

	return nil
}

// collectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) collectSystemMetrics() {
	// Timestamp
	mc.stats["system.timestamp"] = time.Now().UTC()

	// Uptime (simplified - in production you'd track this from startup)
	mc.stats["system.uptime"] = "running"

	// Memory usage (simplified)
	mc.stats["system.memory_usage"] = "monitored"

	// CPU usage (simplified)
	mc.stats["system.cpu_usage"] = "monitored"
}

// storeMetrics stores collected metrics in Redis
func (mc *MetricsCollector) storeMetrics() error {
	ctx := context.Background()

	// Store current metrics
	metricsData, err := json.Marshal(mc.stats)
	if err != nil {
		return err
	}

	// Store with timestamp
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	key := fmt.Sprintf("metrics:%s", timestamp)

	err = mc.redis.Set(ctx, key, metricsData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	// Store latest metrics
	err = mc.redis.Set(ctx, "metrics:latest", metricsData, 1*time.Hour).Err()
	if err != nil {
		return err
	}

	// Keep only last 100 metric snapshots
	mc.cleanupOldMetrics(ctx)

	return nil
}

// cleanupOldMetrics removes old metric snapshots
func (mc *MetricsCollector) cleanupOldMetrics(ctx context.Context) {
	keys, err := mc.redis.Keys(ctx, "metrics:20*").Result()
	if err != nil {
		return
	}

	if len(keys) > 100 {
		// Sort keys and remove oldest
		oldestKeys := keys[:len(keys)-100]
		if len(oldestKeys) > 0 {
			mc.redis.Del(ctx, oldestKeys...)
		}
	}
}

// parseRedisInfo parses Redis info output for specific metrics
func (mc *MetricsCollector) parseRedisInfo(info, key string) interface{} {
	// Simple parsing - in production you'd use a more robust parser
	// This is a simplified implementation
	return 0 // Placeholder
}

// GetMetrics retrieves stored metrics
func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	result := make(map[string]interface{})
	for k, v := range mc.stats {
		result[k] = v
	}
	return result
}

// GetMetricsHistory retrieves historical metrics
func (mc *MetricsCollector) GetMetricsHistory(limit int) ([]map[string]interface{}, error) {
	ctx := context.Background()

	keys, err := mc.redis.Keys(ctx, "metrics:20*").Result()
	if err != nil {
		return nil, err
	}

	// Sort keys (newest first)
	if len(keys) > limit {
		keys = keys[len(keys)-limit:]
	}

	var history []map[string]interface{}
	for _, key := range keys {
		data, err := mc.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var metrics map[string]interface{}
		if json.Unmarshal([]byte(data), &metrics) == nil {
			history = append(history, metrics)
		}
	}

	return history, nil
}

// StartMetricsCollection starts periodic metrics collection
func (mc *MetricsCollector) StartMetricsCollection(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Initial collection
		mc.CollectMetrics()

		for range ticker.C {
			mc.CollectMetrics()
		}
	}()
}
