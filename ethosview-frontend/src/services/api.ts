export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const memoryCache = new Map<string, { expiresAt: number; data: unknown }>();
const backoffUntil = new Map<number, number>();
const backoffUntilByKey = new Map<string, number>();
const requestQueue: Array<() => void> = [];
let activeRequests = 0;
const MAX_CONCURRENCY = 4;

async function acquireSlot() {
  if (activeRequests < MAX_CONCURRENCY) {
    activeRequests++;
    return;
  }
  await new Promise<void>((resolve) => {
    requestQueue.push(() => {
      activeRequests++;
      resolve();
    });
  });
}

function releaseSlot() {
  activeRequests = Math.max(0, activeRequests - 1);
  const next = requestQueue.shift();
  if (next) next();
}
const inflight = new Map<string, Promise<unknown>>();

function nowMs() {
  return Date.now();
}

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function getJson<T>(path: string, init?: RequestInit, ttlMs = 15000): Promise<T> {
  const url = `${API_BASE_URL}${path}`;
  const key = url;
  const backoff = backoffUntil.get(429) || 0;
  const pathBackoff = backoffUntilByKey.get(key) || 0;
  if (backoff > nowMs() || pathBackoff > nowMs()) {
    const cached = memoryCache.get(key);
    if (cached) return cached.data as T; // serve stale during backoff
  }
  const cached = memoryCache.get(key);
  if (cached && cached.expiresAt > nowMs()) {
    return cached.data as T;
  }
  if (inflight.has(key)) {
    return (await inflight.get(key)) as T;
  }

  const fetchPromise = (async () => {
    // Gentle jitter for low-priority endpoints to avoid synchronized bursts
    const lowPriority = path.startsWith("/alerts") || path.startsWith("/api/v1/ws/status") || path.startsWith("/metrics") || path.startsWith("/api/v1/esg/scores");
    if (lowPriority) {
      const jitter = 50 + Math.floor(Math.random() * 100);
      await sleep(jitter);
    }

    await acquireSlot();
    const maxAttempts = 3;
    let attempt = 0;
    let lastErr: unknown = null;
    while (attempt < maxAttempts) {
      attempt++;
      const res = await fetch(url, {
        ...init,
        headers: {
          Accept: "application/json",
          ...(init?.headers || {}),
        },
      });
      if (res.ok) {
        const json = (await res.json()) as T;
        if (ttlMs > 0) memoryCache.set(key, { data: json, expiresAt: nowMs() + ttlMs });
        return json;
      }

      // Build error text for diagnostics
      let text = "";
      try {
        text = await res.text();
      } catch {}

      if (res.status === 429) {
        // Global and per-path backoff with jitter; prefer serving stale
        const retryAfter = Number(res.headers.get("Retry-After") || "") || 0;
        const base = retryAfter > 0 ? Math.min(retryAfter * 1000, 3000) : 900;
        const jitter = 200 + Math.floor(Math.random() * 400);
        const delayMs = Math.max(base, 900) + jitter;
        const until = nowMs() + delayMs;
        backoffUntil.set(429, until);
        backoffUntilByKey.set(key, until);
        const anyCache = memoryCache.get(key);
        if (anyCache) {
          return anyCache.data as T; // serve stale immediately on 429
        }
        if (attempt < maxAttempts) {
          await sleep(delayMs);
          continue;
        }
      }

      lastErr = new Error(`HTTP ${res.status} ${res.statusText}${text ? `: ${text}` : ""}`);
      break;
    }
    // As a last resort, if we have any stale cache, serve it instead of throwing
    const stale = memoryCache.get(key);
    if (stale) return stale.data as T;
    throw lastErr ?? new Error("Request failed");
  })().finally(() => {
    inflight.delete(key);
    releaseSlot();
  });

  inflight.set(key, fetchPromise);
  return (await fetchPromise) as T;
}

export const api = {
  dashboard: () => getJson("/api/v1/dashboard"),
  analyticsSummary: () => getJson("/api/v1/analytics/summary", undefined, 120000),
  marketLatest: () => getJson("/api/v1/financial/market", undefined, 15000),
  latestPrice: (companyId: number) => getJson(`/api/v1/financial/companies/${companyId}/price/latest`),
  stockPrices: (companyId: number, limit = 30) => getJson(`/api/v1/financial/companies/${companyId}/prices?limit=${limit}`, undefined, 20000),
  latestESG: (companyId: number) => getJson(`/api/v1/esg/companies/${companyId}/latest`),
  companyBySymbol: (symbol: string) => getJson(`/api/v1/companies/symbol/${encodeURIComponent(symbol)}`),
  companyById: (companyId: number) => getJson(`/api/v1/companies/${companyId}`),
  sectorComparisons: () => getJson("/api/v1/analytics/sectors/comparisons", undefined, 180000),
  marketHistory: (start: string, end: string, limit = 30) =>
    getJson(`/api/v1/financial/market/history?start_date=${start}&end_date=${end}&limit=${limit}`, undefined, 30000),
  correlation: () => getJson("/api/v1/analytics/correlation/esg-financial", undefined, 60000),
  topPERatio: (limit = 5) => getJson(`/api/v1/analytics/top-performers/pe_ratio?limit=${limit}`, undefined, 30000),
  alerts: (all = false) => getJson(`/alerts?all=${all ? "true" : "false"}`, undefined, 15000),
  wsStatus: () => getJson("/api/v1/ws/status", undefined, 5000),
  esgTrends: (companyId: number, days = 30) => getJson(`/api/v1/analytics/companies/${companyId}/esg-trends?days=${days}`, undefined, 15000),
  financialIndicators: (companyId: number) => getJson(`/api/v1/financial/companies/${companyId}/indicators`, undefined, 30000),
  advancedSummary: () => getJson(`/api/v1/advanced/summary`, undefined, 60000),
  riskAssessment: (companyId: number) => getJson(`/api/v1/advanced/companies/${companyId}/risk-assessment`, undefined, 60000),
  esgScores: (limit = 20, min = 0, offset = 0) => getJson(`/api/v1/esg/scores?limit=${limit}&min_score=${min}&offset=${offset}`, undefined, 10000),
  companyFinancialSummary: (companyId: number) => getJson(`/api/v1/financial/companies/${companyId}/summary`, undefined, 15000),
};

export function isUnderGlobalBackoff(): number {
  const until = backoffUntil.get(429) || 0;
  const now = nowMs();
  return Math.max(0, until - now);
}
