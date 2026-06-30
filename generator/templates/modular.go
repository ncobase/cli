package templates

import "fmt"

func ModularServerHTTPTemplate(packagePath string) string {
	return fmt.Sprintf(`package server

import (
	"net/http"
	"strings"

	"%s/internal/middleware"
	"%s/internal/version"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/ecode"
	ext "github.com/ncobase/ncore/extension/types"
	"github.com/ncobase/ncore/net/resp"
	"github.com/ncobase/ncore/security/jwt"
)

// HTTPConfig holds HTTP server configuration.
type HTTPConfig struct {
	Mode        string
	Middlewares []gin.HandlerFunc
}

type guardedExtensionManager interface {
	ManageRoutesWithGuard(router *gin.RouterGroup, guards ...gin.HandlerFunc)
}

// newHTTPHandler creates the Gin HTTP handler.
func newHTTPHandler(conf *config.Config, em ext.ManagerInterface) (http.Handler, error) {
	httpConf := &HTTPConfig{
		Mode: validateGinMode(conf),
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

func createGinEngine(conf *HTTPConfig, em ext.ManagerInterface, appConf *config.Config) (*gin.Engine, error) {
	gin.SetMode(conf.Mode)

	engine := gin.New()
	for _, mw := range conf.Middlewares {
		engine.Use(mw)
	}

	registerMachineRoutes(engine, appConf)
	em.RegisterRoutes(engine)

	if appConf.Extension != nil && (appConf.Extension.HotReload || (appConf.Extension.Metrics != nil && appConf.Extension.Metrics.Enabled)) {
		registerNCoreManagementRoutes(engine, em, appConf)
	}

	setupNoRouteHandler(engine)
	return engine, nil
}

func registerMachineRoutes(engine *gin.Engine, conf *config.Config) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "running",
			"name":    conf.AppName,
			"version": version.GetVersionInfo().Version,
		})
	})
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"name":   conf.AppName,
		})
	})
	engine.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"name":   conf.AppName,
		})
	})
	engine.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.GetVersionInfo())
	})
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

func setupNoRouteHandler(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		resp.Fail(c.Writer, resp.NotFound(ecode.Text(http.StatusNotFound)))
	})
	engine.NoMethod()
}

func validateGinMode(conf *config.Config) string {
	if conf == nil {
		return gin.ReleaseMode
	}
	switch conf.Environment {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:
		return conf.Environment
	default:
		if conf.IsProd() {
			return gin.ReleaseMode
		}
		return gin.DebugMode
	}
}
`, packagePath, packagePath)
}

func ModularReadmeTemplate() string {
	return `# {{ .Name }}

{{ .Name }} is a modular ncore backend application.

## Structure

    .
    |-- cmd/{{ .Name }}
    |-- core
    |-- biz
    |-- plugin
    |-- internal
    |   |-- middleware
    |   |-- server
    |   |-- version
    |-- migrations
    |-- docs
    |-- config.yaml
    |-- go.mod
    |-- Makefile

Use core for foundational domains, biz for product domains, and plugin for optional
integrations. Keep domain routes in module-local router packages and keep internal/server
focused on process startup, middleware, health, version, extension registration, and ncore
management route protection.

## Add Modules

    nco create core auth --path . --use-ent --with-test
    nco create biz content --path . --use-ent --with-test
    nco create plugin resource --path . --use-ent --with-test

Module implementation order is structs, data/schema, data/repository, service, handler, router,
extension entrypoint, optional wrapper or event, then tests.

## Development

    go test ./...
    go run ./cmd/{{ .Name }} -conf ./config.yaml

## Configuration

config.yaml is for infrastructure, bootstrap, secrets, and security boundaries. Runtime product
policy belongs in application options or settings owned by the product.
`
}

func PackageDocTemplate(packageName, description string) string {
	return fmt.Sprintf("// Package %s %s.\npackage %s\n", packageName, description, packageName)
}

func MigrationReadmeTemplate() string {
	return `# Migrations

Store database migration files here when the application owns migration artifacts outside module
Ent generation.
`
}

func DocsReadmeTemplate() string {
	return `# Documentation

Keep API, operation, and architecture documents for this application here.
`
}
