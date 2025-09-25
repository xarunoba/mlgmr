# mlgmr - micro lambda, go, mongodb, redis

This is a GitHub template for creating micro Lambda functions using Go, MongoDB, and/or Redis. It provides a clean monorepository that's easy to customize and extend, with built-in support for local testing using AWS SAM.

**🎯 Use this template**: Click "Use this template" button above to create a new repository from this template.

## Project Structure

```
.
├── functions/                # Lambda functions directory
│   └── greeter/              # Example function
│       ├── main.go          # Function entry point
│       ├── handler.go       # Function logic
│       └── events/
│           └── event.json   # Sample test event
├── shared/                  # Shared code across functions
│   ├── types.go             # Common types and structs
│   ├── db/
│   │   ├── mongodb.go       # MongoDB client
│   │   └── redis.go         # Redis client
│   └── middleware/
│       └── logger.go         # Structured logging middleware (slog)
├── template.yaml             # SAM template for deployment
├── samconfig.template.toml   # SAM configuration template (rename to samconfig.toml)
├── Makefile                  # Build commands
├── go.mod
├── go.sum
├── .gitignore
├── README.md
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

# Rename samconfig template file
mv samconfig.template.toml samconfig.toml

# Tidy up Go modules
go mod tidy
```
  - Rename the `module` name in `go.mod` and all the import paths in the project files.

2. **Configure Environment Variables**:
  - Edit `samconfig.toml` to set your MongoDB and Redis connection strings, and logging level.

3. **Local Testing**:
```bash
# Use a Docker container for building the function (recommended for runtime compatibility)
sam build --use-container

# Invoke GreeterFunction directly with test event
sam local invoke GreeterFunction -e ./functions/greeter/events/event.json
```

4. **Deploy to AWS**:
```bash
# First-time deployment with guided setup
sam deploy --guided

# Subsequent deployments
sam deploy
```

## Usage Examples

- Refer to the [GreeterFunction](./functions/greeter/main.go) for a simple example of handling an event, connecting to MongoDB and Redis, and using middleware for structured logging.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
