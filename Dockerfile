FROM golang:1.25-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/exe cmd/goTodo/main.go

CMD ["/app/exe"]