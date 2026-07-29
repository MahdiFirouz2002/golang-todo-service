# Task Manager Microservice

A small Go backend that manages to-do tasks over a REST API. Built with clean architecture and the Gin framework.

> **Current stage:** `feature/database` — PostgreSQL persistence, migrations, and repository adapter.

## Architecture

```
cmd/api/                         Application entrypoint (composition root)
internal/
  config/                        Environment-based configuration
  domain/                        Entities and repository ports (interfaces)
  infrastructure/postgres/       PostgreSQL pool, migrations, repository
  delivery/http/                 Gin HTTP adapters (handlers, router, server)
```

Dependency rule: outer layers depend inward. Domain has no framework or infrastructure imports. Infrastructure implements domain ports; HTTP delivery stays thin and delegates to use cases in later stages.

```mermaid
flowchart TB
  Client[HTTP Client] --> Delivery[delivery/http]
  Delivery --> Domain[domain]
  Delivery -.->|future| UseCase[usecase]
  UseCase -.-> Domain
  Postgres[infrastructure/postgres] --> Domain
```

## Prerequisites

- Go 1.24+
- PostgreSQL 15+

## Setup

```bash
cp .env.example .env
# Start PostgreSQL and create the todos database, then:
go mod download
```

## Run locally

```bash
go run ./cmd/api
```

The server listens on `:8080` by default (`HTTP_PORT`).

### Health checks

```bash
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

## Tests

```bash
go test ./...
```

## Roadmap

| Branch | Scope |
|--------|--------|
| `feature/bootstrap` | Skeleton, config, Gin server, health |
| `feature/database` | PostgreSQL persistence |
| `feature/task-api` | Task CRUD REST API |
| `feature/tests` | Coverage ≥ 70%, mocks |
| `feature/swagger` | OpenAPI / Swagger UI |
| `feature/docker` | Dockerfile + Compose |
| `feature/observability` | Prometheus + tracing |
| `feature/cache` | Redis cache-aside (optional) |
| `feature/filtering` | Pagination & filters (optional) |
| `feature/performance` | Load test / pprof (optional) |

## Design decisions (bootstrap)

- **Clean architecture** with an explicit `domain` package and repository ports so persistence and HTTP can evolve independently.
- **Gin** as the required HTTP framework, isolated under `delivery/http`.
- **Env-based config** with defaults suitable for local development; validated at startup.
- **Graceful shutdown** on `SIGINT` / `SIGTERM` so in-flight requests can finish.
- **JSON structured logging** via `log/slog` for production-friendly observability from day one.

## Trade-offs

- Migrations run automatically at startup from embedded SQL files; no separate migration CLI yet (keeps local setup simple for the assessment).
- Task CRUD HTTP endpoints land in `feature/task-api`; this stage wires persistence only.
- Module path matches the GitHub repository (`github.com/MahdiFirouz2002/golang-todo-service`).
