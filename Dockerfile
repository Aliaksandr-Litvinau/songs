# Start from a small, secure base image
FROM golang:1.23-alpine AS builder

# Install swag and required build dependencies
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.3

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Generate Swagger documentation
RUN swag init -g internal/app/transport/handlers.go

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

# Create a minimal production image
FROM alpine:latest

# Install required runtime dependencies
RUN apk update && apk upgrade && apk add bash

# Reduce image size
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Create non-root user
RUN adduser -D appuser
USER appuser

# Set the working directory
WORKDIR /app

# Copy binary and docs from builder
COPY --from=builder /app/app .
COPY --from=builder /app/cmd/wait-for-it.sh .
COPY --from=builder /app/internal/app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Expose the port that the application listens on
EXPOSE 8080

# Run the binary when the container starts
CMD ["./app"]
