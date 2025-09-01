package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ethosview-backend/internal/models"

	"github.com/redis/go-redis/v9"
)

// CacheWarmer handles cache warming operations
type CacheWarmer struct {
	redis *redis.Client
	db    *sql.DB
}

// NewCacheWarmer creates a new cache warmer instance
func NewCacheWarmer(redis *redis.Client, db *sql.DB) *CacheWarmer {
	return &CacheWarmer{
		redis: redis,
		db:    db,
	}
}

// WarmCache performs cache warming for frequently accessed data
func (cw *CacheWarmer) WarmCache() error {
	log.Println("🔥 Starting cache warming...")

	// Warm company data
	if err := cw.warmCompanies(); err != nil {
		log.Printf("Error warming companies: %v", err)
	}

	// Warm ESG scores
	if err := cw.warmESGScores(); err != nil {
		log.Printf("Error warming ESG scores: %v", err)
	}

	// Warm sector data
	if err := cw.warmSectors(); err != nil {
		log.Printf("Error warming sectors: %v", err)
	}

	// Warm analytics data
	if err := cw.warmAnalytics(); err != nil {
		log.Printf("Error warming analytics: %v", err)
	}

	log.Println("✅ Cache warming completed")
	return nil
}

// warmCompanies warms company-related cache
func (cw *CacheWarmer) warmCompanies() error {
	ctx := context.Background()

	// Warm all companies
	companies, err := cw.getCompanies()
	if err != nil {
		return err
	}

	// Cache companies list
	companiesData, _ := json.Marshal(companies)
	cw.redis.Set(ctx, "cache:companies:all", companiesData, 30*time.Minute)

	// Cache companies by sector
	sectors := make(map[string][]models.Company)
	for _, company := range companies {
		if company.Sector != "" {
			sectors[company.Sector] = append(sectors[company.Sector], company)
		}
	}

	for sector, sectorCompanies := range sectors {
		sectorData, _ := json.Marshal(sectorCompanies)
		cw.redis.Set(ctx, fmt.Sprintf("cache:companies:sector:%s", sector), sectorData, 30*time.Minute)
	}

	// Cache individual companies
	for _, company := range companies {
		companyData, _ := json.Marshal(company)
		cw.redis.Set(ctx, fmt.Sprintf("cache:company:%d", company.ID), companyData, 30*time.Minute)
		cw.redis.Set(ctx, fmt.Sprintf("cache:company:symbol:%s", company.Symbol), companyData, 30*time.Minute)
	}

	log.Printf("Warmed %d companies", len(companies))
	return nil
}

// warmESGScores warms ESG score-related cache
func (cw *CacheWarmer) warmESGScores() error {
	ctx := context.Background()

	// Get latest ESG scores for all companies
	query := `
		SELECT DISTINCT ON (company_id) 
			es.id, es.company_id, es.environmental_score, es.social_score, 
			es.governance_score, es.overall_score, es.score_date, es.data_source
		FROM esg_scores es
		ORDER BY company_id, score_date DESC
	`

	rows, err := cw.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var scores []models.ESGScore
	for rows.Next() {
		var score models.ESGScore
		err := rows.Scan(
			&score.ID, &score.CompanyID, &score.EnvironmentalScore,
			&score.SocialScore, &score.GovernanceScore, &score.OverallScore,
			&score.ScoreDate, &score.DataSource,
		)
		if err != nil {
			continue
		}
		scores = append(scores, score)
	}

	// Cache all ESG scores
	scoresData, _ := json.Marshal(scores)
	cw.redis.Set(ctx, "cache:esg:scores:all", scoresData, 15*time.Minute)

	// Cache individual company ESG scores
	for _, score := range scores {
		scoreData, _ := json.Marshal(score)
		cw.redis.Set(ctx, fmt.Sprintf("cache:esg:company:%d:latest", score.CompanyID), scoreData, 15*time.Minute)
	}

	// Cache top performers
	cw.warmTopPerformers(scores)

	log.Printf("Warmed %d ESG scores", len(scores))
	return nil
}

