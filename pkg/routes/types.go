package routes

// WARNING: This is a core package. Do NOT modify unless you're changing route configuration structure.
// For adding routes, edit JSON files in config/routes/ instead.

// MiddlewareConfig represents middleware configuration with phases
type MiddlewareConfig struct {
	OnRequest  []string `json:"onRequest,omitempty"`
	OnResponse []string `json:"onResponse,omitempty"`
}

// Route represents a single route configuration
type Route struct {
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	Target      string            `json:"target"`
	Middlewares *MiddlewareConfig `json:"middlewares,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// RouteConfig represents the configuration for a domain
type RouteConfig struct {
	Domain      string            `json:"domain"`
	BasePath    string            `json:"basePath"`
	Routes      []Route           `json:"routes"`
	Middlewares *MiddlewareConfig `json:"middlewares,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// GlobalConfig represents global configuration for all routes
type GlobalConfig struct {
	Middlewares *MiddlewareConfig `json:"middlewares,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// CompiledRoutes represents the compiled route configuration
type CompiledRoutes struct {
	Middlewares *MiddlewareConfig `json:"middlewares,omitempty"`
	Routes      []Route           `json:"routes"`
}
