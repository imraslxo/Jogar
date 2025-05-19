FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o futbic ./main.go

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/futbic .
COPY --from=builder /app/.env .
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./futbic"]