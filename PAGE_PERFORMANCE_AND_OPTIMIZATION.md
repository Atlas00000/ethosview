# Page Performance and Optimization

Purpose: Improve perceived and actual performance of the EthosView homepage without changing functionality, staying within scope, and hugging the live backend. Work is divided into weekly, incremental, low-risk optimizations.

## Guiding Principles
- Keep it simple; prefer small, measurable wins
- Avoid scope creep; no functional changes
- Hug the backend: rely on real endpoints, respect rate limits
- Measure before/after (LCP, TTI, CLS, total requests, 429 rate)

## KPIs (track weekly)
- LCP (Largest Contentful Paint)
- TTI (Time to Interactive)
- Total JS transferred for homepage
- Number of initial network requests
- 429 rate and error log noise (client)

---

## Week 1: Network + Request Discipline (Quick Wins)
- Increase client TTLs for stable analytics endpoints to 3–5 minutes
- Serve cached responses during backoff (already) and extend stale-while-error window
- Add small jitter (50–150ms) for noncritical requests to avoid synchronized bursts
- Cap per-host concurrent requests at 4 (lightweight request queue)
- Maintain IntersectionObserver lazy details for cards; add per-section concurrency limits
- Add dns-prefetch + preconnect tags for API host (reduce handshake time)
- Audit and remove duplicate/in-flight redundant requests

Deliverables:
- Client request queue utility
- TTL and jitter tuning in API client
- <head> link hints for API host

---

## Week 2: Rendering Cost and Charts
- Slice chart series to last N points (desktop ~30, mobile ~20)
- Defer chart animations until section is revealed (reveal-on-scroll hook)
- Memoize transformed series and computed domains; avoid redundant parses
- Prefer BarChart/Pie where rendering cost is lower and visibility is higher
- Disable tooltips for offscreen charts; enable on reveal

Deliverables:
- Reveal hook toggling chart animation/tooltip
- Memoization for series/domain builders

---

## Week 3: Lazy Loading and Code Splitting
- Dynamic import below-the-fold sections (ESG Feed, Advanced Insights, Financial Snapshot)
- Show skeletons until a section enters viewport
- Split out heavy chart components to separate chunks; preload when near viewport
- Prefetch data on hover/focus for section nav anchors

Deliverables:
- dynamic() wrappers for selected sections with suspense fallbacks
- Preload triggers via IntersectionObserver and link hover

---

## Week 4: List Virtualization and Pagination Polish
- Virtualize long lists (ESG Feed) when item count > 50
- Keep “Load more”; prefetch next page when near end
- De-duplicate items on append (prevent re-renders and extra fetches)
- Maintain lazy per-card details with concurrency cap

Deliverables:
- Lightweight virtualization for ESG Feed
- Next-page prefetch logic

---

## Week 5: Error Budget & Resilience
- Centralize expected 404 handling as “no data” (optional dev-only quiet logs)
- Improve stale-on-429 fallback messages (UI badges only; no functional change)
- Add soft-degrade paths for chart tiles when data is missing

Deliverables:
- Small helper wrapper for gracefully handling 404 in dev
- UI badge for “stale” data state

---

## Optional (Backend-Friendly Enhancements)
Note: Only if/when in scope.
- Filter analytics/top lists to companies with financial coverage (SQL join/exists)
- Batch endpoints: latest price, summary, latest ESG for multiple IDs to reduce N+1
- Add ETag/Last-Modified and Cache-Control on analytics/sectors

---

## Measurement Plan
- Add a lightweight console summary in dev: nav timing, initial request count
- Record weekly KPI snapshots in this file (delta from previous week)

Template:
```
Week N KPIs:
- LCP:  
- TTI:  
- JS (KB):  
- Requests (initial):  
- 429 rate:  
Notes:
```


