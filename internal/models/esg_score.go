package models

import (
	"database/sql"
	"time"
)

// ESGScore represents an ESG score for a company
type ESGScore struct {
	ID                 int       `json:"id"`
	CompanyID          int       `json:"company_id"`
	EnvironmentalScore float64   `json:"environmental_score"`
	SocialScore        float64   `json:"social_score"`
	GovernanceScore    float64   `json:"governance_score"`
	OverallScore       float64   `json:"overall_score"`
	ScoreDate          time.Time `json:"score_date"`
	DataSource         string    `json:"data_source"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Joined data
	CompanyName   string `json:"company_name,omitempty"`
	CompanySymbol string `json:"company_symbol,omitempty"`
}

// ESGScoreRepository handles database operations for ESG scores
type ESGScoreRepository struct {
	db *sql.DB
}

// NewESGScoreRepository creates a new ESG score repository
func NewESGScoreRepository(db *sql.DB) *ESGScoreRepository {
	return &ESGScoreRepository{db: db}
}

// CreateESGScore creates a new ESG score
func (r *ESGScoreRepository) CreateESGScore(score *ESGScore) error {
	query := `
		INSERT INTO esg_scores (company_id, environmental_score, social_score, governance_score, overall_score, score_date, data_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		score.CompanyID,
		score.EnvironmentalScore,
		score.SocialScore,
		score.GovernanceScore,
		score.OverallScore,
		score.ScoreDate,
		score.DataSource,
	).Scan(&score.ID, &score.CreatedAt, &score.UpdatedAt)
}

// GetESGScoreByID retrieves an ESG score by ID
func (r *ESGScoreRepository) GetESGScoreByID(id int) (*ESGScore, error) {
	score := &ESGScore{}
	query := `
		SELECT es.id, es.company_id, es.environmental_score, es.social_score, es.governance_score, 
		       es.overall_score, es.score_date, es.data_source, es.created_at, es.updated_at,
		       c.name as company_name, c.symbol as company_symbol
		FROM esg_scores es
		JOIN companies c ON es.company_id = c.id
		WHERE es.id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&score.ID,
		&score.CompanyID,
		&score.EnvironmentalScore,
		&score.SocialScore,
		&score.GovernanceScore,
		&score.OverallScore,
		&score.ScoreDate,
		&score.DataSource,
		&score.CreatedAt,
		&score.UpdatedAt,
		&score.CompanyName,
		&score.CompanySymbol,
	)

	if err != nil {
		return nil, err
	}

	return score, nil
}

// GetLatestESGScoreByCompany retrieves the latest ESG score for a company
func (r *ESGScoreRepository) GetLatestESGScoreByCompany(companyID int) (*ESGScore, error) {
	score := &ESGScore{}
	query := `
		SELECT es.id, es.company_id, es.environmental_score, es.social_score, es.governance_score, 
		       es.overall_score, es.score_date, es.data_source, es.created_at, es.updated_at,
		       c.name as company_name, c.symbol as company_symbol
		FROM esg_scores es
		JOIN companies c ON es.company_id = c.id
		WHERE es.company_id = $1
		ORDER BY es.score_date DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, companyID).Scan(
		&score.ID,
		&score.CompanyID,
		&score.EnvironmentalScore,
		&score.SocialScore,
		&score.GovernanceScore,
		&score.OverallScore,
		&score.ScoreDate,
		&score.DataSource,
		&score.CreatedAt,
		&score.UpdatedAt,
		&score.CompanyName,
		&score.CompanySymbol,
	)

	if err != nil {
		return nil, err
	}

	return score, nil
}

// GetESGScoresByCompany retrieves all ESG scores for a company
func (r *ESGScoreRepository) GetESGScoresByCompany(companyID int, limit, offset int) ([]*ESGScore, error) {
	query := `
		SELECT es.id, es.company_id, es.environmental_score, es.social_score, es.governance_score, 
		       es.overall_score, es.score_date, es.data_source, es.created_at, es.updated_at,
		       c.name as company_name, c.symbol as company_symbol
		FROM esg_scores es
		JOIN companies c ON es.company_id = c.id
		WHERE es.company_id = $1
		ORDER BY es.score_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []*ESGScore
	for rows.Next() {
		score := &ESGScore{}
		err := rows.Scan(
			&score.ID,
			&score.CompanyID,
			&score.EnvironmentalScore,
			&score.SocialScore,
			&score.GovernanceScore,
			&score.OverallScore,
			&score.ScoreDate,
			&score.DataSource,
			&score.CreatedAt,
			&score.UpdatedAt,
			&score.CompanyName,
			&score.CompanySymbol,
		)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}

// ListESGScores retrieves all ESG scores with pagination and optional filtering
func (r *ESGScoreRepository) ListESGScores(limit, offset int, minScore float64) ([]*ESGScore, error) {
	var query string
	var args []interface{}

	if minScore > 0 {
		query = `
			SELECT es.id, es.company_id, es.environmental_score, es.social_score, es.governance_score, 
			       es.overall_score, es.score_date, es.data_source, es.created_at, es.updated_at,
			       c.name as company_name, c.symbol as company_symbol
			FROM esg_scores es
			JOIN companies c ON es.company_id = c.id
			WHERE es.overall_score >= $1
			ORDER BY es.overall_score DESC, es.score_date DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{minScore, limit, offset}
	} else {
		query = `
			SELECT es.id, es.company_id, es.environmental_score, es.social_score, es.governance_score, 
			       es.overall_score, es.score_date, es.data_source, es.created_at, es.updated_at,
			       c.name as company_name, c.symbol as company_symbol
			FROM esg_scores es
			JOIN companies c ON es.company_id = c.id
			ORDER BY es.overall_score DESC, es.score_date DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []*ESGScore
	for rows.Next() {
		score := &ESGScore{}
		err := rows.Scan(
			&score.ID,
			&score.CompanyID,
			&score.EnvironmentalScore,
			&score.SocialScore,
			&score.GovernanceScore,
			&score.OverallScore,
			&score.ScoreDate,
			&score.DataSource,
			&score.CreatedAt,
			&score.UpdatedAt,
			&score.CompanyName,
			&score.CompanySymbol,
		)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}

// UpdateESGScore updates an existing ESG score
func (r *ESGScoreRepository) UpdateESGScore(score *ESGScore) error {
	query := `
		UPDATE esg_scores 
		SET environmental_score = $1, social_score = $2, governance_score = $3, 
		    overall_score = $4, score_date = $5, data_source = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`

	return r.db.QueryRow(
		query,
		score.EnvironmentalScore,
		score.SocialScore,
		score.GovernanceScore,
		score.OverallScore,
		score.ScoreDate,
		score.DataSource,
		score.ID,
	).Scan(&score.UpdatedAt)
}

// DeleteESGScore deletes an ESG score by ID
func (r *ESGScoreRepository) DeleteESGScore(id int) error {
	query := `DELETE FROM esg_scores WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
