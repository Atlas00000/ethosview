package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"ethosview-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// FinancialHandler handles financial data-related HTTP requests
type FinancialHandler struct {
	stockPriceRepo         *models.StockPriceRepository
	financialIndicatorRepo *models.FinancialIndicatorRepository
	marketDataRepo         *models.MarketDataRepository
}

// NewFinancialHandler creates a new financial handler
func NewFinancialHandler(db *sql.DB) *FinancialHandler {
	return &FinancialHandler{
		stockPriceRepo:         models.NewStockPriceRepository(db),
		financialIndicatorRepo: models.NewFinancialIndicatorRepository(db),
		marketDataRepo:         models.NewMarketDataRepository(db),
	}
}

// GetStockPrices retrieves stock prices for a company
func (h *FinancialHandler) GetStockPrices(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "30")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 30
	}

	prices, err := h.stockPriceRepo.GetByCompanyID(companyID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stock prices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_id": companyID,
		"prices":     prices,
		"count":      len(prices),
	})
}

// GetLatestStockPrice retrieves the latest stock price for a company
func (h *FinancialHandler) GetLatestStockPrice(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	price, err := h.stockPriceRepo.GetLatestByCompanyID(companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock price not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_id": companyID,
		"price":      price,
	})
}

// GetFinancialIndicators retrieves financial indicators for a company
func (h *FinancialHandler) GetFinancialIndicators(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	indicators, err := h.financialIndicatorRepo.GetByCompanyID(companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Financial indicators not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_id": companyID,
		"indicators": indicators,
	})
}

// GetMarketData retrieves the latest market data
func (h *FinancialHandler) GetMarketData(c *gin.Context) {
	data, err := h.marketDataRepo.GetLatest()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Market data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"market_data": data,
	})
}

// GetMarketDataHistory retrieves market data for a date range
func (h *FinancialHandler) GetMarketDataHistory(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format (YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format (YYYY-MM-DD)"})
		return
	}

	limitStr := c.DefaultQuery("limit", "30")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 30
	}

	data, err := h.marketDataRepo.GetByDateRange(startDate, endDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve market data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"data":       data,
		"count":      len(data),
	})
}

// GetCompanyFinancialSummary retrieves a comprehensive financial summary for a company
func (h *FinancialHandler) GetCompanyFinancialSummary(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	// Get latest stock price
	stockPrice, err := h.stockPriceRepo.GetLatestByCompanyID(companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company financial data not found"})
		return
	}

	// Get financial indicators
	indicators, err := h.financialIndicatorRepo.GetByCompanyID(companyID)
	if err != nil {
		// Continue without indicators if not available
		indicators = nil
	}

	// Calculate price change (simplified - would need previous day data for real calculation)
	priceChange := 0.0
	priceChangePercent := 0.0

	c.JSON(http.StatusOK, gin.H{
		"company_id": companyID,
		"summary": gin.H{
			"current_price":        stockPrice.ClosePrice,
			"price_change":         priceChange,
			"price_change_percent": priceChangePercent,
			"volume":               stockPrice.Volume,
			"date":                 stockPrice.Date.Format("2006-01-02"),
		},
		"indicators": indicators,
	})
}
