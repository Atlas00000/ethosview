-- Supabase Migration for EthosView
-- This file contains the database schema for Supabase

-- Enable UUID extension for Supabase
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table (using UUID for Supabase compatibility)
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create companies table
CREATE TABLE IF NOT EXISTS companies (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(20) UNIQUE NOT NULL,
    sector VARCHAR(100),
    industry VARCHAR(100),
    country VARCHAR(100),
    market_cap DECIMAL(20,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create ESG scores table
CREATE TABLE IF NOT EXISTS esg_scores (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    environmental_score DECIMAL(5,2),
    social_score DECIMAL(5,2),
    governance_score DECIMAL(5,2),
    overall_score DECIMAL(5,2),
    score_date DATE NOT NULL,
    data_source VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_companies_symbol ON companies(symbol);
CREATE INDEX IF NOT EXISTS idx_companies_sector ON companies(sector);
CREATE INDEX IF NOT EXISTS idx_esg_scores_company_id ON esg_scores(company_id);
CREATE INDEX IF NOT EXISTS idx_esg_scores_date ON esg_scores(score_date);
CREATE INDEX IF NOT EXISTS idx_esg_scores_overall ON esg_scores(overall_score);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_esg_scores_updated_at BEFORE UPDATE ON esg_scores
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Enable Row Level Security (RLS) for Supabase
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE esg_scores ENABLE ROW LEVEL SECURITY;

-- Create policies for public access to companies and esg_scores
CREATE POLICY "Allow public read access to companies" ON companies
    FOR SELECT USING (true);

CREATE POLICY "Allow public read access to esg_scores" ON esg_scores
    FOR SELECT USING (true);

-- Create policy for users (only authenticated users can access their own data)
CREATE POLICY "Users can view own profile" ON users
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can update own profile" ON users
    FOR UPDATE USING (auth.uid() = id);

-- Insert sample data
INSERT INTO companies (name, symbol, sector, industry, country, market_cap) VALUES
('Apple Inc.', 'AAPL', 'Technology', 'Consumer Electronics', 'United States', 3000000000000),
('Microsoft Corporation', 'MSFT', 'Technology', 'Software', 'United States', 2800000000000),
('Amazon.com Inc.', 'AMZN', 'Consumer Discretionary', 'E-commerce', 'United States', 1500000000000),
('Tesla Inc.', 'TSLA', 'Consumer Discretionary', 'Electric Vehicles', 'United States', 800000000000),
('Alphabet Inc.', 'GOOGL', 'Technology', 'Internet Services', 'United States', 1800000000000),
('Meta Platforms Inc.', 'META', 'Technology', 'Social Media', 'United States', 900000000000),
('NVIDIA Corporation', 'NVDA', 'Technology', 'Semiconductors', 'United States', 1200000000000),
('Johnson & Johnson', 'JNJ', 'Healthcare', 'Pharmaceuticals', 'United States', 450000000000),
('Procter & Gamble', 'PG', 'Consumer Staples', 'Household Products', 'United States', 380000000000),
('Coca-Cola Company', 'KO', 'Consumer Staples', 'Beverages', 'United States', 260000000000)
ON CONFLICT (symbol) DO NOTHING;

-- Insert sample ESG scores
INSERT INTO esg_scores (company_id, environmental_score, social_score, governance_score, overall_score, date, data_source)
SELECT 
    c.id,
    ROUND((70 + (RANDOM() * 30))::numeric, 2) as environmental_score,
    ROUND((70 + (RANDOM() * 30))::numeric, 2) as social_score,
    ROUND((70 + (RANDOM() * 30))::numeric, 2) as governance_score,
    ROUND((70 + (RANDOM() * 30))::numeric, 2) as overall_score,
    CURRENT_DATE - INTERVAL '1 day' * (RANDOM() * 30)::int as date,
    'Sample Data' as data_source
FROM companies c
LIMIT 50;
