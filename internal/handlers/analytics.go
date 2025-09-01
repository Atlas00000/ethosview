package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"ethosview-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsRepo *models.AnalyticsRepository
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(db *sql.DB) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsRepo: models.NewAnalyticsRepository(db),
	}
}

// GetESGTrends retrieves ESG score trends for a company
func (h *AnalyticsHandler) GetESGTrends(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		days = 30
	}

	trends, err := h.analyticsRepo.GetESGTrends(companyID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ESG trends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_id": companyID,
		"trends":     trends,
		"count":      len(trends),
		"days":       days,
	})
}

// GetSectorComparisons retrieves sector-level ESG and financial comparisons
func (h *AnalyticsHandler) GetSectorComparisons(c *gin.Context) {
	comparisons, err := h.analyticsRepo.GetSectorComparisons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sector comparisons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sector_comparisons": comparisons,
		"count":              len(comparisons),
	})
}

// GetFinancialComparisons retrieves financial performance comparisons
func (h *AnalyticsHandler) GetFinancialComparisons(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	comparisons, err := h.analyticsRepo.GetFinancialComparisons(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve financial comparisons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"financial_comparisons": comparisons,
		"count":                 len(comparisons),
		"limit":                 limit,
	})
}

// GetTopPerformers retrieves top performing companies by various metrics
func (h *AnalyticsHandler) GetTopPerformers(c *gin.Context) {
	metric := c.Param("metric")
	if metric == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric parameter is required"})
		return
	}

	// Validate metric
	validMetrics := map[string]bool{
		"esg_score":  true,
		"market_cap": true,
		"pe_ratio":   true,
	}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric. Valid metrics: esg_score, market_cap, pe_ratio"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	performers, err := h.analyticsRepo.GetTopPerformers(metric, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No data found for the specified metric"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top performers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metric":         metric,
		"top_performers": performers,
		"count":          len(performers),
		"limit":          limit,
	})
}

// GetESGvsFinancialCorrelation retrieves correlation analysis between ESG and financial metrics
func (h *AnalyticsHandler) GetESGvsFinancialCorrelation(c *gin.Context) {
	correlation, err := h.analyticsRepo.GetESGvsFinancialCorrelation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate correlations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"correlation_analysis": correlation,
	})
}

// GetAnalyticsSummary retrieves a comprehensive analytics summary
func (h *AnalyticsHandler) GetAnalyticsSummary(c *gin.Context) {
	// Get sector comparisons
	sectorComparisons, err := h.analyticsRepo.GetSectorComparisons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics summary"})
		return
	}

	// Get top ESG performers
	topESG, err := h.analyticsRepo.GetTopPerformers("esg_score", 5)
	if err != nil {
		topESG = []models.PerformanceMetric{}
	}

	// Get top market cap performers
	topMarketCap, err := h.analyticsRepo.GetTopPerformers("market_cap", 5)
	if err != nil {
		topMarketCap = []models.PerformanceMetric{}
	}

	// Get correlation analysis
	correlation, err := h.analyticsRepo.GetESGvsFinancialCorrelation()
	if err != nil {
		correlation = map[string]interface{}{}
	}

	// Calculate summary statistics
	totalCompanies := 0
	totalSectors := len(sectorComparisons)
	avgESGScore := 0.0

	for _, sector := range sectorComparisons {
		totalCompanies += sector.CompanyCount
		avgESGScore += sector.AvgESGScore * float64(sector.CompanyCount)
	}

	if totalCompanies > 0 {
		avgESGScore = avgESGScore / float64(totalCompanies)
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"total_companies": totalCompanies,
			"total_sectors":   totalSectors,
			"avg_esg_score":   avgESGScore,
		},
		"sector_comparisons": sectorComparisons,
		"top_esg_performers": topESG,
		"top_market_cap":     topMarketCap,
		"correlation":        correlation,
	})
}
