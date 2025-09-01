-- Sample data for EthosView
-- Week 2: Seed data for testing

-- Insert sample companies
INSERT INTO companies (name, symbol, sector, industry, country, market_cap) VALUES
('Apple Inc.', 'AAPL', 'Technology', 'Consumer Electronics', 'United States', 3000000000000.00),
('Microsoft Corporation', 'MSFT', 'Technology', 'Software', 'United States', 2800000000000.00),
('Alphabet Inc.', 'GOOGL', 'Technology', 'Internet Services', 'United States', 1800000000000.00),
('Amazon.com Inc.', 'AMZN', 'Consumer Cyclical', 'Internet Retail', 'United States', 1600000000000.00),
('Tesla Inc.', 'TSLA', 'Consumer Cyclical', 'Auto Manufacturers', 'United States', 800000000000.00),
('Johnson & Johnson', 'JNJ', 'Healthcare', 'Drug Manufacturers', 'United States', 400000000000.00),
('Procter & Gamble Co.', 'PG', 'Consumer Defensive', 'Household & Personal Products', 'United States', 350000000000.00),
('Coca-Cola Company', 'KO', 'Consumer Defensive', 'Beverages', 'United States', 250000000000.00),
('Walmart Inc.', 'WMT', 'Consumer Defensive', 'Discount Stores', 'United States', 450000000000.00),
('JPMorgan Chase & Co.', 'JPM', 'Financial Services', 'Banks', 'United States', 500000000000.00)
ON CONFLICT (symbol) DO NOTHING;

-- Insert sample ESG scores
INSERT INTO esg_scores (company_id, environmental_score, social_score, governance_score, overall_score, score_date, data_source) VALUES
((SELECT id FROM companies WHERE symbol = 'AAPL'), 85.5, 78.2, 82.1, 82.1, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'MSFT'), 88.3, 81.5, 85.7, 85.2, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'GOOGL'), 82.1, 76.8, 79.4, 79.4, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'AMZN'), 75.6, 72.3, 68.9, 72.3, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'TSLA'), 92.4, 65.2, 58.7, 72.1, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'JNJ'), 78.9, 85.6, 88.2, 84.2, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'PG'), 81.2, 83.7, 86.5, 83.8, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'KO'), 76.4, 79.1, 82.8, 79.4, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'WMT'), 72.8, 81.5, 75.3, 76.5, '2024-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'JPM'), 68.9, 74.2, 71.6, 71.6, '2024-01-15', 'MSCI');

-- Insert historical ESG scores for trend analysis
INSERT INTO esg_scores (company_id, environmental_score, social_score, governance_score, overall_score, score_date, data_source) VALUES
((SELECT id FROM companies WHERE symbol = 'AAPL'), 83.2, 76.1, 80.5, 80.0, '2023-07-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'MSFT'), 86.1, 79.8, 83.2, 83.0, '2023-07-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'GOOGL'), 80.5, 74.9, 77.8, 77.7, '2023-07-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'AAPL'), 81.8, 74.3, 78.9, 78.3, '2023-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'MSFT'), 84.7, 77.6, 81.1, 81.1, '2023-01-15', 'MSCI'),
((SELECT id FROM companies WHERE symbol = 'GOOGL'), 78.9, 73.2, 76.1, 76.1, '2023-01-15', 'MSCI');
