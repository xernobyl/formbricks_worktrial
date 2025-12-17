# Formbricks Hub

A backend service for collecting and managing experience data, built with Go and PostgreSQL.

## Features

- RESTful API for experience data CRUD operations
- PostgreSQL for data persistence
- Redis for caching and background jobs
- Clean architecture with repository, service, and handler layers
- Docker Compose for local development

## Tech Stack

- **Language**: Go 1.23
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Driver**: pgx/v5
- **HTTP**: Standard library `net/http`

## Project Structure

```
.
├── cmd/
│   ├── api/              # API server entrypoint
│   ├── worker/           # Background jobs (future)
│   └── migrate/          # Migration runner
├── internal/
│   ├── api/
│   │   ├── handlers/     # HTTP request handlers
│   │   └── middleware/   # HTTP middleware
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   └── models/           # Domain models
├── pkg/
│   └── database/         # Database utilities
├── migrations/           # SQL migrations
└── tests/               # Integration tests
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- Make (optional, for convenience)

### Quick Start

1. Clone the repository and navigate to the project directory

2. Set up the development environment:
```bash
make setup
```

This will:
- Start PostgreSQL and Redis containers
- Create a `.env` file from `.env.example`
- Run database migrations

3. Start the API server:
```bash
make run
```

The server will start on `http://localhost:8080`

### Manual Setup

If you prefer not to use Make:

1. Start Docker containers:
```bash
docker-compose up -d
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Run migrations:
```bash
go run ./cmd/migrate/main.go
```

4. Start the server:
```bash
go run ./cmd/api/main.go
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Experience Data

#### Create Experience
```bash
POST /v1/experiences
Content-Type: application/json

{
  "source_type": "survey",
  "source_id": "survey-123",
  "source_name": "Customer Feedback Survey",
  "field_id": "question-1",
  "field_label": "How satisfied are you?",
  "field_type": "rating",
  "value_number": 5,
  "metadata": {"campaign": "summer-2025"},
  "language": "en",
  "user_identifier": "user-abc"
}
```

#### Get Experience by ID
```bash
GET /v1/experiences/{id}
```

#### List Experiences
```bash
GET /v1/experiences?source_type=survey&limit=50&offset=0
```

Query parameters:
- `source_type` - Filter by source type
- `source_id` - Filter by source ID
- `field_id` - Filter by field ID
- `user_identifier` - Filter by user identifier
- `limit` - Number of results (default: 100, max: 1000)
- `offset` - Pagination offset

#### Update Experience
```bash
PATCH /v1/experiences/{id}
Content-Type: application/json

{
  "value_number": 4,
  "metadata": {"updated": true}
}
```

#### Delete Experience
```bash
DELETE /v1/experiences/{id}
```

## Development

### Available Make Commands

```bash
make help         # Show all available commands
make build        # Build all binaries
make run          # Run the API server
make test         # Run tests
make migrate      # Run database migrations
make docker-up    # Start Docker containers
make docker-down  # Stop Docker containers
make clean        # Clean build artifacts
```

### Running Tests

```bash
make test
```

### Database Migrations

Migrations are stored in the `migrations/` directory and are executed in alphabetical order.

To run migrations:
```bash
make migrate
```

Or manually:
```bash
go run ./cmd/migrate/main.go
```

## Environment Variables

See [.env.example](.env.example) for all available configuration options:

- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `PORT` - HTTP server port (default: 8080)
- `ENV` - Environment (development/production)

## Example Requests

### Create a text response
```bash
curl -X POST http://localhost:8080/v1/experiences \
  -H "Content-Type: application/json" \
  -d '{
    "source_type": "form",
    "field_id": "email",
    "field_type": "text",
    "value_text": "user@example.com"
  }'
```

### Get all experiences
```bash
curl http://localhost:8080/v1/experiences
```

### Update an experience
```bash
curl -X PATCH http://localhost:8080/v1/experiences/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {"verified": true}
  }'
```

### Delete an experience
```bash
curl -X DELETE http://localhost:8080/v1/experiences/{id}
```

## Architecture

The application follows a clean architecture pattern:

1. **Handlers** - Handle HTTP requests and responses
2. **Service** - Contain business logic and validation
3. **Repository** - Handle data access and queries
4. **Models** - Define domain entities and DTOs

This separation allows for:
- Easy testing and mocking
- Clear separation of concerns
- Simple maintenance and refactoring
