## EthosView Frontend Integration — Homepage (Weekly, Prioritized)

Status
- In Progress: Week 1 underway. Frontend scaffolded with Next.js via pnpm, env configured, and Core Homepage sections (Hero header + KPI band, Market snapshot bar) integrated to live Go/Gin backend.

Principles
- Keep it simple, in scope, and modular.
- Hug the Go/Gin backend: live endpoints only, no mock data.
- Prefer pnpm; future UI will be Next.js (App Router) when greenlit.
- API base (when implemented): `http://localhost:8080`.

Weekly Implementation Plan (Homepage first)

Week 1 — Core Homepage (Highest priority)
- Hero header + KPI band
  - Content: tagline, CTA, KPIs (total companies, sectors, avg ESG, last updated)
  - Backend: GET `/api/v1/dashboard`, GET `/api/v1/analytics/summary`, GET `/health/live`
  - Details: refresh every 60s; show skeletons; degrade gracefully on health failures
- Market snapshot bar
  - Content: market level, breadth (advancers/decliners), top gainers/losers, intraday sparkline
  - Backend: GET `/api/v1/financial/market`, GET `/api/v1/financial/market/history?range=1d`
  - Details: client cache 30–60s; tiny sparkline from history response
- ESG highlights (baseline)
  - Content: avg ESG by sector, top ESG performers
  - Backend: GET `/api/v1/analytics/summary`, GET `/api/v1/analytics/top-performers/esg_score?limit=5`
  - Details: small bar/tiles; tooltips with sector mean vs market

Week 2 — Discovery & Featured
- Featured companies carousel
  - Content: logo, name, symbol, latest price, % change, ESG badge
  - Backend: GET `/api/v1/companies?sort=market_cap&limit=10`, GET `/api/v1/financial/companies/:id/price/latest`, GET `/api/v1/esg/companies/:id/latest`
  - Details: lazy-load per slide; prefetch next slide details
- Company quick search
  - Content: autocomplete by name/symbol; quick jump
  - Backend: GET `/api/v1/companies?query=...&limit=10`, GET `/api/v1/companies/symbol/:symbol`
  - Details: debounce 250ms; keyboard navigation; highlight matches

Week 3 — Deeper Insights
- Sector performance heatmap
  - Content: sector tiles colored by daily return; click to drill in
  - Backend: GET `/api/v1/analytics/sectors/comparisons`
  - Details: color scale anchored to session min/max to avoid flicker
- Business dashboard preview
  - Content: mini-KPIs from the business dashboard
  - Backend: GET `/api/v1/dashboard`
  - Details: link to full dashboard page (later)

Week 4 — Realtime & Quality
- Live ticker and alerts (non-blocking)
  - Content: scrolling price/ESG alerts, connection status dot
  - Backend: WS `/api/v1/ws`, GET `/api/v1/ws/status`, GET `/alerts`
  - Details: WS optional; fallback to polling `/alerts` every 60s
- Mini ESG vs Financial correlation teaser
  - Content: small scatter with correlation coefficient
  - Backend: GET `/api/v1/analytics/correlation/esg-financial`
  - Details: fixed, small sample size to keep payload tiny
- Health/latency badge + Footer freshness
  - Content: API up/down, p50 latency; “Data updated Xm ago”
  - Backend: GET `/health/ready`, GET `/metrics`, timestamps from `/api/v1/dashboard` and `/api/v1/financial/market`
  - Details: compute latency via timed ping; hide if noisy

Priorities
1. Use existing live endpoints; no mock data.
2. Consolidate requests per section; simple loading/error states.
3. Small, reusable components; minimal shared state.
4. Accessibility and responsive layout.

