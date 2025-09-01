package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"ethosview-backend/internal/models"
	"ethosview-backend/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ESGHandler handles ESG score-related HTTP requests
type ESGHandler struct {
	repo *models.ESGScoreRepository
}

// NewESGHandler creates a new ESG handler
func NewESGHandler(db *sql.DB) *ESGHandler {
	return &ESGHandler{
		repo: models.NewESGScoreRepository(db),
	}
}

// CreateESGScore handles POST /api/v1/esg/scores
func (h *ESGHandler) CreateESGScore(c *gin.Context) {
	var score models.ESGScore
	if err := c.ShouldBindJSON(&score); err != nil {
		errors.HandleValidationError(c, err)
		return
	}

	if err := h.repo.CreateESGScore(&score); err != nil {
		errors.HandleDatabaseError(c, err, "ESG score")
		return
	}

	errors.SuccessResponse(c, score)
}

// GetESGScore handles GET /api/v1/esg/scores/:id
func (h *ESGHandler) GetESGScore(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errors.HandleValidationError(c, errors.ErrInvalidInput)
		return
	}

	score, err := h.repo.GetESGScoreByID(id)
	if err != nil {
		errors.HandleDatabaseError(c, err, "ESG score")
		return
	}

	errors.SuccessResponse(c, score)
}

// GetLatestESGScoreByCompany handles GET /api/v1/esg/companies/:id/latest
func (h *ESGHandler) GetLatestESGScoreByCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errors.HandleValidationError(c, errors.ErrInvalidInput)
		return
	}

	score, err := h.repo.GetLatestESGScoreByCompany(companyID)
	if err != nil {
		errors.HandleDatabaseError(c, err, "ESG score")
		return
	}

	errors.SuccessResponse(c, score)
}

// GetESGScoresByCompany handles GET /api/v1/esg/companies/:id/scores
func (h *ESGHandler) GetESGScoresByCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	limit := 20 // Default limit
	offset := 0 // Default offset

	// Parse pagination parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	scores, err := h.repo.GetESGScoresByCompany(companyID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ESG scores"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scores": scores,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(scores),
		},
	})
}

// ListESGScores handles GET /api/v1/esg/scores
func (h *ESGHandler) ListESGScores(c *gin.Context) {
	limit := 20 // Default limit
	offset := 0 // Default offset
	minScore := 0.0

	// Parse pagination parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Parse minimum score filter
	if minScoreStr := c.Query("min_score"); minScoreStr != "" {
		if ms, err := strconv.ParseFloat(minScoreStr, 64); err == nil && ms >= 0 {
			minScore = ms
		}
	}

	scores, err := h.repo.ListESGScores(limit, offset, minScore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ESG scores"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scores": scores,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(scores),
		},
		"filters": gin.H{
			"min_score": minScore,
		},
	})
}

// UpdateESGScore handles PUT /api/v1/esg/scores/:id
func (h *ESGHandler) UpdateESGScore(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ESG score ID"})
		return
	}

	var score models.ESGScore
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	score.ID = id
	if err := h.repo.UpdateESGScore(&score); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ESG score not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ESG score"})
		return
	}

	c.JSON(http.StatusOK, score)
}

// DeleteESGScore handles DELETE /api/v1/esg/scores/:id
func (h *ESGHandler) DeleteESGScore(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ESG score ID"})
		return
	}

	if err := h.repo.DeleteESGScore(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ESG score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ESG score deleted successfully"})
}
