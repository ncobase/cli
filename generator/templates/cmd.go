package templates

import "fmt"

// CmdMainTemplate generates the main.go file for the cmd directory
func CmdMainTemplate(d *Data) string {
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
	if d.UseS3Storage {
		imports += "\t_ \"github.com/ncobase/ncore/data/storage/s3\"\n"
	}
	if d.UseMinio {
		imports += "\t_ \"github.com/ncobase/ncore/data/storage/minio\"\n"
	}
	if d.UseAliyun {
		imports += "\t_ \"github.com/ncobase/ncore/data/storage/aliyun\"\n"
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
	"strings"
	"syscall"
	"time"

	"%s/internal/server"
	"%s/internal/version"

	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/logging/logger"
	"github.com/ncobase/ncore/logging/observes"

%s)

const (
	shutdownTimeout = 3 * time.Second // service shutdown timeout
)

func main() {
	flag.Parse()
	// handle version flags
	version.Flags()
	
	// load config
	conf := loadConfig()
	
	// set logger version
	logger.SetVersion(version.GetVersionInfo().Version)

	appName := strings.ToLower(conf.AppName)

	// init tracer
	initTracer(conf, appName)

	// init sentry
	initSentry(conf, appName)

	// initialize logger
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
		return fmt.Errorf("failed to create server: %%w", err)
	}
	defer s.Cleanup()

	// create listener
	listener, err := createListener(conf)
	if err != nil {
		return fmt.Errorf("failed to create listener: %%w", err)
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	// create server instance
	srvInstance := &http.Server{
		Addr:    fmt.Sprintf("%%s:%%d", conf.Host, conf.Port),
		Handler: s.Handler(),
	}

	// create error channel
	errChan := make(chan error, 1)

	// start server
	go func() {
		logger.Infof(context.Background(), "Listening and serving HTTP on: %%s", srvInstance.Addr)
		if err := srvInstance.Serve(listener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
				logger.Errorf(context.Background(), "Listen error: %%s", err)
			} else {
				logger.Infof(context.Background(), "Server closed")
			}
		}
	}()

	return gracefulShutdown(srvInstance, errChan)
}

// createListener creates network listener
func createListener(conf *config.Config) (net.Listener, error) {
	addr := fmt.Sprintf("%%s:%%d", conf.Host, conf.Port)
	if conf.Port == 0 {
		addr = fmt.Sprintf("%%s:0", conf.Host)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// If port is in use, try to find a random available port
		if strings.Contains(err.Error(), "address already in use") && conf.Port != 0 {
			// Try with port 0 to let system assign a random port
			originalPort := conf.Port
			tmpAddr := fmt.Sprintf("%%s:0", conf.Host)
			listener, err = net.Listen("tcp", tmpAddr)
			if err != nil {
				return nil, fmt.Errorf("error starting server with random port: %%w", err)
			}
			// Update config with the new port
			conf.Port = listener.Addr().(*net.TCPAddr).Port
			logger.Infof(context.Background(), "Port %%d was in use, switched to port %%d", originalPort, conf.Port)
			return listener, nil
		}
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

// initTracer initializes the tracer
func initTracer(conf *config.Config, appName string) {
	if conf.Observes != nil && conf.Observes.Tracer != nil && conf.Observes.Tracer.Endpoint != "" {
		err := observes.NewTracer(&observes.TracerOption{
			URL:                conf.Observes.Tracer.Endpoint,
			Name:               strings.ToLower(conf.AppName),
			Version:            version.Version,
			Branch:             version.Branch,
			Revision:           version.Revision,
			Environment:        conf.Environment,
			SamplingRate:       1.0,
			MaxAttributes:      100,
			BatchTimeout:       5 * time.Second,
			ExportTimeout:      30 * time.Second,
			MaxExportBatchSize: 512,
		})
		if err != nil {
			logger.Errorf(context.Background(), "tracer.Init: %%s", err)
		}
	}
}

// initSentry initializes the sentry
func initSentry(conf *config.Config, appName string) {
	if conf.Observes != nil && conf.Observes.Sentry != nil && conf.Observes.Sentry.Endpoint != "" {
		if err := observes.NewSentry(&observes.SentryOptions{
			Dsn:         conf.Observes.Sentry.Endpoint,
			Name:        appName,
			Release:     version.Version,
			Environment: conf.Environment,
		}); err != nil {
			logger.Errorf(context.Background(), "sentry.Init: %%s", err)
		}
	}
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
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Debugf(context.Background(), "Shutdown timed out after %%s", shutdownTimeout)
		} else {
			logger.Debugf(context.Background(), "Shutdown completed within %%s", shutdownTimeout)
		}

		return nil
	}
}
`, d.PackagePath, d.PackagePath, imports)
}
