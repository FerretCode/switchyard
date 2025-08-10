FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -ldflags "-w -s" -o main ./cmd/dashboard

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/main ./

COPY --from=builder /app/views ./views

ENTRYPOINT ["/app/main"]