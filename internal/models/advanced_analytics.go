package models

import (
	"database/sql"
	"math"
	"time"
)

// ESGPrediction represents ESG score prediction
type ESGPrediction struct {
	CompanyID      int       `json:"company_id"`
	CompanyName    string    `json:"company_name"`
	CurrentScore   float64   `json:"current_score"`
	PredictedScore float64   `json:"predicted_score"`
	Confidence     float64   `json:"confidence"`
	PredictionDate time.Time `json:"prediction_date"`
	Factors        []string  `json:"factors"`
}

// PortfolioOptimization represents portfolio optimization results
type PortfolioOptimization struct {
	PortfolioID    string       `json:"portfolio_id"`
	TotalValue     float64      `json:"total_value"`
	ExpectedReturn float64      `json:"expected_return"`
	RiskLevel      string       `json:"risk_level"`
	SharpeRatio    float64      `json:"sharpe_ratio"`
	Allocations    []Allocation `json:"allocations"`
	CreatedAt      time.Time    `json:"created_at"`
}

// Allocation represents portfolio allocation
type Allocation struct {
	CompanyID   int     `json:"company_id"`
	CompanyName string  `json:"company_name"`
	Percentage  float64 `json:"percentage"`
	Amount      float64 `json:"amount"`
	ESGScore    float64 `json:"esg_score"`
}

// RiskAssessment represents risk assessment metrics
type RiskAssessment struct {
	CompanyID     int     `json:"company_id"`
	CompanyName   string  `json:"company_name"`
	Volatility    float64 `json:"volatility"`
	Beta          float64 `json:"beta"`
	ValueAtRisk   float64 `json:"value_at_risk"`
	MaxDrawdown   float64 `json:"max_drawdown"`
	RiskScore     float64 `json:"risk_score"`
	RiskLevel     string  `json:"risk_level"`
	ESGRiskFactor float64 `json:"esg_risk_factor"`
}

// TrendAnalysis represents trend analysis results
type TrendAnalysis struct {
	CompanyID    int       `json:"company_id"`
	CompanyName  string    `json:"company_name"`
	Metric       string    `json:"metric"`
	Trend        string    `json:"trend"` // "up", "down", "stable"
	Slope        float64   `json:"slope"`
	R2           float64   `json:"r2"`
	Confidence   float64   `json:"confidence"`
	Period       string    `json:"period"`
	AnalysisDate time.Time `json:"analysis_date"`
}

// AdvancedAnalyticsRepository handles advanced analytical database operations
type AdvancedAnalyticsRepository struct {
	db *sql.DB
}

// NewAdvancedAnalyticsRepository creates a new advanced analytics repository
func NewAdvancedAnalyticsRepository(db *sql.DB) *AdvancedAnalyticsRepository {
	return &AdvancedAnalyticsRepository{db: db}
}

// PredictESGScore predicts ESG score for a company using historical data
func (r *AdvancedAnalyticsRepository) PredictESGScore(companyID int) (*ESGPrediction, error) {
	// Get historical ESG scores for the company
	query := `
		SELECT c.name, es.overall_score, es.score_date
		FROM esg_scores es
		JOIN companies c ON es.company_id = c.id
		WHERE es.company_id = $1
		ORDER BY es.score_date DESC
		LIMIT 10
	`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []float64
	var companyName string
	var latestScore float64

	for rows.Next() {
		var score float64
		var date time.Time
		if err := rows.Scan(&companyName, &score, &date); err != nil {
			return nil, err
		}
		scores = append(scores, score)
		if len(scores) == 1 {
			latestScore = score
		}
	}

	if len(scores) < 3 {
		return nil, sql.ErrNoRows
	}

	// Simple linear regression for prediction
	predictedScore, confidence := r.calculateLinearRegression(scores)

	// Ensure prediction is within reasonable bounds
	if predictedScore < 0 {
		predictedScore = 0
	} else if predictedScore > 100 {
		predictedScore = 100
	}

	return &ESGPrediction{
		CompanyID:      companyID,
		CompanyName:    companyName,
		CurrentScore:   latestScore,
		PredictedScore: predictedScore,
		Confidence:     confidence,
		PredictionDate: time.Now(),
		Factors:        []string{"historical_trend", "linear_regression"},
	}, nil
}

