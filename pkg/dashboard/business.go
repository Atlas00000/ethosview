package dashboard

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// BusinessDashboard handles business metrics collection and dashboard data
type BusinessDashboard struct {
	db    *sql.DB
	redis *redis.Client
}

// DashboardData represents the complete dashboard response
type DashboardData struct {
	Summary     BusinessSummary     `json:"summary"`
	ESGMetrics  ESGMetrics         `json:"esg_metrics"`
	Sectors     SectorMetrics      `json:"sectors"`
	Trends      TrendMetrics       `json:"trends"`
	Performance PerformanceMetrics `json:"performance"`
	LastUpdated time.Time          `json:"last_updated"`
}

// BusinessSummary provides high-level business metrics
type BusinessSummary struct {
	TotalCompanies       int     `json:"total_companies"`
	TotalESGScores       int     `json:"total_esg_scores"`
	AvgESGScore          float64 `json:"avg_esg_score"`
	ActiveSectors        int     `json:"active_sectors"`
	RecentScoreUpdates   int     `json:"recent_score_updates"`
	MarketCapTotal       float64 `json:"market_cap_total"`
	TopPerformingCompany string  `json:"top_performing_company"`
}

// ESGMetrics provides ESG-specific metrics
type ESGMetrics struct {
	OverallAverage      float64              `json:"overall_average"`
	EnvironmentalAvg    float64              `json:"environmental_avg"`
	SocialAvg           float64              `json:"social_avg"`
	GovernanceAvg       float64              `json:"governance_avg"`
	ScoreDistribution   ScoreDistribution    `json:"score_distribution"`
	TopPerformers       []TopPerformer       `json:"top_performers"`
	ImprovementTrends   []ImprovementTrend   `json:"improvement_trends"`
}

// SectorMetrics provides sector-specific analysis
type SectorMetrics struct {
	Distribution    map[string]int        `json:"distribution"`
	ESGAverages     map[string]float64    `json:"esg_averages"`
	MarketCapBySector map[string]float64 `json:"market_cap_by_sector"`
	TopSectors      []SectorRanking       `json:"top_sectors"`
}

// TrendMetrics provides trend analysis
type TrendMetrics struct {
	ScoreChanges     []ScoreChange     `json:"score_changes"`
	MonthlyGrowth    []MonthlyGrowth   `json:"monthly_growth"`
	SectorTrends     []SectorTrend     `json:"sector_trends"`
	PredictedGrowth  float64           `json:"predicted_growth"`
}

// PerformanceMetrics provides system performance insights
type PerformanceMetrics struct {
	DataFreshness      time.Duration `json:"data_freshness_minutes"`
	CacheHitRate       float64       `json:"cache_hit_rate"`
	QueryPerformance   float64       `json:"avg_query_time_ms"`
	DataCompleteness   float64       `json:"data_completeness_percent"`
}

// Supporting types
type ScoreDistribution struct {
	Excellent int `json:"excellent"` // 80-100
	Good      int `json:"good"`      // 60-79
	Average   int `json:"average"`   // 40-59
	Poor      int `json:"poor"`      // 0-39
}

type TopPerformer struct {
	CompanyName  string  `json:"company_name"`
	Symbol       string  `json:"symbol"`
	ESGScore     float64 `json:"esg_score"`
	Sector       string  `json:"sector"`
	MarketCap    float64 `json:"market_cap"`
}

type ImprovementTrend struct {
	CompanyName     string  `json:"company_name"`
	Symbol          string  `json:"symbol"`
	ScoreImprovement float64 `json:"score_improvement"`
	TimeFrame       string  `json:"time_frame"`
}

type SectorRanking struct {
	Sector      string  `json:"sector"`
	AvgESGScore float64 `json:"avg_esg_score"`
	CompanyCount int    `json:"company_count"`
	Rank        int     `json:"rank"`
}

type ScoreChange struct {
	Date        time.Time `json:"date"`
	AvgChange   float64   `json:"avg_change"`
	CompanyCount int      `json:"company_count"`
}

type MonthlyGrowth struct {
	Month           string  `json:"month"`
	NewCompanies    int     `json:"new_companies"`
	NewESGScores    int     `json:"new_esg_scores"`
	GrowthRate      float64 `json:"growth_rate"`
}

type SectorTrend struct {
	Sector         string  `json:"sector"`
	TrendDirection string  `json:"trend_direction"` // "up", "down", "stable"
	ChangePercent  float64 `json:"change_percent"`
}

// NewBusinessDashboard creates a new business dashboard instance
func NewBusinessDashboard(db *sql.DB, redis *redis.Client) *BusinessDashboard {
	return &BusinessDashboard{
		db:    db,
		redis: redis,
	}
}

