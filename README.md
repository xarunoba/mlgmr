# mlgmr - micro lambda, go, mongodb, redis

This is a GitHub template for creating micro Lambda functions using Go, MongoDB, and/or Redis. It provides a clean, modular architecture that's easy to customize and extend, with built-in support for local testing using AWS SAM.

**🎯 Use this template**: Click "Use this template" button above to create a new repository from this template.

## Project Structure

```
.
├── main.go                 # Entry point - Lambda startup
├── go.mod
├── go.sum
├── README.md
├── db/                     # Database clients
│   ├── mongodb.go          # MongoDB client
│   └── redis.go            # Redis client
├── handler/
│   └── handler.go          # Core Lambda function logic
├── middlewares/            # Middlewares for the handler
│   └── logger.go           # Structured logging middleware (slog)
├── events/                 # Test events for local development
│   └── event.json          # Sample test event
├── template.yaml           # SAM template for local testing & deployment
├── Makefile                # Makefile for AWS SAM build command
├── .gitignore
└── LICENSE

```

## Prerequisites
- [Go](https://golang.org/dl/) (version 1.25.1 or later)
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
- [Docker](https://www.docker.com/get-started) (for local testing with SAM)

## Quick Start

1. **Create from template**:
  - Click "Use this template" button above
  - Clone your new repository locally

```bash
git clone https://github.com/yourusername/your-repo-name
cd your-repo-name
go mod tidy
```

2. **Configure Environment Variables**:
  - Edit `template.yaml` to set your MongoDB and Redis connection strings, and logging level.

3. **Local Testing**:
```bash
# Use a Docker container for building the function (recommended for runtime compatibility)
sam build --use-container

# Invoke directly with test event
sam local invoke MLGMRFunction -e events/event.json
```

4. **Deploy to AWS**:
```bash
# First-time deployment with guided setup
sam deploy --guided

# Subsequent deployments
sam deploy
```

## Usage Examples

### Using MongoDB

```go
func LambdaFunction(ctx context.Context, input Input) (Output, error) {
    client, err := db.GetMongoClient()
    if err != nil {
        return Output{}, fmt.Errorf("MongoDB connection failed: %w", err)
    }

    database := client.Database("myapp")
    collection := database.Collection("users")

    // Your MongoDB operations here
    // ...

    return Output{}, nil
}
```

### Using Redis

```go
func LambdaFunction(ctx context.Context, input Input) (Output, error) {
    client, err := db.GetRedisClient()
    if err != nil {
        return Output{}, fmt.Errorf("Redis connection failed: %w", err)
    }

    // Set a value
    err = client.Set(ctx, "key", "value", time.Hour).Err()
    if err != nil {
        return Output{}, err
    }

    // Get a value
    val, err := client.Get(ctx, "key").Result()
    if err != nil {
        return Output{}, err
    }

    return Output{Data: val}, nil
}
```

## Customization Guide

### 1. Modify Handler

Edit `handler/handler.go` according to your needs:

```go
type Input struct {
    UserID string `json:"user_id"`
    Action string `json:"action"`
}

// Adjust the handler function signature if needed
type HandlerFunc func(ctx context.Context, input Input) (error)

func LambdaFunction(ctx context.Context, input Input) (error) {
	// Your logic here
	return nil
}
```

### 2. Add New Middleware

Create new middleware in `middlewares/`:

```go
// middlewares/auth.go
func Auth(next handler.HandlerFunc) handler.HandlerFunc {
    return func(ctx context.Context, input handler.Input) (handler.Output, error) {
        // Authentication logic here
        return next(ctx, input)
    }
}
```

Then apply it in `main.go`:

```go
wrappedHandler := middlewares.Auth(middlewares.Logger(handler.LambdaFunction))
```

### 3. Modify `template.yaml`

Update the `template.yaml` file to adjust your Lambda function's configuration:

```yaml
...
Resources:
  YourFunction:
    Type: AWS::Serverless::Function
    Properties:
      MemorySize: 256
      Timeout: 10
      Environment:
        Variables:
          MONGODB_URI: "your-mongodb-uri"
          REDIS_ADDR: "your-redis-address"
          LOG_LEVEL: "INFO"
...
```

```bash
sam build --use-container
sam local invoke YourFunction -e events/event.json
```

### 4. Remove Unused Components

Remove unused files and dependencies to keep the project clean. Don't forget to perform `go mod tidy` and update `template.yaml` after making changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
