-- Performance Optimization Migration
-- Phase 1: Database Index Optimization
-- Week 10: Immediate optimizations

-- Composite indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_esg_scores_company_date ON esg_scores(company_id, score_date DESC);
CREATE INDEX IF NOT EXISTS idx_esg_scores_overall_date ON esg_scores(overall_score DESC, score_date DESC);
CREATE INDEX IF NOT EXISTS idx_esg_scores_environmental ON esg_scores(environmental_score DESC);
CREATE INDEX IF NOT EXISTS idx_esg_scores_social ON esg_scores(social_score DESC);
CREATE INDEX IF NOT EXISTS idx_esg_scores_governance ON esg_scores(governance_score DESC);

-- Composite indexes for financial data
CREATE INDEX IF NOT EXISTS idx_stock_prices_company_date_desc ON stock_prices(company_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_stock_prices_close_price ON stock_prices(close_price DESC);
CREATE INDEX IF NOT EXISTS idx_stock_prices_volume ON stock_prices(volume DESC);
CREATE INDEX IF NOT EXISTS idx_financial_indicators_company_date_desc ON financial_indicators(company_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_financial_indicators_pe_ratio ON financial_indicators(pe_ratio);
CREATE INDEX IF NOT EXISTS idx_financial_indicators_market_cap ON financial_indicators(market_cap DESC);

-- Partial indexes for active data
CREATE INDEX IF NOT EXISTS idx_companies_active ON companies(id) WHERE market_cap > 0;
CREATE INDEX IF NOT EXISTS idx_esg_scores_recent ON esg_scores(company_id, score_date DESC) WHERE score_date >= CURRENT_DATE - INTERVAL '1 year';

-- Text search indexes for company search
CREATE INDEX IF NOT EXISTS idx_companies_name_gin ON companies USING gin(to_tsvector('english', name));
CREATE INDEX IF NOT EXISTS idx_companies_symbol_gin ON companies USING gin(to_tsvector('english', symbol));

-- Indexes for analytics queries
CREATE INDEX IF NOT EXISTS idx_esg_scores_sector_date ON esg_scores(score_date DESC) INCLUDE (overall_score, environmental_score, social_score, governance_score);
CREATE INDEX IF NOT EXISTS idx_stock_prices_date_range ON stock_prices(date) WHERE date >= CURRENT_DATE - INTERVAL '1 year';

-- Performance optimization: Add covering indexes for common queries
CREATE INDEX IF NOT EXISTS idx_esg_scores_company_covering ON esg_scores(company_id, score_date DESC) 
    INCLUDE (overall_score, environmental_score, social_score, governance_score, data_source);

-- Index for sector-based queries
CREATE INDEX IF NOT EXISTS idx_companies_sector_market_cap ON companies(sector, market_cap DESC) WHERE market_cap > 0;

-- Add comments for documentation
COMMENT ON INDEX idx_esg_scores_company_date IS 'Optimized for company ESG score history queries';
COMMENT ON INDEX idx_esg_scores_overall_date IS 'Optimized for top ESG performers queries';
COMMENT ON INDEX idx_stock_prices_company_date_desc IS 'Optimized for company stock price history queries';
COMMENT ON INDEX idx_companies_name_gin IS 'Full-text search index for company names';
COMMENT ON INDEX idx_esg_scores_company_covering IS 'Covering index for ESG score queries with all score components';
