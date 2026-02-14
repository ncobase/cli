package templates

import "fmt"

// StandaloneMainTemplate generates the main.go file for standalone applications
func StandaloneMainTemplate(d *Data) string {
	imports := ""
	if d.DBDriver != "" && d.DBDriver != "none" {
		imports += fmt.Sprintf("\t_ \"github.com/ncobase/ncore/data/%s\"\n", d.DBDriver)
	}
	if d.UseRedis {
		imports += "\t_ \"github.com/ncobase/ncore/data/redis\"\n"
	}
	if d.UseElastic {
		imports += "\t_ \"github.com/ncobase/ncore/data/elasticsearch\"\n"
	}
	if d.UseOpenSearch {
		imports += "\t_ \"github.com/ncobase/ncore/data/opensearch\"\n"
	}
	if d.UseMeili {
		imports += "\t_ \"github.com/ncobase/ncore/data/meilisearch\"\n"
	}
	if d.UseKafka {
		imports += "\t_ \"github.com/ncobase/ncore/data/kafka\"\n"
	}
	if d.UseRabbitMQ {
		imports += "\t_ \"github.com/ncobase/ncore/data/rabbitmq\"\n"
	}

	return fmt.Sprintf(`package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"%s/internal/server"
	"%s/internal/version"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/logging/logger"
%s)

const (
	shutdownTimeout = 5 * time.Second // service shutdown timeout
)

func main() {
	flag.Parse()

	// Handle version flags
	version.Flags()

	// Load config
	conf := loadConfig()

	// Set logger version
	logger.SetVersion(version.GetVersionInfo().Version)

	// Initialize logger
	cleanupLogger := initializeLogger(conf)
	defer cleanupLogger()

	logger.Infof(context.Background(), "Starting %%s", conf.AppName)

	if err := runServer(conf); err != nil {
		logger.Fatalf(context.Background(), "Server error: %%v", err)
	}
}

// runServer creates and runs HTTP server
func runServer(conf *config.Config) error {
	// Create server
	s, err := server.New(conf)
	if err != nil {
		return fmt.Errorf("failed to create server: %%w", err)
	}
	defer s.Cleanup()

	// Create listener
	listener, err := createListener(conf)
	if err != nil {
		return fmt.Errorf("failed to create listener: %%w", err)
	}
	defer listener.Close()

	// Create HTTP server instance
	srv := &http.Server{
		Addr:         fmt.Sprintf("%%s:%%d", conf.Host, conf.Port),
		Handler:      s.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		logger.Infof(context.Background(), "Listening and serving HTTP on: %%s", srv.Addr)
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server error: %%w", err)
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
		return nil, fmt.Errorf("failed to create listener: %%w", err)
	}

	// Update port if dynamically allocated
	if conf.Port == 0 {
		conf.Port = listener.Addr().(*net.TCPAddr).Port
		logger.Infof(context.Background(), "Using dynamically allocated port: %%d", conf.Port)
	}

	return listener, nil
}

// loadConfig loads the application configuration
func loadConfig() *config.Config {
	conf, err := config.Init()
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to load config: %%+v", err)
	}
	return conf
}

// initializeLogger initializes the logger
func initializeLogger(conf *config.Config) func() {
	cleanup, err := logger.New(conf.Logger)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to initialize logger: %%+v", err)
	}
	return cleanup
}

// gracefulShutdown gracefully shuts down the server
func gracefulShutdown(srv *http.Server, errChan chan error) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err

	case sig := <-quit:
		logger.Infof(context.Background(), "Received signal: %%s, shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %%w", err)
		}

		logger.Info(context.Background(), "Server shutdown completed")
		return nil
	}
}
`, d.PackagePath, d.PackagePath, imports)
}

// StandaloneServerTemplate generates the server.go file for internal/server
func StandaloneServerTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ncobase/ncore/config"
	extm "github.com/ncobase/ncore/extension/manager"
	"github.com/ncobase/ncore/logging/logger"

	"%s/data"
	"%s/data/repository"
	"%s/handler"
	"%s/service"
)

