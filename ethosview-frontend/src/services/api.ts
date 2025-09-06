export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const memoryCache = new Map<string, { expiresAt: number; data: unknown }>();
const backoffUntil = new Map<number, number>();
const inflight = new Map<string, Promise<unknown>>();

function nowMs() {
  return Date.now();
}

export async function getJson<T>(path: string, init?: RequestInit, ttlMs = 15000): Promise<T> {
  const url = `${API_BASE_URL}${path}`;
  const key = url;
  const backoff = backoffUntil.get(429) || 0;
  if (backoff > nowMs()) {
    const cached = memoryCache.get(key);
    if (cached && cached.expiresAt > nowMs()) return cached.data as T;
  }
  const cached = memoryCache.get(key);
  if (cached && cached.expiresAt > nowMs()) {
    return cached.data as T;
  }
  if (inflight.has(key)) {
    return (await inflight.get(key)) as T;
  }

  const fetchPromise = fetch(url, {
    ...init,
    headers: {
      Accept: "application/json",
      ...(init?.headers || {}),
    },
  }).then(async (res) => {
    if (!res.ok) {
      let text = "";
      try {
        text = await res.text();
      } catch {}
      if (res.status === 429) {
        backoffUntil.set(429, nowMs() + 60_000);
        const anyCache = memoryCache.get(key);
        if (anyCache) {
          return anyCache.data as T; // serve stale on 429
        }
      }
      throw new Error(`HTTP ${res.status} ${res.statusText}${text ? `: ${text}` : ""}`);
    }
    const json = (await res.json()) as T;
    if (ttlMs > 0) memoryCache.set(key, { data: json, expiresAt: nowMs() + ttlMs });
    return json;
  }).finally(() => {
    inflight.delete(key);
  });

  inflight.set(key, fetchPromise);
  return (await fetchPromise) as T;
}

export const api = {
  dashboard: () => getJson("/api/v1/dashboard"),
  analyticsSummary: () => getJson("/api/v1/analytics/summary", undefined, 30000),
  marketLatest: () => getJson("/api/v1/financial/market", undefined, 15000),
  latestPrice: (companyId: number) => getJson(`/api/v1/financial/companies/${companyId}/price/latest`),
  latestESG: (companyId: number) => getJson(`/api/v1/esg/companies/${companyId}/latest`),
  companyBySymbol: (symbol: string) => getJson(`/api/v1/companies/symbol/${encodeURIComponent(symbol)}`),
  sectorComparisons: () => getJson("/api/v1/analytics/sectors/comparisons", undefined, 60000),
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
  esgScores: (limit = 10, min = 0) => getJson(`/api/v1/esg/scores?limit=${limit}&min_score=${min}`, undefined, 10000),
};
