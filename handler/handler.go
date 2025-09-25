package handler

import "context"

// Compile-time check to ensure LambdaFunction implements HandlerFunc
var _ HandlerFunc = LambdaFunction

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

// HandlerFunc defines the function signature for the Lambda handler.
// Change this if you want to use a different signature for your LambdaFunction.
// Refer to https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start for more details.
type HandlerFunc func(ctx context.Context, input Input) (Output, error)

// LambdaFunction is the main handler function for the AWS Lambda.
func LambdaFunction(ctx context.Context, input Input) (Output, error) {
	return Output{
		Success: true,
		Message: "Hello, " + input.Name + "!",
	}, nil
}
