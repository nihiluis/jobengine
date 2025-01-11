# jobengine
## Requirements
- PostgreSQL

## OpenAPI
See `http://localhost:<port>/api/v1/docs` when running.

## How to
.env
```
ADDRESS=0.0.0.0:8080
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postgres
```

```bash
# Build the application
go build main.go

# Run the application
./main
# or
go run main.go
```

### With Docker
```bash
# Build the Docker image
docker build -t jobengine .

# Run the container
docker run -p 8080:8080 jobengine
```