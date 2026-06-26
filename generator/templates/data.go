package templates

func DataTemplate(name, extType string) string {
	return `package data

import (
	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/data"
{{- if .DBDriver }}
	_ "github.com/ncobase/ncore/data/{{ .DBDriver }}"
{{- end }}
{{- if .UseRedis }}
	_ "github.com/ncobase/ncore/data/redis"
{{- end }}
{{- if .UseElastic }}
	_ "github.com/ncobase/ncore/data/elasticsearch"
{{- end }}
{{- if .UseOpenSearch }}
	_ "github.com/ncobase/ncore/data/opensearch"
{{- end }}
{{- if .UseMeili }}
	_ "github.com/ncobase/ncore/data/meilisearch"
{{- end }}
{{- if .UseKafka }}
	_ "github.com/ncobase/ncore/data/kafka"
{{- end }}
{{- if .UseRabbitMQ }}
	_ "github.com/ncobase/ncore/data/rabbitmq"
{{- end }}
)

// Data wraps ncore data resources for the generated extension.
type Data struct {
	*data.Data
}

// New creates the data layer.
func New(conf *config.Data, env ...string) (*Data, func(name ...string), error) {
	d, cleanup, err := data.New(conf)
	if err != nil {
		return nil, nil, err
	}
	return &Data{Data: d}, cleanup, nil
}

// Close closes all data resources.
func (d *Data) Close() (errs []error) {
	if d == nil || d.Data == nil {
		return nil
	}
	if baseErrs := d.Data.Close(); len(baseErrs) > 0 {
		errs = append(errs, baseErrs...)
	}
	return errs
}
`
}

func DataTemplateWithEnt(name, extType string) string {
	return `package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/data"
	"github.com/ncobase/ncore/logging/logger"
{{- if .DBDriver }}
	_ "github.com/ncobase/ncore/data/{{ .DBDriver }}"
{{- end }}
{{- if .UseRedis }}
	_ "github.com/ncobase/ncore/data/redis"
{{- end }}
{{- if .UseElastic }}
	_ "github.com/ncobase/ncore/data/elasticsearch"
{{- end }}
{{- if .UseOpenSearch }}
	_ "github.com/ncobase/ncore/data/opensearch"
{{- end }}
{{- if .UseMeili }}
	_ "github.com/ncobase/ncore/data/meilisearch"
{{- end }}
{{- if .UseKafka }}
	_ "github.com/ncobase/ncore/data/kafka"
{{- end }}
{{- if .UseRabbitMQ }}
	_ "github.com/ncobase/ncore/data/rabbitmq"
{{- end }}

	"{{ .PackagePath }}/data/ent"
	"{{ .PackagePath }}/data/ent/migrate"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
)

// Data wraps ncore and Ent data resources for the generated extension.
type Data struct {
	*data.Data
	EC     *ent.Client
	ECRead *ent.Client
}

// New creates the data layer.
func New(conf *config.Data, env ...string) (*Data, func(name ...string), error) {
	d, cleanup, err := data.New(conf)
	if err != nil {
		return nil, nil, err
	}

	masterDB := d.GetMasterDB()
	if masterDB == nil {
		return nil, cleanup, fmt.Errorf("master database connection is nil")
	}

	entClient, err := newEntClient(masterDB, conf.Database.Master, conf.Database.Migrate, env...)
	if err != nil {
		return nil, cleanup, fmt.Errorf("create master Ent client: %w", err)
	}

	entClientRead := entClient
	if readDB, err := d.GetSlaveDB(); err == nil && readDB != nil && readDB != masterDB {
		if readClient, err := newEntClient(readDB, conf.Database.Master, false, env...); err == nil {
			entClientRead = readClient
		} else {
			logger.Warnf(context.Background(), "Failed to create read Ent client, using master: %v", err)
		}
	}

	return &Data{
		Data:   d,
		EC:     entClient,
		ECRead: entClientRead,
	}, cleanup, nil
}

func newEntClient(db *sql.DB, conf *config.DBNode, enableMigrate bool, env ...string) (*ent.Client, error) {
	client := ent.NewClient(ent.Driver(dialect.DebugWithContext(
		entsql.OpenDB(conf.Driver, db),
		func(ctx context.Context, args ...any) {
			if conf.Logging {
				logger.Infof(ctx, "%v", args)
			}
		},
	)))

	if conf.Logging {
		client = client.Debug()
	}

	if enableMigrate {
		opts := []schema.MigrateOption{
			migrate.WithForeignKeys(false),
		}
		if len(env) == 0 || env[0] != "production" {
			opts = append(opts, migrate.WithDropIndex(true), migrate.WithDropColumn(true))
		}
		if err := client.Schema.Create(context.Background(), opts...); err != nil {
			return nil, fmt.Errorf("migrate Ent schema: %w", err)
		}
	}

	return client, nil
}

// GetMasterEntClient returns the write Ent client.
func (d *Data) GetMasterEntClient() *ent.Client {
	return d.EC
}

// GetSlaveEntClient returns the read Ent client.
func (d *Data) GetSlaveEntClient() *ent.Client {
	if d.ECRead != nil {
		return d.ECRead
	}
	return d.EC
}

// GetEntClientWithFallback returns a read or write Ent client.
func (d *Data) GetEntClientWithFallback(ctx context.Context, readOnly ...bool) *ent.Client {
	if len(readOnly) == 0 || !readOnly[0] {
		return d.GetMasterEntClient()
	}
	if d.ECRead != nil && d.ECRead != d.EC {
		return d.ECRead
	}
	if d.IsReadOnlyMode(ctx) {
		logger.Warnf(ctx, "System is in read-only mode, using master Ent client for reads")
	}
	return d.EC
}

// WithEntTx executes fn inside a write transaction.
func (d *Data) WithEntTx(ctx context.Context, fn func(ctx context.Context, tx *ent.Tx) error) error {
	client := d.GetEntClientWithFallback(ctx)
	if client == nil {
		return fmt.Errorf("Ent client is nil")
	}
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin Ent transaction: %w", err)
	}
	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback Ent transaction after %v: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// WithEntTxRead executes fn inside a read transaction.
func (d *Data) WithEntTxRead(ctx context.Context, fn func(ctx context.Context, tx *ent.Tx) error) error {
	client := d.GetEntClientWithFallback(ctx, true)
	if client == nil {
		return fmt.Errorf("Ent read client is nil")
	}
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin Ent read transaction: %w", err)
	}
	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback Ent read transaction after %v: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// Close closes all data resources.
func (d *Data) Close() (errs []error) {
	if d == nil {
		return nil
	}
	if d.EC != nil {
		if err := d.EC.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close master Ent client: %w", err))
		}
	}
	if d.ECRead != nil && d.ECRead != d.EC {
		if err := d.ECRead.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close read Ent client: %w", err))
		}
	}
	if d.Data != nil {
		if baseErrs := d.Data.Close(); len(baseErrs) > 0 {
			errs = append(errs, baseErrs...)
		}
	}
	return errs
}
`
}

