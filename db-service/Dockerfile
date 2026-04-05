FROM golang:1.26-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/db-service ./cmd/main.go

CMD ["/app/db-service"]
