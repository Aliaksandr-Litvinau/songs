# Songs Library API

[![Go Version](https://img.shields.io/github/go-mod/go-version/Aliaksandr-Litvinau/songs)](https://go.dev/)
[![License](https://img.shields.io/github/license/Aliaksandr-Litvinau/songs)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aliaksandr-Litvinau/songs)](https://goreportcard.com/report/github.com/Aliaksandr-Litvinau/songs)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)
[![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=flat&logo=Prometheus&logoColor=white)](https://prometheus.io/)

This project built using **Clean Architecture** principles with **Gin Gonic** (a Go web framework). The API allows you to manage a library of songs, including CRUD operations for songs and advanced lyrics management with pagination support.

## Features

- **Song Management**:
  - Create new songs with metadata and lyrics
  - Update existing songs
  - Delete songs
  - Get detailed song information
- **Advanced Queries**:
  - Filter songs by various parameters
  - Pagination support
  - Search functionality
- **Lyrics Management**:
  - Fetch lyrics with verse pagination
  - Format and structure lyrics
- **Monitoring**:
  - Prometheus metrics
  - Request tracking
  - Performance monitoring

## Technologies

- **Go (Golang)**: The core programming language
- **Gin Gonic**: Web framework for building APIs
- **PostgreSQL**: Primary database
- **Docker & Docker Compose**: Containerization and orchestration
- **Prometheus**: Metrics and monitoring
- **Swagger**: API documentation
- **Standard Go Logger**: Simple and efficient logging
- **Testify**: Testing framework
- **Golang-migrate**: Database migrations

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21 or higher (if running outside Docker)
- PostgreSQL 14 or higher (if running without Docker)

### Installation

1. **Clone the repository:**

```bash
git clone git@github.com:Aliaksandr-Litvinau/songs.git
cd songs
```

2. **Set up environment variables:**

Create a `.env` file in the root directory of the project, with the following contents:

```bash
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/music_library
```

3. **Run the application with Docker:**

This project uses Docker for containerization. The `docker-compose.yml` file provided will set up the database and the API.

```bash
docker-compose up --build
```

![Run project](images/make-run.png)

> After the containers are built and running, the API and Swagger documentation will be accessible.
Swagger Documentation:
Visit ```/swagger/index.html``` to explore and interact with the API.

4. **Database Migration:**

The project uses `golang-migrate` for database migrations. All migrations are versioned and stored in the database. This ensures:
- Consistent database schema across all environments
- Version control for database changes
- Safe and repeatable deployments
- Rollback capability

To apply migrations manually (if needed):
```bash
make migrate-up
```

To rollback migrations:
```bash
make migrate-down
```

### Project Architecture

This project follows the Clean Architecture principles, which provides several benefits:
- Clear separation of concerns
- Independence of frameworks
- Testability
- Maintainability
- Scalability

The project is structured into the following layers:

#### 🎯 Domain Layer (`internal/app/domain`)
- Contains business logic and entities
- Defines interfaces for repositories and use cases
- Independent of external frameworks and libraries

#### 💼 Use Case Layer (`internal/app/usecase`)
- Implements application-specific business rules
- Orchestrates the flow of data between layers
- Contains the core business logic

#### 🔌 Interface Adapters Layer
- **Repository** (`internal/app/repository`): Implements data storage interfaces
- **Transport** (`internal/app/transport`): Handles HTTP requests and responses
- **Middleware** (`internal/app/middleware`): Contains HTTP middleware components

#### 🔧 Infrastructure Layer
- Database configurations
- External services integration
- Framework-specific code

### Directory Structure
```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   └── app/
│       ├── domain/      # Business entities and interfaces
│       ├── usecase/     # Business logic implementation
│       ├── repository/  # Data access layer
│       ├── transport/   # HTTP handlers and routes
│       └── middleware/  # HTTP middleware
├── pkg/                 # Shared packages
└── docker/             # Docker configurations
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
