package models

import (
	"database/sql"
	"time"
)

// ESGTrend represents ESG score trends over time
type ESGTrend struct {
	CompanyID   int       `json:"company_id"`
	CompanyName string    `json:"company_name"`
	Date        time.Time `json:"date"`
	ESGScore    float64   `json:"esg_score"`
	EScore      float64   `json:"e_score"`
	SScore      float64   `json:"s_score"`
	GScore      float64   `json:"g_score"`
}

// SectorComparison represents sector-level ESG and financial comparisons
type SectorComparison struct {
	Sector          string  `json:"sector"`
	CompanyCount    int     `json:"company_count"`
	AvgESGScore     float64 `json:"avg_esg_score"`
	AvgPERatio      float64 `json:"avg_pe_ratio"`
	AvgMarketCap    float64 `json:"avg_market_cap"`
	TotalMarketCap  float64 `json:"total_market_cap"`
	BestESGCompany  string  `json:"best_esg_company"`
	WorstESGCompany string  `json:"worst_esg_company"`
}

// FinancialComparison represents financial performance comparisons
type FinancialComparison struct {
	CompanyID      int     `json:"company_id"`
	CompanyName    string  `json:"company_name"`
	CurrentPrice   float64 `json:"current_price"`
	PriceChange    float64 `json:"price_change"`
	PriceChangePct float64 `json:"price_change_pct"`
	MarketCap      float64 `json:"market_cap"`
	PERatio        float64 `json:"pe_ratio"`
	ESGScore       float64 `json:"esg_score"`
	ESGPercentile  float64 `json:"esg_percentile"`
}

