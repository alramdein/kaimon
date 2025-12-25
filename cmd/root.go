package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/alramdein/kaimon/pkg/framework"
	"github.com/alramdein/kaimon/pkg/middleware"
	"github.com/alramdein/kaimon/pkg/routes"
	"github.com/spf13/cobra"

	// Import middlewares to trigger init() registration
	_ "github.com/alramdein/kaimon/internal/middlewares/onRequest"
	_ "github.com/alramdein/kaimon/internal/middlewares/onResponse"
)

var rootCmd = &cobra.Command{
	Use:   "kaimon",
	Short: "Kaimon API Gateway",
	Long:  "A simple, modular API gateway with framework abstraction",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(serveCmd)
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile route configurations",
	Long:  "Compile all route configuration files from config/routes into a single routes.json",
	Run: func(cmd *cobra.Command, args []string) {
		compiler := routes.NewCompiler("config/routes", "build", "config/global.json")
		if err := compiler.Compile(); err != nil {
			log.Fatalf("Failed to compile routes: %v", err)
		}
		log.Println("Routes compiled successfully to build/routes.json")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API gateway server",
	Long:  "Start the API gateway server with compiled routes",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize framework (Echo by default)
		fw := framework.NewEchoFramework()
		router := fw.Router()

		// Initialize middleware manager
		mwManager := middleware.NewManager()

		// Load all middlewares from registry (dynamic)
		mwManager.LoadFromRegistry()

		// Auto-discover middleware files (for logging purposes)
		if err := mwManager.AutoDiscoverMiddlewares(
			"internal/middlewares/onRequest",
			"internal/middlewares/onResponse",
		); err != nil {
			log.Printf("Warning: Failed to auto-discover middlewares: %v", err)
		}

		// Load routes
		loader := routes.NewLoader(router, mwManager)
		if err := loader.LoadFromFile("build/routes.json"); err != nil {
			log.Fatalf("Failed to load routes: %v", err)
		}

		// Start server
		log.Println("Starting Kaimon API Gateway on :8080")
		if err := fw.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}
