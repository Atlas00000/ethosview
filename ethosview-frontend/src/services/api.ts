export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function getJson<T>(path: string, init?: RequestInit): Promise<T> {
  const url = `${API_BASE_URL}${path}`;
  const res = await fetch(url, {
    ...init,
    headers: {
      Accept: "application/json",
      ...(init?.headers || {}),
    },
  });
  if (!res.ok) {
    let text = "";
    try {
      text = await res.text();
    } catch {}
    throw new Error(`HTTP ${res.status} ${res.statusText}${text ? `: ${text}` : ""}`);
  }
  return (await res.json()) as T;
}

export const api = {
  dashboard: () => getJson("/api/v1/dashboard"),
  analyticsSummary: () => getJson("/api/v1/analytics/summary"),
  marketLatest: () => getJson("/api/v1/financial/market"),
  latestPrice: (companyId: number) => getJson(`/api/v1/financial/companies/${companyId}/price/latest`),
  latestESG: (companyId: number) => getJson(`/api/v1/esg/companies/${companyId}/latest`),
  companyBySymbol: (symbol: string) => getJson(`/api/v1/companies/symbol/${encodeURIComponent(symbol)}`),
  sectorComparisons: () => getJson("/api/v1/analytics/sectors/comparisons"),
  marketHistory: (start: string, end: string, limit = 30) =>
    getJson(`/api/v1/financial/market/history?start_date=${start}&end_date=${end}&limit=${limit}`),
  correlation: () => getJson("/api/v1/analytics/correlation/esg-financial"),
  topPERatio: (limit = 5) => getJson(`/api/v1/analytics/top-performers/pe_ratio?limit=${limit}`),
  alerts: (all = false) => getJson(`/alerts?all=${all ? "true" : "false"}`),
  wsStatus: () => getJson("/api/v1/ws/status"),
  esgTrends: (companyId: number, days = 30) => getJson(`/api/v1/analytics/companies/${companyId}/esg-trends?days=${days}`),
  financialIndicators: (companyId: number) => getJson(`/api/v1/financial/companies/${companyId}/indicators`),
  advancedSummary: () => getJson(`/api/v1/advanced/summary`),
  riskAssessment: (companyId: number) => getJson(`/api/v1/advanced/companies/${companyId}/risk-assessment`),
  esgScores: (limit = 10, min = 0) => getJson(`/api/v1/esg/scores?limit=${limit}&min_score=${min}`),
};
