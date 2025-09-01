package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// CursorPagination handles cursor-based pagination
type CursorPagination struct {
	Limit     int    `json:"limit"`
	Cursor    string `json:"cursor,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore   bool   `json:"has_more"`
}

// CursorData represents the internal cursor data
type CursorData struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

// NewCursorPagination creates a new cursor pagination instance
func NewCursorPagination(limit int, cursor string) *CursorPagination {
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}

	return &CursorPagination{
		Limit:  limit,
		Cursor: cursor,
	}
}

// EncodeCursor encodes cursor data to a base64 string
func EncodeCursor(id int, timestamp time.Time, cursorType string) string {
	data := CursorData{
		ID:        id,
		Timestamp: timestamp,
		Type:      cursorType,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(jsonData)
}

// DecodeCursor decodes a base64 cursor string to cursor data
func DecodeCursor(cursor string) (*CursorData, error) {
	if cursor == "" {
		return nil, nil
	}

	jsonData, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format")
	}

	var data CursorData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("invalid cursor data")
	}

	return &data, nil
}

// SetNextCursor sets the next cursor for pagination
func (cp *CursorPagination) SetNextCursor(id int, timestamp time.Time, cursorType string) {
	cp.NextCursor = EncodeCursor(id, timestamp, cursorType)
}

// GetSQLCondition returns SQL WHERE condition for cursor pagination
func (cp *CursorPagination) GetSQLCondition(tableName string) (string, []interface{}, error) {
	if cp.Cursor == "" {
		return "", nil, nil
	}

	data, err := DecodeCursor(cp.Cursor)
	if err != nil {
		return "", nil, err
	}

	// For cursor pagination, we use (timestamp, id) as composite cursor
	// This ensures consistent ordering even with same timestamps
	condition := fmt.Sprintf("(%s.created_at, %s.id) < ($1, $2)", tableName, tableName)
	params := []interface{}{data.Timestamp, data.ID}

	return condition, params, nil
}

// GetSQLOrderBy returns SQL ORDER BY clause for cursor pagination
func (cp *CursorPagination) GetSQLOrderBy(tableName string) string {
	return fmt.Sprintf("%s.created_at DESC, %s.id DESC", tableName, tableName)
}

// CompanyPaginationResponse represents paginated company response
type CompanyPaginationResponse struct {
	Companies  []interface{}    `json:"companies"`
	Pagination CursorPagination `json:"pagination"`
}

// ESGPaginationResponse represents paginated ESG response
type ESGPaginationResponse struct {
	Scores     []interface{}    `json:"scores"`
	Pagination CursorPagination `json:"pagination"`
}

// PaginationParams represents common pagination parameters
type PaginationParams struct {
	Limit  int    `form:"limit"`
	Cursor string `form:"cursor"`
}

// Validate validates pagination parameters
func (p *PaginationParams) Validate() error {
	if p.Limit < 0 || p.Limit > 100 {
		p.Limit = 20
	}

	if p.Limit == 0 {
		p.Limit = 20
	}

	return nil
}

// ParseCursorFromQuery parses cursor from query parameters
func ParseCursorFromQuery(cursor string, limit string) (*CursorPagination, error) {
	// Parse limit
	limitInt := 20
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitInt = l
		}
	}

	return NewCursorPagination(limitInt, cursor), nil
}

// BuildResponse builds a standardized pagination response
func BuildResponse(data []interface{}, pagination *CursorPagination, dataKey string) map[string]interface{} {
	return map[string]interface{}{
		dataKey:      data,
		"pagination": pagination,
	}
}
