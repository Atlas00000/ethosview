package middleware

import (
	"regexp"
	"strconv"
	"strings"

	"ethosview-backend/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationRules defines validation rules for different endpoints
type ValidationRules struct {
	RequiredFields []string
	StringRules    map[string]StringRule
	NumberRules    map[string]NumberRule
	EmailRules     map[string]EmailRule
	DateRules      map[string]DateRule
}

// StringRule defines validation rules for string fields
type StringRule struct {
	MinLength int
	MaxLength int
	Pattern   string
	Required  bool
}

// NumberRule defines validation rules for numeric fields
type NumberRule struct {
	Min      float64
	Max      float64
	Required bool
	Integer  bool
}

// EmailRule defines validation rules for email fields
type EmailRule struct {
	Required bool
}

// DateRule defines validation rules for date fields
type DateRule struct {
	Required bool
	Format   string
}

// ValidationMiddleware creates validation middleware for specific rules
func ValidationMiddleware(rules ValidationRules) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate required fields
		if err := validateRequiredFields(c, rules.RequiredFields); err != nil {
			errors.HandleValidationError(c, err)
			c.Abort()
			return
		}

		// Validate string fields
		if err := validateStringFields(c, rules.StringRules); err != nil {
			errors.HandleValidationError(c, err)
			c.Abort()
			return
		}

		// Validate number fields
		if err := validateNumberFields(c, rules.NumberRules); err != nil {
			errors.HandleValidationError(c, err)
			c.Abort()
			return
		}

		// Validate email fields
		if err := validateEmailFields(c, rules.EmailRules); err != nil {
			errors.HandleValidationError(c, err)
			c.Abort()
			return
		}

		// Validate date fields
		if err := validateDateFields(c, rules.DateRules); err != nil {
			errors.HandleValidationError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}

// validateRequiredFields checks if all required fields are present
func validateRequiredFields(c *gin.Context, requiredFields []string) error {
	for _, field := range requiredFields {
		value := c.Param(field)
		if value == "" {
			value = c.Query(field)
		}
		if value == "" {
			return errors.ErrInvalidInput
		}
	}
	return nil
}

// validateStringFields validates string fields according to rules
func validateStringFields(c *gin.Context, rules map[string]StringRule) error {
	for field, rule := range rules {
		value := c.Param(field)
		if value == "" {
			value = c.Query(field)
		}
		if value == "" {
			value = c.PostForm(field)
		}

		// Check if required
		if rule.Required && value == "" {
			return errors.ErrInvalidInput
		}

		// Skip validation if not required and empty
		if !rule.Required && value == "" {
			continue
		}

		// Validate length
		if rule.MinLength > 0 && len(value) < rule.MinLength {
			return errors.ErrInvalidInput
		}
		if rule.MaxLength > 0 && len(value) > rule.MaxLength {
			return errors.ErrInvalidInput
		}

		// Validate pattern
		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, value)
			if err != nil || !matched {
				return errors.ErrInvalidInput
			}
		}

		// Sanitize input
		sanitizedValue := sanitizeString(value)
		if field != "" {
			c.Set(field, sanitizedValue)
		}
	}
	return nil
}

// validateNumberFields validates numeric fields according to rules
func validateNumberFields(c *gin.Context, rules map[string]NumberRule) error {
	for field, rule := range rules {
		value := c.Param(field)
		if value == "" {
			value = c.Query(field)
		}
		if value == "" {
			value = c.PostForm(field)
		}

		// Check if required
		if rule.Required && value == "" {
			return errors.ErrInvalidInput
		}

		// Skip validation if not required and empty
		if !rule.Required && value == "" {
			continue
		}

		// Parse number
		var num float64
		var err error
		if rule.Integer {
			var intVal int64
			intVal, err = strconv.ParseInt(value, 10, 64)
			num = float64(intVal)
		} else {
			num, err = strconv.ParseFloat(value, 64)
		}

		if err != nil {
			return errors.ErrInvalidInput
		}

		// Validate range
		if rule.Min != 0 && num < rule.Min {
			return errors.ErrInvalidInput
		}
		if rule.Max != 0 && num > rule.Max {
			return errors.ErrInvalidInput
		}

		// Store sanitized value
		if field != "" {
			c.Set(field, num)
		}
	}
	return nil
}

