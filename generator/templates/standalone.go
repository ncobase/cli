package templates

import "fmt"

// StandaloneMainTemplate generates the main.go file for standalone applications
func StandaloneMainTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{ .PackagePath }}/internal/server"
	"{{ .PackagePath }}/internal/version"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/logging/logger"
)

const (
	shutdownTimeout = 3 * time.Second // service shutdown timeout
)

func main() {
	// Handle version flags
	version.Flags()

	// Set logger version
	logger.SetVersion(version.GetVersionInfo().Version)

	// load config
	conf := loadConfig()

	// Application name
	appName := conf.AppName

	// Initialize logger
	cleanupLogger := initializeLogger(conf)
	defer cleanupLogger()

	logger.Infof(context.Background(), "Starting %%s", appName)

	if err := runServer(conf); err != nil {
		logger.Fatalf(context.Background(), "Server error: %%v", err)
	}
}

// runServer creates and runs HTTP server
func runServer(conf *config.Config) error {
	// create server
	s, err := server.New(conf)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	defer s.Cleanup()

	// create listener
	listener, err := createListener(conf)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	// create server instance
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: s.Handler(),
	}

	// create error channel
	errChan := make(chan error, 1)

	// start server
	go func() {
		logger.Infof(context.Background(), "Listening and serving HTTP on: %%s", srv.Addr)
		if err := srv.Serve(listener); err != nil {
			logger.Errorf(context.Background(), "Listen error: %%s", err)
		} else {
			logger.Infof(context.Background(), "Server closed")
		}
	}()

	return gracefulShutdown(srv, errChan)
}

// createListener creates network listener
func createListener(conf *config.Config) (net.Listener, error) {
	addr := fmt.Sprintf("%%s:%%d", conf.Host, conf.Port)
	if conf.Port == 0 {
		addr = fmt.Sprintf("%%s:0", conf.Host)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %%w", err)
	}

	// update port if dynamically allocated
	if conf.Port == 0 {
		conf.Port = listener.Addr().(*net.TCPAddr).Port
	}

	return listener, nil
}

// loadConfig loads the application configuration
func loadConfig() *config.Config {
	conf, err := config.Init()
	if err != nil {
		logger.Fatalf(context.Background(), "[Config] Initialization error: %%+v", err)
	}
	return conf
}

// initializeLogger initializes the logger
func initializeLogger(conf *config.Config) func() {
	l, err := logger.New(conf.Logger)
	if err != nil {
		logger.Fatalf(context.Background(), "[Logger] Initialization error: %%+v", err)
	}
	return l
}

// gracefulShutdown gracefully shuts down the server
func gracefulShutdown(srv *http.Server, errChan chan error) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %%w", err)

	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Execute shutdown logic
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(context.Background(), "Shutdown error: %%v", err)
			return fmt.Errorf("shutdown error: %%w", err)
		}

		// wait for server to shutdown
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			logger.Debugf(context.Background(), "Shutdown timed out after %%s", shutdownTimeout)
		} else {
			logger.Debugf(context.Background(), "Shutdown completed within %%s", shutdownTimeout)
		}

		return nil
	}
}
`, name)
}

// StandaloneServerTemplate generates the server.go file for internal/server
func StandaloneServerTemplate(name, moduleName string) string {
	return `package server

import (
	"context"
	"net/http"

	"github.com/ncobase/ncore/config"
	extm "github.com/ncobase/ncore/extension/manager"
	"github.com/ncobase/ncore/logging/logger"

	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/data/repository"
	"{{ .PackagePath }}/handler"
	"{{ .PackagePath }}/service"
)

// Server represents the application server
type Server struct {
	config  *config.Config
	handler http.Handler
	cleanup func()
}

