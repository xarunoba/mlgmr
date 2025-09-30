package shared

import "context"

// HandlerFunc defines the signature for AWS Lambda handler functions.
// It takes a context and an input of type TIn, and returns an output of type TOut and an error.
type HandlerFunc[TIn, TOut any] func(ctx context.Context, input TIn) (TOut, error)

// MiddlewareFunc defines a function that wraps a HandlerFunc with additional functionality.
// It takes a HandlerFunc as input and returns a new HandlerFunc.
type MiddlewareFunc[TIn, TOut any] func(next HandlerFunc[TIn, TOut]) HandlerFunc[TIn, TOut]
