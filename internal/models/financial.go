package models

import (
	"database/sql"
	"time"
)

// StockPrice represents a daily stock price record
type StockPrice struct {
	ID            int       `json:"id"`
	CompanyID     int       `json:"company_id"`
	Date          time.Time `json:"date"`
	OpenPrice     float64   `json:"open_price"`
	HighPrice     float64   `json:"high_price"`
	LowPrice      float64   `json:"low_price"`
	ClosePrice    float64   `json:"close_price"`
	Volume        int64     `json:"volume"`
	AdjustedClose float64   `json:"adjusted_close"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// FinancialIndicator represents financial metrics for a company
type FinancialIndicator struct {
	ID             int       `json:"id"`
	CompanyID      int       `json:"company_id"`
	Date           time.Time `json:"date"`
	MarketCap      *float64  `json:"market_cap,omitempty"`
	PERatio        *float64  `json:"pe_ratio,omitempty"`
	PBRatio        *float64  `json:"pb_ratio,omitempty"`
	DebtToEquity   *float64  `json:"debt_to_equity,omitempty"`
	ReturnOnEquity *float64  `json:"return_on_equity,omitempty"`
	ProfitMargin   *float64  `json:"profit_margin,omitempty"`
	RevenueGrowth  *float64  `json:"revenue_growth,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// MarketData represents broader market indicators
type MarketData struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	SP500Close  *float64  `json:"sp500_close,omitempty"`
	NasdaqClose *float64  `json:"nasdaq_close,omitempty"`
	DowClose    *float64  `json:"dow_close,omitempty"`
	VIXClose    *float64  `json:"vix_close,omitempty"`
	Treasury10Y *float64  `json:"treasury_10y,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockPriceRepository handles database operations for stock prices
type StockPriceRepository struct {
	db *sql.DB
}

// NewStockPriceRepository creates a new stock price repository
func NewStockPriceRepository(db *sql.DB) *StockPriceRepository {
	return &StockPriceRepository{db: db}
}

// GetByCompanyID retrieves stock prices for a specific company
func (r *StockPriceRepository) GetByCompanyID(companyID int, limit int) ([]StockPrice, error) {
	query := `
		SELECT id, company_id, date, open_price, high_price, low_price, close_price, volume, adjusted_close, created_at, updated_at
		FROM stock_prices 
		WHERE company_id = $1 
		ORDER BY date DESC 
		LIMIT $2
	`

	rows, err := r.db.Query(query, companyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []StockPrice
	for rows.Next() {
		var price StockPrice
		err := rows.Scan(
			&price.ID, &price.CompanyID, &price.Date, &price.OpenPrice, &price.HighPrice,
			&price.LowPrice, &price.ClosePrice, &price.Volume, &price.AdjustedClose,
			&price.CreatedAt, &price.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}

	return prices, nil
}

// GetLatestByCompanyID gets the most recent stock price for a company
func (r *StockPriceRepository) GetLatestByCompanyID(companyID int) (*StockPrice, error) {
	query := `
		SELECT id, company_id, date, open_price, high_price, low_price, close_price, volume, adjusted_close, created_at, updated_at
		FROM stock_prices 
		WHERE company_id = $1 
		ORDER BY date DESC 
		LIMIT 1
	`

	var price StockPrice
	err := r.db.QueryRow(query, companyID).Scan(
		&price.ID, &price.CompanyID, &price.Date, &price.OpenPrice, &price.HighPrice,
		&price.LowPrice, &price.ClosePrice, &price.Volume, &price.AdjustedClose,
		&price.CreatedAt, &price.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &price, nil
}

// FinancialIndicatorRepository handles database operations for financial indicators
type FinancialIndicatorRepository struct {
	db *sql.DB
}

// NewFinancialIndicatorRepository creates a new financial indicator repository
func NewFinancialIndicatorRepository(db *sql.DB) *FinancialIndicatorRepository {
	return &FinancialIndicatorRepository{db: db}
}

// GetByCompanyID retrieves financial indicators for a specific company
func (r *FinancialIndicatorRepository) GetByCompanyID(companyID int) (*FinancialIndicator, error) {
	query := `
		SELECT id, company_id, date, market_cap, pe_ratio, pb_ratio, debt_to_equity, 
		       return_on_equity, profit_margin, revenue_growth, created_at, updated_at
		FROM financial_indicators 
		WHERE company_id = $1 
		ORDER BY date DESC 
		LIMIT 1
	`

	var indicator FinancialIndicator
	err := r.db.QueryRow(query, companyID).Scan(
		&indicator.ID, &indicator.CompanyID, &indicator.Date, &indicator.MarketCap,
		&indicator.PERatio, &indicator.PBRatio, &indicator.DebtToEquity,
		&indicator.ReturnOnEquity, &indicator.ProfitMargin, &indicator.RevenueGrowth,
		&indicator.CreatedAt, &indicator.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &indicator, nil
}

// MarketDataRepository handles database operations for market data
type MarketDataRepository struct {
	db *sql.DB
}

// NewMarketDataRepository creates a new market data repository
func NewMarketDataRepository(db *sql.DB) *MarketDataRepository {
	return &MarketDataRepository{db: db}
}

// GetLatest retrieves the most recent market data
func (r *MarketDataRepository) GetLatest() (*MarketData, error) {
	query := `
		SELECT id, date, sp500_close, nasdaq_close, dow_close, vix_close, treasury_10y, created_at, updated_at
		FROM market_data 
		ORDER BY date DESC 
		LIMIT 1
	`

	var data MarketData
	err := r.db.QueryRow(query).Scan(
		&data.ID, &data.Date, &data.SP500Close, &data.NasdaqClose, &data.DowClose,
		&data.VIXClose, &data.Treasury10Y, &data.CreatedAt, &data.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetByDateRange retrieves market data for a date range
func (r *MarketDataRepository) GetByDateRange(startDate, endDate time.Time, limit int) ([]MarketData, error) {
	query := `
		SELECT id, date, sp500_close, nasdaq_close, dow_close, vix_close, treasury_10y, created_at, updated_at
		FROM market_data 
		WHERE date BETWEEN $1 AND $2
		ORDER BY date DESC 
		LIMIT $3
	`

	rows, err := r.db.Query(query, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []MarketData
	for rows.Next() {
		var marketData MarketData
		err := rows.Scan(
			&marketData.ID, &marketData.Date, &marketData.SP500Close, &marketData.NasdaqClose,
			&marketData.DowClose, &marketData.VIXClose, &marketData.Treasury10Y,
			&marketData.CreatedAt, &marketData.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, marketData)
	}

	return data, nil
}
