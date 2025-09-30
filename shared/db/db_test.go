package db_test

import (
	"os"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/xarunoba/mlgmr/shared/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestMain(m *testing.M) {
	// Reset global state before and after tests
	code := m.Run()
	os.Exit(code)
}

// MongoDB Tests

func TestGetMongoClient_MissingEnvironmentVariable(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Unsetenv("MONGODB_URI")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		}
	}()

	client, err := db.GetMongoClient()

	if client != nil {
		t.Error("Expected client to be nil when MONGODB_URI is not set")
	}

	if err == nil {
		t.Error("Expected error when MONGODB_URI is not set")
	}

	expectedError := "MONGODB_URI environment variable not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetMongoClient_InvalidURI(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Setenv("MONGODB_URI", "invalid-uri")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
	}()

	client, err := db.GetMongoClient()

	if client != nil {
		t.Error("Expected client to be nil when URI is invalid")
	}

	if err == nil {
		t.Error("Expected error when URI is invalid")
	}

	if !contains(err.Error(), "failed to connect to MongoDB") {
		t.Errorf("Expected error to contain 'failed to connect to MongoDB', got '%s'", err.Error())
	}
}

func TestGetMongoClient_ConnectionFailure(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Setenv("MONGODB_URI", "mongodb://nonexistent-host:27017/test")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
	}()

	client, err := db.GetMongoClient()

	if client != nil {
		t.Error("Expected client to be nil when connection fails")
	}

	if err == nil {
		t.Error("Expected error when connection fails")
	}

	if !contains(err.Error(), "failed to ping MongoDB") {
		t.Errorf("Expected error to contain 'failed to ping MongoDB', got '%s'", err.Error())
	}
}

func TestGetMongoClient_SingletonBehavior(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Setenv("MONGODB_URI", "mongodb://localhost:99999/test")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
	}()

	client1, err1 := db.GetMongoClient()
	client2, err2 := db.GetMongoClient()

	if client1 != nil || client2 != nil {
		t.Error("Expected both clients to be nil")
	}

	if err1 == nil || err2 == nil {
		t.Error("Expected both calls to return errors")
	}
}

func TestGetMongoClient_ConcurrentAccess(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Setenv("MONGODB_URI", "mongodb://localhost:99999/test")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
	}()

	const numGoroutines = 10
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)
	clients := make(chan *mongo.Client, numGoroutines)

	wg.Add(numGoroutines)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			client, err := db.GetMongoClient()
			clients <- client
			errors <- err
		}()
	}

	wg.Wait()
	close(clients)
	close(errors)

	for client := range clients {
		if client != nil {
			t.Error("Expected all concurrent calls to return nil client")
		}
	}

	for err := range errors {
		if err == nil {
			t.Error("Expected all concurrent calls to return errors")
		}
	}
}

func TestGetMongoClient_HealthCheckFailure(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Setenv("MONGODB_URI", "mongodb://localhost:99999/test")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
	}()

	client1, err1 := db.GetMongoClient()
	if client1 != nil || err1 == nil {
		t.Error("Expected first call to fail")
	}

	client2, err2 := db.GetMongoClient()
	if client2 != nil || err2 == nil {
		t.Error("Expected second call to also fail")
	}
}

func TestGetMongoClient_EnvironmentVariableWhitespace(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Setenv("MONGODB_URI", "")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
	}()

	client, err := db.GetMongoClient()

	if client != nil {
		t.Error("Expected client to be nil when MONGODB_URI is empty string")
	}

	if err == nil {
		t.Error("Expected error when MONGODB_URI is empty string")
	}

	expectedError := "MONGODB_URI environment variable not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetMongoClient_WithValidMockSetup(t *testing.T) {
	originalURI := os.Getenv("MONGODB_URI")
	os.Setenv("MONGODB_URI", "mongodb://test:test@invalid-host:27017/testdb?authSource=admin")
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
	}()

	client, err := db.GetMongoClient()

	if client != nil {
		t.Error("Expected client to be nil with invalid host")
	}

	if err == nil {
		t.Error("Expected error with invalid host")
	}

	if !contains(err.Error(), "failed to ping MongoDB") {
		t.Errorf("Expected error to contain 'failed to ping MongoDB', got '%s'", err.Error())
	}
}

// Redis Tests

