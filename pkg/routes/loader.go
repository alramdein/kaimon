package routes

// WARNING: This is a core package. Do NOT modify unless you're changing route loading logic.
// For adding routes, edit JSON files in config/routes/ instead.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/alramdein/kaimon/pkg/framework"
	"github.com/alramdein/kaimon/pkg/middleware"
)

// Loader loads and registers routes
type Loader struct {
	router            framework.Router
	middlewareManager *middleware.Manager
}

// NewLoader creates a new route loader
func NewLoader(router framework.Router, middlewareManager *middleware.Manager) *Loader {
	return &Loader{
		router:            router,
		middlewareManager: middlewareManager,
	}
}

// LoadFromFile loads routes from a compiled routes file
func (l *Loader) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read routes file: %w", err)
	}

	var compiled CompiledRoutes
	if err := json.Unmarshal(data, &compiled); err != nil {
		return fmt.Errorf("failed to parse routes file: %w", err)
	}

	return l.Load(&compiled)
}

// Load loads routes from compiled configuration
func (l *Loader) Load(compiled *CompiledRoutes) error {
	// Apply global middlewares
	if compiled.Middlewares != nil {
		globalMiddlewares := make([]framework.MiddlewareFunc, 0)

		// Add onRequest middlewares
		if len(compiled.Middlewares.OnRequest) > 0 {
			onRequestMws := l.middlewareManager.GetMiddlewares(compiled.Middlewares.OnRequest)
			globalMiddlewares = append(globalMiddlewares, onRequestMws...)
		}

		// Add onResponse middlewares
		if len(compiled.Middlewares.OnResponse) > 0 {
			onResponseMws := l.middlewareManager.GetMiddlewares(compiled.Middlewares.OnResponse)
			globalMiddlewares = append(globalMiddlewares, onResponseMws...)
		}

		if len(globalMiddlewares) > 0 {
			l.router.Use(globalMiddlewares...)
		}
	}

	// Register routes
	for _, route := range compiled.Routes {
		handler := l.createProxyHandler(route)

		// Wrap with route-specific middlewares
		if route.Middlewares != nil {
			routeMiddlewares := make([]framework.MiddlewareFunc, 0)

			// Add route onRequest middlewares
			if len(route.Middlewares.OnRequest) > 0 {
				onRequestMws := l.middlewareManager.GetMiddlewares(route.Middlewares.OnRequest)
				routeMiddlewares = append(routeMiddlewares, onRequestMws...)
			}

			// Add route onResponse middlewares
			if len(route.Middlewares.OnResponse) > 0 {
				onResponseMws := l.middlewareManager.GetMiddlewares(route.Middlewares.OnResponse)
				routeMiddlewares = append(routeMiddlewares, onResponseMws...)
			}

			// Apply middlewares in reverse order
			for i := len(routeMiddlewares) - 1; i >= 0; i-- {
				handler = routeMiddlewares[i](handler)
			}
		}

		// Register based on method
		switch strings.ToUpper(route.Method) {
		case "GET":
			l.router.GET(route.Path, handler)
		case "POST":
			l.router.POST(route.Path, handler)
		case "PUT":
			l.router.PUT(route.Path, handler)
		case "DELETE":
			l.router.DELETE(route.Path, handler)
		case "PATCH":
			l.router.PATCH(route.Path, handler)
		default:
			return fmt.Errorf("unsupported method: %s", route.Method)
		}
	}

	return nil
}

// createProxyHandler creates a proxy handler for the route
func (l *Loader) createProxyHandler(route Route) framework.HandlerFunc {
	return func(ctx framework.Context) error {
		// Parse target URL
		targetURL, err := url.Parse(route.Target)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "invalid target URL",
			})
		}

		// Create proxy request
		req := ctx.Request()
		proxyURL := targetURL.Scheme + "://" + targetURL.Host + targetURL.Path
		if req.URL.RawQuery != "" {
			proxyURL += "?" + req.URL.RawQuery
		}

		proxyReq, err := http.NewRequest(req.Method, proxyURL, req.Body)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to create proxy request",
			})
		}

		// Copy headers
		for key, values := range req.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Add custom headers from route config
		for key, value := range route.Headers {
			proxyReq.Header.Set(key, value)
		}

		// Execute proxy request
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, map[string]string{
				"error": "failed to proxy request",
			})
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				ctx.Response().Header().Add(key, value)
			}
		}

		// Copy response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to read response",
			})
		}

		ctx.Response().WriteHeader(resp.StatusCode)
		ctx.Response().Write(body)

		return nil
	}
}
