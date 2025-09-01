# EthosView Backend

ESG/Financial Analytics Platform Backend API

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make (optional, for convenience)

### Development Setup

1. **Clone and setup**
```bash
# Install dependencies
make deps

# Or manually:
go mod download
go mod tidy
```

2. **Start services with Docker**
```bash
# Start PostgreSQL and Redis
make docker-up

# Or manually:
docker-compose up -d
```

3. **Run the backend**
```bash
# With hot reload (recommended for development)
make dev

# Or run directly:
make run
```

4. **Test the API**
```bash
# Health check
curl http://localhost:8080/health

# API health check
curl http://localhost:8080/api/v1/health
```

### Environment Variables

Copy the example environment file:
```bash
cp .env.example .env
```

Key environment variables:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL connection
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` - Redis connection
- `PORT` - Server port (default: 8080)

### Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make dev           # Run with hot reload
make test          # Run tests
make docker-up     # Start all services
make docker-down   # Stop all services
make clean         # Clean build artifacts
```

## API Endpoints

### Health Checks
- `GET /health` - Basic health check
- `GET /api/v1/health` - API health check with database connectivity

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/profile` - Get user profile (requires authentication)
- `PUT /api/v1/auth/profile` - Update user profile (requires authentication)

### Dashboard
- `GET /api/v1/dashboard` - Get ESG analytics dashboard overview

### Companies
- `GET /api/v1/companies` - List all companies (with pagination)
- `GET /api/v1/companies/:id` - Get company by ID
- `GET /api/v1/companies/symbol/:symbol` - Get company by symbol
- `GET /api/v1/companies/sectors` - Get all unique sectors
- `POST /api/v1/companies` - Create a new company
- `PUT /api/v1/companies/:id` - Update a company
- `DELETE /api/v1/companies/:id` - Delete a company

### ESG Scores
- `GET /api/v1/esg/scores` - List all ESG scores (with pagination and filtering)
- `GET /api/v1/esg/scores/:id` - Get ESG score by ID
- `GET /api/v1/esg/companies/:id/latest` - Get latest ESG score for a company
- `GET /api/v1/esg/companies/:id/scores` - Get all ESG scores for a company
- `POST /api/v1/esg/scores` - Create a new ESG score
- `PUT /api/v1/esg/scores/:id` - Update an ESG score
- `DELETE /api/v1/esg/scores/:id` - Delete an ESG score

### Query Parameters
- `limit` - Number of items per page (default: 20, max: 100)
- `offset` - Number of items to skip (default: 0)
- `sector` - Filter companies by sector
- `min_score` - Filter ESG scores by minimum overall score

### Authentication
Protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Project Structure

```
ethosview-backend/
├── cmd/
│   └── server/           # Main application entry
├── internal/
│   └── server/           # HTTP server and routes
├── pkg/
│   └── database/         # Database connections
├── Dockerfile            # Docker image
├── docker-compose.yml    # Services orchestration
├── Makefile              # Development commands
└── .air.toml            # Hot reload configuration
```

## Development

The backend uses:
- **Gin** for HTTP routing
- **PostgreSQL** for primary data storage
- **Redis** for caching
- **Docker** for containerization
- **Air** for hot reload during development

## Next Steps

This is Week 3 of the EthosView development roadmap. Next phases will include:
- Financial data APIs (Week 4)
- Analytics and reporting (Week 5)
- Real-time streaming (Week 7)
- Frontend development (Week 10)