// Server represents the application server
type Server struct {
	config  *config.Config
	handler http.Handler
	cleanup func()
}

// New creates a new server instance
func New(conf *config.Config) (*Server, error) {
	// Initialize Extension Manager
	em, err := extm.NewManager(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize extension manager: %%w", err)
	}

	// Initialize Data Layer (Database)
	d, cleanupData, err := data.New(conf.Data, conf.Environment)
	if err != nil {
		em.Cleanup()
		return nil, fmt.Errorf("failed to initialize data layer: %%w", err)
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
		if cleanupData != nil {
			cleanupData()
		}
		em.Cleanup()
		return nil, fmt.Errorf("failed to initialize HTTP server: %%w", err)
	}

	return &Server{
		config:  conf,
		handler: router,
		cleanup: func() {
			logger.Debug(context.Background(), "Cleaning up server resources...")

			if err := svc.Close(); err != nil {
				logger.Errorf(context.Background(), "Error during service cleanup: %%v", err)
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
`, moduleName, moduleName, moduleName, moduleName)
}

// StandaloneGinTemplate generates the http.go file for internal/server
func StandaloneGinTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/ecode"
	"github.com/ncobase/ncore/net/resp"

	"%s/handler"
	appConfig "%s/internal/config"
)

// ginServer creates and initializes the Gin engine
func ginServer(conf *config.Config, h handler.Interface) (*gin.Engine, error) {
	// Set gin mode
	if appConfig.IsProd(conf) {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create gin engine
	engine := gin.New()

	// Recovery middleware
	engine.Use(gin.Recovery())

	// Register API routes
	registerRest(engine, conf, h)

	// No route handler
	engine.NoRoute(func(c *gin.Context) {
		resp.Fail(c.Writer, resp.NotFound(ecode.Text(http.StatusNotFound)))
	})

	return engine, nil
}
`, moduleName, moduleName)
}

// StandaloneRestTemplate generates the rest.go file for internal/server
func StandaloneRestTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/config"

	"%s/handler"
	"%s/internal/version"
)

// registerRest registers the REST routes
func registerRest(e *gin.Engine, conf *config.Config, h handler.Interface) {
	// Root endpoint
	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Service is running",
			"name":    conf.AppName,
			"version": version.GetVersionInfo().Version,
		})
	})

	// Health check endpoint
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"name":    conf.AppName,
			"version": version.GetVersionInfo().Version,
		})
	})

	// Version endpoint
	e.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.GetVersionInfo())
	})

	// API v1 routes
	v1 := e.Group("/api/v1")
	{
		// Example endpoint
		v1.GET("/example", h.GetExample)

		// Add your API routes here
		// v1.GET("/users", h.GetUsers)
		// v1.POST("/users", h.CreateUser)
	}
}
`, moduleName, moduleName)
}

// StandaloneConfigTemplate generates the config.go file for internal/config
func StandaloneConfigTemplate(name, moduleName string) string {
	return `package config

import (
	"github.com/ncobase/ncore/config"
)

// IsProd returns true if the current environment is production
func IsProd(c *config.Config) bool {
	if c == nil {
		return false
	}
	return c.IsProd()
}

// IsDev returns true if the current environment is development
func IsDev(c *config.Config) bool {
	if c == nil {
		return true
	}
	return !c.IsProd()
}
`
}

// StandaloneHandlerProviderTemplate generates the provider.go file for handler layer
func StandaloneHandlerProviderTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package handler

import (
	"github.com/gin-gonic/gin"
	"%s/service"
)

// Interface defines the handler interface
type Interface interface {
	// Example endpoint
	GetExample(c *gin.Context)

	// Add your handler methods here
}

// handler implements the Interface
type handler struct {
	svc service.Interface
}

// New creates a new handler instance
func New(svc service.Interface) Interface {
	return &handler{
		svc: svc,
	}
}
`, moduleName)
}

