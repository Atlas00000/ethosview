export type DashboardResponse = {
  summary: { total_companies: number; total_sectors: number; avg_esg_score: number; };
  top_esg_scores: Array<{ id: number; company_id: number; overall_score: number; company_name?: string; company_symbol?: string; score_date: string; }>;\n  sectors: string[];\n  sector_stats: Record<string, number>;\n};\n\nexport type PerformanceMetric = {\n  company_id: number;\n  company_name: string;\n  metric: string;\n  value: number;\n  rank: number;\n  total_count: number;\n  percentile: number;\n  date: string;\n};\n\nexport type SectorComparison = {\n  sector: string;\n  company_count: number;\n  avg_esg_score: number;\n  avg_pe_ratio: number;\n  avg_market_cap: number;\n  total_market_cap: number;\n  best_esg_company: string;\n  worst_esg_company: string;\n};\n\nexport type AnalyticsSummaryResponse = {\n  summary: { total_companies: number; total_sectors: number; avg_esg_score: number; };\n  sector_comparisons: SectorComparison[];\n  top_esg_performers: PerformanceMetric[];\n  top_market_cap: PerformanceMetric[];\n  correlation: Record<string, unknown>;\n};\n\nexport type MarketLatestResponse = {\n  market_data: {\n    date: string;\n    sp500_close?: number;\n    nasdaq_close?: number;\n    dow_close?: number;\n    vix_close?: number;\n    treasury_10y?: number;\n  };\n};\n\nexport type LatestPriceResponse = {\n  company_id: number;\n  price: { close_price: number; date: string; volume: number } | null;\n};\n\nexport type LatestESGResponse = {\n  id: number;\n  company_id: number;\n  overall_score: number;\n  score_date: string;\n  company_name?: string;\n  company_symbol?: string;\n};\n\nexport type CompanyResponse = {\n  id: number;\n  name: string;\n  symbol: string;\n  sector: string;\n  industry: string;\n  country: string;\n  market_cap: number;\n};

export type MarketHistoryItem = {
  id: number;
  date: string;
  sp500_close?: number;
  nasdaq_close?: number;
  dow_close?: number;
  vix_close?: number;
  treasury_10y?: number;
};

export type MarketHistoryResponse = {
  start_date: string;
  end_date: string;
  data: MarketHistoryItem[];
  count: number;
};

export type CorrelationResponse = {
  sample_size: number;
  avg_esg_score: number;
  avg_market_cap: number;
  avg_pe_ratio: number;
  avg_roe: number;
  avg_profit_margin: number;
  esg_market_cap_corr: number;
  esg_pe_corr: number;
  esg_roe_corr: number;
  esg_profit_corr: number;
};

export type TopPEResponse = {
  metric: string;
  top_performers: PerformanceMetric[];
  count: number;
  limit: number;
};

export type AlertsResponse = {
  alerts: unknown[];
  count: number;
  active_only: boolean;
};

export type ESGTrendPoint = { date: string; esg_score: number; e_score: number; s_score: number; g_score: number };
export type ESGTrendsResponse = { company_id: number; trends: ESGTrendPoint[]; count: number; days: number };

export type FinancialIndicatorsResponse = {
  company_id: number;
  indicators: {
    market_cap?: number;
    pe_ratio?: number;
    pb_ratio?: number;
    debt_to_equity?: number;
    return_on_equity?: number;
    profit_margin?: number;
    revenue_growth?: number;
  } | null;
};

export type AdvancedSummaryResponse = {
  summary: {
    portfolio_optimization: unknown;
    risk_summary: { total_companies_assessed: number; average_risk_score: number; risk_distribution: { low: number; medium: number; high: number } };
    trend_summary: { esg_trends: { improving: number; declining: number; stable: number }; price_trends: { up: number; down: number; stable: number } };
    message: string;
  };
};

export type RiskAssessmentResponse = { assessment: unknown; message: string };

export type ESGScoreItem = {
  id: number;
  company_id: number;
  overall_score: number;
  score_date: string;
  company_name?: string;
  company_symbol?: string;
};
export type ESGScoresListResponse = {
  scores: ESGScoreItem[];
  pagination: { limit: number; offset: number; count: number };
  filters?: { min_score?: number };
};
