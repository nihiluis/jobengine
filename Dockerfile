FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .
COPY migrations ./migrations
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM golang:1.23-alpine

WORKDIR /app

COPY --from=builder /main .
COPY --from=builder /migrations ./migrations

EXPOSE 8080

CMD ["/app/main"]