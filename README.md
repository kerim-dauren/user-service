# User Service

## Description

The User Service is a Go-based microservice that manages user accounts.
It provides functionalities for creating, retrieving, updating, and deleting user profiles.
This service uses a PostgreSQL database for persistent storage and leverages Prometheus for metrics collection.

## Features

- **User Management**: Create, read, update, and delete user accounts.
- **Password Hashing**: Uses Argon2 for secure password hashing.
- **API Documentation**: Swagger-generated API documentation.
- **Metrics**: Exports service metrics via Prometheus.
- **Configuration**: Uses a configuration file for easy setup and customization.
- **Logging**: Structured logging with `slog`.

## Technologies Used

- **Go**: Programming language.
- **PostgreSQL**: Database.
- **Goose**: Database migration tool.
- **Swagger**: API documentation generator.
- **Prometheus**: Metrics monitoring.
- **Argon2**: Password hashing algorithm.
- **slog**: Structured logging.
- **Docker**: Containerization.
- **Makefile**: Task automation.

## Getting Started

### Prerequisites

- Go (version 1.24 or higher)
- Docker
- Docker Compose
- PostgreSQL

### Installation

1. **Clone the repository:**

   ```bash
   git clone <repository_url>
   cd user-service
   ```

2. **Set up the database:**

    - Ensure PostgreSQL is installed and running.
    - Update the database connection string in the `.env` file.

3. **Run database migrations:**

   ```bash
   make migration-up
   ```

### Configuration

The service is configured via a `.env` file. Example configurations:

```env
HTTP_PORT=8080
DB_URL="postgres://user:password@host:port/database"
LOG_LEVEL=info
LOG_HANDLER=text
LOG_WRITER=stdout
```

### Running the Service

1. **Build the project:**

   ```bash
   make project-build
   ```

2. **Run with Docker Compose:**

   ```bash
   make compose-run
   ```

3. **Access the API:**

   The service will be accessible at `http://localhost:8080/api/v1`.

### API Documentation

API documentation is generated using Swagger. To generate and view the documentation:

1. **Generate the Swagger documentation:**

   ```bash
   make godoc
   ```

2. **Access the Swagger UI:**

   Open your browser and navigate to `http://localhost:8080/swagger-ui/index.html`.

## Makefile Targets

- `make migration-up`: Run database migrations up.
- `make migration-down`: Run database migrations down.
- `make migration-create`: Create a new migration script.
- `make project-build`: Build the Go project.
- `make compose-run`: Run the service using Docker Compose.
- `make godoc`: Generate Swagger API documentation.
- `make go-test`: Run all unit tests in the project with verbose output and coverage report.
- `make go-proto-gen`: Generate Go code from .proto files

## Directory Structure

.
├── cmd # Main application entrypoint
├── internal
│ ├── api # HTTP API handlers and routing
│ ├── configs # Configuration loading and management
│ ├── domain # Domain models and interfaces
│ ├── services # Business logic and service implementations
│ └── storages # Data storage implementations (PostgreSQL)
├── pkg # Reusable packages
├── db
│ └── migration # Database migration scripts
├── .env # Environment configuration
├── Dockerfile # Dockerfile for building the service image
├── docker-compose.yml # Docker Compose configuration
├── go.mod # Go module definition
├── go.sum # Go module checksums
└── Makefile # Task automation