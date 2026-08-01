# Task Manager Microservice

A Go backend that manages to-do tasks over a REST API. Built with clean architecture, Gin, PostgreSQL, and optional Redis caching.

## Architecture

```
cmd/api/                         Application entrypoint (composition root)
internal/
  config/                        Environment-based configuration
  domain/                        Entities and repository ports
  usecase/task/                  Task business rules
  infrastructure/
    postgres/                    PostgreSQL pool, migrations, repository
    redis/                       Redis client
    cache/                       Cache-aside list decorator
  delivery/http/                 Gin handlers, router, server
  observability/                 Prometheus metrics and OpenTelemetry tracing
docs/                            OpenAPI spec and performance notes
```

```mermaid
flowchart TB
  Client[HTTP Client] --> Delivery[delivery/http]
  Delivery --> UseCase[usecase/task]
  UseCase --> Cache[infrastructure/cache]
  Cache --> Postgres[infrastructure/postgres]
  Delivery --> Metrics[observability/metrics]
  Delivery --> Tracing[observability/tracing]
```

## Prerequisites

- Go 1.24+
- Docker & Docker Compose (recommended)

## Quick start (Docker Compose)

```bash
docker compose up --build
```

| Endpoint | URL |
|----------|-----|
| API | http://localhost:8080/api/v1/tasks |
| Swagger UI | http://localhost:8080/swagger/index.html |
| OpenAPI spec | http://localhost:8080/openapi.yaml |
| Metrics | http://localhost:8080/metrics |
| pprof | http://localhost:8080/debug/pprof/ |

## Local development

```bash
cp .env.example .env
go mod download
go run ./cmd/api
```

## API examples

Create a task:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Write docs","status":"todo","assignee":"mahdi"}'
```

List with pagination and filters:

```bash
curl "http://localhost:8080/api/v1/tasks?page=1&page_size=10&status=todo&assignee=mahdi"
```

Update and delete:

```bash
curl -X PUT http://localhost:8080/api/v1/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","status":"in_progress"}'
curl -X DELETE http://localhost:8080/api/v1/tasks/{id}
```

### List response format

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "page_size": 20,
  "total_pages": 0
}
```

Valid `status` values: `todo`, `in_progress`, `done`.

## Tests

```bash
make test
make coverage
make bench
```

Integration tests (requires PostgreSQL):

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/todos?sslmode=disable make test-integration
```

## Performance

See [docs/PERFORMANCE.md](docs/PERFORMANCE.md) for benchmarks, pprof usage, and load testing.

```bash
make bench
make load-test
```

## Branch history

| Branch | Scope |
|--------|--------|
| `feature/bootstrap` | Skeleton, config, Gin server, health |
| `feature/database` | PostgreSQL persistence |
| `feature/task-api` | Task CRUD REST API |
| `feature/tests` | Coverage ≥ 70%, mocks |
| `feature/swagger` | OpenAPI / Swagger UI |
| `feature/docker` | Dockerfile + Compose |
| `feature/observability` | Prometheus + tracing |
| `feature/filtering` | Pagination & filters |
| `feature/cache` | Redis cache-aside |
| `feature/performance` | Benchmarks / pprof |

## Design decisions

- **Clean architecture** keeps domain and use cases independent from frameworks.
- **Repository port** enables mocking and cache decoration without changing business logic.
- **Embedded migrations** simplify local and container startup.
- **Cache-aside** wraps the repository to keep HTTP handlers unaware of Redis.
- **Observability** uses Prometheus counters/histograms/gauges and stdout OpenTelemetry traces for development.

## Trade-offs

- Migrations run at startup instead of a separate CLI.
- Redis caching is opt-in via `REDIS_ENABLED`.
- pprof is exposed for assessment; restrict in production deployments.
- Module path: `github.com/MahdiFirouz2002/golang-todo-service`.
