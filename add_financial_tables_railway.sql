-- Add financial tables to Railway PostgreSQL
-- This script adds the missing financial tables that the charts need

-- Stock Prices table
CREATE TABLE IF NOT EXISTS stock_prices (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    volume BIGINT NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

-- Financial Indicators table
CREATE TABLE IF NOT EXISTS financial_indicators (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    pe_ratio DECIMAL(8,2),
    roe DECIMAL(8,2),
    profit_margin DECIMAL(8,2),
    revenue BIGINT,
    net_income BIGINT,
    market_cap BIGINT,
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

-- Market Data table
CREATE TABLE IF NOT EXISTS market_data (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    sp500_price DECIMAL(10,2),
    nasdaq_price DECIMAL(10,2),
    dow_price DECIMAL(10,2),
    vix DECIMAL(5,2),
    treasury_10y DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stock_prices_company_date ON stock_prices(company_id, date);
CREATE INDEX IF NOT EXISTS idx_stock_prices_date ON stock_prices(date);
CREATE INDEX IF NOT EXISTS idx_financial_indicators_company_date ON financial_indicators(company_id, date);
CREATE INDEX IF NOT EXISTS idx_financial_indicators_date ON financial_indicators(date);
CREATE INDEX IF NOT EXISTS idx_market_data_date ON market_data(date);

-- Create triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_stock_prices_updated_at BEFORE UPDATE ON stock_prices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_financial_indicators_updated_at BEFORE UPDATE ON financial_indicators FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_market_data_updated_at BEFORE UPDATE ON market_data FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample stock prices for companies
INSERT INTO stock_prices (company_id, price, volume, date) VALUES
(1, 175.50, 45000000, '2025-09-13'),
(1, 174.20, 42000000, '2025-09-12'),
(1, 176.80, 48000000, '2025-09-11'),
(2, 420.30, 25000000, '2025-09-13'),
(2, 418.90, 23000000, '2025-09-12'),
(2, 422.10, 27000000, '2025-09-11'),
(3, 280.45, 35000000, '2025-09-13'),
(3, 278.90, 33000000, '2025-09-12'),
(3, 282.15, 37000000, '2025-09-11'),
(4, 250.20, 28000000, '2025-09-13'),
(4, 248.50, 26000000, '2025-09-12'),
(4, 252.80, 30000000, '2025-09-11'),
(5, 145.75, 20000000, '2025-09-13'),
(5, 144.30, 18000000, '2025-09-12'),
(5, 147.20, 22000000, '2025-09-11');

-- Insert sample financial indicators for companies
INSERT INTO financial_indicators (company_id, pe_ratio, roe, profit_margin, revenue, net_income, market_cap, date) VALUES
(1, 28.5, 22.3, 24.1, 365000000000, 88000000000, 3000000000000, '2025-09-13'),
(2, 31.2, 18.7, 21.5, 211000000000, 45300000000, 2800000000000, '2025-09-13'),
(3, 25.8, 19.4, 23.8, 185000000000, 44000000000, 1500000000000, '2025-09-13'),
(4, 45.2, 15.6, 8.9, 96000000000, 8500000000, 800000000000, '2025-09-13'),
(5, 26.3, 16.8, 19.2, 282000000000, 54000000000, 1800000000000, '2025-09-13'),
(6, 22.1, 20.5, 26.3, 134000000000, 35200000000, 900000000000, '2025-09-13'),
(7, 35.7, 14.2, 12.8, 61000000000, 7800000000, 1200000000000, '2025-09-13'),
(8, 18.9, 17.3, 18.7, 89000000000, 16600000000, 450000000000, '2025-09-13'),
(9, 21.4, 16.9, 17.2, 76000000000, 13100000000, 380000000000, '2025-09-13'),
(10, 25.6, 18.1, 20.4, 42000000000, 8600000000, 260000000000, '2025-09-13');

-- Insert sample market data
INSERT INTO market_data (date, sp500_price, nasdaq_price, dow_price, vix, treasury_10y) VALUES
('2025-09-13', 5450.25, 17450.80, 39550.30, 18.5, 4.25),
('2025-09-12', 5432.10, 17380.45, 39420.80, 19.2, 4.30),
('2025-09-11', 5465.80, 17520.90, 39680.50, 17.8, 4.22),
('2025-09-10', 5448.30, 17420.60, 39510.20, 18.9, 4.28),
('2025-09-09', 5420.15, 17350.25, 39380.90, 20.1, 4.35),
('2025-09-08', 5480.75, 17580.40, 39720.60, 16.9, 4.18),
('2025-09-07', 5455.20, 17480.15, 39620.30, 18.3, 4.24),
('2025-09-06', 5438.90, 17410.85, 39480.70, 19.6, 4.31);

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
