package middleware

// WARNING: This is a core package. Do NOT modify unless you're changing middleware architecture.
// For adding middlewares, create files in internal/middlewares/onRequest or onResponse instead.

import "github.com/alramdein/kaimon/pkg/framework"

// Phase represents middleware execution phase
type Phase string

const (
	PhaseOnRequest  Phase = "onRequest"
	PhaseOnResponse Phase = "onResponse"
)

// Middleware represents a middleware with metadata
type Middleware struct {
	Name    string
	Phase   Phase
	Handler framework.MiddlewareFunc
}

// OnRequestMiddleware is the interface for onRequest middlewares
type OnRequestMiddleware interface {
	Name() string
	Handle(next framework.HandlerFunc) framework.HandlerFunc
}

// OnResponseMiddleware is the interface for onResponse middlewares
type OnResponseMiddleware interface {
	Name() string
	Handle(next framework.HandlerFunc) framework.HandlerFunc
}
