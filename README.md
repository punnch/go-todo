# ToDo (Study Project)

A microservices-based REST API for task management, written in Go.

## Architecture

```
Client
  └── api-service (:8080)   ← HTTP REST API, input validation
        └── db-service (:8010)   ← internal HTTP, business logic
              └── PostgreSQL (:5432)
```

Two services communicate over HTTP. `api-service` is the public-facing gateway; `db-service` owns the database layer.

| Service | Port | Responsibility |
|---|---|---|
| `api-service` | 8080 | Public REST API, request validation |
| `db-service` | 8010 | Business logic, PostgreSQL access |
| `postgres` | 5432 | Data persistence |

## Stack

- **Go** — `gorilla/mux`, `pgx/v5`, `go.uber.org/zap`
- **PostgreSQL 18** — primary storage
- **golang-migrate** — schema migrations
- **Docker / Docker Compose** — containerized deployment

## Prerequisites

- Docker + Docker Compose
- Go 1.21+
- Make

## Getting Started

**1. Configure environment**

Create a `.env` file in the project root:

```env
LOG_LEVEL=INFO

POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
```

**2. Start the database**

```bash
make env-up
```

**3. Run migrations**

```bash
make migrate-up
```

**4. Start services**

```bash
make services-up
```

The API is now available at `http://localhost:8080`.

## API Reference

All requests and responses use JSON.

### Create a task

```
POST /tasks
```

```json
{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread"
}
```

Response `201 Created`:

```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2024-11-01T12:00:00Z"
}
```

---

### Get tasks

```
GET /tasks
GET /tasks?id=1
GET /tasks?completed=false
GET /tasks?id=1&completed=true
```

Response `200 OK`: array of task objects.

---

### Get a task

```
GET /tasks/{id}
```

Response `200 OK` or `404 Not Found`.

---

### Update task completion

```
PATCH /tasks/{id}
```

```json
{
  "completed": true
}
```

Response `200 OK`: updated task object.

---

### Delete a task

```
DELETE /tasks/{id}
```

Response `204 No Content` or `404 Not Found`.

## Makefile Reference

### Environment

```bash
make env-up            # start postgres
make env-down          # stop postgres
make env-cleanup       # stop postgres and wipe data (destructive)
make env-port-forward  # expose postgres on localhost:5432
make env-port-close    # close the port forwarder
```

### Services

```bash
make services-up       # start api-service and db-service
make services-down     # stop both services
make services-rebuild  # rebuild images without cache
```

### Migrations

```bash
make migrate-up                    # apply all pending migrations
make migrate-down                  # roll back last migration
make migrate-create seq=<name>     # create a new migration pair
```

### Other

```bash
make ps   # show running containers
```

## Project Structure

```
go-todo/
├── api-service/
│   ├── cmd/main.go
│   ├── Dockerfile
│   └── internal/
│       ├── core/
│       │   ├── domain/        # Task domain model
│       │   ├── errors/
│       │   └── logger/
│       └── features/tasks/
│           ├── api/           # HTTP handlers, server, DTOs
│           └── client/        # db-service HTTP client
├── db-service/
│   ├── cmd/main.go
│   ├── Dockerfile
│   └── internal/
│       ├── core/
│       │   ├── domain/        # Task domain model
│       │   ├── errors/
│       │   └── logger/
│       └── features/tasks/
│           ├── repository/postgres/  # pgx queries
│           ├── service/             # business logic
│           └── transport/           # HTTP handlers, server
├── migrations/
│   └── 000001_init.up.sql
├── docker-compose.yaml
├── Makefile
└── .env
```

## Logs

Each service writes logs to its own directory, mounted from the host:

```
out/
├── logs/
│   ├── api-service/
│   └── db-service/
└── pgdata/
```
