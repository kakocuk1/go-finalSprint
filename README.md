# Task Manager API

A RESTful API for task management written in Go.

The application allows users to create, update, delete, search, and complete tasks, supports recurring events, JWT authentication, and persistent data storage using SQLite.

> Originally developed as the final project for the Yandex Practicum Go Developer course and further improved as part of my backend portfolio.

---

## Features

- User registration and authentication
- JWT-based authorization
- CRUD operations for tasks
- Recurring task support
- Task search and filtering
- SQLite database
- REST API
- Docker support
- Automated testing with GitHub Actions

---

## Tech Stack

- Go
- Chi Router
- SQLite
- JWT
- Docker
- GitHub Actions
- net/http
- JSON

---

## Project Structure

```
.
├── pkg/
│   ├── api/          # HTTP handlers
│   ├── db/           # Database layer
│   └── server/       # HTTP server
├── scheduler/        # Task scheduler
├── tests/            # Integration tests
├── Dockerfile
├── docker-compose.yml
└── main.go
```

> The project is planned to be migrated to the `cmd/` + `internal/` layout as it evolves.

---

## Getting Started

### Run locally

```bash
git clone https://github.com/kakocuk1/go-finalSprint.git

cd go-finalSprint

go run .
```

---

### Run with Docker

```bash
docker compose up --build
```

The application will be available at:

```
http://localhost:7540
```

---

## Configuration

The application uses environment variables.

| Variable | Description | Default |
|----------|-------------|---------|
| `TODO_DBFILE` | SQLite database path | `scheduler.db` |
| `TODO_PORT` | HTTP server port | `7540` |
| `TODO_PASSWORD` | Admin password | not set |

---

## API Endpoints

### Authentication

| Method | Endpoint |
|---------|----------|
| POST | `/api/register` |
| POST | `/api/login` |

### Tasks

| Method | Endpoint |
|---------|----------|
| GET | `/api/tasks` |
| GET | `/api/task` |
| POST | `/api/task` |
| PUT | `/api/task` |
| DELETE | `/api/task` |
| POST | `/api/task/done` |

---

## Running Tests

Run all tests:

```bash
go test ./...
```

Verbose mode:

```bash
go test -v ./...
```

---

## Continuous Integration

GitHub Actions automatically runs on every push and pull request.

Pipeline includes:

- Build
- Unit tests

---

## Future Improvements

- PostgreSQL support
- Swagger / OpenAPI documentation
- Structured logging
- Graceful shutdown
- YAML configuration
- Migration to `cmd/` + `internal/`
- Refresh token support
- Integration tests with Testcontainers

---

## Author

**Vitaliy**

Backend Developer (Go)

GitHub: https://github.com/kakocuk1

---

## License

This project is published for educational and portfolio purposes.