// New creates a new server instance.
func New(conf *config.Config) (*Server, error) {
	// Initialize Extension Manager
	em, err := extm.NewManager(conf)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed initializing extension manager: %+v", err)
		return nil, err
	}

	// Initialize Data Layer (Database)
	d, cleanupData, err := data.New(conf.Data, conf.Environment)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed initializing data layer: %+v", err)
		return nil, err
	}

	// Initialize Repository Layer
	repo := repository.New(d)

	// Initialize Service Layer
	svc := service.New(repo)

	// Initialize Handler Layer
	h := handler.New(svc)

	// Initialize HTTP server
	router, err := ginServer(conf, h)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed initializing http server: %+v", err)
		cleanupData()
		return nil, err
	}

	return &Server{
		config:  conf,
		handler: router,
		cleanup: func() {
			// Cleanup when server shuts down
			if err := svc.Close(); err != nil {
				logger.Errorf(context.Background(), "Error during service cleanup: %+v", err)
			}
			if cleanupData != nil {
				cleanupData()
			}
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
`
}

// StandaloneGinTemplate generates the gin.go file for internal/server
func StandaloneGinTemplate(name, moduleName string) string {
	return `package server

import (
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/ecode"
	"github.com/ncobase/ncore/net/resp"
	"net/http"
	"{{ .PackagePath }}/handler"
	appConfig "{{ .PackagePath }}/internal/config"

	"github.com/gin-gonic/gin"
)

// ginServer creates and initializes the server.
func ginServer(conf *config.Config, h handler.Interface) (*gin.Engine, error) {
	// Set gin mode
	if appConfig.IsProd(conf) {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create gin engine
	engine := gin.New()

	// Add middleware here if needed
	// engine.Use(middleware.Logger)
	// engine.Use(middleware.CORSHandler)

	// Register API routes
	registerRest(engine, conf, h)

	// No route
	engine.NoRoute(func(c *gin.Context) {
		resp.Fail(c.Writer, resp.NotFound(ecode.Text(http.StatusNotFound)))
	})

	return engine, nil
}
`
}

// StandaloneRestTemplate generates the rest.go file for internal/server
func StandaloneRestTemplate(name, moduleName string) string {
	return `package server

import (
	"net/http"
	"github.com/ncobase/ncore/config"
	"{{ .PackagePath }}/internal/version"
	"{{ .PackagePath }}/handler"

	"github.com/gin-gonic/gin"
)

// registerRest registers the REST routes.
func registerRest(e *gin.Engine, conf *config.Config, h handler.Interface) {
	// Root endpoint
	e.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Service is running.")
	})

	// Health check endpoint
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"name":   conf.AppName,
			"version": version.GetVersionInfo().Version,
		})
	})

	// API v1 routes
	v1 := e.Group("/api/v1")
	{
		// Add your API routes here
		v1.GET("/example", h.GetExample)
	}

	// Version endpoint
	e.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.GetVersionInfo())
	})
}
`
}

// StandaloneConfigTemplate generates the config.go file for standalone applications
func StandaloneConfigTemplate(name, moduleName string) string {
	return `package config

import (
	"github.com/ncobase/ncore/config"
)

// GetAppConfig returns the application-specific configuration
func GetAppConfig(cfg *config.Config) *config.Config {
	if cfg == nil {
		return &config.Config{}
	}
	return cfg
}

// IsProd returns current environment is production
func IsProd(c *config.Config) bool {
	if c == nil {
		return false
	}
	return c.IsProd()
}
`
}

// StandaloneHandlerProviderTemplate generates the provider.go file for handler layer
func StandaloneHandlerProviderTemplate(name, moduleName string) string {
	return `package handler

import (
	"github.com/gin-gonic/gin"
	"{{ .PackagePath }}/service"
)

// Interface defines the handler interface
type Interface interface {
	GetExample(c *gin.Context)
}

// handler implements the Interface
type handler struct {
	svc service.Interface
}

// New creates a new handler
func New(svc service.Interface) Interface {
	return &handler{
		svc: svc,
	}
}
`
}

// StandaloneHandlerTemplate generates the handler.go file (implementation)
func StandaloneHandlerTemplate(name, moduleName string) string {
	return `package handler

import (
	"github.com/ncobase/ncore/net/resp"
	"github.com/gin-gonic/gin"
)

