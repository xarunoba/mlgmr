package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient      *redis.Client
	redisOnce        sync.Once
	redisErr         error
	redisInitialized bool
	redisMu          sync.Mutex
)

// GetRedisClient returns a singleton Redis client instance optimized for AWS Lambda.
// It reads the Redis URI from the REDIS_URI environment variable.
// If the connection fails, it resets the singleton to allow retries on subsequent calls.
// The client persists across Lambda invocations for connection reuse.
func GetRedisClient() (*redis.Client, error) {
	redisOnce.Do(func() {
		uri := os.Getenv("REDIS_URI")
		if uri == "" {
			redisErr = fmt.Errorf("REDIS_URI environment variable not set")
			return
		}

		opt, err := redis.ParseURL(uri)
		if err != nil {
			redisErr = fmt.Errorf("failed to parse REDIS_URI: %w", err)
			return
		}

		redisClient = redis.NewClient(opt)

		_, redisErr = redisClient.Ping(context.Background()).Result()
		if redisErr != nil {
			redisClient = nil
			return
		}

		redisInitialized = true
	})

	if redisErr != nil {
		resetRedisClient()
		return nil, redisErr
	}

	if redisClient != nil && redisInitialized {
		if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
			resetRedisClient()
			return nil, err
		}
	}

	return redisClient, nil
}

// resetRedisClient safely resets the singleton to allow retry on next call
func resetRedisClient() {
	redisMu.Lock()
	defer redisMu.Unlock()

	redisClient = nil
	redisErr = nil
	redisInitialized = false
	redisOnce = sync.Once{}
}