func DataTemplateWithGorm(name, extType string) string {
	return `package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/data"
	"github.com/ncobase/ncore/logging/logger"
	"{{ .PackagePath }}/data/model"
{{- if .DBDriver }}
	_ "github.com/ncobase/ncore/data/{{ .DBDriver }}"
{{- end }}
{{- if .UseRedis }}
	_ "github.com/ncobase/ncore/data/redis"
{{- end }}
{{- if .UseElastic }}
	_ "github.com/ncobase/ncore/data/elasticsearch"
{{- end }}
{{- if .UseOpenSearch }}
	_ "github.com/ncobase/ncore/data/opensearch"
{{- end }}
{{- if .UseMeili }}
	_ "github.com/ncobase/ncore/data/meilisearch"
{{- end }}
{{- if .UseKafka }}
	_ "github.com/ncobase/ncore/data/kafka"
{{- end }}
{{- if .UseRabbitMQ }}
	_ "github.com/ncobase/ncore/data/rabbitmq"
{{- end }}

{{- if eq .DBDriver "mysql" }}
	"gorm.io/driver/mysql"
{{- end }}
{{- if eq .DBDriver "postgres" }}
	"gorm.io/driver/postgres"
{{- end }}
{{- if eq .DBDriver "sqlite" }}
	"gorm.io/driver/sqlite"
{{- end }}
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Data wraps ncore and GORM data resources for the generated extension.
type Data struct {
	*data.Data
	GormClient *gorm.DB
	GormRead   *gorm.DB
}

// New creates the data layer.
func New(conf *config.Data, env ...string) (*Data, func(name ...string), error) {
	d, cleanup, err := data.New(conf)
	if err != nil {
		return nil, nil, err
	}

	masterDB := d.GetMasterDB()
	if masterDB == nil {
		return nil, cleanup, fmt.Errorf("master database connection is nil")
	}

	gormClient, err := newGormClient(masterDB, conf.Database.Master, conf.Database.Migrate)
	if err != nil {
		return nil, cleanup, fmt.Errorf("create master GORM client: %w", err)
	}

	gormRead := gormClient
	if readDB, err := d.GetSlaveDB(); err == nil && readDB != nil && readDB != masterDB {
		if readClient, err := newGormClient(readDB, conf.Database.Master, false); err == nil {
			gormRead = readClient
		} else {
			logger.Warnf(context.Background(), "Failed to create read GORM client, using master: %v", err)
		}
	}

	return &Data{
		Data:       d,
		GormClient: gormClient,
		GormRead:   gormRead,
	}, cleanup, nil
}

func newGormClient(db *sql.DB, conf *config.DBNode, enableMigrate bool) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch conf.Driver {
	case "postgres", "postgresql":
{{- if eq .DBDriver "postgres" }}
		dialector = postgres.New(postgres.Config{Conn: db})
{{- else }}
		return nil, fmt.Errorf("postgres support was not generated")
{{- end }}
	case "mysql":
{{- if eq .DBDriver "mysql" }}
		dialector = mysql.New(mysql.Config{Conn: db})
{{- else }}
		return nil, fmt.Errorf("mysql support was not generated")
{{- end }}
	case "sqlite", "sqlite3":
{{- if eq .DBDriver "sqlite" }}
		dialector = sqlite.Open(conf.Source)
{{- else }}
		return nil, fmt.Errorf("sqlite support was not generated")
{{- end }}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", conf.Driver)
	}

	logLevel := gormlogger.Silent
	if conf.Logging {
		logLevel = gormlogger.Info
	}
	client, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("open GORM client: %w", err)
	}

	if enableMigrate {
		if err := client.AutoMigrate(&model.Item{}); err != nil {
			return nil, fmt.Errorf("migrate GORM schema: %w", err)
		}
	}

	return client, nil
}

// GetMasterGormClient returns the write GORM client.
func (d *Data) GetMasterGormClient() *gorm.DB {
	return d.GormClient
}

// GetSlaveGormClient returns the read GORM client.
func (d *Data) GetSlaveGormClient() *gorm.DB {
	if d.GormRead != nil {
		return d.GormRead
	}
	return d.GormClient
}

// GetGormClientWithFallback returns a read or write GORM client.
func (d *Data) GetGormClientWithFallback(ctx context.Context, readOnly ...bool) *gorm.DB {
	if len(readOnly) == 0 || !readOnly[0] {
		return d.GetMasterGormClient()
	}
	if d.GormRead != nil && d.GormRead != d.GormClient {
		return d.GormRead
	}
	if d.IsReadOnlyMode(ctx) {
		logger.Warnf(ctx, "System is in read-only mode, using master GORM client for reads")
	}
	return d.GormClient
}

// WithGormTx executes fn inside a write transaction.
func (d *Data) WithGormTx(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	client := d.GetGormClientWithFallback(ctx)
	if client == nil {
		return fmt.Errorf("GORM client is nil")
	}
	return client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

// Close closes all data resources.
func (d *Data) Close() (errs []error) {
	if d == nil || d.Data == nil {
		return nil
	}
	if baseErrs := d.Data.Close(); len(baseErrs) > 0 {
		errs = append(errs, baseErrs...)
	}
	return errs
}
`
}

