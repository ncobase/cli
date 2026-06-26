# Ncobase CLI

基于 [ncore](https://github.com/ncobase/ncore) 的 Go 应用和扩展模块脚手架命令。

[![Go 版本](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![许可证](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**[中文文档](README_zh-CN.md)** | **[English](README.md)**

## 安装

```bash
go install github.com/ncobase/cli/cmd/nco@latest
nco -v
```

请确保 `$(go env GOPATH)/bin` 或 `GOBIN` 已加入 `PATH`。

## 源码构建

```bash
git clone https://github.com/ncobase/cli.git
cd cli
make build
./bin/nco -v
```

## 命令

```bash
nco init <名称> [参数]
nco create [类型] <名称> [参数]
nco migrate <命令>
nco schema <命令>
nco -v
nco --version
nco -h
nco --help
```

`nco init` 创建独立 ncore 应用。`nco create` 为现有项目创建扩展模块。
迁移和 schema 命令依赖 Atlas。

## `nco init`

创建独立应用。

```bash
nco init myapp --db postgres --use-ent --use-redis --with-test
```

### 参数

| 参数 | 说明 |
| --- | --- |
| `--db` | 数据库驱动：`postgres`、`mysql`、`sqlite`、`mongodb` |
| `--use-ent` | SQL 数据库使用 Ent |
| `--use-gorm` | SQL 数据库使用 GORM |
| `--use-mongo` | 使用 MongoDB 驱动 |
| `--use-redis` | 启用 Redis 支持 |
| `--use-elastic` | 启用 Elasticsearch 支持 |
| `--use-opensearch` | 启用 OpenSearch 支持 |
| `--use-meilisearch` | 启用 Meilisearch 支持 |
| `--use-kafka` | 启用 Kafka 支持 |
| `--use-rabbitmq` | 启用 RabbitMQ 支持 |
| `--use-s3` | 启用 AWS S3 存储支持 |
| `--use-minio` | 启用 MinIO 存储支持 |
| `--use-aliyun` | 启用阿里云 OSS 存储支持 |
| `--with-grpc` | 生成 gRPC 服务装配 |
| `--with-tracing` | 生成 OpenTelemetry 追踪装配 |
| `--with-test` | 生成测试文件 |
| `-p, --path` | 输出路径 |
| `-m, --module` | Go 模块名 |
| `--dry-run` | 只输出生成计划，不写入文件 |
| `--output` | 输出格式：`text`、`json` |

### 示例

```bash
nco init api --db postgres --use-ent --with-test
nco init service --db postgres --use-ent --use-kafka --use-redis
nco init files --db postgres --use-ent --use-s3
nco init analytics --db mongodb --use-mongo
```

## `nco create`

创建扩展模块。

```bash
nco create core auth --use-ent --with-test
```

### 扩展类型

| 类型 | 目标路径 | 用途 |
| --- | --- | --- |
| `core` | `core/<name>` | 基础领域模块 |
| `business` | `business/<name>` | 应用业务模块 |
| `plugin` | `plugin/<name>` | 可选集成模块 |
| 自定义 | `<type>/<name>` | 项目自定义模块分组 |

### 参数

| 参数 | 说明 |
| --- | --- |
| `--db` | 数据库驱动：`postgres`、`mysql`、`sqlite`、`mongodb` |
| `--use-ent` | SQL 数据库使用 Ent |
| `--use-gorm` | SQL 数据库使用 GORM |
| `--use-mongo` | 使用 MongoDB 驱动 |
| `--with-test` | 生成测试文件 |
| `--with-cmd` | 生成可运行服务装配 |
| `-p, --path` | 输出路径 |
| `-m, --module` | Go 模块名 |
| `--group` | 可选领域分组 |
| `--dry-run` | 只输出生成计划，不写入文件 |
| `--output` | 输出格式：`text`、`json` |

### 示例

```bash
nco create core auth --use-ent --with-test
nco create business order --use-gorm --with-cmd
nco create plugin payment --use-mongo
nco create features notification --use-ent
nco create core audit --use-gorm --db sqlite --with-cmd --dry-run --output json
```

## 生成应用结构

```text
myapp/
├── cmd/myapp/
├── internal/
│   ├── config/
│   ├── middleware/
│   ├── server/
│   └── version/
├── handler/
├── service/
├── data/
│   ├── ent/
│   ├── model/
│   ├── repository/
│   └── schema/
├── tests/
├── config.yaml
├── go.mod
├── Makefile
└── README.md
```

生成项目遵循 ncore 服务常用的 handler、service、repository 和 data 分层。

## 配置

生成的 `config.yaml` 遵循 ncore 配置约定。

```yaml
app_name: myapp
environment: debug

server:
  protocol: http
  host: 127.0.0.1
  port: 8080

data:
  database:
    migrate: true
    strategy: random
    max_retry: 3
    master:
      driver: postgres
      source: postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable
      max_open_conn: 32
      max_life_time: 7200
      max_idle_conn: 8
      logging: false
```

Redis、搜索、消息队列、对象存储、gRPC 和追踪配置只在启用对应参数时生成。

## 生成项目命令

```bash
make build
make run
make test
make fmt
make lint
make clean
```

## 开发验证

```bash
go test ./...
go vet ./...
make build
```

## 许可证

MIT
