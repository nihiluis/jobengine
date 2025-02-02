# jobengine
Simple Go service that manages job state with a PostgreSQL backend.

## Requirements
- Go 1.22+
- PostgreSQL 12+
- sqlc (for code generation)
- Docker (optional)

## Getting started
Setting up env variables (e.g., .env)
```
ADDRESS=0.0.0.0:8080
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postgres
```

Running the app
```bash
# Build the application
go build main.go

# Run the application
./main
# or
go run main.go
```

With Docker
```bash
# Build the Docker image
docker build -t jobengine .

# Run the container
docker run -p 8080:8080 jobengine
```

## Code generation
```bash
sqlc generate
```

## OpenAPI
See `http://localhost:<port>/api/v1/docs` when running.

## Tests
```bash
go test ./...
```