// GetDashboardData retrieves comprehensive dashboard data
func (bd *BusinessDashboard) GetDashboardData() (*DashboardData, error) {
	// Check cache first
	cached, err := bd.getCachedDashboard()
	if err == nil && cached != nil {
		return cached, nil
	}

	// Generate fresh dashboard data
	dashboard := &DashboardData{
		LastUpdated: time.Now().UTC(),
	}

	// Collect all metrics
	if err := bd.collectSummary(&dashboard.Summary); err != nil {
		return nil, fmt.Errorf("failed to collect summary: %v", err)
	}

	if err := bd.collectESGMetrics(&dashboard.ESGMetrics); err != nil {
		return nil, fmt.Errorf("failed to collect ESG metrics: %v", err)
	}

	if err := bd.collectSectorMetrics(&dashboard.Sectors); err != nil {
		return nil, fmt.Errorf("failed to collect sector metrics: %v", err)
	}

	if err := bd.collectTrendMetrics(&dashboard.Trends); err != nil {
		return nil, fmt.Errorf("failed to collect trend metrics: %v", err)
	}

	if err := bd.collectPerformanceMetrics(&dashboard.Performance); err != nil {
		return nil, fmt.Errorf("failed to collect performance metrics: %v", err)
	}

	// Cache the result
	bd.cacheDashboard(dashboard)

	return dashboard, nil
}

// collectSummary gathers business summary metrics
func (bd *BusinessDashboard) collectSummary(summary *BusinessSummary) error {
	// Total companies
	err := bd.db.QueryRow("SELECT COUNT(*) FROM companies").Scan(&summary.TotalCompanies)
	if err != nil {
		return err
	}

	// Total ESG scores
	err = bd.db.QueryRow("SELECT COUNT(*) FROM esg_scores").Scan(&summary.TotalESGScores)
	if err != nil {
		return err
	}

	// Average ESG score
	err = bd.db.QueryRow("SELECT AVG(overall_score) FROM esg_scores WHERE overall_score IS NOT NULL").Scan(&summary.AvgESGScore)
	if err != nil {
		return err
	}

	// Active sectors
	err = bd.db.QueryRow("SELECT COUNT(DISTINCT sector) FROM companies WHERE sector IS NOT NULL AND sector != ''").Scan(&summary.ActiveSectors)
	if err != nil {
		return err
	}

	// Recent score updates (last 7 days)
	err = bd.db.QueryRow("SELECT COUNT(*) FROM esg_scores WHERE created_at > NOW() - INTERVAL '7 days'").Scan(&summary.RecentScoreUpdates)
	if err != nil {
		return err
	}

	// Total market cap
	err = bd.db.QueryRow("SELECT COALESCE(SUM(market_cap), 0) FROM companies WHERE market_cap > 0").Scan(&summary.MarketCapTotal)
	if err != nil {
		return err
	}

	// Top performing company
	err = bd.db.QueryRow(`
		SELECT c.name 
		FROM companies c 
		JOIN esg_scores e ON c.id = e.company_id 
		WHERE e.overall_score IS NOT NULL 
		ORDER BY e.overall_score DESC 
		LIMIT 1
	`).Scan(&summary.TopPerformingCompany)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

// collectESGMetrics gathers ESG-specific metrics
func (bd *BusinessDashboard) collectESGMetrics(metrics *ESGMetrics) error {
	// Average scores
	err := bd.db.QueryRow(`
		SELECT 
			AVG(overall_score),
			AVG(environmental_score),
			AVG(social_score),
			AVG(governance_score)
		FROM esg_scores 
		WHERE overall_score IS NOT NULL
	`).Scan(&metrics.OverallAverage, &metrics.EnvironmentalAvg, &metrics.SocialAvg, &metrics.GovernanceAvg)
	if err != nil {
		return err
	}

	// Score distribution
	err = bd.collectScoreDistribution(&metrics.ScoreDistribution)
	if err != nil {
		return err
	}

	// Top performers
	err = bd.collectTopPerformers(&metrics.TopPerformers)
	if err != nil {
		return err
	}

	return nil
}

// collectSectorMetrics gathers sector-specific metrics
func (bd *BusinessDashboard) collectSectorMetrics(metrics *SectorMetrics) error {
	metrics.Distribution = make(map[string]int)
	metrics.ESGAverages = make(map[string]float64)
	metrics.MarketCapBySector = make(map[string]float64)

	// Sector distribution
	rows, err := bd.db.Query("SELECT sector, COUNT(*) FROM companies WHERE sector IS NOT NULL GROUP BY sector")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sector string
		var count int
		if err := rows.Scan(&sector, &count); err != nil {
			continue
		}
		metrics.Distribution[sector] = count
	}

	// ESG averages by sector
	rows, err = bd.db.Query(`
		SELECT c.sector, AVG(e.overall_score)
		FROM companies c
		JOIN esg_scores e ON c.id = e.company_id
		WHERE c.sector IS NOT NULL AND e.overall_score IS NOT NULL
		GROUP BY c.sector
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sector string
		var avgScore float64
		if err := rows.Scan(&sector, &avgScore); err != nil {
			continue
		}
		metrics.ESGAverages[sector] = avgScore
	}

	return nil
}

// collectTrendMetrics gathers trend analysis data
func (bd *BusinessDashboard) collectTrendMetrics(metrics *TrendMetrics) error {
	// Monthly growth for last 6 months
	rows, err := bd.db.Query(`
		SELECT 
			DATE_TRUNC('month', created_at) as month,
			COUNT(DISTINCT company_id) as companies,
			COUNT(*) as scores
		FROM esg_scores 
		WHERE created_at > NOW() - INTERVAL '6 months'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY month DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var month time.Time
		var companies, scores int
		if err := rows.Scan(&month, &companies, &scores); err != nil {
			continue
		}

		growth := MonthlyGrowth{
			Month:        month.Format("2006-01"),
			NewCompanies: companies,
			NewESGScores: scores,
			GrowthRate:   float64(scores) / float64(companies), // Simple metric
		}
		metrics.MonthlyGrowth = append(metrics.MonthlyGrowth, growth)
	}

	return nil
}

