package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"ethosview-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// CompanyHandler handles company-related HTTP requests
type CompanyHandler struct {
	repo *models.CompanyRepository
}

// NewCompanyHandler creates a new company handler
func NewCompanyHandler(db *sql.DB) *CompanyHandler {
	return &CompanyHandler{
		repo: models.NewCompanyRepository(db),
	}
}

// CreateCompany handles POST /api/v1/companies
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.repo.CreateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

// GetCompany handles GET /api/v1/companies/:id
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	company, err := h.repo.GetCompanyByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve company"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// GetCompanyBySymbol handles GET /api/v1/companies/symbol/:symbol
func (h *CompanyHandler) GetCompanyBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	company, err := h.repo.GetCompanyBySymbol(symbol)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve company"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// ListCompanies handles GET /api/v1/companies
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	limit := 20 // Default limit
	offset := 0 // Default offset
	sector := c.Query("sector")

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

	companies, err := h.repo.ListCompanies(limit, offset, sector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve companies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": companies,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(companies),
		},
	})
}

// UpdateCompany handles PUT /api/v1/companies/:id
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	company.ID = id
	if err := h.repo.UpdateCompany(&company); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// DeleteCompany handles DELETE /api/v1/companies/:id
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	if err := h.repo.DeleteCompany(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}

// GetSectors handles GET /api/v1/companies/sectors
func (h *CompanyHandler) GetSectors(c *gin.Context) {
	sectors, err := h.repo.GetSectors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sectors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sectors": sectors})
}