// PerformanceMetric represents performance metrics for analysis
type PerformanceMetric struct {
	CompanyID   int       `json:"company_id"`
	CompanyName string    `json:"company_name"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Rank        int       `json:"rank"`
	TotalCount  int       `json:"total_count"`
	Percentile  float64   `json:"percentile"`
	Date        time.Time `json:"date"`
}

// AnalyticsRepository handles complex analytical database operations
type AnalyticsRepository struct {
	db *sql.DB
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// GetESGTrends retrieves ESG score trends for a company
func (r *AnalyticsRepository) GetESGTrends(companyID int, days int) ([]ESGTrend, error) {
	query := `
		SELECT es.company_id, c.name as company_name, es.score_date, es.overall_score, es.environmental_score, es.social_score, es.governance_score
		FROM esg_scores es
		JOIN companies c ON es.company_id = c.id
		WHERE es.company_id = $1
		ORDER BY es.score_date DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, companyID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []ESGTrend
	for rows.Next() {
		var trend ESGTrend
		err := rows.Scan(
			&trend.CompanyID, &trend.CompanyName, &trend.Date,
			&trend.ESGScore, &trend.EScore, &trend.SScore, &trend.GScore,
		)
		if err != nil {
			return nil, err
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

// GetSectorComparisons retrieves sector-level ESG and financial comparisons
func (r *AnalyticsRepository) GetSectorComparisons() ([]SectorComparison, error) {
	query := `
		WITH latest_esg AS (
			SELECT DISTINCT ON (company_id) company_id, overall_score, score_date
			FROM esg_scores
			ORDER BY company_id, score_date DESC
		),
		latest_financial AS (
			SELECT DISTINCT ON (company_id) company_id, market_cap, pe_ratio
			FROM financial_indicators
			ORDER BY company_id, date DESC
		),
		sector_stats AS (
			SELECT 
				c.sector,
				COUNT(c.id) as company_count,
				AVG(le.overall_score) as avg_esg_score,
				AVG(lf.pe_ratio) as avg_pe_ratio,
				AVG(lf.market_cap) as avg_market_cap,
				SUM(lf.market_cap) as total_market_cap,
				MAX(le.overall_score) as max_esg_score,
				MIN(le.overall_score) as min_esg_score
			FROM companies c
			LEFT JOIN latest_esg le ON c.id = le.company_id
			LEFT JOIN latest_financial lf ON c.id = lf.company_id
			GROUP BY c.sector
		),
		best_esg AS (
			SELECT DISTINCT ON (c.sector) c.sector, c.name as company_name, le.overall_score
			FROM companies c
			JOIN latest_esg le ON c.id = le.company_id
			ORDER BY c.sector, le.overall_score DESC
		),
		worst_esg AS (
			SELECT DISTINCT ON (c.sector) c.sector, c.name as company_name, le.overall_score
			FROM companies c
			JOIN latest_esg le ON c.id = le.company_id
			ORDER BY c.sector, le.overall_score ASC
		)
		SELECT 
			ss.sector,
			ss.company_count,
			ROUND(ss.avg_esg_score::numeric, 2) as avg_esg_score,
			ROUND(ss.avg_pe_ratio::numeric, 2) as avg_pe_ratio,
			ROUND(ss.avg_market_cap::numeric, 2) as avg_market_cap,
			ROUND(ss.total_market_cap::numeric, 2) as total_market_cap,
			be.company_name as best_esg_company,
			we.company_name as worst_esg_company
		FROM sector_stats ss
		LEFT JOIN best_esg be ON ss.sector = be.sector
		LEFT JOIN worst_esg we ON ss.sector = we.sector
		ORDER BY ss.avg_esg_score DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comparisons []SectorComparison
	for rows.Next() {
		var comp SectorComparison
		err := rows.Scan(
			&comp.Sector, &comp.CompanyCount, &comp.AvgESGScore,
			&comp.AvgPERatio, &comp.AvgMarketCap, &comp.TotalMarketCap,
			&comp.BestESGCompany, &comp.WorstESGCompany,
		)
		if err != nil {
			return nil, err
		}
		comparisons = append(comparisons, comp)
	}

	return comparisons, nil
}

// GetFinancialComparisons retrieves financial performance comparisons
func (r *AnalyticsRepository) GetFinancialComparisons(limit int) ([]FinancialComparison, error) {
	query := `
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
			SELECT 
				company_id,
				overall_score,
				PERCENT_RANK() OVER (ORDER BY overall_score) * 100 as percentile
			FROM latest_esg
		)
		SELECT 
			c.id as company_id,
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
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comparisons []FinancialComparison
	for rows.Next() {
		var comp FinancialComparison
		err := rows.Scan(
			&comp.CompanyID, &comp.CompanyName, &comp.CurrentPrice,
			&comp.PriceChange, &comp.PriceChangePct, &comp.MarketCap,
			&comp.PERatio, &comp.ESGScore, &comp.ESGPercentile,
		)
		if err != nil {
			return nil, err
		}
		comparisons = append(comparisons, comp)
	}

	return comparisons, nil
}

// GetTopPerformers retrieves top performing companies by various metrics
func (r *AnalyticsRepository) GetTopPerformers(metric string, limit int) ([]PerformanceMetric, error) {
	var query string

	switch metric {
	case "esg_score":
		query = `
			WITH latest_esg AS (
				SELECT DISTINCT ON (company_id) company_id, overall_score, score_date
				FROM esg_scores
				ORDER BY company_id, score_date DESC
			),
			ranked AS (
				SELECT 
					c.id as company_id,
					c.name as company_name,
					le.overall_score as value,
					le.score_date,
					RANK() OVER (ORDER BY le.overall_score DESC) as rank,
					COUNT(*) OVER () as total_count,
					PERCENT_RANK() OVER (ORDER BY le.overall_score) * 100 as percentile
				FROM companies c
				JOIN latest_esg le ON c.id = le.company_id
			)
			SELECT company_id, company_name, 'ESG Score' as metric, value, rank, total_count, percentile, score_date
			FROM ranked
			ORDER BY rank
			LIMIT $1
		`
	case "market_cap":
		query = `
			WITH latest_financial AS (
				SELECT DISTINCT ON (company_id) company_id, market_cap, date
				FROM financial_indicators
				ORDER BY company_id, date DESC
			),
			ranked AS (
				SELECT 
					c.id as company_id,
					c.name as company_name,
					lf.market_cap as value,
					lf.date,
					RANK() OVER (ORDER BY lf.market_cap DESC) as rank,
					COUNT(*) OVER () as total_count,
					PERCENT_RANK() OVER (ORDER BY lf.market_cap) * 100 as percentile
				FROM companies c
				JOIN latest_financial lf ON c.id = lf.company_id
			)
			SELECT company_id, company_name, 'Market Cap' as metric, value, rank, total_count, percentile, date
			FROM ranked
			ORDER BY rank
			LIMIT $1
		`
	case "pe_ratio":
		query = `
			WITH latest_financial AS (
				SELECT DISTINCT ON (company_id) company_id, pe_ratio, date
				FROM financial_indicators
				WHERE pe_ratio > 0
				ORDER BY company_id, date DESC
			),
			ranked AS (
				SELECT 
					c.id as company_id,
					c.name as company_name,
					lf.pe_ratio as value,
					lf.date,
					RANK() OVER (ORDER BY lf.pe_ratio ASC) as rank,
					COUNT(*) OVER () as total_count,
					PERCENT_RANK() OVER (ORDER BY lf.pe_ratio) * 100 as percentile
				FROM companies c
				JOIN latest_financial lf ON c.id = lf.company_id
			)
			SELECT company_id, company_name, 'P/E Ratio' as metric, value, rank, total_count, percentile, date
			FROM ranked
			ORDER BY rank
			LIMIT $1
		`
	default:
		return nil, sql.ErrNoRows
	}

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []PerformanceMetric
	for rows.Next() {
		var metric PerformanceMetric
		err := rows.Scan(
			&metric.CompanyID, &metric.CompanyName, &metric.Metric,
			&metric.Value, &metric.Rank, &metric.TotalCount,
			&metric.Percentile, &metric.Date,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// GetESGvsFinancialCorrelation calculates correlation between ESG scores and financial metrics
func (r *AnalyticsRepository) GetESGvsFinancialCorrelation() (map[string]interface{}, error) {
	query := `
		WITH latest_data AS (
			SELECT 
				c.id as company_id,
				c.name as company_name,
				le.overall_score,
				lf.market_cap,
				lf.pe_ratio,
				lf.return_on_equity,
				lf.profit_margin
			FROM companies c
			LEFT JOIN (
				SELECT DISTINCT ON (company_id) company_id, overall_score
				FROM esg_scores
				ORDER BY company_id, score_date DESC
			) le ON c.id = le.company_id
			LEFT JOIN (
				SELECT DISTINCT ON (company_id) company_id, market_cap, pe_ratio, return_on_equity, profit_margin
				FROM financial_indicators
				ORDER BY company_id, date DESC
			) lf ON c.id = lf.company_id
			WHERE le.overall_score IS NOT NULL AND lf.market_cap IS NOT NULL
		),
		correlations AS (
			SELECT 
				COUNT(*) as sample_size,
				AVG(overall_score) as avg_esg_score,
				AVG(market_cap) as avg_market_cap,
				AVG(pe_ratio) as avg_pe_ratio,
				AVG(return_on_equity) as avg_roe,
				AVG(profit_margin) as avg_profit_margin,
				CORR(overall_score, market_cap) as esg_market_cap_corr,
				CORR(overall_score, pe_ratio) as esg_pe_corr,
				CORR(overall_score, return_on_equity) as esg_roe_corr,
				CORR(overall_score, profit_margin) as esg_profit_corr
			FROM latest_data
		)
		SELECT 
			sample_size,
			ROUND(avg_esg_score::numeric, 2) as avg_esg_score,
			ROUND(avg_market_cap::numeric, 2) as avg_market_cap,
			ROUND(avg_pe_ratio::numeric, 2) as avg_pe_ratio,
			ROUND(avg_roe::numeric, 4) as avg_roe,
			ROUND(avg_profit_margin::numeric, 4) as avg_profit_margin,
			ROUND(esg_market_cap_corr::numeric, 4) as esg_market_cap_corr,
			ROUND(esg_pe_corr::numeric, 4) as esg_pe_corr,
			ROUND(esg_roe_corr::numeric, 4) as esg_roe_corr,
			ROUND(esg_profit_corr::numeric, 4) as esg_profit_corr
		FROM correlations
	`

	var sampleSize int
	var avgESGScore, avgMarketCap, avgPERatio, avgROE, avgProfitMargin float64
	var esgMarketCapCorr, esgPECorr, esgROECorr, esgProfitCorr float64

	err := r.db.QueryRow(query).Scan(
		&sampleSize, &avgESGScore, &avgMarketCap, &avgPERatio, &avgROE, &avgProfitMargin,
		&esgMarketCapCorr, &esgPECorr, &esgROECorr, &esgProfitCorr,
	)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"sample_size":         sampleSize,
		"avg_esg_score":       avgESGScore,
		"avg_market_cap":      avgMarketCap,
		"avg_pe_ratio":        avgPERatio,
		"avg_roe":             avgROE,
		"avg_profit_margin":   avgProfitMargin,
		"esg_market_cap_corr": esgMarketCapCorr,
		"esg_pe_corr":         esgPECorr,
		"esg_roe_corr":        esgROECorr,
		"esg_profit_corr":     esgProfitCorr,
	}

	return result, nil
}