Key Backend Endpoints (Homepage)
- GET `/api/v1/dashboard`
- GET `/api/v1/analytics/summary`
- GET `/api/v1/financial/market`
- GET `/api/v1/financial/market/history`
- GET `/api/v1/companies?sort=market_cap&limit=10`
- GET `/api/v1/financial/companies/:id/price/latest`
- GET `/api/v1/esg/companies/:id/latest`
- GET `/api/v1/analytics/top-performers/esg_score`
- GET `/api/v1/companies?query=...&limit=10`
- GET `/api/v1/companies/symbol/:symbol`
- GET `/api/v1/analytics/sectors/comparisons`
- WS `/api/v1/ws`, GET `/api/v1/ws/status`
- GET `/alerts`, GET `/metrics`, GET `/health/*`

---

## Additions — Homepage Enhancements (Proposed)
Note: These are additive, in-scope enhancements that hug existing APIs. Implement incrementally.

- Market sparkline with range select
  - Backend: GET `/api/v1/financial/market/history` (e.g., 1W/1M)
  - UI: Tiny sparkline under Market Snapshot; smooth crossfade on range change

- ESG vs Financial correlation teaser
  - Backend: GET `/api/v1/analytics/correlation/esg-financial`
  - UI: Mini scatter with current R values; tooltip shows metrics

- Top P/E leaders/laggards
  - Backend: GET `/api/v1/analytics/top-performers/pe_ratio?limit=5`
  - UI: Two compact lists (low/high P/E) with company names and values

- Alerts strip (non-blocking)
  - Backend: GET `/alerts`
  - UI: Subtle scrolling banner; click to expand details; hides when empty

- WS status indicator
  - Backend: GET `/api/v1/ws/status`
  - UI: Status dot (green/amber/red) near header or ticker

- Data freshness/latency badges
  - Backend: GET `/health/ready`, GET `/metrics`, timestamps from `/api/v1/dashboard` and `/api/v1/financial/market`
  - UI: “Updated Xm ago” + p50 latency; auto-hides if noisy

- ESG highlights (enhanced tooltips)
  - Backend: existing summary + `/api/v1/analytics/sectors/comparisons`
  - UI: Tooltip shows sector mean vs company score deltas

- Quick jump after symbol lookup
  - Backend: GET `/api/v1/companies/symbol/:symbol`
  - UI: On success, CTA to open company detail (future), for now scroll to ESG/Market sections

---

## Additions — Homepage Enhancements (Proposed, Batch 2)
All items hug existing endpoints; implement incrementally, no mock data.

- ESG trend mini for a top ESG company
  - Backend: GET `/api/v1/analytics/companies/:id/esg-trends?days=30`, GET `/api/v1/esg/companies/:id/latest`
  - UI: Small line chart with latest score badge

- Advanced insights teaser
  - Backend: GET `/api/v1/advanced/summary`
  - UI: 3–4 KPI tiles (e.g., avg risk, portfolio lift), link to dashboard

- Company financial snapshot (spotlight)
  - Backend: GET `/api/v1/financial/companies/:id/indicators`, GET `/api/v1/financial/companies/:id/price/latest`
  - UI: Compact metric tiles (P/E, ROE, margin, price)

- Sector distribution pie
  - Backend: GET `/api/v1/analytics/sectors/comparisons`
  - UI: Pie/donut of company counts by sector; legend with counts

- Risk assessment teaser
  - Backend: GET `/api/v1/advanced/companies/:id/risk-assessment`
  - UI: Gauge/badge (low/med/high) with tooltip context

- ESG scores feed
  - Backend: GET `/api/v1/esg/scores?limit=10&min_score=0`
  - UI: Recent/top scores list with company and date

- Business alerts counters
  - Backend: GET `/alerts`
  - UI: Small counters (active alerts) with link to details

- Live ticker (optional, non-blocking)
  - Backend: WS `/api/v1/ws`, fallback to GET `/alerts`
  - UI: Scrolling events; auto-pause on hover; hides when empty

Notes
- No frontend scaffolding will begin until explicitly requested.
- When approved, use pnpm and Next.js; integrate directly with the running Docker backend.
