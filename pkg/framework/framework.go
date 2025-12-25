package framework

// WARNING: This is a core package. Do NOT modify unless you're changing the framework abstraction.
// For adding features, work in internal/ directory instead.

import "net/http"

// Context represents an HTTP request context
type Context interface {
	Request() *http.Request
	Response() http.ResponseWriter
	Param(key string) string
	QueryParam(key string) string
	Body() ([]byte, error)
	JSON(code int, data interface{}) error
	String(code int, data string) error
	Set(key string, value interface{})
	Get(key string) interface{}
}

// HandlerFunc represents a handler function
type HandlerFunc func(Context) error

// MiddlewareFunc represents a middleware function
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Router represents a router interface
type Router interface {
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	PATCH(path string, handler HandlerFunc)
	Group(prefix string) Router
	Use(middleware ...MiddlewareFunc)
}

// Framework represents a web framework interface
type Framework interface {
	Router() Router
	Start(address string) error
	Shutdown() error
}
