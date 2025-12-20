package templates

import "fmt"

// CmdMainTemplate generates the main.go file for the cmd directory
func CmdMainTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package main

import (
	"context"
	"fmt"
	"os"

	"{{ .PackagePath }}"
	"{{ .PackagePath }}/cmd/provider"

	"github.com/ncobase/ncore/logging/logger"
	"github.com/ncobase/ncore/version"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

func main() {
	logger.SetVersion(version.GetVersionInfo().Version)

	rootCmd := &cobra.Command{
		Use:	"%s",
		Short: "%s service",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize extension
			cleanup, err := provider.NewExtension(configFile)
			if err != nil {
				return err
			}
			defer cleanup()
			
			// This module is already registered in registry by importing the package
			// The registry will handle lifecycle management
			
			// Wait for interrupt signal
			<-cmd.Context().Done()
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "config file path")

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Import the module to register it
var _ = %s.New
`, name, name, name)
}

// CmdServerTemplate generates the server.go file for the provider directory
func CmdServerTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package provider

import (
	"context"
	extm "github.com/ncobase/ncore/extension/manager"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/logging/logger"
	"net/http"
)

// NewServer creates a new server.
func NewServer(conf *config.Config) (http.Handler, func(), error) {
	// Initialize Extension Manager
	em, err := extm.NewManager(conf)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed initializing extension manager: %%+v", err)
		return nil, nil, err
	}

	// Register extensions
	registerExtensions(em)
	if err := em.LoadPlugins(); err != nil {
		logger.Fatalf(context.Background(), "Failed loading plugins: %%+v", err)
	}

	// New server
	h, err := ginServer(conf, em)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed initializing http: %%+v", err)
	}

	return h, func() {
		em.Cleanup()
	}, nil
}
`)
}

// CmdExtensionTemplate generates the extension.go file for the provider directory
func CmdExtensionTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package provider

import (
	"context"
	"fmt"
	
	"github.com/ncobase/ncore/config"
	exr "github.com/ncobase/ncore/extension/registry"
	"github.com/ncobase/ncore/logging/logger"
)

// NewExtension initializes the extension
func NewExtension(configFile string) (func(), error) {
	// Initialize config
	conf, err := config.Init(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %%w", err)
	}

	// Initialize logger
	cleanupLogger, err := logger.New(conf.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %%w", err)
	}

	// Initialize extension manager
	em := exr.NewManager(conf)
	
	// Initialize all registered extensions (including this one)
	if err := em.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize extensions: %%w", err)
	}
	
	// Start extensions
	if err := em.PostInit(); err != nil {
		return nil, fmt.Errorf("failed to start extensions: %%w", err)
	}

	return func() {
		// Cleanup extensions
		em.Cleanup()
		cleanupLogger()
	}, nil
}
`)
}

// CmdGinTemplate generates the gin.go file for the provider directory
func CmdGinTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package provider

import (
	ext "github.com/ncobase/ncore/extension/types"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/ecode"
	"github.com/ncobase/ncore/net/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ginServer creates and initializes the server.
func ginServer(conf *config.Config, em ext.ManagerInterface) (*gin.Engine, error) {
	// Set gin mode
	if conf.RunMode == "" {
		conf.RunMode = gin.ReleaseMode
	}
	// Set mode before creating engine
	gin.SetMode(conf.RunMode)
	// Create gin engine
	engine := gin.New()

	// Add middleware here if needed
	// engine.Use(middleware.Logger)
	// engine.Use(middleware.CORSHandler)

	// Register REST
	registerRest(engine, conf)

	// Register extension / plugin routes
	em.RegisterRoutes(engine)

	// Register extension management routes
	if conf.Extension.HotReload {
		g := engine.Group("/sys")
		em.ManageRoutes(g)
	}

	// No route
	engine.NoRoute(func(c *gin.Context) {
		resp.Fail(c.Writer, resp.NotFound(ecode.Text(http.StatusNotFound)))
	})

	return engine, nil
}
`)
}

// CmdRestTemplate generates the rest.go file for the provider directory
func CmdRestTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package provider

import (
	"github.com/ncobase/ncore/helper"
	"net/http"

	"github.com/ncobase/ncore/config"

	"github.com/gin-gonic/gin"
)

// registerRest registers the REST routes.
func registerRest(e *gin.Engine, conf *config.Config) {
	// Root endpoint
	e.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Service is running.")
	})

	// Health check endpoint
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"name":   conf.AppName,
			"version": version.Version,
		})
	})

	// Add your API routes here
}
`)
}
