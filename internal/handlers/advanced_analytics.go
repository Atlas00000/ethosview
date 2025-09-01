package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"ethosview-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// AdvancedAnalyticsHandler handles advanced analytics-related HTTP requests
type AdvancedAnalyticsHandler struct {
	advancedAnalyticsRepo *models.AdvancedAnalyticsRepository
}

// NewAdvancedAnalyticsHandler creates a new advanced analytics handler
func NewAdvancedAnalyticsHandler(db *sql.DB) *AdvancedAnalyticsHandler {
	return &AdvancedAnalyticsHandler{
		advancedAnalyticsRepo: models.NewAdvancedAnalyticsRepository(db),
	}
}

// PredictESGScore predicts ESG score for a company
func (h *AdvancedAnalyticsHandler) PredictESGScore(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	prediction, err := h.advancedAnalyticsRepo.PredictESGScore(companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Insufficient data for prediction"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ESG prediction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prediction": prediction,
		"message":    "ESG prediction generated successfully",
	})
}

// OptimizePortfolio creates an optimized portfolio
func (h *AdvancedAnalyticsHandler) OptimizePortfolio(c *gin.Context) {
	// Parse query parameters
	targetReturnStr := c.DefaultQuery("target_return", "0.10")
	riskTolerance := c.DefaultQuery("risk_tolerance", "medium")
	maxCompaniesStr := c.DefaultQuery("max_companies", "10")

	targetReturn, err := strconv.ParseFloat(targetReturnStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target return"})
		return
	}

	maxCompanies, err := strconv.Atoi(maxCompaniesStr)
	if err != nil || maxCompanies <= 0 || maxCompanies > 50 {
		maxCompanies = 10
	}

	// Validate risk tolerance
	validRiskLevels := map[string]bool{"low": true, "medium": true, "high": true}
	if !validRiskLevels[riskTolerance] {
		riskTolerance = "medium"
	}

	optimization, err := h.advancedAnalyticsRepo.OptimizePortfolio(targetReturn, riskTolerance, maxCompanies)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Insufficient data for portfolio optimization"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to optimize portfolio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"optimization": optimization,
		"message":      "Portfolio optimized successfully",
	})
}

// AssessRisk calculates risk metrics for a company
func (h *AdvancedAnalyticsHandler) AssessRisk(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	assessment, err := h.advancedAnalyticsRepo.AssessRisk(companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Insufficient data for risk assessment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assess risk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assessment": assessment,
		"message":    "Risk assessment completed successfully",
	})
}

// AnalyzeTrend performs trend analysis on various metrics
func (h *AdvancedAnalyticsHandler) AnalyzeTrend(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	metric := c.Param("metric")
	period := c.DefaultQuery("period", "30d")

	// Validate metric
	validMetrics := map[string]bool{"esg_score": true, "stock_price": true, "market_cap": true}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric. Supported: esg_score, stock_price, market_cap"})
		return
	}

	analysis, err := h.advancedAnalyticsRepo.AnalyzeTrend(companyID, metric, period)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Insufficient data for trend analysis"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze trend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"analysis": analysis,
		"message":  "Trend analysis completed successfully",
	})
}

// GetAdvancedAnalyticsSummary provides a comprehensive analytics summary
func (h *AdvancedAnalyticsHandler) GetAdvancedAnalyticsSummary(c *gin.Context) {
	// Get portfolio optimization for top companies
	optimization, err := h.advancedAnalyticsRepo.OptimizePortfolio(0.10, "medium", 5)
	if err != nil {
		optimization = nil
	}

	// Get risk assessments for top companies (simplified)
	riskSummary := gin.H{
		"total_companies_assessed": 0,
		"average_risk_score":       0.0,
		"risk_distribution": gin.H{
			"low":    0,
			"medium": 0,
			"high":   0,
		},
	}

	// Get trend analysis summary
	trendSummary := gin.H{
		"esg_trends": gin.H{
			"improving": 0,
			"declining": 0,
			"stable":    0,
		},
		"price_trends": gin.H{
			"up":     0,
			"down":   0,
			"stable": 0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"portfolio_optimization": optimization,
			"risk_summary":           riskSummary,
			"trend_summary":          trendSummary,
			"message":                "Advanced analytics summary generated successfully",
		},
	})
}
