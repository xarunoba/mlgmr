package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xarunoba/mlgmr/shared/middleware"
)

func main() {
	// Wrap the lambdaFn with the Logger middleware
	wrappedHandler := middleware.Logger(LambdaFunction)

	// Start the Lambda with the wrapped handler
	lambda.Start(wrappedHandler)
}
