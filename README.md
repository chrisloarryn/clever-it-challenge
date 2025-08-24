# Beer API - Hexagonal Architecture

A robust, production-ready Beer API service built with Go, implementing Clean Architecture principles and industry best practices.

[![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-green.svg)](https://alistair.cockburn.us/hexagonal-architecture/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 🚀 Quick Start

### Prerequisites
- Go 1.24.5 or higher
- Docker & Docker Compose (optional)
- Make (optional, for convenience commands)

### Run Locally (In-Memory Database)
```bash
# Clone the repository
git clone <repository-url>
cd clever-it-challenge

# Install dependencies
go mod tidy

# Run in development mode
make dev
# or
export DB_TYPE=inmemory && go run cmd/main.go
```

### Run with PostgreSQL
```bash
# Start PostgreSQL with Docker
make db-up

# Run the application
export DB_TYPE=postgres && make run
```

### Run with Docker Compose
```bash
# Build and run everything
make docker-run-postgres
```

The API will be available at `http://localhost:8080`

## 🏗️ Architecture

This project implements **Hexagonal Architecture** (Ports & Adapters) with the following benefits:

- **Clean separation** between business logic and infrastructure
- **Easy testing** with dependency injection
- **Database flexibility** - switch between PostgreSQL and In-Memory
- **SOLID principles** implementation
- **Domain-driven design** with rich domain models

### Architecture Diagram

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

### Project Structure
```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── core/                   # Core business logic
│   │   ├── domain/             # Domain entities & business rules
│   │   ├── ports/              # Interface definitions
│   │   │   ├── primary/        # Use case interfaces
│   │   │   └── secondary/      # Infrastructure interfaces
│   │   └── services/           # Business logic implementation
│   ├── adapters/               # External interface adapters
│   │   └── http/              # HTTP adapter (REST API)
│   └── infrastructure/         # Infrastructure implementations
│       ├── config/            # Configuration management
│       ├── logger/            # Structured logging
│       ├── storage/           # Repository implementations
│       ├── external/          # External service implementations
│       └── dependencies/      # Dependency injection container
├── docs/                       # Documentation
├── scripts/                    # Utility scripts
└── sql/                        # Database migration scripts
```

## � API Documentation

### Health Check
```bash
curl http://localhost:8080/ping
```

### Beer Operations
```bash
# Get all beers
curl http://localhost:8080/api/v1/beers

# Get beer by ID
curl http://localhost:8080/api/v1/beers/1

# Create a new beer
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 100,
    "name": "IPA Craft",
    "brewery": "Local Brewery",
    "country": "USA",
    "price": 25.99,
    "currency": "USD"
  }'

# Calculate box price with currency conversion
curl "http://localhost:8080/api/v1/beers/1/boxprice?quantity=6&currency=EUR"
```

### Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/ping` | Health check |
| `GET` | `/api/v1/beers` | Get all beers |
| `GET` | `/api/v1/beers/{id}` | Get beer by ID |
| `POST` | `/api/v1/beers` | Create new beer |
| `GET` | `/api/v1/beers/{id}/boxprice` | Calculate box price |

Legacy routes are also supported for backward compatibility:
- `/beers` (same functionality as `/api/v1/beers`)

## 🛠️ Development

### Available Make Commands
```bash
make help           # Show all available commands
make build          # Build the application
make test           # Run unit tests
make test-coverage  # Run tests with coverage
make lint           # Run linter
make fmt            # Format code
make setup          # Setup development environment
```

### Testing

#### Unit Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

#### API Testing

**Option 1: REST Client (VS Code) - Recommended ⭐**
Use the properly formatted REST files with the REST Client extension:

1. Install the "REST Client" extension in VS Code
2. Open `test.REST` (complete test suite with variables)
3. Click "Send Request" on any endpoint

See detailed guide: [docs/REST_CLIENT_GUIDE.md](docs/REST_CLIENT_GUIDE.md)

**Option 2: Automated Script**
```bash
# Start the server first
make dev

# Run comprehensive API tests (in another terminal)
make api-test-full
# or
./scripts/test-api.sh
```

**Option 3: Manual cURL Commands**
```bash
# Quick test
make api-test

# See all cURL examples
make api-test-manual
# or check docs/API_TESTING.md
```

## 🔧 Configuration

Configure the application using environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `ENVIRONMENT` | Environment (dev/staging/prod) | `development` | No |
| `LOG_LEVEL` | Logging level | `info` | No |
| `LOG_FORMAT` | Log format (json/text) | `json` | No |
| `DB_TYPE` | Database type (`postgres`/`inmemory`) | `inmemory` | No |
| `DB_HOST` | PostgreSQL host | `localhost` | No |
| `DB_PORT` | PostgreSQL port | `5432` | No |
| `DB_NAME` | Database name | `beers_db` | No |
| `DB_USER` | Database user | `postgres` | No |
| `DB_PASSWORD` | Database password | `password` | No |
| `CURRENCY_API_KEY` | CurrencyLayer API key | - | No* |

*Required when using currency conversion features

## 🐳 Docker

### Build and Run
```bash
# Build Docker image
make docker-build

# Run with Docker
make docker-run

# Run with PostgreSQL using docker-compose
make docker-run-postgres
```

### Docker Compose Services
- **app**: Go application
- **postgres**: PostgreSQL database
- **adminer**: Database management UI (optional)
- **redis**: Redis cache (optional, for future use)

## 📈 Features

- ✅ **RESTful API** for beer management
- ✅ **Box price calculation** with tax and discount support
- ✅ **Multi-currency support** with real-time conversion
- ✅ **Multiple database backends** (PostgreSQL, In-Memory)
- ✅ **Structured logging** (JSON/Text formats)
- ✅ **Graceful shutdown** handling
- ✅ **CORS support** for cross-origin requests
- ✅ **Health check endpoint** for monitoring
- ✅ **Docker support** with multi-stage builds
- ✅ **Comprehensive testing** (unit + integration)
- ✅ **API documentation** with OpenAPI/Swagger
- ✅ **Development tools** (Makefile, scripts)

## 🔍 Code Quality & Best Practices

This project follows Go best practices and enterprise patterns:

### Architecture Principles
- **Hexagonal Architecture** for clean separation of concerns
- **SOLID principles** implementation
- **Dependency Injection** for loose coupling
- **Interface-driven design** for testability
- **Domain-driven design** with rich domain models

### Code Quality
- **Comprehensive error handling** with typed errors
- **Input validation** at domain and API levels
- **Thread-safe operations** for concurrent access
- **Structured logging** with contextual information
- **Configuration management** with environment variables
- **Database abstraction** for easy switching between storage backends

### Testing Strategy
- **Unit tests** for domain logic and services
- **Integration tests** for database operations
- **API tests** for HTTP endpoints
- **Mocking** for external dependencies
- **Test coverage** reporting

## 📝 API Examples

### Beer Entity
```json
{
  "id": 1,
  "name": "Corona Extra",
  "brewery": "Cervecería Modelo",
  "country": "Mexico",
  "price": 1200,
  "currency": "CLP",
  "created_at": "2025-08-24T10:30:00Z",
  "updated_at": "2025-08-24T10:30:00Z"
}
```

### Box Price Calculation
```bash
GET /api/v1/beers/1/boxprice?quantity=12&currency=USD
```

Response:
```json
{
  "beer_id": 1,
  "beer_name": "Corona Extra",
  "quantity": 12,
  "unit_price": 1200,
  "unit_currency": "CLP",
  "total_price": 14400,
  "target_currency": "USD",
  "exchange_rate": 0.0012,
  "converted_total": 17.28
}
```

### Error Response
```json
{
  "error": "VALIDATION_ERROR",
  "message": "Price must be greater than 0",
  "code": "INVALID_PRICE"
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/awesome-feature`)
3. Make your changes following the coding standards
4. Add tests for your changes
5. Run `make lint` and `make test`
6. Commit your changes (`git commit -m 'Add awesome feature'`)
7. Push to the branch (`git push origin feature/awesome-feature`)
8. Open a Pull Request

