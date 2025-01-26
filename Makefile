APP_NAME=songs
DOCKER_COMPOSE=docker-compose.yml

.PHONY: build run test clean docker-up docker-down migrate lint help

# Docker up and down
up:
	docker-compose -f $(DOCKER_COMPOSE) up -d

down:
	docker-compose -f $(DOCKER_COMPOSE) down

# Local running
run:
	go run ./cmd/main.go

# Local build
build:
	go build -o $(APP_NAME) ./cmd/main.go

# Run tests
test:
	go test -v ./...

# Run lint tests
lint:
	golangci-lint run

# Clean binary files
clean:
	rm -f $(APP_NAME)
	go clean

# Apply migrations
migrate:
	@echo "Applying database migrations..."
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/songs?sslmode=disable" up

# Rollback migrations
migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/songs?sslmode=disable" down