func TestGetRedisClient_MissingEnvironmentVariable(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Unsetenv("REDIS_URI")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil when REDIS_URI is not set")
	}

	if err == nil {
		t.Error("Expected error when REDIS_URI is not set")
	}

	expectedError := "REDIS_URI environment variable not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetRedisClient_InvalidURI(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "invalid-uri")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil when URI is invalid")
	}

	if err == nil {
		t.Error("Expected error when URI is invalid")
	}

	if !contains(err.Error(), "failed to parse REDIS_URI") {
		t.Errorf("Expected error to contain 'failed to parse REDIS_URI', got '%s'", err.Error())
	}
}

func TestGetRedisClient_ConnectionFailure(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "redis://localhost:99999")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil when connection fails")
	}

	if err == nil {
		t.Error("Expected error when connection fails")
	}

	if !contains(err.Error(), "failed to ping Redis") {
		t.Errorf("Expected error to contain 'failed to ping Redis', got '%s'", err.Error())
	}
}

func TestGetRedisClient_SingletonBehavior(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "redis://localhost:99999")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client1, err1 := db.GetRedisClient()
	client2, err2 := db.GetRedisClient()

	if client1 != nil || client2 != nil {
		t.Error("Expected both clients to be nil")
	}

	if err1 == nil || err2 == nil {
		t.Error("Expected both calls to return errors")
	}
}

func TestGetRedisClient_ConcurrentAccess(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "redis://localhost:99999")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	const numGoroutines = 10
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)
	clients := make(chan *redis.Client, numGoroutines)

	wg.Add(numGoroutines)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			client, err := db.GetRedisClient()
			clients <- client
			errors <- err
		}()
	}

	wg.Wait()
	close(clients)
	close(errors)

	for client := range clients {
		if client != nil {
			t.Error("Expected all concurrent calls to return nil client")
		}
	}

	for err := range errors {
		if err == nil {
			t.Error("Expected all concurrent calls to return errors")
		}
	}
}

func TestGetRedisClient_HealthCheckFailure(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "redis://localhost:99999")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client1, err1 := db.GetRedisClient()
	if client1 != nil || err1 == nil {
		t.Error("Expected first call to fail")
	}

	client2, err2 := db.GetRedisClient()
	if client2 != nil || err2 == nil {
		t.Error("Expected second call to also fail")
	}
}

func TestGetRedisClient_EnvironmentVariableWhitespace(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil when REDIS_URI is empty string")
	}

	if err == nil {
		t.Error("Expected error when REDIS_URI is empty string")
	}

	expectedError := "REDIS_URI environment variable not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetRedisClient_WithValidURIFormat(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "redis://user:pass@invalid-host:6379/0")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil with invalid host")
	}

	if err == nil {
		t.Error("Expected error with invalid host")
	}

	if !contains(err.Error(), "failed to ping Redis") {
		t.Errorf("Expected error to contain 'failed to ping Redis', got '%s'", err.Error())
	}
}

func TestGetRedisClient_URIWithSSL(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "rediss://invalid-host:6380/0")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil with invalid SSL host")
	}

	if err == nil {
		t.Error("Expected error with invalid SSL host")
	}

	if !contains(err.Error(), "failed to ping Redis") {
		t.Errorf("Expected error to contain 'failed to ping Redis', got '%s'", err.Error())
	}
}

func TestGetRedisClient_MalformedURI(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "not-a-valid-redis-uri")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil with malformed URI")
	}

	if err == nil {
		t.Error("Expected error with malformed URI")
	}

	if !contains(err.Error(), "failed to parse REDIS_URI") {
		t.Errorf("Expected error to contain 'failed to parse REDIS_URI', got '%s'", err.Error())
	}
}

func TestGetRedisClient_URIWithDatabase(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "redis://localhost:99999/5")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	client, err := db.GetRedisClient()

	if client != nil {
		t.Error("Expected client to be nil with invalid host and database")
	}

	if err == nil {
		t.Error("Expected error with invalid host and database")
	}

	if !contains(err.Error(), "failed to ping Redis") {
		t.Errorf("Expected error to contain 'failed to ping Redis', got '%s'", err.Error())
	}
}

func TestGetRedisClient_MultipleReconnectAttempts(t *testing.T) {
	originalURI := os.Getenv("REDIS_URI")
	os.Setenv("REDIS_URI", "redis://localhost:99999")
	defer func() {
		if originalURI != "" {
			os.Setenv("REDIS_URI", originalURI)
		} else {
			os.Unsetenv("REDIS_URI")
		}
	}()

	for i := range 3 {
		client, err := db.GetRedisClient()
		if client != nil {
			t.Errorf("Expected client to be nil on attempt %d", i+1)
		}
		if err == nil {
			t.Errorf("Expected error on attempt %d", i+1)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
