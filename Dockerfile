FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM golang:1.23-alpine

WORKDIR /app

COPY --from=builder /main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["/app/main"]