### Development Guidelines
- Follow Go conventions and idioms
- Write tests for new functionality
- Update documentation for API changes
- Use meaningful commit messages
- Keep PRs focused and small

## 📚 Documentation

- [Architecture Guide](docs/ARCHITECTURE.md) - Detailed architecture documentation
- [API Testing Guide](docs/REST_CLIENT_GUIDE.md) - How to test the API
- [cURL Examples](docs/API_TESTING.md) - Command-line testing examples
- [Cleanup Report](docs/CLEANUP_REPORT.md) - Code cleanup and refactoring notes

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) web framework
- Database integration with [lib/pq](https://github.com/lib/pq) PostgreSQL driver
- Testing with [testify](https://github.com/stretchr/testify)
- Architecture inspired by Hexagonal Architecture principles

---

**Built with ❤️ and Go**

## Especificaciones
Para ejecutar este aplicativo se necesita:

### Desarrollo
Solo es necesario el api token:
````shell
API_TOKEN=<YOUR_API_TOKEN>
````

### Producción
Se necesita una DB postgres configurada, por default toma los siguientes valores:
````shell
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=postgres
DB_USER=postgres
DB_PASSWORD=postgres
````

## Problema

Bender es fanático de las cervezas y quiere tener un registro de todas las cervezas que prueba y como calcular el precio que necesita para comprar una caja de algún tipo especifico de cervezas. Para esto necesita una API REST con esta información que posteriormente compartirá con sus amigos.

### Descripción

Se solicita crear un API REST basándonos en la definición que se encuentra en el archivo **openapi.yaml**.

#### Funcionalidad

- GET /Beers: Lista todas las cervezas que se encuentran en el sistema.
- POST /Beers: Permite ingresar una nueva cerveza.
- GET /beers/{beerID}: Lista un detalle de una cerveza especifica.
- GET /beets/{beerID}/boxprice: Entrega el valor que cuesta una caja específica de cerveza dependiendo de los parámetros ingresados, esto quiere decir que multiplique el precio por la cantidad una vez se homologara la moneda a lo que se ingreso por parámetro.
    - Quantity: Cantidad de cervezas a comprar (valor por defecto 6).
    - Currency: Tipo de moneda con la que desea pagar, para este caso se recomienda que utilice esta API <https://currencylayer.com/>

### Requisitos

- Puede usar alguno de los siguientes lenguajes Java, NodeJS, Go o Python. Aunque valoramos el uso de GO.
- Usar Docker y Docker Compose para los diferentes servicios.
- Se puede usar librarías externas y frameworks
- Requisito un 70% de cobertura de código
- Completa libertad para agregar nuevas funcionalidades.

### Entrega

- Enviar el link del repositorio donde se realiza este ejercicio.