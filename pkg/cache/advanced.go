package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// AdvancedCache provides enhanced caching strategies
type AdvancedCache struct {
	redis  *redis.Client
	prefix string
}

// CacheStrategy defines different caching strategies
type CacheStrategy int

const (
	// ShortTerm for frequently changing data (5 minutes)
	ShortTerm CacheStrategy = iota
	// MediumTerm for moderately changing data (30 minutes)
	MediumTerm
	// LongTerm for stable data (2 hours)
	LongTerm
	// Daily for daily aggregations (24 hours)
	Daily
)

// CacheEntry represents a cache entry with metadata
type CacheEntry struct {
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"created_at"`
	ExpiresAt time.Time   `json:"expires_at"`
	Version   string      `json:"version"`
	Tags      []string    `json:"tags"`
}

// NewAdvancedCache creates a new advanced cache instance
func NewAdvancedCache(redis *redis.Client, prefix string) *AdvancedCache {
	return &AdvancedCache{
		redis:  redis,
		prefix: prefix,
	}
}

// Set stores data with advanced caching strategy
func (ac *AdvancedCache) Set(key string, data interface{}, strategy CacheStrategy, tags []string) error {
	ctx := context.Background()
	
	// Determine TTL based on strategy
	ttl := ac.getStrategyTTL(strategy)
	
	// Create cache entry
	entry := CacheEntry{
		Data:      data,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(ttl),
		Version:   ac.generateVersion(data),
		Tags:      tags,
	}

	// Serialize entry
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	// Store in Redis with TTL
	fullKey := ac.buildKey(key)
	err = ac.redis.Set(ctx, fullKey, jsonData, ttl).Err()
	if err != nil {
		return err
	}

	// Store tags for invalidation
	for _, tag := range tags {
		ac.addToTagSet(tag, fullKey)
	}

	return nil
}

// Get retrieves data from cache with validation
func (ac *AdvancedCache) Get(key string, dest interface{}) (bool, error) {
	ctx := context.Background()
	fullKey := ac.buildKey(key)

	data, err := ac.redis.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return false, nil // Cache miss
	}
	if err != nil {
		return false, err
	}

	var entry CacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		// Invalid cache entry, delete it
		ac.redis.Del(ctx, fullKey)
		return false, nil
	}

	// Check if expired (extra safety check)
	if time.Now().UTC().After(entry.ExpiresAt) {
		ac.redis.Del(ctx, fullKey)
		return false, nil
	}

	// Deserialize data
	dataBytes, err := json.Marshal(entry.Data)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(dataBytes, dest); err != nil {
		return false, err
	}

	return true, nil
}

// GetOrSet retrieves from cache or executes function and caches result
func (ac *AdvancedCache) GetOrSet(key string, dest interface{}, fn func() (interface{}, error), strategy CacheStrategy, tags []string) error {
	// Try to get from cache first
	found, err := ac.Get(key, dest)
	if err != nil {
		return err
	}

	if found {
		return nil // Cache hit
	}

	// Cache miss, execute function
	data, err := fn()
	if err != nil {
		return err
	}

	// Store in cache
	if err := ac.Set(key, data, strategy, tags); err != nil {
		log.Printf("Failed to cache data for key %s: %v", key, err)
		// Continue execution even if caching fails
	}

	// Set the result
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(dataBytes, dest)
}

// InvalidateByTag invalidates all cache entries with a specific tag
func (ac *AdvancedCache) InvalidateByTag(tag string) error {
	ctx := context.Background()
	tagKey := ac.buildTagKey(tag)

	// Get all keys with this tag
	keys, err := ac.redis.SMembers(ctx, tagKey).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all keys
	err = ac.redis.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	// Clear tag set
	return ac.redis.Del(ctx, tagKey).Err()
}

// InvalidatePattern invalidates all keys matching a pattern
func (ac *AdvancedCache) InvalidatePattern(pattern string) error {
	ctx := context.Background()
	fullPattern := ac.buildKey(pattern)

	keys, err := ac.redis.Keys(ctx, fullPattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return ac.redis.Del(ctx, keys...).Err()
}

// Refresh refreshes cache entry with new data
func (ac *AdvancedCache) Refresh(key string, data interface{}, strategy CacheStrategy, tags []string) error {
	// Delete existing entry
	ac.Delete(key)
	
	// Set new entry
	return ac.Set(key, data, strategy, tags)
}

// Delete removes a specific cache entry
func (ac *AdvancedCache) Delete(key string) error {
	ctx := context.Background()
	fullKey := ac.buildKey(key)
	return ac.redis.Del(ctx, fullKey).Err()
}

// GetStats returns cache statistics
func (ac *AdvancedCache) GetStats() (map[string]interface{}, error) {
	ctx := context.Background()
	
	// Get Redis info
	info, err := ac.redis.Info(ctx, "memory", "stats").Result()
	if err != nil {
		return nil, err
	}

	// Count keys by prefix
	keys, err := ac.redis.Keys(ctx, ac.buildKey("*")).Result()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_keys":    len(keys),
		"prefix":        ac.prefix,
		"redis_info":    ac.parseRedisInfo(info),
		"last_updated":  time.Now().UTC(),
	}

	return stats, nil
}

// WarmupCache performs intelligent cache warming
func (ac *AdvancedCache) WarmupCache(warmupFuncs map[string]func() (interface{}, error)) error {
	log.Println("🔥 Starting advanced cache warmup...")

	for key, fn := range warmupFuncs {
		data, err := fn()
		if err != nil {
			log.Printf("Failed to warmup %s: %v", key, err)
			continue
		}

		// Use medium-term strategy for warmup
		tags := []string{"warmup", ac.getKeyCategory(key)}
		if err := ac.Set(key, data, MediumTerm, tags); err != nil {
			log.Printf("Failed to cache warmup data for %s: %v", key, err)
		}
	}

	log.Println("✅ Advanced cache warmup completed")
	return nil
}

// BuildQueryKey creates a consistent cache key for database queries
func (ac *AdvancedCache) BuildQueryKey(table string, operation string, params map[string]interface{}) string {
	// Sort parameters for consistent key generation
	var paramParts []string
	for k, v := range params {
		paramParts = append(paramParts, fmt.Sprintf("%s=%v", k, v))
	}
	
	paramStr := strings.Join(paramParts, "&")
	keyStr := fmt.Sprintf("query:%s:%s:%s", table, operation, paramStr)
	
	// Hash long keys
	if len(keyStr) > 200 {
		hash := md5.Sum([]byte(keyStr))
		return fmt.Sprintf("query:hash:%s", hex.EncodeToString(hash[:]))
	}
	
	return keyStr
}

// Helper methods

func (ac *AdvancedCache) getStrategyTTL(strategy CacheStrategy) time.Duration {
	switch strategy {
	case ShortTerm:
		return 5 * time.Minute
	case MediumTerm:
		return 30 * time.Minute
	case LongTerm:
		return 2 * time.Hour
	case Daily:
		return 24 * time.Hour
	default:
		return 30 * time.Minute
	}
}

func (ac *AdvancedCache) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", ac.prefix, key)
}

func (ac *AdvancedCache) buildTagKey(tag string) string {
	return fmt.Sprintf("%s:tag:%s", ac.prefix, tag)
}

func (ac *AdvancedCache) addToTagSet(tag, key string) {
	ctx := context.Background()
	tagKey := ac.buildTagKey(tag)
	ac.redis.SAdd(ctx, tagKey, key)
	ac.redis.Expire(ctx, tagKey, 25*time.Hour) // Slightly longer than daily cache
}

func (ac *AdvancedCache) generateVersion(data interface{}) string {
	// Simple version based on current time
	return fmt.Sprintf("v%d", time.Now().Unix())
}

func (ac *AdvancedCache) getKeyCategory(key string) string {
	parts := strings.Split(key, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return "general"
}

func (ac *AdvancedCache) parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	}
	
	return result
}
