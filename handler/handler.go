package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/xarunoba/mlgmr/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Compile-time check to ensure LambdaFunction implements HandlerFunc
var _ HandlerFunc = LambdaFunction

// HandlerFunc defines the function signature for the Lambda handler.
// Change this if you want to use a different signature for your LambdaFunction.
// Refer to https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start for more details.
type HandlerFunc func(ctx context.Context, input Input) (*Output, error)

// Input represents the input structure for the Lambda function. (The Event)
type Input struct {
	// Define your input fields here
	Name string `json:"name"`
}

// Output represents the output structure for the Lambda function.
type Output struct {
	// Define your output fields here
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type NameDocument struct {
	Name      string `bson:"name"`
	CreatedAt int64  `bson:"createdAt"`
}

// LambdaFunction is the main handler function for the AWS Lambda.
func LambdaFunction(ctx context.Context, input Input) (*Output, error) {
	mongoClient, err := db.GetMongoClient()
	if err != nil {
		return nil, err
	}
	mongoCollection := mongoClient.Database("mlgmr").Collection("name")

	redisClient, err := db.GetRedisClient()
	if err != nil {
		return nil, err
	}

	var doc NameDocument
	if check := mongoCollection.FindOne(ctx, bson.M{
		"name": input.Name,
	}); check.Err() == mongo.ErrNoDocuments {
		_, err := mongoCollection.InsertOne(ctx, NameDocument{
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
