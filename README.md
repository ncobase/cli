# Ncobase CLI

A powerful scaffolding tool for building Go applications with the [ncore](https://github.com/ncobase/ncore) framework.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**[中文文档](README_zh-CN.md)** | **[English](README.md)**

## Features

- 🚀 **Quick Scaffolding** - Production-ready applications in seconds
- 🏗️ **Clean Architecture** - Layered architecture with handler → service → repository
- 🔌 **Database Support** - PostgreSQL, MySQL, SQLite, MongoDB, Redis, Elasticsearch, etc.
- 📦 **ORM Integration** - Ent, GORM, or MongoDB driver
- 🧪 **Test Ready** - Unit, integration, and e2e test templates
- 🔄 **Version Injection** - Automatic versioning via Makefile
- 🎯 **Middleware** - CORS, Logger, Tracer, Security Headers, Rate Limiting
- 📂 **Extension System** - Core, Business, and Plugin modules
- 🔐 **Security** - JWT authentication, RBAC templates
- 📡 **Real-time** - WebSocket and notification support
- 📤 **File Upload** - Storage abstraction (Local, S3, MinIO, Aliyun OSS)

## Installation

```bash
# From source
git clone https://github.com/ncobase/cli.git
cd cli && make build
sudo mv bin/nco /usr/local/bin/

# Quick build
go build -o nco cmd/nco/main.go
```

## Quick Start

### Two Commands, Two Purposes

| Command      | Use For                        | Output                                                |
| ------------ | ------------------------------ | ----------------------------------------------------- |
| `nco init`   | **New standalone application** | Complete project with cmd/, data/, handler/, service/ |
| `nco create` | **Extension modules**          | core/, business/, or plugin/ module                   |

### Create Your First App

```bash
# Basic app
nco init myapp

# Full-featured microservice
nco init myapp \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-kafka \
  --use-s3 \
  --with-test

cd myapp && make run
```

Visit `http://localhost:8080` to see your app running.

## Commands

### `nco init` - Initialize Application

Create a standalone application with complete structure.

```bash
nco init <name> [flags]
```

**Key Flags:**

| Category          | Flags                                                                         | Description                             |
| ----------------- | ----------------------------------------------------------------------------- | --------------------------------------- |
| **ORM**           | `--use-ent`<br>`--use-gorm`<br>`--use-mongo`                                  | Ent (SQL), GORM (SQL), or MongoDB       |
| **Database**      | `--db`                                                                        | postgres, mysql, sqlite, mongodb, neo4j |
| **Cache/Search**  | `--use-redis`<br>`--use-elastic`<br>`--use-opensearch`<br>`--use-meilisearch` | Caching and search engines              |
| **Message Queue** | `--use-kafka`<br>`--use-rabbitmq`                                             | Messaging systems                       |
| **Storage**       | `--use-s3`<br>`--use-minio`<br>`--use-aliyun`                                 | Object storage                          |
| **Other**         | `--with-test`<br>`-m, --module`                                               | Generate tests<br>Custom module name    |

**Examples:**

```bash
# REST API with PostgreSQL
nco init blog --db postgres --use-ent --with-test

# Microservice with full stack
nco init orders \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-kafka \
  --use-s3 \
  --with-test

# Data service with MongoDB
nco init analytics --db mongodb --use-mongo --use-elastic
```

### `nco create` - Create Extensions

Add modules to existing projects.

```bash
nco create [type] <name> [flags]
```

**Extension Types:**

| Type       | Purpose               | Example                     |
| ---------- | --------------------- | --------------------------- |
| `core`     | Fundamental logic     | `nco create core auth`      |
| `business` | App-specific features | `nco create business order` |
| `plugin`   | Optional features     | `nco create plugin payment` |
| Custom     | Your directory        | `nco create myext user`     |

**Flags:** `--use-ent`, `--use-gorm`, `--use-mongo`, `--with-test`, `--with-cmd`

### Other Commands

```bash
nco version              # Show version
nco migrate <command>    # Database migrations (requires atlas)
nco schema <command>     # Schema management (requires atlas)
```

## Project Structure

```text
myapp/
├── cmd/myapp/           # Application entry point
│   └── main.go
├── internal/            # Private application code
│   ├── config/          # Configuration helpers
│   ├── middleware/      # HTTP middleware (CORS, auth, logging)
│   ├── server/          # Server setup (HTTP, routes)
│   └── version/         # Version info
├── handler/             # HTTP handlers (controllers)
├── service/             # Business logic
├── data/                # Data access layer
│   ├── repository/      # Repositories
│   ├── schema/          # Database schemas (Ent)
│   └── model/           # Data models
├── tests/               # Test files (if --with-test)
├── config.yaml          # Configuration
├── Makefile             # Build commands
└── README.md            # Documentation
```

## Configuration

Generated `config.yaml`:

```yaml
app_name: myapp
environment: debug # debug, release

server:
  host: 127.0.0.1
  port: 8080

data:
  database:
    master:
      driver: postgres
      source: postgres://user:pass@localhost/db?sslmode=disable
      maxOpenConns: 10
      maxIdleConns: 5
      logging: true
    # Optional slaves for read replicas
    slaves: []

  # Optional data sources (if enabled via flags)
  redis:
    addr: localhost:6379

  elasticsearch:
    addresses: ["http://localhost:9200"]

logger:
  level: 4 # 1:fatal, 2:error, 3:warn, 4:info, 5:debug
  format: text # text, json
```

## Makefile Commands

Every generated project includes:

```bash
make build      # Build with version injection
make run        # Run in dev mode
make test       # Run all tests
make clean      # Clean artifacts
make lint       # Run linters
make fmt        # Format code
make help       # Show all commands
```

**Version Injection:**

```bash
make build
./bin/myapp --version
# Version: v0.1.0-3-g1a2b3c4
# Branch:  main
# Built At: 2026-02-14T10:30:00Z
```

## Built-in Features

### Middleware

- **CORS** - Configurable cross-origin resource sharing
- **Logger** - Request/response logging with context
- **Trace** - OpenTelemetry distributed tracing
- **Security Headers** - HSTS, CSP, X-Frame-Options, etc.
- **Client Info** - IP, User-Agent extraction
- **Auth** - JWT authentication middleware (template)
- **Rate Limit** - Token bucket rate limiting (template)

### Advanced Features (Templates)

- **Pagination** - Cursor and offset-based pagination
- **Filtering** - Advanced query filters with operators
- **WebSocket** - Real-time communication with rooms
- **Notifications** - Push notification system
- **File Upload** - Validation, thumbnails, MD5 hashing
- **Storage** - Abstraction for Local/S3/MinIO/Aliyun

## Common Use Cases

### REST API Server

```bash
nco init api \
  --db postgres \
  --use-ent \
  --use-redis \
  --with-test
```

### Microservice

```bash
nco init service \
  --db postgres \
  --use-ent \
  --use-kafka \
  --use-redis
```

### File Service

```bash
nco init files \
  --db postgres \
  --use-ent \
  --use-s3
```

### Real-time App

```bash
nco init chat \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-rabbitmq
```

## Database Support Matrix

| Database      | Flag                       | ORM    | Use Case          |
| ------------- | -------------------------- | ------ | ----------------- |
| PostgreSQL    | `--db postgres --use-ent`  | Ent    | Production SQL    |
| MySQL         | `--db mysql --use-gorm`    | GORM   | Legacy systems    |
| SQLite        | `--db sqlite --use-gorm`   | GORM   | Local dev/testing |
| MongoDB       | `--db mongodb --use-mongo` | Native | Document store    |
| Redis         | `--use-redis`              | Native | Cache/Queue       |
| Elasticsearch | `--use-elastic`            | Native | Search            |
| Neo4j         | `--db neo4j`               | Native | Graph database    |

## Template Capabilities

### 1. Authentication & Authorization

Generated projects include JWT authentication middleware templates:

```go
// middleware/auth.go.tmpl
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // JWT verification logic
    }
}

func RequireRoles(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // RBAC role checking
    }
}
```

### 2. Pagination & Filtering

Two pagination strategies:

**Cursor-based:**

```go
// features/pagination.go.tmpl
type CursorPagination struct {
    Cursor    string
    Limit     int
    Direction string  // "next" or "prev"
}
```

**Offset-based:**

```go
type OffsetPagination struct {
    Page     int
    PageSize int
}
```

**Advanced Filtering:**

```go
// features/filter.go.tmpl
// Operators: eq, ne, gt, gte, lt, lte, in, like, between
type Filter struct {
    Field    string
    Operator string
    Value    interface{}
}
```

### 3. WebSocket Real-time

Complete WebSocket support:

```go
// features/websocket.go.tmpl
type WSHub struct {
    clients    map[*WSClient]bool
    rooms      map[string]map[*WSClient]bool
    broadcast  chan *WSBroadcast
}

// Usage
hub.BroadcastToRoom("chat-room-1", message, nil)
```

### 4. Notification System

Integrated push notifications:

```go
// features/notification.go.tmpl
type NotificationService struct {
    hub     *WSHub
    storage NotificationStorage
}

// Send notification
notifService.SendToUser(ctx, userID, NotificationInfo, "Title", "Message", data)
```

### 5. File Upload

Complete file upload handling:

```go
// features/upload.go.tmpl
type UploadConfig struct {
    MaxFileSize       int64
    AllowedTypes      []string
    GenerateThumbnail bool
    ThumbnailSizes    []ThumbnailSize
}

// Features:
// - File type validation
// - Size limits
// - Auto thumbnail generation
// - MD5 checksums
// - Single and multi-file upload
```

### 6. Storage Abstraction

Unified storage interface:

```go
// features/storage.go.tmpl
type StorageProvider interface {
    Put(ctx, path, reader, size, contentType) error
    Get(ctx, path) (io.ReadCloser, error)
    Delete(ctx, path) error
    GetURL(ctx, path) (string, error)
    GetSignedURL(ctx, path, expiry) (string, error)
}

// Implementations:
// - LocalStorageProvider (local filesystem)
// - S3StorageProvider (AWS S3)
// Extensible to MinIO, Aliyun OSS, etc.
```

### 7. Test Templates

Three types of tests generated:

**Handler Tests:**

```go
// tests/handler_test.go.tmpl
func TestHandler_Create(t *testing.T) {
    tests := []struct {
        name       string
        input      interface{}
        wantStatus int
    }{
        // Test cases
    }
}
```

**Service Tests:**

```go
// tests/service_test.go.tmpl
// Using mock repositories
```

**Integration Tests:**

```go
// tests/integration_test.go.tmpl
// Using testify/suite
// Full CRUD workflow tests
```

## Advanced Usage

### 1. Master-Slave Database

```yaml
data:
  database:
    master:
      driver: postgres
      source: postgres://user:pass@master:5432/db
    slaves:
      - driver: postgres
        source: postgres://user:pass@slave1:5432/db
      - driver: postgres
        source: postgres://user:pass@slave2:5432/db
```

Generated code handles read/write splitting automatically:

```go
// Write operations use master
d.GetMasterEntClient()

// Read operations use slave (with automatic fallback)
d.GetSlaveEntClient()
```

### 2. Environment Configuration

Multiple environment support:

```bash
# Development
export APP_ENV=debug
make run

# Production
export APP_ENV=release
./bin/myapp
```

### 3. Custom Extension Directories

```bash
# Create custom directory structures
nco create features user-management --use-ent
nco create services notification --use-redis
nco create modules analytics --use-mongo
```

### 4. Module Organization

Use `--group` flag to organize related modules:

```bash
nco create business order --group ecommerce
nco create business product --group ecommerce
nco create business payment --group ecommerce
```

## Performance Optimization

### 1. Database Connection Pool

Generated config includes pool optimization:

```yaml
database:
  master:
    maxOpenConns: 25 # Maximum open connections
    maxIdleConns: 10 # Maximum idle connections
    connMaxLifetime: 3600 # Connection max lifetime (seconds)
```

### 2. Redis Caching

If Redis is enabled:

```go
// Auto-generated caching methods
func (r *Repository) GetCached(ctx context.Context, key string) (*Entity, error) {
    // Try Redis first
    // Fallback to database and cache result
}
```

### 3. Query Optimization

Generated code includes optimization tips:

```go
// Eager loading
client.User.Query().
    WithOrders().
    WithProfile().
    All(ctx)

// Select specific fields
client.User.Query().
    Select(user.FieldName, user.FieldEmail).
    All(ctx)
```

## Security Best Practices

### 1. Security Headers

Auto-generated security header middleware:

```go
// middleware/security_headers.go.tmpl
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
w.Header().Set("Strict-Transport-Security", "max-age=31536000")
w.Header().Set("Content-Security-Policy", "default-src 'self'")
```

### 2. Input Validation

Recommended validator usage:

```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### 3. SQL Injection Protection

Parameterized queries:

```go
// Ent automatic protection
client.User.Query().Where(user.EmailEQ(email))

// GORM automatic protection
db.Where("email = ?", email).Find(&users)
```

## Production Deployment

### 1. Docker Deployment

Create Dockerfile:

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
COPY --from=builder /app/bin/myapp /myapp
COPY config.yaml /config.yaml
EXPOSE 8080
CMD ["/myapp"]
```

### 2. Environment Variables

Config override via environment variables:

```bash
export DB_HOST=prod-db.example.com
export DB_PASSWORD=secure-password
export REDIS_ADDR=redis:6379
./bin/myapp
```

### 3. Health Checks

Auto-generated health check endpoint:

```bash
curl http://localhost:8080/health
# {"status": "ok", "database": "connected", "redis": "connected"}
```

## Development Workflow

```bash
# 1. Create project
nco init myapp --db postgres --use-ent --with-test

# 2. Navigate and setup
cd myapp
vim config.yaml  # Configure database

# 3. Generate schemas (if using Ent)
vim data/schema/user.go
go generate ./...

# 4. Implement features
vim handler/handler.go
vim service/service.go
vim data/repository/repository.go

# 5. Test and run
make test
make run
```

## Troubleshooting

**`go mod tidy` fails:**

```bash
# Use go.work for local ncore development
go work init
go work use /path/to/ncore/config
go work use /path/to/ncore/logging
# ... or just continue - go.mod is already correct
```

**Port conflicts:**
The server automatically finds available ports if 8080 is busy.

**Version not showing:**
Always use `make build` instead of `go build` for version injection.

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```bash
git clone https://github.com/ncobase/cli.git
cd cli && make build
./bin/nco init test-app --with-test
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Links

- [Ncore Framework](https://github.com/ncobase/ncore)
- [Documentation](https://github.com/ncobase/ncobase)
- [Issues](https://github.com/ncobase/cli/issues)

## Acknowledgments

Thanks to all contributors to the Ncobase project.

Special thanks to:

MIT License - see [LICENSE](LICENSE) file for details.

## Links

- [Ncore Framework](https://github.com/ncobase/ncore)
- [Documentation](https://github.com/ncobase/ncobase)
- [Issues](https://github.com/ncobase/cli/issues)

## Acknowledgments

Thanks to all contributors to the Ncobase project.

Special thanks to:

- [Ent](https://entgo.io/) - Powerful Go ORM framework
- [GORM](https://gorm.io/) - Popular Go ORM library
- [Gin](https://gin-gonic.com/) - High-performance HTTP framework
- [Cobra](https://cobra.dev/) - CLI framework

---

Built with ❤️ by the Ncobase team.
