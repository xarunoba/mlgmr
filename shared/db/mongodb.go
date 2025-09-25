package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
	clientErr      error
	mu             sync.Mutex
)

// GetMongoClient returns a singleton MongoDB client instance.
// It reads the MongoDB URI from the MONGODB_URI environment variable.
// If the connection fails, it resets the singleton to allow retries on subsequent calls.
// The client is safe for concurrent use by multiple goroutines.
func GetMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			clientErr = fmt.Errorf("MONGODB_URI environment variable not set")
			return
		}
		clientInstance, clientErr = mongo.Connect(options.Client().ApplyURI(uri))
		if clientErr != nil {
			return
		}
		clientErr = clientInstance.Ping(context.Background(), nil)
	})

	if clientErr != nil {
		mu.Lock()
		if clientErr != nil {
			clientOnce = sync.Once{}
			clientInstance = nil
		}
		mu.Unlock()
	}

	return clientInstance, clientErr
}
