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

# 功能完整的微服务
nco init myapp \
  --db postgres \
  --use-ent \
  --use-redis \
  --use-kafka \
  --use-s3 \
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
  --with-test

# MongoDB 数据服务
nco init analytics --db mongodb --use-mongo --use-elastic
```

### `nco create` - 创建扩展

向现有项目添加模块。

```bash
nco create [类型] <名称> [标志]
```

**扩展类型：**

| 类型       | 用途       | 示例                        |
| ---------- | ---------- | --------------------------- |
| `core`     | 核心逻辑   | `nco create core auth`      |
| `business` | 业务功能   | `nco create business order` |
| `plugin`   | 可选功能   | `nco create plugin payment` |
| 自定义     | 自定义目录 | `nco create myext user`     |

**标志：** `--use-ent`、`--use-gorm`、`--use-mongo`、`--with-test`、`--with-cmd`

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

生成的 `config.yaml`：

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
    # 可选的从库配置（读副本）
    slaves: []

  # 可选的数据源（通过标志启用）
  redis:
    addr: localhost:6379

  elasticsearch:
    addresses: ["http://localhost:9200"]

logger:
  level: 4 # 1:fatal, 2:error, 3:warn, 4:info, 5:debug
  format: text # text, json
```

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
- [文档](https://github.com/ncobase/ncobase)
- [问题反馈](https://github.com/ncobase/cli/issues)

## 社区

- 微信群：[添加管理员]
- Discord：[链接]
- 论坛：[链接]

## 致谢

感谢所有为 Ncobase 做出贡献的开发者。

特别感谢：

- [Ent](https://entgo.io/) - 强大的 Go ORM 框架
- [GORM](https://gorm.io/) - 流行的 Go ORM 库
- [Gin](https://gin-gonic.com/) - 高性能 HTTP 框架
- [Cobra](https://cobra.dev/) - CLI 框架

---

用 ❤️ 构建，由 Ncobase 团队出品。