func DataTemplateWithMongo(name, extType string) string {
	return `package data

import (
	"context"
	"fmt"

	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/data"
	"github.com/ncobase/ncore/logging/logger"
	_ "github.com/ncobase/ncore/data/mongodb"
{{- if .UseRedis }}
	_ "github.com/ncobase/ncore/data/redis"
{{- end }}
{{- if .UseElastic }}
	_ "github.com/ncobase/ncore/data/elasticsearch"
{{- end }}
{{- if .UseOpenSearch }}
	_ "github.com/ncobase/ncore/data/opensearch"
{{- end }}
{{- if .UseMeili }}
	_ "github.com/ncobase/ncore/data/meilisearch"
{{- end }}
{{- if .UseKafka }}
	_ "github.com/ncobase/ncore/data/kafka"
{{- end }}
{{- if .UseRabbitMQ }}
	_ "github.com/ncobase/ncore/data/rabbitmq"
{{- end }}

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Data wraps ncore MongoDB data resources for the generated extension.
type Data struct {
	*data.Data
}

// New creates the data layer.
func New(conf *config.Data, env ...string) (*Data, func(name ...string), error) {
	d, cleanup, err := data.New(conf)
	if err != nil {
		return nil, nil, err
	}

	if err := d.MongoHealthCheck(context.Background()); err != nil {
		if cleanup != nil {
			cleanup()
		}
		return nil, nil, fmt.Errorf("mongodb health check failed: %w", err)
	}

	logger.Debug(context.Background(), "MongoDB data layer initialized")
	return &Data{Data: d}, cleanup, nil
}

// Collection returns a typed MongoDB collection from the configured manager.
func (d *Data) Collection(database, collection string, readOnly bool) (*mongo.Collection, error) {
	coll, err := d.GetMongoCollection(database, collection, readOnly)
	if err != nil {
		return nil, err
	}
	typed, ok := coll.(*mongo.Collection)
	if !ok || typed == nil {
		return nil, fmt.Errorf("mongodb collection %s.%s has unexpected type %T", database, collection, coll)
	}
	return typed, nil
}

// Close closes all data resources.
func (d *Data) Close() (errs []error) {
	if d == nil || d.Data == nil {
		return nil
	}
	if baseErrs := d.Data.Close(); len(baseErrs) > 0 {
		errs = append(errs, baseErrs...)
	}
	return errs
}
`
}
