package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoMutex  sync.Mutex
)

const (
	connectionTimeout = 10 * time.Second
	pingTimeout       = 5 * time.Second
)

// GetMongoClient returns a singleton MongoDB client instance optimized for AWS Lambda.
// It reads the MongoDB URI from the MONGODB_URI environment variable.
// The client persists across Lambda invocations for connection reuse.
// Automatically handles health checking and reconnection transparently.
func GetMongoClient() (*mongo.Client, error) {
	mongoMutex.Lock()
	defer mongoMutex.Unlock()

	// If we have a client, check if it's still healthy
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		defer cancel()

		if err := mongoClient.Ping(ctx, nil); err != nil {
			// Client is unhealthy, disconnect and reset
			mongoClient.Disconnect(context.Background())
			mongoClient = nil
		} else {
			// Client is healthy, return it
			return mongoClient, nil
		}
	}

	// Need to create new client
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable not set")
	}

	// Create connection
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	mongoClient = client
	return mongoClient, nil
}
