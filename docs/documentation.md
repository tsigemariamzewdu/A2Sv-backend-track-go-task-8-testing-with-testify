# Task Management System - Clean Architecture Implementation

## Project Structure Overview

This project follows Clean Architecture principles, separating concerns into distinct layers with clear dependencies. The structure includes comprehensive test coverage at each architectural layer.

### Core Layers
TASK_MANAGEMENT_WITH_TEST/
├── Domain/ # Enterprise business rules
│ └── domain.go # Core entities and interfaces
├── Repositories/ # Database operations
│ ├── task_repository.go
│ ├── task_repository_test.go
│ ├── user_repository.go
│ └── user_repository_test.go
├── usecases/ # Application business rules
│ ├── interface.go
│ ├── task_usecases.go
│ ├── task_usecases_test.go
│ ├── user_usecases.go
│ └── user_usecase_test.go
├── infrastructure/ # Framework & drivers
│ ├── auth_middleware.go
│ ├── auth_middleware_test.go
│ ├── jwt_service.go
│ ├── jwt-service_test.go
│ ├── password_service.go
│ └── password_service_test.go
├── Delivery/ # Interface adapters
│ ├── controllers/
│ │ └── task_controller.go
│ └── router/
│ └── router.go
├── db/ # Database connection
│ └── db.go
├── docs/ # Documentation
│ └── documentation.md
├── .env # Environment variables
├── .gitignore
├── go.mod # Go dependencies
└── main.go # Entry point


## Testing Strategy

The project implements a robust testing strategy with test files for each component:

### Infrastructure Layer Tests
- **Authentication**: `auth_middleware_test.go`
- **JWT Services**: `jwt-service_test.go`
- **Password Services**: `password_service_test.go`

### Repository Layer Tests
- **Task Repository**: `task_repository_test.go`
- **User Repository**: `user_repository_test.go`

### Use Case Layer Tests
- **Task Use Cases**: `task_usecases_test.go`
- **User Use Cases**: `user_usecase_test.go`

## Clean Architecture Implementation

### 1. Domain Layer
- Contains enterprise-wide business rules
- Defines core entities and interfaces
- No dependencies on other layers

### 2. Use Case Layer
- Implements application-specific business rules
- Depends only on Domain layer interfaces
- Contains all business logic

### 3. Interface Adapters
- `Delivery/`: Converts data between layers
- `Repositories/`: Implements persistence
- Both depend inward toward Use Cases

### 4. Frameworks & Drivers
- `infrastructure/`: External frameworks
- `db/`: Database connection
- Outer layer with most volatility

## Development Setup

1. Clone the repository
2. Set up environment variables in `.env`
3. Install dependencies: `go mod tidy`
4. Run tests: `go test ./...`

## Key Features

- Clear separation of concerns
- Independent testable components
- Database-agnostic design
- Authentication middleware
- Comprehensive test coverage

The architecture ensures that business rules remain independent of frameworks, databases, or external interfaces, making the core logic more maintainable and testable.