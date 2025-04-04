FROM golang:1.24

RUN apt-get update && apt-get install -y postgresql-client

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./main.go

ENV PORT=8080
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=sekretik123
ENV DB_NAME=futbikSecond

CMD ["./main"]