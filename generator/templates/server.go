package templates

import "fmt"

// ServerTemplate generates the server.go file
func ServerTemplate(packagePath string) string {
	return fmt.Sprintf(`package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ncobase/ncore/config"
	extm "github.com/ncobase/ncore/extension/manager"
	"github.com/ncobase/ncore/logging/logger"
)

// Server represents the application server
type Server struct {
	config  *config.Config
	handler http.Handler
	cleanup func()
}

// New creates a new server instance
func New(conf *config.Config) (*Server, error) {
	ctx := context.Background()

	// Initialize components in dependency order
	em, err := initExtensionManager(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize extension manager: %%w", err)
	}

	// Create HTTP handler
	httpHandler, err := newHTTPHandler(conf, em)
	if err != nil {
		em.Cleanup()
		return nil, fmt.Errorf("failed to create HTTP handler: %%w", err)
	}

	return &Server{
		config:  conf,
		handler: httpHandler,
		cleanup: func() {
			em.Cleanup()
		},
	}, nil
}

// Handler returns the HTTP handler
func (s *Server) Handler() http.Handler {
	return s.handler
}

// Cleanup performs server cleanup
func (s *Server) Cleanup() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

// initExtensionManager initializes the extension manager
func initExtensionManager(ctx context.Context, conf *config.Config) (*extm.Manager, error) {
	em, err := extm.NewManager(conf)
	if err != nil {
		logger.Errorf(ctx, "Failed initializing extension manager: %%+v", err)
		return nil, err
	}

	// Register built-in extensions
	registerExtensions(em)

	// Load plugins
	if err = em.LoadPlugins(); err != nil {
		logger.Errorf(ctx, "Failed loading plugins: %%+v", err)
		return nil, err
	}

	return em, nil
}
`)
}

// ServerHTTPTemplate generates the http.go file
func ServerHTTPTemplate(packagePath string) string {
	return fmt.Sprintf(`package server

import (
	"net/http"

	"%s/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/ecode"
	ext "github.com/ncobase/ncore/extension/types"
	"github.com/ncobase/ncore/net/resp"
)

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Mode        string
	Middlewares []gin.HandlerFunc
}

// newHTTPHandler creates HTTP handler with middleware chain
func newHTTPHandler(conf *config.Config, em ext.ManagerInterface) (http.Handler, error) {
	ginMode := validateGinMode(conf)

	httpConf := &HTTPConfig{
		Mode:   ginMode,
		Middlewares: []gin.HandlerFunc{
			gin.Recovery(),
			middleware.CORSHandler(conf),
			middleware.SecurityHeaders(nil),
			middleware.InputValidation(nil),
			middleware.InputSanitization(nil),
			middleware.ClientInfo,
			middleware.Logger,
		},
	}

	return createGinEngine(httpConf, em, conf)
}

// createGinEngine creates and configures Gin engine
func createGinEngine(conf *HTTPConfig, em ext.ManagerInterface, config *config.Config) (*gin.Engine, error) {
	gin.SetMode(conf.Mode)

	engine := gin.New()

	for _, mw := range conf.Middlewares {
		engine.Use(mw)
	}

	em.RegisterRoutes(engine)

	if config.Extension != nil && (config.Extension.HotReload || (config.Extension.Metrics != nil && config.Extension.Metrics.Enabled)) {
		em.ManageRoutes(engine.Group("/ncore"))
	}

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "running",
		})
	})

	setupNoRouteHandler(engine)
	return engine, nil
}

// setupNoRouteHandler configures 404 handler
func setupNoRouteHandler(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		resp.Fail(c.Writer, resp.NotFound(ecode.Text(http.StatusNotFound)))
	})
	engine.NoMethod()
}

// validateGinMode validates and returns appropriate gin mode
func validateGinMode(conf *config.Config) string {
	// Check if Environment is one of the valid gin modes
	switch conf.Environment {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:
		return conf.Environment
	default:
		// Fallback based on production flag
		if conf.IsProd() {
			return gin.ReleaseMode
		}
		return gin.DebugMode
	}
}
`, packagePath)
}

// ServerExtsTemplate generates the exts.go file
func ServerExtsTemplate(packagePath string) string {
	return `package server

import (
	"context"

	ext "github.com/ncobase/ncore/extension/types"
	"github.com/ncobase/ncore/logging/logger"
)

// registerExtensions registers all built-in extensions
func registerExtensions(em ext.ManagerInterface) {
	// Registration is handled by the registry system through init() functions
	if err := em.InitExtensions(); err != nil {
		logger.Errorf(context.Background(), "Failed to initialize extensions: %v", err)
		return
	}
}
`
}
