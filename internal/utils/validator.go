package utils

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hatena-bookmark-mcp/internal/types"
)

// Validator provides parameter validation functions
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateGetBookmarksParams validates the parameters for GetBookmarks
func (v *Validator) ValidateGetBookmarksParams(params types.GetHatenaBookmarksParams) error {
	// Validate username
	if err := v.ValidateUsername(params.Username); err != nil {
		return err
	}

	// Validate optional parameters
	if params.Tag != "" {
		if err := v.ValidateTag(params.Tag); err != nil {
			return err
		}
	}

	if params.Date != "" {
		if err := v.ValidateDate(params.Date); err != nil {
			return err
		}
	}

	if params.URL != "" {
		if err := v.ValidateURL(params.URL); err != nil {
			return err
		}
	}

	if err := v.ValidatePage(params.Page); err != nil {
		return err
	}

	return nil
}

// ValidateUsername validates the username parameter
func (v *Validator) ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	
	if username == "" {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Username is required",
			Details: map[string]interface{}{"field": "username"},
		}
	}

	// Username should be 1-50 characters
	if len(username) > 50 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Username must be 50 characters or less",
			Details: map[string]interface{}{"username": username, "length": len(username)},
		}
	}

	// Username should contain only alphanumeric characters and hyphens
	validUsernameRegex := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
	if !validUsernameRegex.MatchString(username) {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Username must contain only alphanumeric characters and hyphens",
			Details: map[string]interface{}{"username": username},
		}
	}

	return nil
}

// ValidateTag validates the tag parameter
func (v *Validator) ValidateTag(tag string) error {
	tag = strings.TrimSpace(tag)
	
	if len(tag) > 100 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Tag must be 100 characters or less",
			Details: map[string]interface{}{"tag": tag, "length": len(tag)},
		}
	}

	// Tags should not contain certain special characters
	invalidChars := []string{"<", ">", "\"", "'", "&"}
	for _, char := range invalidChars {
		if strings.Contains(tag, char) {
			return &types.MCPError{
				Code:    types.ErrorCodeValidation,
				Message: "Tag contains invalid characters",
				Details: map[string]interface{}{"tag": tag, "invalid_char": char},
			}
		}
	}

	return nil
}

// ValidateDate validates the date parameter (YYYYMMDD format)
func (v *Validator) ValidateDate(date string) error {
	date = strings.TrimSpace(date)
	
	// Check format: YYYYMMDD
	if len(date) != 8 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Date must be in YYYYMMDD format",
			Details: map[string]interface{}{"date": date, "expected_format": "YYYYMMDD"},
		}
	}

	// Check if all characters are digits
	for _, r := range date {
		if r < '0' || r > '9' {
			return &types.MCPError{
				Code:    types.ErrorCodeValidation,
				Message: "Date must contain only numeric characters",
				Details: map[string]interface{}{"date": date},
			}
		}
	}

	// Validate actual date values
	year, err := strconv.Atoi(date[:4])
	if err != nil || year < 1900 || year > time.Now().Year()+1 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Invalid year in date",
			Details: map[string]interface{}{"date": date, "year": year},
		}
	}

	month, err := strconv.Atoi(date[4:6])
	if err != nil || month < 1 || month > 12 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Invalid month in date",
			Details: map[string]interface{}{"date": date, "month": month},
		}
	}

	day, err := strconv.Atoi(date[6:8])
	if err != nil || day < 1 || day > 31 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Invalid day in date",
			Details: map[string]interface{}{"date": date, "day": day},
		}
	}

	// Additional validation: check if the date is actually valid
	_, err = time.Parse("20060102", date)
	if err != nil {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Invalid date",
			Details: map[string]interface{}{"date": date, "error": err.Error()},
		}
	}

	return nil
}

// ValidateURL validates the URL parameter
func (v *Validator) ValidateURL(urlStr string) error {
	urlStr = strings.TrimSpace(urlStr)
	
	if len(urlStr) > 2000 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "URL must be 2000 characters or less",
			Details: map[string]interface{}{"url": urlStr, "length": len(urlStr)},
		}
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Invalid URL format",
			Details: map[string]interface{}{"url": urlStr, "error": err.Error()},
		}
	}

	// URL must have scheme and host
	if parsedURL.Scheme == "" {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "URL must include scheme (http:// or https://)",
			Details: map[string]interface{}{"url": urlStr},
		}
	}

	if parsedURL.Host == "" {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "URL must include host",
			Details: map[string]interface{}{"url": urlStr},
		}
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "URL scheme must be http or https",
			Details: map[string]interface{}{"url": urlStr, "scheme": parsedURL.Scheme},
		}
	}

	return nil
}

// ValidatePage validates the page parameter
func (v *Validator) ValidatePage(page int) error {
	if page < 0 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Page number must be positive",
			Details: map[string]interface{}{"page": page},
		}
	}

	// Reasonable upper limit for page numbers
	if page > 10000 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Page number is too large (maximum: 10000)",
			Details: map[string]interface{}{"page": page},
		}
	}

	return nil
}