package middleware

// WARNING: This is a core package. Do NOT modify unless you're changing middleware registration.
// For adding middlewares, create files in internal/middlewares/onRequest or onResponse instead.
// Middlewares auto-register via init() functions.

// Global registry for self-registration
var globalRegistry = struct {
	onRequest  map[string]func() OnRequestMiddleware
	onResponse map[string]func() OnResponseMiddleware
}{
	onRequest:  make(map[string]func() OnRequestMiddleware),
	onResponse: make(map[string]func() OnResponseMiddleware),
}

// RegisterOnRequestFactory registers an onRequest middleware factory
// This should be called in init() functions of middleware files
func RegisterOnRequestFactory(name string, factory func() OnRequestMiddleware) {
	globalRegistry.onRequest[name] = factory
}

// RegisterOnResponseFactory registers an onResponse middleware factory
// This should be called in init() functions of middleware files
func RegisterOnResponseFactory(name string, factory func() OnResponseMiddleware) {
	globalRegistry.onResponse[name] = factory
}

// LoadFromRegistry loads all middlewares from the global registry
func (m *Manager) LoadFromRegistry() {
	// Load all onRequest middlewares
	for _, factory := range globalRegistry.onRequest {
		m.RegisterOnRequest(factory())
	}

	// Load all onResponse middlewares
	for _, factory := range globalRegistry.onResponse {
		m.RegisterOnResponse(factory())
	}
}
