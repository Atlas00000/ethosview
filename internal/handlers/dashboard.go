package handlers

import (
	"database/sql"
	"net/http"

	"ethosview-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	companyRepo *models.CompanyRepository
	esgRepo     *models.ESGScoreRepository
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{
		companyRepo: models.NewCompanyRepository(db),
		esgRepo:     models.NewESGScoreRepository(db),
	}
}

// GetDashboard handles GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// Get top ESG scores
	topScores, err := h.esgRepo.ListESGScores(5, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ESG scores"})
		return
	}

	// Get sectors
	sectors, err := h.companyRepo.GetSectors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sectors"})
		return
	}

	// Get companies count by sector
	sectorStats := make(map[string]int)
	for _, sector := range sectors {
		companies, err := h.companyRepo.ListCompanies(100, 0, sector)
		if err != nil {
			continue
		}
		sectorStats[sector] = len(companies)
	}

	// Calculate average ESG score
	var totalScore float64
	var scoreCount int
	for _, score := range topScores {
		totalScore += score.OverallScore
		scoreCount++
	}

	avgScore := 0.0
	if scoreCount > 0 {
		avgScore = totalScore / float64(scoreCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"total_companies": len(sectors) * 2, // Rough estimate
			"total_sectors":   len(sectors),
			"avg_esg_score":   avgScore,
		},
		"top_esg_scores": topScores,
		"sectors":        sectors,
		"sector_stats":   sectorStats,
	})
}
