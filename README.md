# API Quest – Go Backend

## 1. Project Overview

The service is written in **Go** and follows **Clean Architecture**, keeping business logic completely decoupled from infrastructure and delivery concerns. Each layer communicates only through explicit interfaces defined in the domain package, so the HTTP framework, storage engine, or auth library can be swapped out without touching business rules.

| Concern | Choice |
|---|---|
| Language | Go 1.24 |
| HTTP framework | [Fiber v2](https://github.com/gofiber/fiber) |
| Authentication | JWT (HS256) via [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Storage | Thread-safe in-memory (`sync.RWMutex`) |
| IDs | UUIDs via [google/uuid](https://github.com/google/uuid) |

---

## 2. Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go          # Composition root – wires all layers, starts Fiber
├── internal/
│   ├── domain/              # Enterprise layer – entities & interface contracts
│   │   ├── auth.go          #   AuthUseCase interface
│   │   ├── book.go          #   Book entity, BookRepository & BookUseCase interfaces
│   │   └── errors.go        #   Sentinel errors (ErrNotFound, ErrUnauthorized, …)
│   ├── usecase/             # Application layer – pure business logic, no HTTP
│   │   ├── auth_usecase.go  #   JWT generation & validation
│   │   └── book_usecase.go  #   CRUD orchestration, input validation
│   ├── repository/
│   │   └── memory/          # Infrastructure layer – in-memory BookRepository
│   │       ├── book_repository.go
│   │       └── book_repository_test.go
│   ├── handler/             # Delivery layer – Fiber HTTP handlers
│   │   ├── ping_handler.go
│   │   ├── echo_handler.go
│   │   ├── auth_handler.go
│   │   └── book_handler.go
│   └── middleware/
│       └── auth.go          # JWT Bearer token middleware
├── Dockerfile               # Multi-stage build (builder → alpine)
├── docker-compose.yml
├── go.mod
└── go.sum
```

### Layer responsibilities

| Layer | Package | Role |
|---|---|---|
| **Domain** | `internal/domain` | Defines entities and interface contracts. Zero external dependencies. |
| **Use Case** | `internal/usecase` | Implements business rules. Depends only on domain interfaces. |
| **Repository** | `internal/repository/memory` | Satisfies `domain.BookRepository` with a mutex-guarded in-memory map. |
| **Handler** | `internal/handler` | Translates HTTP requests/responses. Calls use-case interfaces. |
| **Middleware** | `internal/middleware` | Cross-cutting concerns (auth). Fiber-specific, but isolated from business logic. |

---

## 3. Prerequisites

| Tool | Minimum version |
|---|---|
| Go | 1.24 |
| Docker | 24.x |
| Docker Compose | v2 (`docker compose`) |

---

## 4. Getting Started

### Run locally (without Docker)

```bash
go run ./cmd/api/main.go
```

The server starts on **http://localhost:8080**.

### Run with Docker Compose

```bash
docker compose up --build
```

Add `-d` to run in the background:

```bash
docker compose up --build -d
docker compose down   # stop and remove containers
```

---

## 5. API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/ping` | Public | Health-check – returns `{"message":"pong"}` |
| `POST` | `/echo` | Public | Echoes the JSON request body back verbatim |
| `POST` | `/auth/token` | Public | Issues a signed JWT for valid credentials |
| `POST` | `/books` | 🔒 Bearer | Create a new book |
| `GET` | `/books` | 🔒 Bearer | List books – supports `?author=`, `?page=`, `?limit=` |
| `GET` | `/books/:id` | 🔒 Bearer | Retrieve a single book by UUID |
| `PUT` | `/books/:id` | 🔒 Bearer | Replace all mutable fields of a book |
| `DELETE` | `/books/:id` | 🔒 Bearer | Delete a book (returns `204 No Content`) |

#### `GET /books` query parameters

| Parameter | Type | Default | Description |
|---|---|---|---|
| `author` | string | — | Filter by exact author name |
| `page` | int | `1` | Page number (1-based) |
| `limit` | int | `10` | Items per page |

Response shape:
```json
{
  "data":  [ /* Book objects */ ],
  "total": 42,
  "page":  1,
  "limit": 10
}
```

---

## 6. Authentication

### Obtain a token

Send any non-empty `username` and `password` to receive a 24-hour JWT:

```bash
curl -s -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "secret"}' | jq .
```

```json
{ "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…" }
```

### Call a protected endpoint

Pass the token as a `Bearer` credential in the `Authorization` header:

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…"

curl -s -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "The Go Programming Language", "author": "Donovan", "year": 2015}' | jq .
```

Requests to protected routes without a valid token return `401 Unauthorized`.

---

## 7. Running Tests

The repository layer ships with concurrency-safety tests that verify correct behaviour under parallel reads, writes, updates, and deletes.

```bash
# Standard run
go test ./internal/repository/memory/...

# With the built-in race detector (recommended)
go test -race ./internal/repository/memory/...
```

To run all tests in the module:

```bash
go test -race ./...
```

