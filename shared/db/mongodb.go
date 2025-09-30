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
	client            *mongo.Client
	clientOnce        sync.Once
	clientErr         error
	clientInitialized bool
	mu                sync.Mutex
)

// GetMongoClient returns a singleton MongoDB client instance optimized for AWS Lambda.
// It reads the MongoDB URI from the MONGODB_URI environment variable.
// If the connection fails, it resets the singleton to allow retries on subsequent calls.
// The client persists across Lambda invocations for connection reuse.
func GetMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			clientErr = fmt.Errorf("MONGODB_URI environment variable not set")
			return
		}

		client, clientErr = mongo.Connect(options.Client().ApplyURI(uri))
		if clientErr != nil {
			return
		}

		clientErr = client.Ping(context.Background(), nil)
		if clientErr != nil {
			client = nil
			return
		}

		clientInitialized = true
	})

	if clientErr != nil {
		resetMongoDBClient()
		return nil, clientErr
	}

	if client != nil && clientInitialized {
		if err := client.Ping(context.Background(), nil); err != nil {
			resetMongoDBClient()
			return nil, err
		}
	}

	return client, nil
}

// resetMongoDBClient safely resets the singleton to allow retry on next call
func resetMongoDBClient() {
	mu.Lock()
	defer mu.Unlock()

	client = nil
	clientErr = nil
	clientInitialized = false
	clientOnce = sync.Once{}
}
