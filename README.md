## EthosView

ESG/Financial analytics platform with a Go/Gin backend and a Next.js App Router frontend.

### At a glance
- **Backend**: Go 1.21, Gin, PostgreSQL, Redis, Docker
- **Frontend**: Next.js 15 (App Router), pnpm, server-side data fetching
- **Ports**: Frontend `3000`, Backend `8080`, Postgres `5432`, Redis `6379`
- **Live health**: `GET /health/live`, `GET /api/v1/health`

### Tech stack highlights
- **Go/Gin**: minimal, fast HTTP server with middleware for security, rate limiting, compression, monitoring
- **PostgreSQL**: relational core for companies, ESG scores, prices, indicators, market aggregates
- **Redis**: response caching and background warming for low latency
- **Next.js App Router**: modern SSR with dynamic rendering for live data; zero client mocks
- **pnpm**: efficient Node package management (front end)
- **Docker Compose**: reproducible local orchestration of all services

### Architecture
```mermaid
flowchart LR
  subgraph Browser
    UI["Next.js UI"]
  end
  subgraph Frontend Container
    FE["Next.js Server (SSR)"]
  end
  subgraph Backend Container
    API["Go/Gin API"]
  end
  subgraph Data Services
    PG["PostgreSQL"]
    R["Redis"]
  end

  UI <--> FE
  FE -->|HTTP (server)| API
  UI -->|HTTP (browser)| API
  API <--> PG
  API <--> R
```

### Backend capabilities
- Health suite: liveness, readiness, detailed metrics endpoints
- Companies API: CRUD, by-id/symbol lookup, sectors
- ESG analytics: latest scores, company trends, sector comparisons, top performers
- Financials: market snapshot/history, indicators, prices, company financial summary
- Advanced analytics (scaffolding): correlation, portfolio/risk summaries
- Middleware: request ID, compression, CORS, security headers, size limits, basic injection/XSS guards
- Rate limiting and simple in-memory monitoring hooks
- Background workers: cache warming, metrics collection, alert scanning

### Frontend aesthetics & UX
- Clean glassmorphism accents with subtle gradients
- Responsive grid layout for sections and cards
- Progressive loading: skeletons for async components; soft-fail for noncritical data
- Micro-interactions: scroll reveals, compact charts (sparkline, pie, bars)
- Accessibility: semantic headings, adequate contrast, keyboard-friendly interactions

### Key homepage sections (live data)
- Hero with KPIs (companies, sectors, avg ESG, last updated)
- Market bar with snapshot and intraday sparkline
- ESG highlights: sector comparisons and top performers
- Featured companies carousel + P/E leaders
- Symbol lookup with quick company preview
- ESG trend mini + Business preview + Financial snapshot
- Alerts strip and WebSocket status indicator

### Quick start (Docker)
1) Start services
```bash
docker compose up -d
```
2) Open the app
```
http://localhost:3000
```
3) Verify API
```bash
curl http://localhost:8080/api/v1/dashboard
```

### Environment configuration
- Backend (container):
  - `DB_HOST=postgres`, `DB_PORT=5432`, `DB_USER=postgres`, `DB_PASSWORD=password`, `DB_NAME=ethosview`
  - `REDIS_HOST=redis`, `REDIS_PORT=6379`
  - `PORT=8080`
- Frontend:
  - `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080` (browser → host mapped port)
  - `INTERNAL_API_BASE_URL=http://backend:8080` (server-side in container → backend service)

These are set in `docker-compose.yml` for a zero-config local run.

### Data seeding
The repo includes SQL migrations and sample data:
- `scripts/migrations/*.sql`
- `scripts/seeds/*.sql`

Seed using the running Postgres container (psql is available inside):
```bash
# Schema/migrations
docker exec -i ethosview-postgres psql -U postgres -d ethosview -f /tmp/001_initial_schema.sql
docker exec -i ethosview-postgres psql -U postgres -d ethosview -f /tmp/002_financial_data.sql
docker exec -i ethosview-postgres psql -U postgres -d ethosview -f /tmp/003_performance_optimization.sql

# Sample data
docker exec -i ethosview-postgres psql -U postgres -d ethosview -f /tmp/sample_data.sql
docker exec -i ethosview-postgres psql -U postgres -d ethosview -f /tmp/financial_data.sql
```

Alternatively, run the helper script (Linux/macOS shells):
```bash
chmod +x scripts/migrate.sh
./scripts/migrate.sh --seed
```

### Development (local)
- Backend
```bash
make deps && make run   # or: make dev (hot reload with Air)
```
- Frontend
```bash
cd ethosview-frontend
pnpm install
pnpm dev
```

### Key API endpoints
- Health: `GET /health`, `GET /health/live`, `GET /api/v1/health`
- Dashboard: `GET /api/v1/dashboard`
- Companies: `GET /api/v1/companies`, `GET /api/v1/companies/:id`, `GET /api/v1/companies/symbol/:symbol`
- ESG: `GET /api/v1/esg/companies/:id/latest`, `GET /api/v1/esg/scores`
- Financial: `GET /api/v1/financial/market`, `GET /api/v1/financial/companies/:id/summary`

### Performance & monitoring
- API client: in-memory TTL cache, max concurrency control, jitter/backoff on 429
- Server: compression, light caching, metrics collection, cache warming
- Containers: small production images (frontend standalone output), healthchecks

### Security
- CORS allowlist (`http://localhost:3000`, `https://ethosview.com`)
- Security headers: HSTS, X-Frame-Options, X-Content-Type-Options, CSP, Permissions-Policy
- Input sanitization and basic injection/XSS guards
- Request size limits and rate limiting

### Troubleshooting
- **Frontend shows empty data**
  - Ensure backend is healthy: `curl http://localhost:8080/health/live`
  - Confirm env split: browser uses `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080`, server-side uses `INTERNAL_API_BASE_URL=http://backend:8080`
  - We use dynamic rendering on the homepage (`force-dynamic`) to avoid stale ISR in containers.
- **Ports busy**
  - Free port 3000/8080 or stop local dev servers, then `docker compose up -d`.
- **Cache issues**
  - Clear Redis: `docker exec -i ethosview-redis redis-cli FLUSHALL`

### Make commands
```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make dev           # Run with hot reload (Air)
make test          # Run tests
make docker-up     # Start all services
make docker-down   # Stop all services
make clean         # Clean build artifacts
```

### Project structure
```
.
├── cmd/server                     # Main application entry
├── internal/                      # HTTP server, handlers, websocket
├── pkg/                           # auth, cache, db, middleware, health, metrics
├── ethosview-frontend/            # Next.js App Router frontend
├── scripts/                       # migrations, seeds, test scripts
├── docker-compose.yml             # Orchestration
├── Dockerfile                     # Backend image
└── README.md
```

### Notes
- Keep it simple and in-scope; prefer pnpm on the frontend.
- Frontend components fetch live data only; no mock data.
