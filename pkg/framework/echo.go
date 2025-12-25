package framework

// WARNING: This is a core package. Do NOT modify unless you're implementing a new framework.
// For adding features, work in internal/ directory instead.

import (
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// EchoFramework wraps Echo framework
type EchoFramework struct {
	e *echo.Echo
}

// NewEchoFramework creates a new Echo framework instance
func NewEchoFramework() Framework {
	return &EchoFramework{
		e: echo.New(),
	}
}

// Router returns the router
func (ef *EchoFramework) Router() Router {
	return &EchoRouter{
		group: ef.e.Group(""),
	}
}

// Start starts the server
func (ef *EchoFramework) Start(address string) error {
	return ef.e.Start(address)
}

// Shutdown gracefully shuts down the server
func (ef *EchoFramework) Shutdown() error {
	return ef.e.Shutdown(context.Background())
}

// EchoRouter wraps Echo router
type EchoRouter struct {
	group *echo.Group
}

// GET registers a GET route
func (er *EchoRouter) GET(path string, handler HandlerFunc) {
	er.group.GET(path, er.wrapHandler(handler))
}

// POST registers a POST route
func (er *EchoRouter) POST(path string, handler HandlerFunc) {
	er.group.POST(path, er.wrapHandler(handler))
}

// PUT registers a PUT route
func (er *EchoRouter) PUT(path string, handler HandlerFunc) {
	er.group.PUT(path, er.wrapHandler(handler))
}

// DELETE registers a DELETE route
func (er *EchoRouter) DELETE(path string, handler HandlerFunc) {
	er.group.DELETE(path, er.wrapHandler(handler))
}

// PATCH registers a PATCH route
func (er *EchoRouter) PATCH(path string, handler HandlerFunc) {
	er.group.PATCH(path, er.wrapHandler(handler))
}

// Group creates a route group
func (er *EchoRouter) Group(prefix string) Router {
	return &EchoRouter{
		group: er.group.Group(prefix),
	}
}

// Use adds middleware
func (er *EchoRouter) Use(middleware ...MiddlewareFunc) {
	for _, m := range middleware {
		er.group.Use(er.wrapMiddleware(m))
	}
}

// wrapHandler wraps our HandlerFunc to Echo handler
func (er *EchoRouter) wrapHandler(handler HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &EchoContext{c: c}
		return handler(ctx)
	}
}

// wrapMiddleware wraps our MiddlewareFunc to Echo middleware
func (er *EchoRouter) wrapMiddleware(middleware MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &EchoContext{c: c}
			wrappedNext := func(ctx Context) error {
				return next(c)
			}
			handler := middleware(wrappedNext)
			return handler(ctx)
		}
	}
}

// EchoContext wraps Echo context
type EchoContext struct {
	c echo.Context
}

// Request returns the HTTP request
func (ec *EchoContext) Request() *http.Request {
	return ec.c.Request()
}

// Response returns the HTTP response writer
func (ec *EchoContext) Response() http.ResponseWriter {
	return ec.c.Response()
}

// Param returns the URL parameter
func (ec *EchoContext) Param(key string) string {
	return ec.c.Param(key)
}

// QueryParam returns the query parameter
func (ec *EchoContext) QueryParam(key string) string {
	return ec.c.QueryParam(key)
}

// Body returns the request body
func (ec *EchoContext) Body() ([]byte, error) {
	return io.ReadAll(ec.c.Request().Body)
}

// JSON sends a JSON response
func (ec *EchoContext) JSON(code int, data interface{}) error {
	return ec.c.JSON(code, data)
}

// String sends a string response
func (ec *EchoContext) String(code int, data string) error {
	return ec.c.String(code, data)
}

// Set stores a value in the context
func (ec *EchoContext) Set(key string, value interface{}) {
	ec.c.Set(key, value)
}

// Get retrieves a value from the context
func (ec *EchoContext) Get(key string) interface{} {
	return ec.c.Get(key)
}
