-- Fix financial tables to match the expected schema
-- Drop and recreate with correct column names

-- Drop existing tables
DROP TABLE IF EXISTS financial_indicators CASCADE;
DROP TABLE IF EXISTS stock_prices CASCADE;
DROP TABLE IF EXISTS market_data CASCADE;

-- Recreate Financial Indicators table with correct column names
CREATE TABLE financial_indicators (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    date DATE NOT NULL,
    market_cap BIGINT,
    pe_ratio DECIMAL(8,2),
    pb_ratio DECIMAL(8,2),
    debt_to_equity DECIMAL(8,2),
    return_on_equity DECIMAL(8,2),
    profit_margin DECIMAL(8,2),
    revenue_growth DECIMAL(8,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

-- Recreate Stock Prices table with correct column names
CREATE TABLE stock_prices (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    date DATE NOT NULL,
    open_price DECIMAL(10,2),
    high_price DECIMAL(10,2),
    low_price DECIMAL(10,2),
    close_price DECIMAL(10,2),
    volume BIGINT,
    adjusted_close DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

-- Recreate Market Data table with correct column names
CREATE TABLE market_data (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    sp500_close DECIMAL(10,2),
    nasdaq_close DECIMAL(10,2),
    dow_close DECIMAL(10,2),
    vix_close DECIMAL(5,2),
    treasury_10y DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_financial_indicators_company_date ON financial_indicators(company_id, date);
CREATE INDEX idx_financial_indicators_date ON financial_indicators(date);
CREATE INDEX idx_stock_prices_company_date ON stock_prices(company_id, date);
CREATE INDEX idx_stock_prices_date ON stock_prices(date);
CREATE INDEX idx_market_data_date ON market_data(date);

-- Create triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_financial_indicators_updated_at BEFORE UPDATE ON financial_indicators FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stock_prices_updated_at BEFORE UPDATE ON stock_prices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_market_data_updated_at BEFORE UPDATE ON market_data FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample financial indicators for companies
INSERT INTO financial_indicators (company_id, date, market_cap, pe_ratio, pb_ratio, debt_to_equity, return_on_equity, profit_margin, revenue_growth) VALUES
(1, '2025-09-13', 3000000000000, 28.5, 5.2, 0.8, 22.3, 24.1, 8.5),
(2, '2025-09-13', 2800000000000, 31.2, 4.8, 1.2, 18.7, 21.5, 12.3),
(3, '2025-09-13', 1500000000000, 25.8, 3.9, 0.6, 19.4, 23.8, 9.7),
(4, '2025-09-13', 800000000000, 45.2, 12.1, 0.3, 15.6, 8.9, 25.4),
(5, '2025-09-13', 1800000000000, 26.3, 4.1, 0.9, 16.8, 19.2, 11.8),
(6, '2025-09-13', 900000000000, 22.1, 3.2, 0.4, 20.5, 26.3, 7.9),
(7, '2025-09-13', 1200000000000, 35.7, 8.5, 0.2, 14.2, 12.8, 18.6),
(8, '2025-09-13', 450000000000, 18.9, 2.8, 0.7, 17.3, 18.7, 5.4),
(9, '2025-09-13', 380000000000, 21.4, 3.1, 0.5, 16.9, 17.2, 4.8),
(10, '2025-09-13', 260000000000, 25.6, 4.3, 0.6, 18.1, 20.4, 6.2);

-- Insert sample stock prices for companies
INSERT INTO stock_prices (company_id, date, open_price, high_price, low_price, close_price, volume, adjusted_close) VALUES
(1, '2025-09-13', 174.20, 176.80, 173.50, 175.50, 45000000, 175.50),
(1, '2025-09-12', 172.80, 175.20, 171.90, 174.20, 42000000, 174.20),
(1, '2025-09-11', 175.90, 178.10, 174.50, 176.80, 48000000, 176.80),
(2, '2025-09-13', 418.90, 422.50, 417.20, 420.30, 25000000, 420.30),
(2, '2025-09-12', 416.50, 420.80, 415.10, 418.90, 23000000, 418.90),
(2, '2025-09-11', 420.10, 424.50, 418.80, 422.10, 27000000, 422.10),
(3, '2025-09-13', 278.90, 282.50, 277.20, 280.45, 35000000, 280.45),
(3, '2025-09-12', 276.80, 280.20, 275.50, 278.90, 33000000, 278.90),
(3, '2025-09-11', 280.50, 284.10, 279.20, 282.15, 37000000, 282.15),
(4, '2025-09-13', 248.50, 252.80, 246.90, 250.20, 28000000, 250.20),
(4, '2025-09-12', 246.20, 250.50, 244.80, 248.50, 26000000, 248.50),
(4, '2025-09-11', 250.80, 255.20, 249.50, 252.80, 30000000, 252.80),
(5, '2025-09-13', 144.30, 147.50, 143.20, 145.75, 20000000, 145.75),
(5, '2025-09-12', 142.80, 145.90, 141.50, 144.30, 18000000, 144.30),
(5, '2025-09-11', 146.20, 148.80, 144.90, 147.20, 22000000, 147.20);

-- Insert sample market data
INSERT INTO market_data (date, sp500_close, nasdaq_close, dow_close, vix_close, treasury_10y) VALUES
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
