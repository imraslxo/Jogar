FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o futbic ./main.go

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/futbic .

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./futbic"]