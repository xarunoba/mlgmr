package middlewares

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/xarunoba/mlgmr/handler"
)

// Compile-time check to ensure Logger implements Middleware
var _ Middleware = Logger

var (
	loggerInstance *slog.Logger
	loggerOnce     sync.Once
	loggerMu       sync.Mutex
)

// GetLogger returns a singleton slog.Logger instance.
// It reads the LOG_LEVEL environment variable to set the log level.
// The logger is safe for concurrent use by multiple goroutines.
func GetLogger() *slog.Logger {
	loggerOnce.Do(func() {
		level := slog.LevelDebug

		if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
			switch logLevel {
			case "DEBUG":
				level = slog.LevelDebug
			case "WARN":
				level = slog.LevelWarn
			case "ERROR":
				level = slog.LevelError
			}
		}

		loggerInstance = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
	})

	return loggerInstance
}

// Logger is a middleware that logs the input and output of the handler.
// It wraps a any and returns a new any
// that logs the input before calling the original handler and logs the output
// after the handler has been called.
func Logger(next handler.HandlerFunc) handler.HandlerFunc {
	logger := GetLogger()

	return func(ctx context.Context, input handler.Input) (*handler.Output, error) {
		// Log the input with structured logging
		logger.DebugContext(ctx, "Lambda invocation started",
			slog.Any("input", input),
		)

		// Call the next handler
		output, err := next(ctx, input)

		// Log the output and error (if any) with structured logging
		if err != nil {
			logger.ErrorContext(ctx, "Lambda invocation failed",
				slog.Any("input", input),
				slog.Any("output", output),
				slog.Any("error", err),
			)
		} else {
			logger.DebugContext(ctx, "Lambda invocation completed",
				slog.Any("input", input),
				slog.Any("output", output),
			)
		}

		return output, err
	}
}
