package middleware

// WARNING: This is a core package. Do NOT modify unless you're changing middleware management logic.
// For adding middlewares, create files in internal/middlewares/onRequest or onResponse instead.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alramdein/kaimon/pkg/framework"
)

// ExecutionOrder represents middleware execution order configuration
type ExecutionOrder struct {
	OnRequest  []string `json:"onRequest"`
	OnResponse []string `json:"onResponse"`
}

// Manager manages middleware registration and execution
type Manager struct {
	onRequestMiddlewares  map[string]OnRequestMiddleware
	onResponseMiddlewares map[string]OnResponseMiddleware
	executionOrder        *ExecutionOrder
}

// NewManager creates a new middleware manager
func NewManager() *Manager {
	return &Manager{
		onRequestMiddlewares:  make(map[string]OnRequestMiddleware),
		onResponseMiddlewares: make(map[string]OnResponseMiddleware),
		executionOrder:        &ExecutionOrder{},
	}
}

// RegisterOnRequest registers an onRequest middleware
func (m *Manager) RegisterOnRequest(middleware OnRequestMiddleware) {
	m.onRequestMiddlewares[middleware.Name()] = middleware
}

// RegisterOnResponse registers an onResponse middleware
func (m *Manager) RegisterOnResponse(middleware OnResponseMiddleware) {
	m.onResponseMiddlewares[middleware.Name()] = middleware
}

// LoadExecutionOrder loads middleware execution order from file
func (m *Manager) LoadExecutionOrder(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read execution order file: %w", err)
	}

	if err := json.Unmarshal(data, m.executionOrder); err != nil {
		return fmt.Errorf("failed to parse execution order: %w", err)
	}

	return nil
}

// AutoDiscoverMiddlewares scans directories and auto-registers middlewares
func (m *Manager) AutoDiscoverMiddlewares(onRequestDir, onResponseDir string) error {
	// Note: Go doesn't support dynamic loading like Node.js
	// This is a placeholder for the concept - in practice, you'd need to
	// either use build tags, code generation, or manual registration
	// For now, we'll scan directories but still require manual registration

	// Scan onRequest directory
	if _, err := os.Stat(onRequestDir); err == nil {
		files, err := os.ReadDir(onRequestDir)
		if err != nil {
			return fmt.Errorf("failed to read onRequest directory: %w", err)
		}

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".go" {
				// Log discovered middleware files
				fmt.Printf("Discovered onRequest middleware: %s\n", file.Name())
			}
		}
	}

	// Scan onResponse directory
	if _, err := os.Stat(onResponseDir); err == nil {
		files, err := os.ReadDir(onResponseDir)
		if err != nil {
			return fmt.Errorf("failed to read onResponse directory: %w", err)
		}

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".go" {
				// Log discovered middleware files
				fmt.Printf("Discovered onResponse middleware: %s\n", file.Name())
			}
		}
	}

	return nil
}

// GetMiddlewares returns middleware functions based on names and phases
func (m *Manager) GetMiddlewares(names []string) []framework.MiddlewareFunc {
	middlewares := make([]framework.MiddlewareFunc, 0)

	for _, name := range names {
		if mw, exists := m.onRequestMiddlewares[name]; exists {
			middlewares = append(middlewares, mw.Handle)
		} else if mw, exists := m.onResponseMiddlewares[name]; exists {
			middlewares = append(middlewares, mw.Handle)
		}
	}

	return middlewares
}

// GetGlobalMiddlewares returns all global middlewares in execution order
func (m *Manager) GetGlobalMiddlewares() []framework.MiddlewareFunc {
	middlewares := make([]framework.MiddlewareFunc, 0)

	// Add onRequest middlewares
	for _, name := range m.executionOrder.OnRequest {
		if mw, exists := m.onRequestMiddlewares[name]; exists {
			middlewares = append(middlewares, mw.Handle)
		}
	}

	// Add onResponse middlewares
	for _, name := range m.executionOrder.OnResponse {
		if mw, exists := m.onResponseMiddlewares[name]; exists {
			middlewares = append(middlewares, mw.Handle)
		}
	}

	return middlewares
}