// warmTopPerformers caches top ESG performers
func (cw *CacheWarmer) warmTopPerformers(scores []models.ESGScore) {
	ctx := context.Background()

	// Top overall performers
	topOverall := cw.getTopPerformers(scores, "overall", 10)
	topOverallData, _ := json.Marshal(topOverall)
	cw.redis.Set(ctx, "cache:esg:top:overall", topOverallData, 15*time.Minute)

	// Top environmental performers
	topEnvironmental := cw.getTopPerformers(scores, "environmental", 10)
	topEnvironmentalData, _ := json.Marshal(topEnvironmental)
	cw.redis.Set(ctx, "cache:esg:top:environmental", topEnvironmentalData, 15*time.Minute)

	// Top social performers
	topSocial := cw.getTopPerformers(scores, "social", 10)
	topSocialData, _ := json.Marshal(topSocial)
	cw.redis.Set(ctx, "cache:esg:top:social", topSocialData, 15*time.Minute)

	// Top governance performers
	topGovernance := cw.getTopPerformers(scores, "governance", 10)
	topGovernanceData, _ := json.Marshal(topGovernance)
	cw.redis.Set(ctx, "cache:esg:top:governance", topGovernanceData, 15*time.Minute)
}

// warmSectors warms sector-related cache
func (cw *CacheWarmer) warmSectors() error {
	ctx := context.Background()

	// Get unique sectors
	query := `SELECT DISTINCT sector FROM companies WHERE sector IS NOT NULL AND sector != ''`
	rows, err := cw.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sectors []string
	for rows.Next() {
		var sector string
		if err := rows.Scan(&sector); err != nil {
			continue
		}
		sectors = append(sectors, sector)
	}

	// Cache sectors list
	sectorsData, _ := json.Marshal(sectors)
	cw.redis.Set(ctx, "cache:sectors:all", sectorsData, 1*time.Hour)

	log.Printf("Warmed %d sectors", len(sectors))
	return nil
}

// warmAnalytics warms analytics-related cache
func (cw *CacheWarmer) warmAnalytics() error {
	ctx := context.Background()

	// Cache analytics summary
	summary := map[string]interface{}{
		"total_companies":  cw.getCompanyCount(),
		"total_esg_scores": cw.getESGScoreCount(),
		"last_updated":     time.Now().UTC(),
	}

	summaryData, _ := json.Marshal(summary)
	cw.redis.Set(ctx, "cache:analytics:summary", summaryData, 10*time.Minute)

	log.Println("Warmed analytics data")
	return nil
}

// getCompanies retrieves all companies from database
func (cw *CacheWarmer) getCompanies() ([]models.Company, error) {
	query := `SELECT id, name, symbol, sector, industry, country, market_cap, created_at, updated_at FROM companies`
	rows, err := cw.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var company models.Company
		err := rows.Scan(
			&company.ID, &company.Name, &company.Symbol, &company.Sector,
			&company.Industry, &company.Country, &company.MarketCap,
			&company.CreatedAt, &company.UpdatedAt,
		)
		if err != nil {
			continue
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// getTopPerformers gets top performers by metric
func (cw *CacheWarmer) getTopPerformers(scores []models.ESGScore, metric string, limit int) []models.ESGScore {
	// This is a simplified implementation
	// In a real scenario, you'd want to join with company data
	if len(scores) <= limit {
		return scores
	}
	return scores[:limit]
}

// getCompanyCount gets total company count
func (cw *CacheWarmer) getCompanyCount() int {
	var count int
	cw.db.QueryRow("SELECT COUNT(*) FROM companies").Scan(&count)
	return count
}

// getESGScoreCount gets total ESG score count
func (cw *CacheWarmer) getESGScoreCount() int {
	var count int
	cw.db.QueryRow("SELECT COUNT(*) FROM esg_scores").Scan(&count)
	return count
}

// StartCacheWarming starts periodic cache warming
func (cw *CacheWarmer) StartCacheWarming(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Initial warming
		cw.WarmCache()

		for range ticker.C {
			cw.WarmCache()
		}
	}()
}
