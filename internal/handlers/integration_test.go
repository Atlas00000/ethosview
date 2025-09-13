package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

// TestIntegration tests the handlers with a real database connection
// This requires the database to be running and seeded with test data
func TestIntegration_DashboardHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Connect to the test database
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/ethosview?sslmode=disable")
	if err != nil {
		t.Skip("Skipping integration test - database not available")
		return
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test - database not accessible")
		return
	}

	handler := NewDashboardHandler(db)
	router := gin.New()
	router.GET("/dashboard", handler.GetDashboard)

	req, _ := http.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "summary")
	assert.Contains(t, response, "top_esg_scores")
	assert.Contains(t, response, "sectors")
	assert.Contains(t, response, "sector_stats")

	// Verify summary structure
	summary, ok := response["summary"].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, summary, "total_companies")
	assert.Contains(t, summary, "total_sectors")
	assert.Contains(t, summary, "avg_esg_score")
}

func TestIntegration_CompanyHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/ethosview?sslmode=disable")
	if err != nil {
		t.Skip("Skipping integration test - database not available")
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test - database not accessible")
		return
	}

	handler := NewCompanyHandler(db)
	router := gin.New()
	router.GET("/companies", handler.ListCompanies)

	req, _ := http.NewRequest("GET", "/companies", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have companies array
	companies, ok := response["companies"].([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(companies), 0)
}

func TestIntegration_CompanyHandler_GetBySymbol(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/ethosview?sslmode=disable")
	if err != nil {
		t.Skip("Skipping integration test - database not available")
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test - database not accessible")
		return
	}

	handler := NewCompanyHandler(db)
	router := gin.New()
	router.GET("/companies/symbol/:symbol", handler.GetCompanyBySymbol)

	// Test with a known symbol
	req, _ := http.NewRequest("GET", "/companies/symbol/AAPL", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return OK or NotFound
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)

	if w.Code == http.StatusOK {
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "symbol")
	}
}

func TestIntegration_AnalyticsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/ethosview?sslmode=disable")
	if err != nil {
		t.Skip("Skipping integration test - database not available")
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test - database not accessible")
		return
	}

	handler := NewAnalyticsHandler(db)
	router := gin.New()
	router.GET("/analytics/summary", handler.GetAnalyticsSummary)

	req, _ := http.NewRequest("GET", "/analytics/summary", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "summary")
	assert.Contains(t, response, "sector_comparisons")
	assert.Contains(t, response, "top_esg_performers")
	assert.Contains(t, response, "top_market_cap")
	assert.Contains(t, response, "correlation")
}

func TestIntegration_ESGHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/ethosview?sslmode=disable")
	if err != nil {
		t.Skip("Skipping integration test - database not available")
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test - database not accessible")
		return
	}

	handler := NewESGHandler(db)
	router := gin.New()
	router.GET("/esg/scores", handler.ListESGScores)

	req, _ := http.NewRequest("GET", "/esg/scores", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have ESG scores array
	scores, ok := response["scores"].([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(scores), 0)
}

func TestIntegration_FinancialHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/ethosview?sslmode=disable")
	if err != nil {
		t.Skip("Skipping integration test - database not available")
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test - database not accessible")
		return
	}

	handler := NewFinancialHandler(db)
	router := gin.New()
	router.GET("/financial/market", handler.GetMarketData)

	req, _ := http.NewRequest("GET", "/financial/market", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "market_data")
}
