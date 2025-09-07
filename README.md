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

### Architecture (ASCII)
```
 [Browser]
    |
    |  HTTP (browser)
    v
 [Go/Gin API] <---- HTTP (server) ---- [Next.js Server (SSR) in Frontend container]
    |   \
    |    \__ metrics/health
    |
    +---- PostgreSQL (relational data)
    |
    +---- Redis (cache)
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
├── bin/
│   └── ethosview-backend                - compiled backend binary (local builds)
├── cmd/
│   └── server/
│       └── main.go                      - backend entrypoint
├── internal/
│   ├── handlers/                        - HTTP handlers (auth, company, esg, financial, analytics)
│   │   ├── advanced_analytics.go
│   │   ├── analytics.go
│   │   ├── auth.go
│   │   ├── company.go
│   │   ├── dashboard.go
│   │   ├── esg.go
│   │   ├── financial.go
│   │   └── websocket.go
│   ├── models/                          - data models shared by handlers
│   │   ├── advanced_analytics.go
│   │   ├── analytics.go
│   │   ├── company.go
│   │   ├── esg_score.go
│   │   ├── financial.go
│   │   └── user.go
│   ├── server/
│   │   └── server.go                    - router, middleware, routes
│   └── websocket/
│       └── manager.go                   - WS manager
├── pkg/                                 - reusable backend packages
│   ├── auth/jwt.go
│   ├── cache/{advanced.go,warming.go}
│   ├── dashboard/business.go
│   ├── database/{postgresql.go,redis.go}
│   ├── errors/errors.go
│   ├── health/health.go
│   ├── metrics/metrics.go
│   ├── middleware/{auth.go,cache.go,compression.go,monitoring.go,rate_limit.go,request_id.go,validation.go}
│   ├── monitoring/alerts.go
│   ├── pagination/cursor.go
│   └── security/security.go
├── ethosview-frontend/                  - Next.js (App Router) frontend
│   ├── Dockerfile
│   ├── next.config.ts
│   ├── package.json
│   ├── pnpm-lock.yaml
│   └── src/
│       ├── app/                         - App Router pages/layouts
│       │   ├── layout.tsx
│       │   ├── page.tsx                 - homepage (SSR, dynamic)
│       │   ├── loading.tsx
│       │   ├── not-found.tsx
│       │   └── favicon.ico
│       ├── components/
│       │   ├── layout/{Header.tsx,Footer.tsx}
│       │   ├── dashboard/{SectorBar.tsx,TopESGList.tsx}
│       │   └── home/                     - homepage widgets (market, ESG, featured, etc.)
│       │       ├── AdvancedInsights.tsx
│       │       ├── AlertsStrip.tsx
│       │       ├── BusinessPreview.tsx
│       │       ├── CompanyQuickView.tsx
│       │       ├── CorrelationTeaser.tsx
│       │       ├── ESGFeed.tsx
│       │       ├── ESGHighlights.tsx
│       │       ├── ESGHighlightsPro.tsx
│       │       ├── ESGTrendMini.tsx
│       │       ├── FeaturedCarousel.tsx
│       │       ├── FinancialSnapshot.tsx
│       │       ├── HeroNew.tsx
│       │       ├── MarketBar.tsx
│       │       ├── MarketSparkline.tsx
│       │       ├── PELeaders.tsx
│       │       ├── QuickWidget.tsx
│       │       ├── RiskTeaser.tsx
│       │       ├── ScrollReveal.tsx
│       │       ├── SectorHeatmap.tsx
│       │       ├── SectorPie.tsx
│       │       ├── SymbolLookup.tsx
│       │       ├── useCountUp.ts
│       │       └── WSStatus.tsx
│       ├── services/api.ts              - API client with caching/backoff
│       └── types/api.ts                 - shared types
├── scripts/                             - migrations, seeds, utilities
│   ├── migrations/{001_initial_schema.sql,002_financial_data.sql,003_performance_optimization.sql}
│   ├── seeds/{sample_data.sql,financial_data.sql}
│   ├── migrate.sh
│   ├── performance_test.sh
│   ├── phase2_test.sh
│   └── phase3_test.sh
├── Dockerfile                           - backend image
├── docker-compose.yml                   - local orchestration (Postgres, Redis, backend, frontend)
├── Makefile                             - dev commands
├── ETHOSVIEW_ROADMAP.md
├── FRONTEND_INTEGRATION.md
├── UI_THEME_SPEC.md
├── PERFORMANCE_OPTIMIZATIONS.md
├── PAGE_PERFORMANCE_AND_OPTIMIZATION.md
├── go.mod
└── go.sum
```

### Notes
- Keep it simple and in-scope; prefer pnpm on the frontend.
- Frontend components fetch live data only; no mock data.
