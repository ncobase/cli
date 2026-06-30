# Ncobase CLI

Command-line scaffolding for Go applications and extension modules based on
[ncore](https://github.com/ncobase/ncore).

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**[中文文档](README_zh-CN.md)** | **[English](README.md)**

## Installation

```bash
go install github.com/ncobase/cli/cmd/nco@latest
nco -v
```

Ensure `$(go env GOPATH)/bin` or `GOBIN` is available in `PATH`.

## Build From Source

```bash
git clone https://github.com/ncobase/cli.git
cd cli
make build
./bin/nco -v
```

## Commands

```bash
nco init <name> [flags]
nco create [type] <name> [flags]
nco migrate <command>
nco schema <command>
nco -v
nco --version
nco -h
nco --help
```

`nco init` creates a standalone service or a modular product backend. `nco create`
adds an extension module to an existing project. Migration and schema commands require Atlas.

## `nco init`

Create an application. The default type is `service`; use `--type modular` for a
product backend with `core`, `biz`, `plugin`, and `internal/server` structure.

```bash
nco init myapp --db postgres --use-ent --use-redis --with-test
nco init product --type modular --db postgres --use-redis --use-meilisearch --with-test
```

### Flags

| Flag | Description |
| --- | --- |
| `-t, --type` | Init type: `service`, `modular` |
| `--db` | Database driver: `postgres`, `mysql`, `sqlite`, `mongodb` |
| `--use-ent` | Use Ent for SQL databases |
| `--use-gorm` | Use GORM for SQL databases |
| `--use-mongo` | Use MongoDB driver |
| `--use-redis` | Include Redis support |
| `--use-elastic` | Include Elasticsearch support |
| `--use-opensearch` | Include OpenSearch support |
| `--use-meilisearch` | Include Meilisearch support |
| `--use-kafka` | Include Kafka support |
| `--use-rabbitmq` | Include RabbitMQ support |
| `--use-s3` | Include AWS S3 storage support |
| `--use-minio` | Include MinIO storage support |
| `--use-aliyun` | Include Aliyun OSS storage support |
| `--with-grpc` | Generate gRPC server wiring |
| `--with-tracing` | Generate OpenTelemetry tracing wiring |
| `--with-test` | Generate test files |
| `-p, --path` | Output path |
| `-m, --module` | Go module name |
| `--dry-run` | Print the generation plan without writing files |
| `--output` | Output format: `text`, `json` |

### Examples

```bash
nco init api --db postgres --use-ent --with-test
nco init product --type modular --db postgres --use-redis --use-meilisearch --with-test
nco init service --db postgres --use-ent --use-kafka --use-redis
nco init files --db postgres --use-ent --use-s3
nco init analytics --db mongodb --use-mongo
```

## `nco create`

Create an extension module.

```bash
nco create core auth --use-ent --with-test
```

### Extension Types

| Type | Target path | Purpose |
| --- | --- | --- |
| `core` | `core/<name>` | Foundational domain modules |
| `biz` | `biz/<name>` | Product business modules |
| `business` | `business/<name>` | Application-specific modules |
| `plugin` | `plugin/<name>` | Optional integration modules |
| Custom | `<type>/<name>` | Project-defined module groups |

### Flags

| Flag | Description |
| --- | --- |
| `--db` | Database driver: `postgres`, `mysql`, `sqlite`, `mongodb` |
| `--use-ent` | Use Ent for SQL databases |
| `--use-gorm` | Use GORM for SQL databases |
| `--use-mongo` | Use MongoDB driver |
| `--with-test` | Generate test files |
| `--with-cmd` | Generate runnable service wiring |
| `-p, --path` | Output path |
| `-m, --module` | Go module name |
| `--group` | Optional domain group name |
| `--dry-run` | Print the generation plan without writing files |
| `--output` | Output format: `text`, `json` |

### Examples

```bash
nco create core auth --use-ent --with-test
nco create biz order --use-ent --with-test
nco create business order --use-gorm --with-cmd
nco create plugin payment --use-mongo
nco create features notification --use-ent
nco create core audit --use-gorm --db sqlite --with-cmd --dry-run --output json
```

## Generated Application Structure

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

Generated projects follow the handler, service, repository, and data layering
used by ncore-based services.

## Generated Modular Structure

```text
myapp/
├── cmd/myapp/
├── core/
├── biz/
├── plugin/
├── internal/
│   ├── middleware/
│   ├── server/
│   └── version/
├── migrations/
├── docs/
├── tests/
├── config.yaml
├── go.mod
├── Makefile
└── README.md
```

Modular applications keep process startup and cross-cutting HTTP composition in
`internal/server`. Domain modules should be added with `nco create` and should keep
their own `structs`, `data/schema`, `data/repository`, `service`, `handler`, and
`router` packages.

## Configuration

Generated `config.yaml` follows ncore configuration conventions.

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

Optional Redis, search, message queue, object storage, gRPC, and tracing sections
are generated only when the corresponding flags are enabled.

## Generated Project Commands

```bash
make build
make run
make test
make fmt
make lint
make clean
```

## Development

```bash
go test ./...
go vet ./...
make build
```

## License

MIT
