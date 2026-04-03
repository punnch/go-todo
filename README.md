# go-todo

A microservices-based Todo REST API written in Go.

## Architecture

The system is split into two independent services:

- **api-service** — HTTP REST API. Handles all incoming requests, validates input, and delegates to `db-service` via an internal TCP client.
- **db-service** — Internal data service. Owns the PostgreSQL connection and all database logic. Not exposed publicly.

Both services run in Docker containers alongside a managed PostgreSQL instance. Migrations are handled by `golang-migrate`.

```
Client → api-service (HTTP :8030) → db-service (TCP :8010) → PostgreSQL
```

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Router | gorilla/mux |
| Database Driver | pgx/v5 |
| Migrations | golang-migrate |
| Database | PostgreSQL 18 |
| Runtime | Docker Compose |
| Logging | zap |

## Project Structure

```
go-todo/
├── cmd/
│   ├── api-service/        # API service entrypoint + Dockerfile
│   └── db-service/         # DB service entrypoint + Dockerfile
├── internal/
│   ├── core/
│   │   ├── domains/todo/   # Task domain model
│   │   ├── apperrors/      # Shared application errors
│   │   └── logger/         # Zap logger setup
│   └── services/
│       ├── api-service/    # HTTP handlers, router, DB client
│       └── db-service/     # Repository, service layer, TCP server, postgres pool
├── migrations/             # SQL migration files
├── docker-compose.yaml
├── Makefile
└── .env
```

## Setup

### 1. Configure environment

```bash
cp .env.example .env
```

Fill in the required values:

```env
LOG_LEVEL=info

POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_db

API_SERVICE_PORT=8030
DB_SERVICE_PORT=8010
```

### 2. Start the database

```bash
make env-up
```

### 3. Run migrations

```bash
make migrate-up
```

### 4. Start the services

```bash
make services-up
```

The API is now available at `http://localhost:8030`.

## API Reference

### Task model

```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2025-01-01T00:00:00Z"
}
```

### Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/tasks` | Create a task |
| `GET` | `/tasks` | Get all tasks (filterable) |
| `GET` | `/tasks/{id}` | Get a task by ID |
| `PATCH` | `/tasks/{id}` | Mark a task complete/incomplete |
| `DELETE` | `/tasks/{id}` | Delete a task |

---

#### `POST /tasks`

Create a new task.

**Request body:**
```json
{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2025-01-01T12:00:00Z"
}
```

---

#### `GET /tasks`

Get all tasks. Supports optional query parameters:

| Param | Type | Description |
|---|---|---|
| `id` | int | Filter by task ID |
| `completed` | bool | Filter by completion status |

**Examples:**
```
GET /tasks
GET /tasks?completed=false
GET /tasks?id=3
```

**Response:** `200 OK` — array of task objects

---

#### `GET /tasks/{id}`

Get a single task by ID.

**Response:** `200 OK` — task object, or `404 Not Found`

---

#### `PATCH /tasks/{id}`

Update the completion status of a task.

**Request body:**
```json
{
  "completed": true
}
```

**Response:** `200 OK` — updated task object

---

#### `DELETE /tasks/{id}`

Delete a task by ID.

**Response:** `204 No Content`, or `404 Not Found`

---

### Error responses

All errors return a JSON body with a timestamp:

```json
{
  "error": "task not found",
  "time": "2025-01-01T12:00:00Z"
}
```

## Makefile Reference

### Environment

| Command | Description |
|---|---|
| `make env-up` | Start the PostgreSQL container |
| `make env-down` | Stop the PostgreSQL container |
| `make env-cleanup` | Destroy PostgreSQL container and data volume |
| `make env-port-forward` | Expose PostgreSQL on `127.0.0.1:5432` (for local tools) |
| `make env-port-close` | Stop the port forwarder |

### Services

| Command | Description |
|---|---|
| `make services-up` | Start api-service and db-service |
| `make services-down` | Stop api-service and db-service |
| `make services-rebuild` | Rebuild service images without cache |

### Migrations

| Command | Description |
|---|---|
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back the last migration |
| `make migrate-create seq=<name>` | Create a new migration file pair |

### Other

| Command | Description |
|---|---|
| `make ps` | Show status of all Compose containers |

## Development

To connect a local SQL client or migration tool directly to PostgreSQL:

```bash
make env-port-forward
```

This forwards `127.0.0.1:5432` to the Postgres container via socat.

To create a new migration:

```bash
make migrate-create seq=add_priority_column
```

This generates `migrations/000002_add_priority_column.up.sql` and `.down.sql`.
