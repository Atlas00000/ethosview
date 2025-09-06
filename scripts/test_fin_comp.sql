WITH latest_esg AS (
  SELECT DISTINCT ON (company_id) company_id, overall_score
  FROM esg_scores
  ORDER BY company_id, score_date DESC
),
latest_price AS (
  SELECT DISTINCT ON (company_id) company_id, close_price
  FROM stock_prices
  ORDER BY company_id, date DESC
),
latest_financial AS (
  SELECT DISTINCT ON (company_id) company_id, market_cap, pe_ratio
  FROM financial_indicators
  ORDER BY company_id, date DESC
),
esg_percentiles AS (
  SELECT company_id, overall_score,
         PERCENT_RANK() OVER (ORDER BY overall_score) * 100 as percentile
  FROM latest_esg
)
SELECT c.id as company_id,
       c.name as company_name,
       COALESCE(lp.close_price, 0) as current_price,
       0 as price_change,
       0 as price_change_pct,
       COALESCE(lf.market_cap, 0) as market_cap,
       COALESCE(lf.pe_ratio, 0) as pe_ratio,
       COALESCE(le.overall_score, 0) as overall_score,
       COALESCE(ep.percentile, 0) as esg_percentile
FROM companies c
LEFT JOIN latest_price lp ON c.id = lp.company_id
LEFT JOIN latest_financial lf ON c.id = lf.company_id
LEFT JOIN latest_esg le ON c.id = le.company_id
LEFT JOIN esg_percentiles ep ON c.id = ep.company_id
ORDER BY COALESCE(le.overall_score, 0) DESC
LIMIT 5;
