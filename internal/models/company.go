package models

import (
	"database/sql"
	"time"
)

// Company represents a company in the system
type Company struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	Sector    string    `json:"sector"`
	Industry  string    `json:"industry"`
	Country   string    `json:"country"`
	MarketCap float64   `json:"market_cap"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CompanyRepository handles database operations for companies
type CompanyRepository struct {
	db *sql.DB
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// CreateCompany creates a new company
func (r *CompanyRepository) CreateCompany(company *Company) error {
	query := `
		INSERT INTO companies (name, symbol, sector, industry, country, market_cap)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		company.Name,
		company.Symbol,
		company.Sector,
		company.Industry,
		company.Country,
		company.MarketCap,
	).Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
}

// GetCompanyByID retrieves a company by ID
func (r *CompanyRepository) GetCompanyByID(id int) (*Company, error) {
	company := &Company{}
	query := `
		SELECT id, name, symbol, sector, industry, country, market_cap, created_at, updated_at
		FROM companies WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&company.ID,
		&company.Name,
		&company.Symbol,
		&company.Sector,
		&company.Industry,
		&company.Country,
		&company.MarketCap,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}

// GetCompanyBySymbol retrieves a company by symbol
func (r *CompanyRepository) GetCompanyBySymbol(symbol string) (*Company, error) {
	company := &Company{}
	query := `
		SELECT id, name, symbol, sector, industry, country, market_cap, created_at, updated_at
		FROM companies WHERE symbol = $1
	`

	err := r.db.QueryRow(query, symbol).Scan(
		&company.ID,
		&company.Name,
		&company.Symbol,
		&company.Sector,
		&company.Industry,
		&company.Country,
		&company.MarketCap,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}

// UpdateCompany updates an existing company
func (r *CompanyRepository) UpdateCompany(company *Company) error {
	query := `
		UPDATE companies 
		SET name = $1, sector = $2, industry = $3, country = $4, market_cap = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at
	`

	return r.db.QueryRow(
		query,
		company.Name,
		company.Sector,
		company.Industry,
		company.Country,
		company.MarketCap,
		company.ID,
	).Scan(&company.UpdatedAt)
}

// DeleteCompany deletes a company by ID
func (r *CompanyRepository) DeleteCompany(id int) error {
	query := `DELETE FROM companies WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ListCompanies retrieves all companies with pagination and optional filtering
func (r *CompanyRepository) ListCompanies(limit, offset int, sector string) ([]*Company, error) {
	var query string
	var args []interface{}

	if sector != "" {
		query = `
			SELECT id, name, symbol, sector, industry, country, market_cap, created_at, updated_at
			FROM companies
			WHERE sector = $1
			ORDER BY name ASC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{sector, limit, offset}
	} else {
		query = `
			SELECT id, name, symbol, sector, industry, country, market_cap, created_at, updated_at
			FROM companies
			ORDER BY name ASC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*Company
	for rows.Next() {
		company := &Company{}
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Symbol,
			&company.Sector,
			&company.Industry,
			&company.Country,
			&company.MarketCap,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// GetSectors retrieves all unique sectors
func (r *CompanyRepository) GetSectors() ([]string, error) {
	query := `SELECT DISTINCT sector FROM companies WHERE sector IS NOT NULL ORDER BY sector`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectors []string
	for rows.Next() {
		var sector string
		if err := rows.Scan(&sector); err != nil {
			return nil, err
		}
		sectors = append(sectors, sector)
	}

	return sectors, nil
}
