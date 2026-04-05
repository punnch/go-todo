FROM golang:1.26-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o /app/api-service ./cmd/main.go

CMD ["/app/api-service"]
