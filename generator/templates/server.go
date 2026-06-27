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
	"strings"

	"%s/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/ecode"
	ext "github.com/ncobase/ncore/extension/types"
	"github.com/ncobase/ncore/net/resp"
	"github.com/ncobase/ncore/security/jwt"
)

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Mode        string
	Middlewares []gin.HandlerFunc
}

type guardedExtensionManager interface {
	ManageRoutesWithGuard(router *gin.RouterGroup, guards ...gin.HandlerFunc)
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
		registerNCoreManagementRoutes(engine, em, config)
	}

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "running",
		})
	})

	setupNoRouteHandler(engine)
	return engine, nil
}

func registerNCoreManagementRoutes(engine *gin.Engine, em ext.ManagerInterface, conf *config.Config) {
	guard := ncoreManagementGuard(conf)
	if guardedManager, ok := em.(guardedExtensionManager); ok {
		guardedManager.ManageRoutesWithGuard(engine.Group("/ncore"), guard)
		return
	}
	em.ManageRoutes(engine.Group("/ncore", guard))
}

func ncoreManagementGuard(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf == nil || conf.Auth == nil || conf.Auth.JWT == nil || conf.Auth.JWT.Secret == "" {
			resp.Fail(c.Writer, resp.ServiceUnavailable("JWT secret is not configured"))
			c.Abort()
			return
		}

		token := extractBearerToken(c)
		if token == "" {
			resp.Fail(c.Writer, resp.UnAuthorized("Bearer token is required"))
			c.Abort()
			return
		}

		payload, err := jwt.NewTokenManager(conf.Auth.JWT.Secret).GetPayload(token)
		if err != nil {
			resp.Fail(c.Writer, resp.UnAuthorized("Invalid or expired token"))
			c.Abort()
			return
		}

		if !payloadAllowsPermission(payload, "manage:ncore") {
			resp.Fail(c.Writer, resp.Forbidden("Permission manage:ncore is required"))
			c.Abort()
			return
		}

		if userID := stringValue(payload["user_id"]); userID != "" {
			c.Set("user_id", userID)
		}
		c.Set("permissions", stringSliceValue(payload["permissions"]))
		c.Next()
	}
}

func extractBearerToken(c *gin.Context) string {
	header := strings.TrimSpace(c.GetHeader("Authorization"))
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func payloadAllowsPermission(payload map[string]any, required string) bool {
	if boolValue(payload["is_admin"]) || boolValue(payload["admin"]) {
		return true
	}
	for _, role := range stringSliceValue(payload["roles"]) {
		if role == "admin" || role == "super_admin" || role == "root" {
			return true
		}
	}
	for _, permission := range stringSliceValue(payload["permissions"]) {
		if permissionMatches(permission, required) {
			return true
		}
	}
	for _, permission := range stringSliceValue(payload["permission_codes"]) {
		if permissionMatches(permission, required) {
			return true
		}
	}
	return false
}

func permissionMatches(granted, required string) bool {
	granted = strings.TrimSpace(granted)
	required = strings.TrimSpace(required)
	if granted == "" || required == "" {
		return false
	}
	if granted == "*" || granted == "*:*" || granted == "admin:*" || granted == "super:*" || strings.EqualFold(granted, required) {
		return true
	}

	grantedParts := strings.Split(granted, ":")
	requiredParts := strings.Split(required, ":")
	if len(grantedParts) != 2 || len(requiredParts) != 2 {
		return false
	}

	actionMatches := grantedParts[0] == "*" ||
		strings.EqualFold(grantedParts[0], "admin") ||
		strings.EqualFold(grantedParts[0], "super") ||
		strings.EqualFold(grantedParts[0], requiredParts[0])
	resourceMatches := grantedParts[1] == "*" || strings.EqualFold(grantedParts[1], requiredParts[1])
	return actionMatches && resourceMatches
}

func stringSliceValue(value any) []string {
	switch v := value.(type) {
	case []string:
		return v
	case []any:
		values := make([]string, 0, len(v))
		for _, item := range v {
			if s := stringValue(item); s != "" {
				values = append(values, s)
			}
		}
		return values
	default:
		if s := stringValue(value); s != "" {
			return []string{s}
		}
		return nil
	}
}

func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return ""
	}
}

func boolValue(value any) bool {
	v, ok := value.(bool)
	return ok && v
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
