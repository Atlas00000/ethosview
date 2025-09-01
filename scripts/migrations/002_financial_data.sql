-- Financial Data Schema Migration
-- Week 4: Financial Data APIs

-- Stock prices table for historical price data
CREATE TABLE IF NOT EXISTS stock_prices (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    open_price DECIMAL(10,2) NOT NULL,
    high_price DECIMAL(10,2) NOT NULL,
    low_price DECIMAL(10,2) NOT NULL,
    close_price DECIMAL(10,2) NOT NULL,
    volume BIGINT NOT NULL,
    adjusted_close DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, date)
);

-- Financial indicators table
CREATE TABLE IF NOT EXISTS financial_indicators (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    market_cap DECIMAL(20,2),
    pe_ratio DECIMAL(10,4),
    pb_ratio DECIMAL(10,4),
    debt_to_equity DECIMAL(10,4),
    return_on_equity DECIMAL(10,4),
    profit_margin DECIMAL(10,4),
    revenue_growth DECIMAL(10,4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, date)
);

-- Market data table for broader market indicators
CREATE TABLE IF NOT EXISTS market_data (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    sp500_close DECIMAL(10,2),
    nasdaq_close DECIMAL(10,2),
    dow_close DECIMAL(10,2),
    vix_close DECIMAL(10,4),
    treasury_10y DECIMAL(10,4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stock_prices_company_date ON stock_prices(company_id, date);
CREATE INDEX IF NOT EXISTS idx_stock_prices_date ON stock_prices(date);
CREATE INDEX IF NOT EXISTS idx_financial_indicators_company_date ON financial_indicators(company_id, date);
CREATE INDEX IF NOT EXISTS idx_financial_indicators_date ON financial_indicators(date);
CREATE INDEX IF NOT EXISTS idx_market_data_date ON market_data(date);

-- Create updated_at trigger function if not exists
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add triggers for updated_at
CREATE TRIGGER update_stock_prices_updated_at BEFORE UPDATE ON stock_prices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_financial_indicators_updated_at BEFORE UPDATE ON financial_indicators FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_market_data_updated_at BEFORE UPDATE ON market_data FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
