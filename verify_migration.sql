-- Verification queries to check if migration was successful
-- Run these in Supabase SQL Editor after running the main migration

-- Check if tables exist
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('users', 'companies', 'esg_scores');

-- Check companies data
SELECT COUNT(*) as company_count FROM companies;
SELECT name, symbol, sector FROM companies LIMIT 5;

-- Check ESG scores data
SELECT COUNT(*) as esg_score_count FROM esg_scores;
SELECT 
    c.name as company_name,
    c.symbol,
    es.overall_score,
    es.environmental_score,
    es.social_score,
    es.governance_score,
    es.date
FROM esg_scores es
JOIN companies c ON es.company_id = c.id
ORDER BY es.overall_score DESC
LIMIT 5;

-- Check indexes
SELECT indexname, tablename 
FROM pg_indexes 
WHERE schemaname = 'public' 
AND tablename IN ('users', 'companies', 'esg_scores');

-- Check RLS policies
SELECT schemaname, tablename, policyname, permissive, roles, cmd, qual
FROM pg_policies 
WHERE schemaname = 'public';
