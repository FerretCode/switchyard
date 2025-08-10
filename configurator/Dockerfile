FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY configurator/go.mod configurator/go.sum ./

RUN go mod download

COPY configurator/ ./

COPY migrations/ ./migrations/

RUN go build -ldflags "-w -s" -o main ./cmd/configurator

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/main ./

COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["/app/main"]

