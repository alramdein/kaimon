package routes

// WARNING: This is a core package. Do NOT modify unless you're changing route compilation logic.
// For adding routes, edit JSON files in config/routes/ instead.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Compiler compiles route configurations from multiple files
type Compiler struct {
	configDir  string
	outputDir  string
	globalFile string
}

// NewCompiler creates a new compiler instance
func NewCompiler(configDir, outputDir, globalFile string) *Compiler {
	return &Compiler{
		configDir:  configDir,
		outputDir:  outputDir,
		globalFile: globalFile,
	}
}

// Compile reads all route configs and compiles them into a single file
func (c *Compiler) Compile() error {
	compiled := &CompiledRoutes{
		Middlewares: &MiddlewareConfig{
			OnRequest:  make([]string, 0),
			OnResponse: make([]string, 0),
		},
		Routes: make([]Route, 0),
	}

	// Load global config first
	if c.globalFile != "" {
		if err := c.loadGlobalConfig(compiled); err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}
	}

	// Read domain route configs
	files, err := os.ReadDir(c.configDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(c.configDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
		}

		var config RouteConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse file %s: %w", file.Name(), err)
		}

		// Process routes
		for _, route := range config.Routes {
			compiledRoute := route

			// Prepend base path if defined
			if config.BasePath != "" {
				compiledRoute.Path = config.BasePath + route.Path
			}

			// Merge domain-level middlewares if route doesn't have its own
			if compiledRoute.Middlewares == nil && config.Middlewares != nil {
				compiledRoute.Middlewares = config.Middlewares
			}

			// Merge headers
			if compiledRoute.Headers == nil {
				compiledRoute.Headers = make(map[string]string)
			}
			for k, v := range config.Headers {
				if _, exists := compiledRoute.Headers[k]; !exists {
					compiledRoute.Headers[k] = v
				}
			}

			compiled.Routes = append(compiled.Routes, compiledRoute)
		}
	}

	// Write compiled routes
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(c.outputDir, "routes.json")
	data, err := json.MarshalIndent(compiled, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal compiled routes: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write compiled routes: %w", err)
	}

	return nil
}

// loadGlobalConfig loads global configuration
func (c *Compiler) loadGlobalConfig(compiled *CompiledRoutes) error {
	data, err := os.ReadFile(c.globalFile)
	if err != nil {
		// If global file doesn't exist, it's okay
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var global GlobalConfig
	if err := json.Unmarshal(data, &global); err != nil {
		return err
	}

	// Set global middlewares
	if global.Middlewares != nil {
		compiled.Middlewares = global.Middlewares
	}

	return nil
}