// collectPerformanceMetrics gathers system performance metrics
func (bd *BusinessDashboard) collectPerformanceMetrics(metrics *PerformanceMetrics) error {
	// Data freshness - check latest ESG score update
	var latestUpdate time.Time
	err := bd.db.QueryRow("SELECT MAX(created_at) FROM esg_scores").Scan(&latestUpdate)
	if err == nil {
		metrics.DataFreshness = time.Since(latestUpdate)
	}

	// Cache hit rate (simplified)
	if bd.redis != nil {
		ctx := context.Background()
		info, err := bd.redis.Info(ctx, "stats").Result()
		if err == nil {
			metrics.CacheHitRate = bd.parseCacheHitRate(info)
		}
	}

	// Data completeness
	var totalCompanies, companiesWithScores int
	bd.db.QueryRow("SELECT COUNT(*) FROM companies").Scan(&totalCompanies)
	bd.db.QueryRow("SELECT COUNT(DISTINCT company_id) FROM esg_scores").Scan(&companiesWithScores)
	
	if totalCompanies > 0 {
		metrics.DataCompleteness = float64(companiesWithScores) / float64(totalCompanies) * 100
	}

	return nil
}

// Helper methods

func (bd *BusinessDashboard) collectScoreDistribution(dist *ScoreDistribution) error {
	err := bd.db.QueryRow(`
		SELECT 
			COUNT(CASE WHEN overall_score >= 80 THEN 1 END) as excellent,
			COUNT(CASE WHEN overall_score >= 60 AND overall_score < 80 THEN 1 END) as good,
			COUNT(CASE WHEN overall_score >= 40 AND overall_score < 60 THEN 1 END) as average,
			COUNT(CASE WHEN overall_score < 40 THEN 1 END) as poor
		FROM esg_scores 
		WHERE overall_score IS NOT NULL
	`).Scan(&dist.Excellent, &dist.Good, &dist.Average, &dist.Poor)
	
	return err
}

func (bd *BusinessDashboard) collectTopPerformers(performers *[]TopPerformer) error {
	rows, err := bd.db.Query(`
		SELECT c.name, c.symbol, e.overall_score, c.sector, c.market_cap
		FROM companies c
		JOIN esg_scores e ON c.id = e.company_id
		WHERE e.overall_score IS NOT NULL
		ORDER BY e.overall_score DESC
		LIMIT 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var performer TopPerformer
		err := rows.Scan(&performer.CompanyName, &performer.Symbol, &performer.ESGScore, &performer.Sector, &performer.MarketCap)
		if err != nil {
			continue
		}
		*performers = append(*performers, performer)
	}

	return nil
}

func (bd *BusinessDashboard) getCachedDashboard() (*DashboardData, error) {
	if bd.redis == nil {
		return nil, fmt.Errorf("redis not available")
	}

	ctx := context.Background()
	data, err := bd.redis.Get(ctx, "dashboard:business:data").Result()
	if err != nil {
		return nil, err
	}

	var dashboard DashboardData
	err = json.Unmarshal([]byte(data), &dashboard)
	if err != nil {
		return nil, err
	}

	// Check if cache is too old (15 minutes)
	if time.Since(dashboard.LastUpdated) > 15*time.Minute {
		return nil, fmt.Errorf("cache expired")
	}

	return &dashboard, nil
}

func (bd *BusinessDashboard) cacheDashboard(dashboard *DashboardData) {
	if bd.redis == nil {
		return
	}

	ctx := context.Background()
	data, err := json.Marshal(dashboard)
	if err != nil {
		return
	}

	// Cache for 15 minutes
	bd.redis.Set(ctx, "dashboard:business:data", data, 15*time.Minute)
}

func (bd *BusinessDashboard) parseCacheHitRate(info string) float64 {
	// Simplified cache hit rate parsing
	// In production, implement proper Redis info parsing
	return 85.0 // Placeholder
}
