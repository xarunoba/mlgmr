package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisInstance *redis.Client
	redisOnce     sync.Once
	redisErr      error
	redisMu       sync.Mutex
)

// GetRedisClient returns a singleton Redis client instance.
// It reads the Redis URI from the REDIS_URI environment variable.
// If the connection fails, it resets the singleton to allow retries on subsequent calls.
// The client is safe for concurrent use by multiple goroutines.
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

		redisInstance = redis.NewClient(opt)

		_, redisErr = redisInstance.Ping(context.Background()).Result()
	})

	if redisErr != nil {
		redisMu.Lock()
		if redisErr != nil {
			redisOnce = sync.Once{}
			redisInstance = nil
		}
		redisMu.Unlock()
	}

	return redisInstance, redisErr
}