// validateEmailFields validates email fields
func validateEmailFields(c *gin.Context, rules map[string]EmailRule) error {
	validate := validator.New()

	for field, rule := range rules {
		value := c.Param(field)
		if value == "" {
			value = c.Query(field)
		}
		if value == "" {
			value = c.PostForm(field)
		}

		// Check if required
		if rule.Required && value == "" {
			return errors.ErrInvalidInput
		}

		// Skip validation if not required and empty
		if !rule.Required && value == "" {
			continue
		}

		// Validate email format
		if err := validate.Var(value, "email"); err != nil {
			return errors.ErrInvalidInput
		}

		// Sanitize email
		sanitizedEmail := strings.ToLower(strings.TrimSpace(value))
		if field != "" {
			c.Set(field, sanitizedEmail)
		}
	}
	return nil
}

// validateDateFields validates date fields
func validateDateFields(c *gin.Context, rules map[string]DateRule) error {
	for field, rule := range rules {
		value := c.Param(field)
		if value == "" {
			value = c.Query(field)
		}
		if value == "" {
			value = c.PostForm(field)
		}

		// Check if required
		if rule.Required && value == "" {
			return errors.ErrInvalidInput
		}

		// Skip validation if not required and empty
		if !rule.Required && value == "" {
			continue
		}

		// Validate date format (basic validation)
		if !isValidDate(value) {
			return errors.ErrInvalidInput
		}

		// Store sanitized value
		if field != "" {
			c.Set(field, value)
		}
	}
	return nil
}

// sanitizeString removes potentially dangerous characters
func sanitizeString(input string) string {
	// Remove null bytes and control characters
	re := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	sanitized := re.ReplaceAllString(input, "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// isValidDate performs basic date validation
func isValidDate(dateStr string) bool {
	// Basic date format validation (YYYY-MM-DD)
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return datePattern.MatchString(dateStr)
}

// Common validation rules
var (
	// ESG Score validation rules
	ESGScoreValidation = ValidationRules{
		RequiredFields: []string{"company_id"},
		NumberRules: map[string]NumberRule{
			"environmental_score": {Min: 0, Max: 100, Required: false},
			"social_score":        {Min: 0, Max: 100, Required: false},
			"governance_score":    {Min: 0, Max: 100, Required: false},
			"overall_score":       {Min: 0, Max: 100, Required: false},
		},
		DateRules: map[string]DateRule{
			"score_date": {Required: true},
		},
	}

	// Company validation rules
	CompanyValidation = ValidationRules{
		RequiredFields: []string{"name", "symbol"},
		StringRules: map[string]StringRule{
			"name":     {MinLength: 1, MaxLength: 255, Required: true},
			"symbol":   {MinLength: 1, MaxLength: 20, Required: true, Pattern: `^[A-Z0-9]+$`},
			"sector":   {MinLength: 0, MaxLength: 100, Required: false},
			"industry": {MinLength: 0, MaxLength: 100, Required: false},
			"country":  {MinLength: 0, MaxLength: 100, Required: false},
		},
		NumberRules: map[string]NumberRule{
			"market_cap": {Min: 0, Required: false},
		},
	}

	// Pagination validation rules
	PaginationValidation = ValidationRules{
		NumberRules: map[string]NumberRule{
			"limit":  {Min: 1, Max: 100, Required: false, Integer: true},
			"offset": {Min: 0, Required: false, Integer: true},
		},
	}

	// User validation rules
	UserValidation = ValidationRules{
		RequiredFields: []string{"email"},
		StringRules: map[string]StringRule{
			"first_name": {MinLength: 1, MaxLength: 100, Required: false},
			"last_name":  {MinLength: 1, MaxLength: 100, Required: false},
		},
		EmailRules: map[string]EmailRule{
			"email": {Required: true},
		},
	}
)
