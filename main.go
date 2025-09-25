package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xarunoba/mlgm/handler"
	"github.com/xarunoba/mlgm/middlewares"
)

func main() {
	// Wrap the LambdaFunction with the Logger middleware
	wrappedHandler := middlewares.Logger(handler.LambdaFunction)

	// Start the Lambda with the wrapped handler
	lambda.Start(wrappedHandler)
}
