package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisMutex  sync.Mutex
)

const (
	redisConnectionTimeout = 10 * time.Second
	redisPingTimeout       = 5 * time.Second
)

// GetRedisClient returns a singleton Redis client instance optimized for AWS Lambda.
// It reads the Redis URI from the REDIS_URI environment variable.
// The client persists across Lambda invocations for connection reuse.
// Automatically handles health checking and reconnection transparently.
func GetRedisClient() (*redis.Client, error) {
	redisMutex.Lock()
	defer redisMutex.Unlock()

	// If we have a client, check if it's still healthy
	if redisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), redisPingTimeout)
		defer cancel()

		if _, err := redisClient.Ping(ctx).Result(); err != nil {
			// Client is unhealthy, close and reset
			redisClient.Close()
			redisClient = nil
		} else {
			// Client is healthy, return it
			return redisClient, nil
		}
	}

	// Need to create new client
	uri := os.Getenv("REDIS_URI")
	if uri == "" {
		return nil, fmt.Errorf("REDIS_URI environment variable not set")
	}

	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REDIS_URI: %w", err)
	}

	// Set connection timeouts
	opt.DialTimeout = redisConnectionTimeout
	opt.ReadTimeout = redisPingTimeout
	opt.WriteTimeout = redisPingTimeout

	client := redis.NewClient(opt)

	// Verify connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), redisPingTimeout)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	redisClient = client
	return redisClient, nil
}
