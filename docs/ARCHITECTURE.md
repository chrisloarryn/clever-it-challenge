# Beer API Documentation

## Architecture Overview

This project implements a **Hexagonal Architecture** (Ports & Adapters) pattern for a Beer API service, following SOLID principles and best practices.

### Architecture Components

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Layer                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │   Handlers  │  │ Middleware  │  │   Server    │       │
│  └─────────────┘  └─────────────┘  └─────────────┘       │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                   Application Core                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                 Use Cases                           │   │
│  │  ┌──────────────┐  ┌──────────────┐                │   │
│  │  │ Beer Service │  │ Price Calc   │                │   │
│  │  └──────────────┘  └──────────────┘                │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Domain                            │   │
│  │  ┌──────────────┐  ┌──────────────┐                │   │
│  │  │    Beer      │  │   Currency   │                │   │
│  │  └──────────────┘  └──────────────┘                │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │ PostgreSQL  │  │  In-Memory  │  │ Currency API│       │
│  └─────────────┘  └─────────────┘  └─────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── core/                   # Core business logic
│   │   ├── domain/            
│   │   │   └── beers/          # Beer domain entity
│   │   ├── ports/              # Interface definitions
│   │   │   ├── primary/        # Use case interfaces
│   │   │   └── secondary/      # Infrastructure interfaces
│   │   └── services/           # Business logic implementation
│   ├── adapters/               # External interface adapters
│   │   └── http/              # HTTP adapter
│   └── infrastructure/         # Infrastructure implementations
│       ├── config/            # Configuration
│       ├── logger/            # Logging
│       ├── storage/           # Repository implementations
│       ├── external/          # External service implementations
│       │   └── currencyLayer/ # Currency API integration
│       └── dependencies/      # Dependency injection
└── sql/
    └── init.sql               # Database initialization
```

## Key Features

### 1. Hexagonal Architecture
- **Domain-driven design** with rich domain models
- **Clear separation** between business logic and infrastructure
- **Testable** components with dependency injection
- **Flexible** infrastructure switching

### 2. SOLID Principles Implementation
- **Single Responsibility**: Each component has one clear purpose
- **Open/Closed**: Extensible through interfaces
- **Liskov Substitution**: Proper interface implementations
- **Interface Segregation**: Specific, focused interfaces
- **Dependency Inversion**: Depends on abstractions, not concretions

### 3. Database Abstraction
- **Repository pattern** for data access
- **Factory pattern** for easy database switching
- Support for **PostgreSQL** and **In-Memory** storage
- Easy to extend with new database types

### 4. Configuration Management
- **Environment-based** configuration
- **Type-safe** configuration with validation
- Support for different environments (dev, staging, prod)

### 5. Structured Logging
- **JSON/Text** format support
- **Contextual** logging with request IDs
- **Configurable** log levels

## API Endpoints

### Health Check
```
GET /ping
```

### Beers Management
```
GET    /api/v1/beers           # Get all beers
GET    /api/v1/beers/{id}      # Get beer by ID
POST   /api/v1/beers           # Create new beer
```

### Box Price Calculation
```
POST   /api/v1/beers/box-price # Calculate beer box price
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `ENVIRONMENT` | Environment (dev/staging/prod) | `development` | No |
| `LOG_LEVEL` | Logging level | `info` | No |
| `LOG_FORMAT` | Log format (json/text) | `json` | No |
| `DB_TYPE` | Database type (postgres/inmemory) | `inmemory` | No |
| `DB_HOST` | PostgreSQL host | `localhost` | No |
| `DB_PORT` | PostgreSQL port | `5432` | No |
| `DB_NAME` | Database name | `beers_db` | No |
| `DB_USER` | Database user | `postgres` | No |
| `DB_PASSWORD` | Database password | `password` | No |
| `CURRENCY_API_KEY` | CurrencyLayer API key | - | Yes* |

*Required when using currency conversion features

## Running the Application

### Development Mode (In-Memory)
```bash
make dev
```

### Production Mode with PostgreSQL
```bash
# Start PostgreSQL
make db-up

# Run application
make run
```

### Using Docker
```bash
# Build and run with PostgreSQL
make docker-run-postgres

# Or just run the app
make docker-build
make docker-run
```

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Test API endpoints (requires running server)
make api-test
```

## Development Tools

```bash
# Setup development environment
make setup

# Format code
make fmt

# Run linter
make lint

# Run static analysis
make vet
```

## Business Rules

### Beer Entity
- **ID**: Unique identifier (positive integer)
- **Name**: Required, non-empty string
- **Brewery**: Required, non-empty string  
- **Country**: Required, non-empty string
- **Price**: Positive decimal value
- **Currency**: ISO currency code (CLP, USD, EUR)

### Box Price Calculation
- **Quantity**: Number of beers in box (minimum 1)
- **Currency**: Target currency for price calculation
- **Tax**: Applied to total price (percentage)
- **Discount**: Applied to total price (percentage)

### Validation Rules
- Beer ID must be unique
- All string fields are trimmed and validated
- Price must be positive
- Currency must be valid ISO code
- Box quantity must be positive
- Tax and discount must be between 0-100%

## Error Handling

The API returns standardized error responses:

```json
{
  "error": "Error type",
  "message": "Detailed error message",
  "details": "Additional context (optional)"
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `409` - Conflict (duplicate ID)
- `500` - Internal Server Error

## Extending the Application

### Adding New Repository
1. Implement the `BeerRepository` interface
2. Add factory case in `storage_factory.go`
3. Update configuration as needed

### Adding New Use Case
1. Define interface in `ports/primary/`
2. Implement service in `core/services/`
3. Add HTTP handler in `adapters/http/handlers/`
4. Register in dependency container

### Adding New External Service
1. Define interface in `ports/secondary/`
2. Implement adapter in `infrastructure/`
3. Register in dependency container

## Performance Considerations

- **Connection pooling** for PostgreSQL
- **Thread-safe** in-memory repository
- **Graceful shutdown** handling
- **Context-aware** operations
- **Middleware** for common concerns (CORS, logging)

## Security Considerations

- **Input validation** at domain level
- **SQL injection** protection through proper queries
- **CORS** configuration for cross-origin requests
- **Environment variable** protection for secrets
- **Structured logging** without sensitive data exposure