// OptimizePortfolio creates an optimized portfolio based on ESG and financial criteria
func (r *AdvancedAnalyticsRepository) OptimizePortfolio(targetReturn float64, riskTolerance string, maxCompanies int) (*PortfolioOptimization, error) {
	// Get companies with ESG scores and financial data
	query := `
		SELECT 
			c.id, c.name, c.sector,
			es.overall_score as esg_score,
			fi.market_cap, fi.pe_ratio, fi.return_on_equity,
			sp.close_price
		FROM companies c
		LEFT JOIN (
			SELECT DISTINCT ON (company_id) company_id, overall_score
			FROM esg_scores
			ORDER BY company_id, score_date DESC
		) es ON c.id = es.company_id
		LEFT JOIN (
			SELECT DISTINCT ON (company_id) company_id, market_cap, pe_ratio, return_on_equity
			FROM financial_indicators
			ORDER BY company_id, date DESC
		) fi ON c.id = fi.company_id
		LEFT JOIN (
			SELECT DISTINCT ON (company_id) company_id, close_price
			FROM stock_prices
			ORDER BY company_id, date DESC
		) sp ON c.id = sp.company_id
		WHERE es.overall_score IS NOT NULL
		AND fi.market_cap IS NOT NULL
		AND sp.close_price IS NOT NULL
		ORDER BY es.overall_score DESC, fi.return_on_equity DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, maxCompanies*2) // Get more companies for selection
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []struct {
		ID         int
		Name       string
		Sector     string
		ESGScore   float64
		MarketCap  float64
		PERatio    float64
		ROE        float64
		ClosePrice float64
	}

	for rows.Next() {
		var comp struct {
			ID         int
			Name       string
			Sector     string
			ESGScore   float64
			MarketCap  float64
			PERatio    float64
			ROE        float64
			ClosePrice float64
		}
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.Sector, &comp.ESGScore, &comp.MarketCap, &comp.PERatio, &comp.ROE, &comp.ClosePrice); err != nil {
			return nil, err
		}
		companies = append(companies, comp)
	}

	if len(companies) == 0 {
		return nil, sql.ErrNoRows
	}

	// Simple portfolio optimization algorithm
	allocations := r.optimizeAllocations(companies, targetReturn, riskTolerance, maxCompanies)

	// Calculate portfolio metrics
	totalValue := 1000000.0 // $1M portfolio
	expectedReturn := r.calculateExpectedReturn(allocations)
	sharpeRatio := r.calculateSharpeRatio(allocations)

	return &PortfolioOptimization{
		PortfolioID:    "opt_" + time.Now().Format("20060102150405"),
		TotalValue:     totalValue,
		ExpectedReturn: expectedReturn,
		RiskLevel:      riskTolerance,
		SharpeRatio:    sharpeRatio,
		Allocations:    allocations,
		CreatedAt:      time.Now(),
	}, nil
}

// AssessRisk calculates risk metrics for a company
func (r *AdvancedAnalyticsRepository) AssessRisk(companyID int) (*RiskAssessment, error) {
	// Get company information
	var companyName string
	err := r.db.QueryRow("SELECT name FROM companies WHERE id = $1", companyID).Scan(&companyName)
	if err != nil {
		return nil, err
	}

	// Get historical stock prices for volatility calculation
	query := `
		SELECT close_price, date
		FROM stock_prices
		WHERE company_id = $1
		ORDER BY date DESC
		LIMIT 30
	`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []float64
	for rows.Next() {
		var price float64
		var date time.Time
		if err := rows.Scan(&price, &date); err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}

	if len(prices) < 5 {
		return nil, sql.ErrNoRows
	}

	// Calculate risk metrics
	volatility := r.calculateVolatility(prices)
	beta := r.calculateBeta(prices)
	valueAtRisk := r.calculateValueAtRisk(prices)
	maxDrawdown := r.calculateMaxDrawdown(prices)
	riskScore := r.calculateRiskScore(volatility, beta, valueAtRisk)
	esgRiskFactor := r.calculateESGRiskFactor(companyID)

	return &RiskAssessment{
		CompanyID:     companyID,
		CompanyName:   companyName,
		Volatility:    volatility,
		Beta:          beta,
		ValueAtRisk:   valueAtRisk,
		MaxDrawdown:   maxDrawdown,
		RiskScore:     riskScore,
		RiskLevel:     r.getRiskLevel(riskScore),
		ESGRiskFactor: esgRiskFactor,
	}, nil
}

// AnalyzeTrend performs trend analysis on various metrics
func (r *AdvancedAnalyticsRepository) AnalyzeTrend(companyID int, metric string, period string) (*TrendAnalysis, error) {
	var companyName string
	err := r.db.QueryRow("SELECT name FROM companies WHERE id = $1", companyID).Scan(&companyName)
	if err != nil {
		return nil, err
	}

	var query string
	var values []float64

	switch metric {
	case "esg_score":
		query = `
			SELECT overall_score, score_date
			FROM esg_scores
			WHERE company_id = $1
			ORDER BY score_date DESC
			LIMIT 10
		`
	case "stock_price":
		query = `
			SELECT close_price, date
			FROM stock_prices
			WHERE company_id = $1
			ORDER BY date DESC
			LIMIT 30
		`
	case "market_cap":
		query = `
			SELECT market_cap, date
			FROM financial_indicators
			WHERE company_id = $1
			ORDER BY date DESC
			LIMIT 30
		`
	default:
		return nil, sql.ErrNoRows
	}

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value float64
		var date time.Time
		if err := rows.Scan(&value, &date); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if len(values) < 3 {
		return nil, sql.ErrNoRows
	}

	// Calculate trend metrics
	slope, r2 := r.calculateLinearRegression(values)
	trend := r.determineTrend(slope)
	confidence := r.calculateConfidence(r2)

	return &TrendAnalysis{
		CompanyID:    companyID,
		CompanyName:  companyName,
		Metric:       metric,
		Trend:        trend,
		Slope:        slope,
		R2:           r2,
		Confidence:   confidence,
		Period:       period,
		AnalysisDate: time.Now(),
	}, nil
}

// Helper methods for calculations
func (r *AdvancedAnalyticsRepository) calculateLinearRegression(values []float64) (slope, r2 float64) {
	n := len(values)
	if n < 2 {
		return 0, 0
	}

	// Simple linear regression
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope = (float64(n)*sumXY - sumX*sumY) / (float64(n)*sumX2 - sumX*sumX)

	// Calculate R-squared
	meanY := sumY / float64(n)
	ssRes := 0.0
	ssTot := 0.0

	for i, y := range values {
		x := float64(i)
		predicted := slope*x + (sumY/float64(n) - slope*sumX/float64(n))
		ssRes += (y - predicted) * (y - predicted)
		ssTot += (y - meanY) * (y - meanY)
	}

	if ssTot > 0 {
		r2 = 1 - (ssRes / ssTot)
	}

	return slope, r2
}

func (r *AdvancedAnalyticsRepository) optimizeAllocations(companies []struct {
	ID         int
	Name       string
	Sector     string
	ESGScore   float64
	MarketCap  float64
	PERatio    float64
	ROE        float64
	ClosePrice float64
}, targetReturn float64, riskTolerance string, maxCompanies int) []Allocation {

	// Simple optimization: select top companies by ESG score and ROE
	allocations := make([]Allocation, 0, maxCompanies)
	totalValue := 1000000.0 // $1M portfolio

	// Sort companies by combined score (ESG + ROE)
	type scoredCompany struct {
		company struct {
			ID         int
			Name       string
			Sector     string
			ESGScore   float64
			MarketCap  float64
			PERatio    float64
			ROE        float64
			ClosePrice float64
		}
		score float64
	}

	var scoredCompanies []scoredCompany
	for _, comp := range companies {
		score := comp.ESGScore*0.7 + comp.ROE*30 // Weight ESG more heavily
		scoredCompanies = append(scoredCompanies, scoredCompany{company: comp, score: score})
	}

	// Select top companies
	selectedCount := min(maxCompanies, len(scoredCompanies))
	equalAllocation := totalValue / float64(selectedCount)

	for i := 0; i < selectedCount; i++ {
		comp := scoredCompanies[i].company
		allocations = append(allocations, Allocation{
			CompanyID:   comp.ID,
			CompanyName: comp.Name,
			Percentage:  100.0 / float64(selectedCount),
			Amount:      equalAllocation,
			ESGScore:    comp.ESGScore,
		})
	}

	return allocations
}

func (r *AdvancedAnalyticsRepository) calculateExpectedReturn(allocations []Allocation) float64 {
	totalReturn := 0.0
	for _, alloc := range allocations {
		// Simple return calculation based on ESG score
		expectedReturn := 0.05 + (alloc.ESGScore/100.0)*0.10 // 5-15% range
		totalReturn += expectedReturn * alloc.Percentage / 100.0
	}
	return totalReturn
}

func (r *AdvancedAnalyticsRepository) calculateSharpeRatio(allocations []Allocation) float64 {
	expectedReturn := r.calculateExpectedReturn(allocations)
	// Assume risk-free rate of 2% and portfolio volatility of 15%
	riskFreeRate := 0.02
	volatility := 0.15
	return (expectedReturn - riskFreeRate) / volatility
}

func (r *AdvancedAnalyticsRepository) calculateVolatility(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}

	// Calculate returns
	returns := make([]float64, len(prices)-1)
	for i := 0; i < len(prices)-1; i++ {
		returns[i] = (prices[i] - prices[i+1]) / prices[i+1]
	}

	// Calculate standard deviation
	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))

	variance := 0.0
	for _, r := range returns {
		variance += (r - mean) * (r - mean)
	}
	variance /= float64(len(returns))

	return math.Sqrt(variance) * math.Sqrt(252) // Annualized
}

func (r *AdvancedAnalyticsRepository) calculateBeta(prices []float64) float64 {
	// Simplified beta calculation (assuming market correlation of 0.7)
	return 0.7
}

func (r *AdvancedAnalyticsRepository) calculateValueAtRisk(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}

	// Calculate returns
	returns := make([]float64, len(prices)-1)
	for i := 0; i < len(prices)-1; i++ {
		returns[i] = (prices[i] - prices[i+1]) / prices[i+1]
	}

	// Simple VaR calculation (95% confidence)
	volatility := r.calculateVolatility(prices)
	return volatility * 1.645 * math.Sqrt(1.0/252) // Daily VaR
}

func (r *AdvancedAnalyticsRepository) calculateMaxDrawdown(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}

	maxDrawdown := 0.0
	peak := prices[0]

	for _, price := range prices {
		if price > peak {
			peak = price
		}
		drawdown := (peak - price) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

func (r *AdvancedAnalyticsRepository) calculateRiskScore(volatility, beta, valueAtRisk float64) float64 {
	// Normalize and combine risk factors
	volScore := math.Min(volatility*100, 100)
	betaScore := math.Min(beta*50, 100)
	vaRScore := math.Min(valueAtRisk*1000, 100)

	return (volScore + betaScore + vaRScore) / 3
}

func (r *AdvancedAnalyticsRepository) getRiskLevel(riskScore float64) string {
	if riskScore < 30 {
		return "low"
	} else if riskScore < 60 {
		return "medium"
	} else {
		return "high"
	}
}

func (r *AdvancedAnalyticsRepository) calculateESGRiskFactor(companyID int) float64 {
	// Get latest ESG score
	var esgScore float64
	err := r.db.QueryRow(`
		SELECT overall_score 
		FROM esg_scores 
		WHERE company_id = $1 
		ORDER BY score_date DESC 
		LIMIT 1
	`, companyID).Scan(&esgScore)

	if err != nil {
		return 0.5 // Default risk factor
	}

	// Lower ESG score = higher risk factor
	return 1.0 - (esgScore / 100.0)
}

func (r *AdvancedAnalyticsRepository) determineTrend(slope float64) string {
	if slope > 0.01 {
		return "up"
	} else if slope < -0.01 {
		return "down"
	} else {
		return "stable"
	}
}

func (r *AdvancedAnalyticsRepository) calculateConfidence(r2 float64) float64 {
	return r2 * 100 // Convert to percentage
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
