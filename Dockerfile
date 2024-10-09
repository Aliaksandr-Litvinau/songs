FROM golang:1.23

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN swag init -g internal/api/handlers.go

RUN go build -o main ./cmd/server

CMD ["./main"]