// GetExample is an example handler
func (h *handler) GetExample(c *gin.Context) {
	result, err := h.svc.GetExample(c.Request.Context())
	if err != nil {
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
		return
	}

	resp.Success(c.Writer, result)
}
`
}

// StandaloneModelTemplate generates the model.go file for standalone applications
func StandaloneModelTemplate(name, moduleName string) string {
	return `package model

import (
	"time"
)

// Example represents an example model
type Example struct {
	ID        string     ` + "`" + `json:"id"` + "`" + `
	Name      string     ` + "`" + `json:"name"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time  ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt *time.Time ` + "`" + `json:"deleted_at,omitempty"` + "`" + `
}

// ExampleResponse represents an example response
type ExampleResponse struct {
	ID   string ` + "`" + `json:"id"` + "`" + `
	Name string ` + "`" + `json:"name"` + "`" + `
}
`
}

// StandaloneServiceProviderTemplate generates the provider.go file for service layer
func StandaloneServiceProviderTemplate(name, moduleName string) string {
	return `package service

import (
	"context"
	"{{ .PackagePath }}/data/repository"
	"{{ .PackagePath }}/data/model"
)

// Interface defines the service interface
type Interface interface {
	GetExample(ctx context.Context) (*model.ExampleResponse, error)
	Close() error
}

// service implements the Interface
type service struct {
	repo repository.Interface
}

// New creates a new service
func New(repo repository.Interface) Interface {
	return &service{
		repo: repo,
	}
}
`
}

// StandaloneServiceTemplate generates the service.go file (implementation)
func StandaloneServiceTemplate(name, moduleName string) string {
	return `package service

import (
	"context"
	"{{ .PackagePath }}/data/model"
	"time"

	"github.com/google/uuid"
)

// Close performs cleanup operations for the service
func (s *service) Close() error {
	return nil
}

// GetExample is an example service method
func (s *service) GetExample(ctx context.Context) (*model.ExampleResponse, error) {
	// Call repository (mock example)
	// result, err := s.repo.FindExample(ctx, "id")

	// This is a mock example
	example := &model.Example{
		ID:        uuid.New().String(),
		Name:      "Example Model",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &model.ExampleResponse{
		ID:   example.ID,
		Name: example.Name,
	},
nil
}
`
}

// StandaloneRepositoryProviderTemplate generates the provider.go file for repository layer
func StandaloneRepositoryProviderTemplate(name, moduleName string) string {
	return `package repository

import (
	"context"
	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/data/model"
)

// Interface defines the repository interface
type Interface interface {
	FindExample(ctx context.Context, id string) (*model.Example, error)
}

// repository implements the Interface
type repository struct {
	d *data.Data
}

// New creates a new repository
func New(d *data.Data) Interface {
	return &repository{
		d: d,
	}
}
`
}

// StandaloneRepositoryTemplate generates the repository.go file (implementation)
func StandaloneRepositoryTemplate(name, moduleName string, useMongo, useEnt, useGorm bool) string {
	var imports string

	// Basic imports
	imports = fmt.Sprintf(`import (
	"context"
	"fmt"
	"{{ .PackagePath }}/data/model"
`)

	// Add DB specific imports
	if useMongo {
		// imports += `	"go.mongodb.org/mongo-driver/mongo"`
	}
	if useEnt {
		imports += fmt.Sprintf(`	"{{ .PackagePath }}/data/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
`)
	}
	if useGorm {
		imports += `	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
`
	}

	imports += `)
`

	// This is just implementation methods, struct is defined in provider.go

	return fmt.Sprintf(`package repository

%s

// FindExample finds an example by ID
func (r *repository) FindExample(ctx context.Context, id string) (*model.Example, error) {
	// Access database via r.d
	// e.g., r.d.Conn.DB (GORM), r.d.Conn.Ent (Ent), r.d.MC (MongoDB)

	// Example: using search API from ncore/data
	// req := &search.IndexRequest{Index: "examples", Document: model.Example{ID: id}}
	// err := r.d.IndexDocument(ctx, req)

	fmt.Println("FindExample called")

	// Mock return
	return nil, nil
}
`, imports)
}