// StandaloneHandlerTemplate generates the handler.go file (implementation)
func StandaloneHandlerTemplate(name, moduleName string) string {
	return `package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/net/resp"
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

// StandaloneHandlerTestTemplate generates the handler test file
func StandaloneHandlerTestTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"%s/data"
	"%s/data/repository"
	"%s/handler"
	"%s/service"
)

func TestGetExample(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Mock dependencies (you should use proper mocks)
	d := &data.Data{}
	repo := repository.New(d)
	svc := service.New(repo)
	h := handler.New(svc)

	// Create test router
	router := gin.New()
	router.GET("/api/v1/example", h.GetExample)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/example", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
`, moduleName, moduleName, moduleName, moduleName)
}

// StandaloneModelTemplate generates the model.go file
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
	return fmt.Sprintf(`package service

import (
	"context"

	"%s/data/model"
	"%s/data/repository"
)

// Interface defines the service interface
type Interface interface {
	// Example method
	GetExample(ctx context.Context) (*model.ExampleResponse, error)

	// Add your service methods here

	// Cleanup
	Close() error
}

// service implements the Interface
type service struct {
	repo repository.Interface
}

// New creates a new service instance
func New(repo repository.Interface) Interface {
	return &service{
		repo: repo,
	}
}
`, moduleName, moduleName)
}

// StandaloneServiceTemplate generates the service.go file (implementation)
func StandaloneServiceTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"%s/data/model"
)

// Close performs cleanup operations for the service
func (s *service) Close() error {
	// Add cleanup logic here if needed
	return nil
}

// GetExample is an example service method
func (s *service) GetExample(ctx context.Context) (*model.ExampleResponse, error) {
	// Example: Call repository
	// result, err := s.repo.FindExample(ctx, "some-id")
	// if err != nil {
	//     return nil, err
	// }

	// Mock example for demonstration
	example := &model.Example{
		ID:        uuid.New().String(),
		Name:      "Example Model",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &model.ExampleResponse{
		ID:   example.ID,
		Name: example.Name,
	}, nil
}
`, moduleName)
}

// StandaloneServiceTestTemplate generates the service test file
func StandaloneServiceTestTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"%s/data"
	"%s/data/repository"
	"%s/service"
)

func TestGetExample(t *testing.T) {
	// Setup
	d := &data.Data{}
	repo := repository.New(d)
	svc := service.New(repo)
	defer svc.Close()

	// Execute
	result, err := svc.GetExample(context.Background())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Name)
}
`, moduleName, moduleName, moduleName)
}

// StandaloneRepositoryProviderTemplate generates the provider.go file for repository layer
func StandaloneRepositoryProviderTemplate(name, moduleName string) string {
	return fmt.Sprintf(`package repository

import (
	"context"

	"%s/data"
	"%s/data/model"
)

// Interface defines the repository interface
type Interface interface {
	// Example method
	FindExample(ctx context.Context, id string) (*model.Example, error)

	// Add your repository methods here
}

// repository implements the Interface
type repository struct {
	d *data.Data
}

// New creates a new repository instance
func New(d *data.Data) Interface {
	return &repository{
		d: d,
	}
}
`, moduleName, moduleName)
}

// StandaloneRepositoryTemplate generates the repository.go file (implementation)
func StandaloneRepositoryTemplate(name, moduleName string, useMongo, useEnt, useGorm bool) string {
	imports := `import (
	"context"
	"fmt"
`

	if useEnt {
		imports += fmt.Sprintf(`	"%s/data/ent"
`, moduleName)
	}

	if useGorm {
		imports += `	"gorm.io/gorm"
`
	}

	if useMongo {
		imports += `	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
`
	}

	imports += fmt.Sprintf(`	"%s/data/model"
)
`, moduleName)

	implementation := `
// FindExample finds an example by ID
func (r *repository) FindExample(ctx context.Context, id string) (*model.Example, error) {
	// Access database via r.d
	//
	// For Ent: r.d.EC (master) or r.d.ECRead (read replica)
	// For GORM: r.d.GormClient (master) or r.d.GormRead (read replica)
	// For MongoDB: r.d.MC (master) or r.d.MCRead (read replica)

	// Example implementation would go here
	_ = id

	return nil, fmt.Errorf("not implemented")
}
`

	return fmt.Sprintf(`package repository

%s
%s`, imports, implementation)
}
