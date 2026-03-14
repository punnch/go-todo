FROM golang:1.25-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o /app/api-service ./cmd/api-service/main.go

CMD ["/app/api-service"]