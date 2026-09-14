# Task API

REST API for task management built with Go.

## Tech Stack

- Go 1.27+
- Gin
- PostgreSQL 16
- pgx
- Docker
- Docker Compose
- golang-migrate

## Features

- Create tasks
- Get task by ID
- Get all tasks
- Pagination
- Update tasks
- Delete tasks
- PostgreSQL migrations
- Graceful server shutdown
- Dockerized application and database

## Project Structure

```text
task-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── model/
│   └── repository/
├── migrations/
│   ├── 001_create_tasks.up.sql
│   └── 001_create_tasks.down.sql
├── .dockerignore
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Running with Docker

Make sure Docker is installed and running.

Start the application:

```bash
docker compose up --build
```

This starts:

1. PostgreSQL
2. Database migrations
3. Go API

The API will be available at:

`http://localhost:8080`

To stop the application:

```bash
docker compose down
```

PostgreSQL data is stored in a Docker volume and persists between container restarts.

## API Endpoints

### Get tasks

```http
GET /api/tasks
```

Pagination:

```http
GET /api/tasks?page=1&limit=20
```

Response:

```json
{
  "items": [],
  "page": 1,
  "limit": 20,
  "total": 0
}
```

### Get task

```http
GET /api/tasks/:id
```

Example:

```http
GET /api/tasks/1
```

### Create task

```http
POST /api/tasks
Content-Type: application/json
```

Request:

```json
{
  "title": "Learn Go",
  "description": "Build a REST API"
}
```

### Update task

```http
PATCH /api/tasks/:id
Content-Type: application/json
```

Request:

```json
{
  "title": "Learn Go properly",
  "completed": true
}
```

### Delete task

```http
DELETE /api/tasks/:id
```

Example:

```http
DELETE /api/tasks/1
```

## Database

The application uses PostgreSQL.

Database migrations are located in:

`migrations/`

Migrations are applied automatically when the application is started with Docker Compose.

## Configuration

For local development outside Docker, create a `.env` file:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/task_api
```

`.env` is not committed to the repository.

## Architecture

The application follows a simple layered structure:

```text
HTTP Request
     ↓
  Handler
     ↓
 Repository
     ↓
 PostgreSQL
```

- **Handler** — HTTP request/response handling and validation.
- **Repository** — database operations.
- **Model** — application data structures.
- **Config** — application configuration.

## Development

Run the application locally:

```bash
go run ./cmd/api
```

Run with Docker:

```bash
docker compose up --build
```
