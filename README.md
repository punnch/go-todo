# go-todo

A todo REST API written in Go. Split into two services that talk to each other over HTTP — one handles the public API, the other owns the database.

## Architecture

```-
client → api-service :8080 → db-service :8010 → postgres :5432
```

**api-service** is the only thing exposed to the outside. It validates requests and forwards them to db-service. **db-service** is internal — it manages the postgres connection and runs all queries. They're separated so the database layer is never directly reachable from outside.

## Stack

- Go 1.25
- PostgreSQL 18
- [gorilla/mux](https://github.com/gorilla/mux) — routing
- [jackc/pgx](https://github.com/jackc/pgx) — postgres driver
- [migrate](https://github.com/golang-migrate/migrate) — database migrations
- Docker + Docker Compose

## Project structure

```-
.
├── cmd/
│   ├── api-service/
│   │   ├── main.go
│   │   └── Dockerfile
│   └── db-service/
│       ├── main.go
│       └── Dockerfile
├── internal/
│   ├── core/
│   │   ├── apperrors/        # shared error types
│   │   └── domains/todo/     # Task domain model
│   └── services/
│       ├── api-service/
│       │   ├── api/          # handlers, DTOs, validation
│       │   └── client/       # HTTP client for db-service
│       └── db-service/
│           ├── migrations/   # SQL migration files
│           ├── postgres/     # connection pool + repository
│           ├── server/       # HTTP handlers
│           └── todo/         # service + repository interface
├── docker-compose.yaml
└── .env.example
```

## Getting started

Copy the example env and fill in your values:

```bash
cp .env.example .env
```

```env
# Connection string used by db-service and the migrate container
DB_URL=postgres://db_user:db_password@postgres:5432/db_name?sslmode=disable

# Postgres credentials
DB_PASS=db_password
DB_USER=db_user
DB_NAME=db_name

# Where api-service reaches db-service (don't use 8080, that's api-service)
DB_SERVICE_URL=http://db-service:8010

# Used for running migrations locally (outside Docker)
LOCAL_DB_URL=postgres://postgres:yourpassword@localhost:5432/go_todo?sslmode=disable

# Used for step migrations (see below)
step=1
```

Then start everything:

```bash
make deploy
```

This builds both services, starts postgres, runs migrations, and brings up the API. To tear it down:

```bash
make undeploy
```

## API

All requests go to `api-service` on port `8080`. Errors come back as JSON with a `message` and `time` field.

### Create a task

```-
POST /tasks
Content-Type: application/json

{
  "title": "buy milk",
  "description": "2% please"
}
```

Both fields are required. Title must be unique (max 50 chars), description max 200 chars.

### Get all tasks

```-
GET /tasks
GET /tasks?completed=true
GET /tasks?id=1
```

### Get a task

```-
GET /tasks/:id
```

### Complete (or uncomplete) a task

```-
PATCH /tasks/:id
Content-Type: application/json

{
  "completed": true
}
```

### Delete a task

```-
DELETE /tasks/:id
```

## Migrations

Migrations live in `internal/services/db-service/migrations/` and run automatically on `make deploy` via the `migrate` container.

To run them manually against a local database:

```bash
make migrate-up                # apply all pending migrations
make migrate-down              # roll back all migrations

make migrate-up-step step=1    # apply N migrations
make migrate-down-step step=1  # roll back N migrations
```

## Running locally without Docker

Make sure postgres is running locally and `LOCAL_DB_URL` in your `.env` points to it, then:

```bash
go run ./cmd/db-service/main.go
go run ./cmd/api-service/main.go
```

db-service needs to be up before api-service, same as in Docker.
