package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xarunoba/mlgmr/shared"
	"github.com/xarunoba/mlgmr/shared/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Compile-time check to ensure LambdaFunction implements HandlerFunc
var _ shared.HandlerFunc[Input, *Output] = LambdaFunction

// Input represents the input structure for the Lambda function. (The Event)
type Input struct {
	Name string `json:"name"`
}

// Output represents the output structure for the Lambda function.
type Output struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// nameDocument represents a document in the MongoDB "name" collection.
type nameDocument struct {
	Name      string `bson:"name"`
	CreatedAt int64  `bson:"createdAt"`
}

// LambdaFunction is the main handler function for the AWS Lambda.
func LambdaFunction(ctx context.Context, input Input) (*Output, error) {
	mongoClient, err := db.GetMongoClient()
	if err != nil {
		return nil, err
	}
	// Use (or create) the "mlgmr" database and "name" collection
	mongoCollection := mongoClient.Database("mlgmr").Collection("name")

	redisClient, err := db.GetRedisClient()
	if err != nil {
		return nil, err
	}

	var doc nameDocument
	if check := mongoCollection.FindOne(ctx, bson.M{
		"name": input.Name,
	}); check.Err() == mongo.ErrNoDocuments {
		_, err := mongoCollection.InsertOne(ctx, nameDocument{
			Name:      input.Name,
			CreatedAt: time.Now().UnixMilli(),
		})
		if err != nil {
			return nil, err
		}
	} else if check.Err() != nil {
		return nil, check.Err()
	} else {
		if err := check.Decode(&doc); err != nil {
			return nil, err
		}
	}
	createdAt := time.UnixMilli(doc.CreatedAt).Format("January 2, 2006 at 3:04 PM MST")

	// Increment the counter in Redis for the given name
	counterKey := fmt.Sprintf("counter:%s", input.Name)
	counter, err := redisClient.Incr(ctx, counterKey).Result()
	if err != nil {
		return nil, err
	}

	return &Output{
		Success: true,
		Message: fmt.Sprintf("Hello, %s! You have been greeted %d times since %s.", input.Name, counter, createdAt),
	}, nil
}
