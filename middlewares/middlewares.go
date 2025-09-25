package middlewares

import "github.com/xarunoba/mlgmr/handler"

// Middleware defines a function that wraps a handler.HandlerFunc with additional functionality.
// It takes a handler.HandlerFunc as input and returns a new handler.HandlerFunc.
type Middleware func(next handler.HandlerFunc) handler.HandlerFunc
