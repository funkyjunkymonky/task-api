FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o task-api ./cmd/api

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/task-api .

CMD ["./task-api"]