# Ncobase CLI

基于 [ncore](https://github.com/ncobase/ncore) 框架的强大 Go 应用脚手架工具。

[![Go 版本](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![许可证](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**[中文文档](README_zh-CN.md)** | **[English](README.md)**

## 特性

- 🚀 **快速脚手架** - 数秒内生成生产级应用
- 🏗️ **清晰架构** - Handler → Service → Repository 分层架构
- 🔌 **数据库支持** - PostgreSQL、MySQL、SQLite、MongoDB、Redis、Elasticsearch 等
- 📦 **ORM 集成** - Ent、GORM 或 MongoDB 驱动
- 🧪 **测试就绪** - 单元测试、集成测试、E2E 测试模板
- 🔄 **版本注入** - 通过 Makefile 自动版本管理
- 🎯 **中间件** - CORS、日志、追踪、安全头、限流
- 📂 **扩展系统** - Core、Business 和 Plugin 模块
- 🔐 **安全** - JWT 认证、RBAC 模板
- 📡 **实时通信** - WebSocket 和通知支持
- 📤 **文件上传** - 存储抽象（本地、S3、MinIO、阿里云 OSS）
- 🔄 **gRPC 支持** - 内置 gRPC 服务器，支持健康检查和反射
- 📊 **分布式追踪** - OpenTelemetry 集成，提供完整的可观测性

## 安装

```bash
# 从源码安装
git clone https://github.com/ncobase/cli.git
cd cli && make build
sudo mv bin/nco /usr/local/bin/

# 快速构建
go build -o nco cmd/nco/main.go
```

## 快速开始

### 两个命令，两种用途

| 命令         | 用于             | 输出                                          |
| ------------ | ---------------- | --------------------------------------------- |
| `nco init`   | **新建独立应用** | 完整项目结构：cmd/、data/、handler/、service/ |
| `nco create` | **创建扩展模块** | core/、business/ 或 plugin/ 模块              |

### 创建第一个应用

```bash
# 基础应用
nco init myapp

# 功能完整的微服务（带 gRPC 和追踪）
nco init myapp \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-kafka \
  --use-s3 \
  --with-grpc \
  --with-tracing \
  --with-test

cd myapp && make run
```

访问 `http://localhost:8080` 查看运行的应用。

## 命令详解

### `nco init` - 初始化应用

创建具有完整结构的独立应用。

```bash
nco init <项目名> [标志]
```

**关键标志：**

| 类别          | 标志                                                                          | 说明                                    |
| ------------- | ----------------------------------------------------------------------------- | --------------------------------------- |
| **ORM**       | `--use-ent`<br>`--use-gorm`<br>`--use-mongo`                                  | Ent (SQL)、GORM (SQL) 或 MongoDB        |
| **数据库**    | `--db`                                                                        | postgres、mysql、sqlite、mongodb、neo4j |
| **缓存/搜索** | `--use-redis`<br>`--use-elastic`<br>`--use-opensearch`<br>`--use-meilisearch` | 缓存和搜索引擎                          |
| **消息队列**  | `--use-kafka`<br>`--use-rabbitmq`                                             | 消息系统                                |
| **存储**      | `--use-s3`<br>`--use-minio`<br>`--use-aliyun`                                 | 对象存储                                |
| **服务**      | `--with-grpc`<br>`--with-tracing`                                             | gRPC 服务器<br>OpenTelemetry 追踪       |
| **其他**      | `--with-test`<br>`-m, --module`                                               | 生成测试<br>自定义模块名                |

**示例：**

```bash
# PostgreSQL REST API
nco init blog --db postgres --use-ent --with-test

# 全栈微服务
nco init orders \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-kafka \
  --use-s3 \
  --with-grpc \
  --with-tracing \
  --with-test

# MongoDB 数据服务
nco init analytics --db mongodb --use-mongo --use-elastic
```

### `nco create` - 创建扩展

在现有的项目中添加模块。扩展遵循与独立应用相同的清晰架构模式。

```bash
nco create [类型] <名称> [标志]
```

**扩展类型：**

| 类型       | 用途            | 路径              | 示例                        |
| ---------- | --------------- | ----------------- | --------------------------- |
| `core`     | 基础业务逻辑    | `core/<name>`     | `nco create core auth`      |
| `business` | 应用特定功能    | `business/<name>` | `nco create business order` |
| `plugin`   | 可选/插件化功能 | `plugin/<name>`   | `nco create plugin payment` |
| 自定义     | 自定义目录名称  | `<dir>/<name>`    | `nco create myext user`     |

**标志：**

| 标志           | 说明                          | 默认值 |
| -------------- | ----------------------------- | ------ |
| `--use-ent`    | 使用 Ent ORM（SQL 数据库）    | false  |
| `--use-gorm`   | 使用 GORM（SQL 数据库）       | false  |
| `--use-mongo`  | 使用 MongoDB 驱动             | false  |
| `--with-test`  | 生成测试文件                  | false  |
| `--with-cmd`   | 生成 cmd/main.go 用于独立运行 | false  |
| `-p, --path`   | 输出路径（默认：当前目录）    | `.`    |
| `-m, --module` | Go 模块名称                   | auto   |
| `--group`      | 可选的域组名称                | -      |

**生成的结构（每个扩展）：**

```text
<type>/<name>/
├── handler/             # HTTP 处理器
│   ├── provider.go
│   └── <name>.go
├── service/             # 业务逻辑
│   ├── provider.go
│   └── <name>.go
├── data/                # 数据访问
│   ├── model/
│   ├── repository/
│   └── schema/          # 如果使用 --use-ent
└── tests/               # 如果使用 --with-test
```

**示例：**

```bash
# 使用 Ent 的核心认证模块
nco create core auth --use-ent --with-test

# 使用 GORM 的业务 CRM 模块
nco create business crm --use-gorm --with-cmd

# 使用 MongoDB 的支付插件
nco create plugin payment --use-mongo

# 在 'features' 目录中的自定义扩展
nco create features notification --use-ent
```

**注意：** 扩展可以无缝集成到现有的 ncobase 项目中，并可以使用 `--with-cmd` 独立开发/测试。

### 其他命令

```bash
nco version              # 显示版本
nco migrate <命令>       # 数据库迁移（需要 atlas）
nco schema <命令>        # Schema 管理（需要 atlas）
```

## 项目结构

```text
myapp/
├── cmd/myapp/           # 应用入口
│   └── main.go
├── internal/            # 私有代码
│   ├── config/          # 配置辅助
│   ├── middleware/      # HTTP 中间件（CORS、认证、日志）
│   ├── server/          # 服务器设置（HTTP、路由）
│   └── version/         # 版本信息
├── handler/             # HTTP 处理器（控制器）
├── service/             # 业务逻辑
├── data/                # 数据访问层
│   ├── repository/      # 仓储
│   ├── schema/          # 数据库模式（Ent）
│   └── model/           # 数据模型
├── tests/               # 测试文件（如果 --with-test）
├── config.yaml          # 配置
├── Makefile             # 构建命令
└── README.md            # 文档
```

## 配置

生成的 `config.yaml`（位于项目根目录）：

```yaml
# 应用名称
app_name: myapp
# 运行环境：production / development / debug
environment: debug

server:
  # 协议类型：http / https
  protocol: http
  # 运行域名
  domain: localhost
  # 应用运行地址
  host: 127.0.0.1
  # 应用运行端口
  port: 8080

# gRPC 服务器（使用 --with-grpc 生成）
grpc:
  enabled: true
  host: 127.0.0.1
  port: 9090

data:
  database:
    # 全局配置
    migrate: true # 自动运行迁移
    strategy: random # 负载均衡策略：round_robin / random
    max_retry: 3 # 连接重试次数
    # 主库配置
    master:
      driver: postgres
      source: postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable
      max_open_conn: 32
      max_life_time: 7200
      max_idle_conn: 8
      logging: false
    # 可选的从库配置（读副本）
    # slaves:
    #   - driver: postgres
    #     source: postgres://postgres:postgres@localhost:5433/myapp?sslmode=disable
    #     max_open_conn: 64
    #     max_idle_conn: 16
    #     width: 1  # 负载均衡权重

  # 可选的数据源（根据标志生成）
  redis:
    addr: localhost:6379
    password: ""
    db: 0
    read_timeout: 0.4s
    write_timeout: 0.6s
    dial_timeout: 1s

  search:
    elasticsearch:
      addresses:
        - http://localhost:9200
      username: ""
      password: ""

  kafka:
    brokers:
      - localhost:9092
    sasl:
      enable: false

auth:
  jwt:
    secret: "生产环境请修改此密钥"
    expire: 48 # 过期时间（小时）
  whitelist:
    - /health
    - /login
    - "*swagger*"

# OpenTelemetry 追踪（使用 --with-tracing 生成）
observes:
  tracer:
    endpoint: localhost:4317 # OTLP gRPC 端点
  # 可选：Sentry 错误追踪
  # sentry:
  #   endpoint: https://your-sentry-dsn@sentry.io/project-id

logger:
  level: 5 # 1:fatal, 2:error, 3:warn, 4:info, 5:debug
  format: text # text / json
  output: stdout # stdout / stderr / file
  output_file: ./logs/runtime.log

storage:
  provider: minio # filesystem / minio / aliyun-oss / aws-s3
  id: minioadmin
  secret: minioadmin
  bucket: myapp
  endpoint: http://localhost:9000
```

**注意**：配置结构遵循 [ncore](https://github.com/ncobase/ncore) 框架标准。所有字段名称和值都是生产就绪的。

## Makefile 命令

每个生成的项目都包含：

```bash
make build      # 构建（带版本注入）
make run        # 开发模式运行
make test       # 运行所有测试
make clean      # 清理构建产物
make lint       # 运行代码检查
make fmt        # 格式化代码
make help       # 显示所有命令
```

**版本注入：**

```bash
make build
./bin/myapp --version
# Version: v0.1.0-3-g1a2b3c4
# Branch:  main
# Built At: 2026-02-14T10:30:00Z
```

## 内置功能

### 中间件

- **CORS** - 可配置的跨域资源共享
- **Logger** - 带上下文的请求/响应日志
- **Trace** - OpenTelemetry 分布式追踪
- **Security Headers** - HSTS、CSP、X-Frame-Options 等
- **Client Info** - IP、User-Agent 提取
- **Auth** - JWT 认证中间件（模板）
- **Rate Limit** - 令牌桶限流（模板）

### 高级功能（模板）

- **Pagination** - 游标和偏移分页
- **Filtering** - 高级查询过滤器
- **WebSocket** - 实时通信与房间
- **Notifications** - 推送通知系统
- **File Upload** - 验证、缩略图、MD5 哈希
- **Storage** - 存储抽象（本地/S3/MinIO/阿里云）

## 常见使用场景

### REST API 服务器

```bash
nco init api \
  --db postgres \
  --use-ent \
  --use-redis \
  --with-test
```

### 微服务

```bash
nco init service \
  --db postgres \
  --use-ent \
  --use-kafka \
  --use-redis
```

### 文件服务

```bash
nco init files \
  --db postgres \
  --use-ent \
  --use-s3
```

### 实时应用

```bash
nco init chat \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-rabbitmq
```

### gRPC 微服务

```bash
nco init grpc-service \
  --db postgres \
  --use-ent \
  --use-redis \
  --with-grpc \
  --with-tracing \
  --with-test

# HTTP 服务在 :8080，gRPC 服务在 :9090
# OpenTelemetry 追踪导出到 localhost:4317
```

**特性：**

- ✅ gRPC 服务器，支持健康检查和反射
- ✅ HTTP 和 gRPC 并发运行
- ✅ 分布式追踪覆盖两种协议
- ✅ 服务注册就绪

### 可观测微服务

```bash
nco init observable-api \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-kafka \
  --with-tracing \
  --with-test

# 完整的可观测性堆栈：
# - OpenTelemetry 追踪（OTLP）
# - 带 trace ID 的结构化日志
# - HTTP 请求追踪
```

### 模块化应用（使用扩展）

```bash
# 1. 初始化基础应用
nco init myapp --db postgres --use-ent --use-redis

cd myapp

# 2. 添加核心认证模块
nco create core auth --use-ent --with-test

# 3. 添加业务模块
nco create business order --use-ent --with-test
nco create business inventory --use-ent

# 4. 添加可选插件
nco create plugin notification --use-ent
nco create plugin analytics --use-mongo

# 项目结构：
# myapp/
# ├── core/auth/           # 认证模块
# ├── business/order/      # 订单管理
# ├── business/inventory/  # 库存跟踪
# └── plugin/notification/ # 通知系统
```

## 数据库支持矩阵

| 数据库        | 标志                       | ORM  | 使用场景      |
| ------------- | -------------------------- | ---- | ------------- |
| PostgreSQL    | `--db postgres --use-ent`  | Ent  | 生产环境 SQL  |
| MySQL         | `--db mysql --use-gorm`    | GORM | 遗留系统      |
| SQLite        | `--db sqlite --use-gorm`   | GORM | 本地开发/测试 |
| MongoDB       | `--db mongodb --use-mongo` | 原生 | 文档存储      |
| Redis         | `--use-redis`              | 原生 | 缓存/队列     |
| Elasticsearch | `--use-elastic`            | 原生 | 搜索          |
| Neo4j         | `--db neo4j`               | 原生 | 图数据库      |

## gRPC 与可观测性

### gRPC 服务器 (`--with-grpc`)

启用后，生成生产就绪的 gRPC 服务器：

**特性：**

- ✅ **健康检查** - gRPC 健康检查协议
- ✅ **反射** - 服务器反射，便于调试
- ✅ **并发运行** - 与 HTTP 服务器同时运行
- ✅ **优雅关闭** - 与 HTTP 服务器协调关闭
- ✅ **拦截器** - 内置日志和追踪

**生成的结构：**

```text
internal/server/
├── server.go  # 初始化 HTTP 和 gRPC
├── http.go    # HTTP 服务器（Gin）
├── grpc.go    # gRPC 服务器包装器
└── rest.go    # REST 路由
```

**配置：**

```yaml
grpc:
  enabled: true
  host: 127.0.0.1
  port: 9090
```

**使用方法：**

```bash
# 启动服务器（HTTP 和 gRPC 都会启动）
make run

# HTTP: http://localhost:8080
# gRPC: localhost:9090

# 使用 grpcurl 测试
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check
```

### 分布式追踪 (`--with-tracing`)

OpenTelemetry 集成，提供完整的可观测性：

**特性：**

- ✅ **OTLP 导出** - gRPC 导出器到收集器
- ✅ **W3C Trace Context** - 标准传播
- ✅ **自动 Span** - HTTP 请求自动追踪
- ✅ **Trace ID** - 注入到日志和响应中
- ✅ **服务元数据** - 名称、版本、环境

**生成的文件：**

```text
internal/middleware/
├── trace.go   # OpenTelemetry 中间件
└── utils.go   # 辅助函数

cmd/myapp/main.go  # Tracer 初始化
```

**配置：**

```yaml
observes:
  tracer:
    endpoint: localhost:4317 # OTLP gRPC 端点
```

**自动追踪上下文：**

```go
// 每个请求自动：
// 1. 创建 span
// 2. 传播 trace context
// 3. 将 trace ID 添加到日志
// 4. 导出到 OTLP 收集器

// 响应头包含：
// X-Trace-Id: 1234567890abcdef
```

**与可观测性堆栈集成：**

```bash
# Jaeger（一体化）
docker run -d --name jaeger \
  -p 4317:4317 \
  -p 16686:16686 \
  jaegertracing/all-in-one:latest

# 访问 Jaeger UI
open http://localhost:16686

# 或使用 Tempo、Zipkin 等
```

## 开发流程

```bash
# 1. 创建项目
nco init myapp --db postgres --use-ent --with-test

# 2. 进入目录并配置
cd myapp
vim config.yaml  # 配置数据库

# 3. 生成 Schema（如果使用 Ent）
vim data/schema/user.go
go generate ./...

# 4. 实现功能
vim handler/handler.go
vim service/service.go
vim data/repository/repository.go

# 5. 测试和运行
make test
make run
```

## 故障排查

**`go mod tidy` 失败：**

```bash
# 使用 go.work 进行本地 ncore 开发
go work init
go work use /path/to/ncore/config
go work use /path/to/ncore/logging
# ... 或者直接继续 - go.mod 已经正确生成
```

**端口冲突：**
如果 8080 端口被占用，服务器会自动寻找可用端口。

**版本信息不显示：**
始终使用 `make build` 而不是 `go build` 来进行版本注入。

## 模板能力详解

### 1. 认证与授权

生成的项目包含 JWT 认证中间件模板：

```go
// middleware/auth.go.tmpl
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // JWT 验证逻辑
    }
}

func RequireRoles(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // RBAC 角色检查
    }
}
```

### 2. 分页与过滤

提供两种分页方式：

**游标分页：**

```go
// features/pagination.go.tmpl
type CursorPagination struct {
    Cursor    string
    Limit     int
    Direction string  // "next" 或 "prev"
}
```

**偏移分页：**

```go
type OffsetPagination struct {
    Page     int
    PageSize int
}
```

**高级过滤：**

```go
// features/filter.go.tmpl
// 支持操作符：eq, ne, gt, gte, lt, lte, in, like, between
type Filter struct {
    Field    string
    Operator string
    Value    interface{}
}
```

### 3. WebSocket 实时通信

完整的 WebSocket 支持：

```go
// features/websocket.go.tmpl
type WSHub struct {
    clients    map[*WSClient]bool
    rooms      map[string]map[*WSClient]bool
    broadcast  chan *WSBroadcast
}

// 使用示例
hub.BroadcastToRoom("chat-room-1", message, nil)
```

### 4. 通知系统

集成的推送通知：

```go
// features/notification.go.tmpl
type NotificationService struct {
    hub     *WSHub
    storage NotificationStorage
}

// 发送通知
notifService.SendToUser(ctx, userID, NotificationInfo, "标题", "消息", data)
```

### 5. 文件上传

完整的文件上传处理：

```go
// features/upload.go.tmpl
type UploadConfig struct {
    MaxFileSize       int64
    AllowedTypes      []string
    GenerateThumbnail bool
    ThumbnailSizes    []ThumbnailSize
}

// 支持：
// - 文件类型验证
// - 大小限制
// - 自动生成缩略图
// - MD5 校验
// - 单文件和多文件上传
```

### 6. 存储抽象

统一的存储接口：

```go
// features/storage.go.tmpl
type StorageProvider interface {
    Put(ctx, path, reader, size, contentType) error
    Get(ctx, path) (io.ReadCloser, error)
    Delete(ctx, path) error
    GetURL(ctx, path) (string, error)
    GetSignedURL(ctx, path, expiry) (string, error)
}

// 实现：
// - LocalStorageProvider（本地文件系统）
// - S3StorageProvider（AWS S3）
// 可扩展至 MinIO、阿里云 OSS 等
```

### 7. 测试模板

生成三种测试：

**Handler 测试：**

```go
// tests/handler_test.go.tmpl
func TestHandler_Create(t *testing.T) {
    tests := []struct {
        name       string
        input      interface{}
        wantStatus int
    }{
        // 测试用例
    }
}
```

**Service 测试：**

```go
// tests/service_test.go.tmpl
// 使用 mock repository
```

**集成测试：**

```go
// tests/integration_test.go.tmpl
// 使用 testify/suite
// 包含完整的 CRUD 工作流测试
```

## 高级用法

### 1. 主从数据库配置

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

生成的代码自动处理读写分离：

```go
// 写操作使用主库
d.GetMasterEntClient()

// 读操作使用从库（带自动降级）
d.GetSlaveEntClient()
```

### 2. 环境配置

支持多环境配置：

```bash
# 开发环境
export APP_ENV=debug
make run

# 生产环境
export APP_ENV=release
./bin/myapp
```

### 3. 自定义扩展目录

```bash
# 创建自定义目录结构
nco create features user-management --use-ent
nco create services notification --use-redis
nco create modules analytics --use-mongo
```

### 4. 模块组织

使用 `--group` 标志组织相关模块：

```bash
nco create business order --group ecommerce
nco create business product --group ecommerce
nco create business payment --group ecommerce
```

## 性能优化

### 1. 数据库连接池

生成的配置包含连接池优化：

```yaml
database:
  master:
    maxOpenConns: 25 # 最大打开连接数
    maxIdleConns: 10 # 最大空闲连接数
    connMaxLifetime: 3600 # 连接最大生命周期（秒）
```

### 2. Redis 缓存

如果启用 Redis：

```go
// 自动生成的缓存方法
func (r *Repository) GetCached(ctx context.Context, key string) (*Entity, error) {
    // 先查 Redis
    // 未命中则查数据库并缓存
}
```

### 3. 查询优化

生成的代码包含查询优化建议：

```go
// 预加载关联数据
client.User.Query().
    WithOrders().
    WithProfile().
    All(ctx)

// 选择特定字段
client.User.Query().
    Select(user.FieldName, user.FieldEmail).
    All(ctx)
```

## 安全最佳实践

### 1. 安全头

自动生成的安全头中间件：

```go
// middleware/security_headers.go.tmpl
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
w.Header().Set("Strict-Transport-Security", "max-age=31536000")
w.Header().Set("Content-Security-Policy", "default-src 'self'")
```

### 2. 输入验证

建议使用 validator：

```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### 3. SQL 注入防护

使用参数化查询：

```go
// Ent 自动防护
client.User.Query().Where(user.EmailEQ(email))

// GORM 自动防护
db.Where("email = ?", email).Find(&users)
```

## 生产部署

### 1. Docker 部署

创建 Dockerfile：

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

### 2. 环境变量

支持环境变量覆盖配置：

```bash
export DB_HOST=prod-db.example.com
export DB_PASSWORD=secure-password
export REDIS_ADDR=redis:6379
./bin/myapp
```

### 3. 健康检查

自动生成的健康检查端点：

```bash
curl http://localhost:8080/health
# {"status": "ok", "database": "connected", "redis": "connected"}
```

## 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献指南。

```bash
git clone https://github.com/ncobase/cli.git
cd cli && make build
./bin/nco init test-app --with-test
```

## 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 链接

- [Ncore 框架](https://github.com/ncobase/ncore)
- [问题反馈](https://github.com/ncobase/cli/issues)

## 致谢

感谢所有为 Ncobase 做出贡献的开发者。

特别感谢：

- [Ent](https://entgo.io/) - 强大的 Go ORM 框架
- [GORM](https://gorm.io/) - 流行的 Go ORM 库
- [Gin](https://gin-gonic.com/) - 高性能 HTTP 框架
- [Cobra](https://cobra.dev/) - CLI 框架

---

用 ❤️ 构建，由 Ncobase 团队出品。
