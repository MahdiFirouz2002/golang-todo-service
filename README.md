# Task Manager Microservice

A small Go backend that manages to-do tasks over a REST API. Built with clean architecture and the Gin framework.

> **Current stage:** `feature/task-api` — RESTful task CRUD backed by PostgreSQL.

## Architecture

```
cmd/api/                         Application entrypoint (composition root)
internal/
  config/                        Environment-based configuration
  domain/                        Entities and repository ports (interfaces)
  usecase/task/                  Task business rules
  infrastructure/postgres/       PostgreSQL pool, migrations, repository
  delivery/http/                 Gin HTTP adapters (handlers, router, server)
```

Dependency rule: outer layers depend inward. Domain has no framework imports. Use cases depend on repository ports; infrastructure implements them; HTTP handlers stay thin.

```mermaid
flowchart TB
  Client[HTTP Client] --> Delivery[delivery/http]
  Delivery --> UseCase[usecase/task]
  UseCase --> Domain[domain]
  Postgres[infrastructure/postgres] --> Domain
```

## Prerequisites

- Go 1.24+
- PostgreSQL 15+ (or Docker Compose)

## Setup

### Option A: Docker Compose (recommended)

```bash
docker compose up --build
```

The API is available at `http://localhost:8080` and Swagger UI at `http://localhost:8080/swagger/index.html`.

### Option B: Local Go

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

### Task API

Create a task:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Write docs","description":"README examples","status":"todo","assignee":"mahdi"}'
```

List tasks:

```bash
curl http://localhost:8080/api/v1/tasks
```

Get, update, and delete:

```bash
curl http://localhost:8080/api/v1/tasks/{id}
curl -X PUT http://localhost:8080/api/v1/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated title","status":"in_progress"}'
curl -X DELETE http://localhost:8080/api/v1/tasks/{id}
```

#### Request / response formats

**Create (POST /api/v1/tasks)**

Request:

```json
{
  "title": "Write docs",
  "description": "optional",
  "status": "todo",
  "assignee": "optional"
}
```

Response `201 Created`:

```json
{
  "id": "uuid",
  "title": "Write docs",
  "description": "optional",
  "status": "todo",
  "assignee": "optional",
  "created_at": "2026-07-29T12:00:00Z",
  "updated_at": "2026-07-29T12:00:00Z"
}
```

**List (GET /api/v1/tasks)** — returns a JSON array of tasks.

**Update (PUT /api/v1/tasks/:id)** — send any fields to change; omitted fields stay unchanged.

**Delete (DELETE /api/v1/tasks/:id)** — returns `204 No Content` on success.

Valid `status` values: `todo`, `in_progress`, `done`.

## Tests

```bash
# Unit tests
make test

# Coverage report
make coverage

# Integration tests (requires running PostgreSQL)
DATABASE_URL=postgres://postgres:postgres@localhost:5432/todos?sslmode=disable make test-integration
```

The project targets **≥ 70% test coverage** using unit tests with mocked repositories and optional PostgreSQL integration tests.

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

## Design decisions

- **Clean architecture** keeps domain and use cases independent from Gin and PostgreSQL.
- **Repository port** in `domain` allows mocking in tests and swapping storage later.
- **Embedded migrations** run at startup to simplify local development before Docker Compose lands.
- **Use case layer** centralises validation (title required, valid status, UUID checks).

## Trade-offs

- Migrations run automatically at startup; no standalone migration CLI yet.
- Docker Compose and Swagger documentation are planned for upcoming branches.
- Module path matches the GitHub repository (`github.com/MahdiFirouz2002/golang-todo-service`).
