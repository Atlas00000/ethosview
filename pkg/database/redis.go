package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitRedis initializes and returns a Redis client
func InitRedis() (*redis.Client, error) {
	// Check if REDIS_URL is provided (for Railway/cloud Redis)
	redisURL := getEnv("REDIS_URL", "")
	if redisURL != "" {
		// Parse the Redis URL
		parsedURL, err := url.Parse(redisURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}

		// Determine if we need SSL/TLS based on the scheme
		var tlsConfig *tls.Config
		if parsedURL.Scheme == "rediss" {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: false, // Use SSL/TLS for rediss://
			}
		}

		// Create Redis client
		client := redis.NewClient(&redis.Options{
			Addr:      parsedURL.Host,
			Password:  parsedURL.User.Username(),
			DB:        0,
			PoolSize:  10,
			TLSConfig: tlsConfig,
		})

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to ping Redis: %w", err)
		}

		return client, nil
	}

	// Fallback to individual environment variables for local Redis
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0, // Default DB
		PoolSize: 10,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}
