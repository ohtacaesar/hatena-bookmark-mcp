package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"hatena-bookmark-mcp/internal/parser"
	"hatena-bookmark-mcp/internal/types"
)

// BookmarkService handles Hatena Bookmark API interactions
type BookmarkService struct {
	baseURL    string
	logger     *slog.Logger
	client     *http.Client
	rssParser  *parser.RSSParser
}

// NewBookmarkService creates a new bookmark service instance
func NewBookmarkService(logger *slog.Logger) *BookmarkService {
	return &BookmarkService{
		baseURL: "https://b.hatena.ne.jp",
		logger:  logger,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		rssParser: parser.NewRSSParser(logger),
	}
}

// GetBookmarks retrieves bookmarks from Hatena Bookmark RSS feed
func (s *BookmarkService) GetBookmarks(ctx context.Context, params types.GetHatenaBookmarksParams) (*types.GetHatenaBookmarksResponse, error) {
	s.logger.Info("Getting bookmarks", 
		"username", params.Username,
		"tag", params.Tag,
		"date", params.Date,
		"url", params.URL,
		"page", params.Page)

	// Validate parameters
	if err := s.validateParams(params); err != nil {
		return nil, err
	}

	// Build request URL
	requestURL := s.buildRequestURL(params)
	s.logger.Debug("Built request URL", "url", requestURL)

	// Make HTTP request
	xmlContent, err := s.fetchRSSFeed(ctx, requestURL)
	if err != nil {
		return nil, err
	}

	// Parse RSS content
	parsedData, err := s.rssParser.ParseRSSFeed(ctx, xmlContent)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &types.GetHatenaBookmarksResponse{
		User:       params.Username,
		Page:       s.getPageOrDefault(params.Page),
		TotalCount: len(parsedData.Items),
		Bookmarks:  parsedData.Items,
	}

	// Add filters if any were applied
	if params.Tag != "" || params.Date != "" || params.URL != "" {
		response.Filters = &types.FilterParams{
			Tag:  params.Tag,
			Date: params.Date,
			URL:  params.URL,
		}
	}

	s.logger.Info("Successfully retrieved bookmarks", 
		"username", params.Username,
		"count", len(parsedData.Items))

	return response, nil
}

// validateParams validates the input parameters
func (s *BookmarkService) validateParams(params types.GetHatenaBookmarksParams) error {
	if strings.TrimSpace(params.Username) == "" {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Username is required",
			Details: map[string]interface{}{"field": "username"},
		}
	}

	// Validate username format (alphanumeric and hyphens only)
	if !isValidUsername(params.Username) {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Username must contain only alphanumeric characters and hyphens",
			Details: map[string]interface{}{"username": params.Username},
		}
	}

	// Validate date format if provided
	if params.Date != "" && !isValidDateFormat(params.Date) {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Date must be in YYYYMMDD format",
			Details: map[string]interface{}{"date": params.Date},
		}
	}

	// Validate URL format if provided
	if params.URL != "" && !isValidURL(params.URL) {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Invalid URL format",
			Details: map[string]interface{}{"url": params.URL},
		}
	}

	// Validate page number
	if params.Page < 0 {
		return &types.MCPError{
			Code:    types.ErrorCodeValidation,
			Message: "Page number must be positive",
			Details: map[string]interface{}{"page": params.Page},
		}
	}

	return nil
}

// buildRequestURL constructs the RSS feed URL with query parameters
func (s *BookmarkService) buildRequestURL(params types.GetHatenaBookmarksParams) string {
	// Base URL: https://b.hatena.ne.jp/{username}/rss
	baseURL := fmt.Sprintf("%s/%s/rss", s.baseURL, params.Username)

	// Build query parameters
	query := url.Values{}

	if params.Tag != "" {
		query.Set("tag", params.Tag)
	}

	if params.Date != "" {
		query.Set("date", params.Date)
	}

	if params.URL != "" {
		query.Set("url", params.URL)
	}

	if params.Page > 1 {
		query.Set("page", strconv.Itoa(params.Page))
	}

	if len(query) > 0 {
		return baseURL + "?" + query.Encode()
	}

	return baseURL
}

// fetchRSSFeed makes HTTP request to get RSS content
func (s *BookmarkService) fetchRSSFeed(ctx context.Context, requestURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, &types.MCPError{
			Code:    types.ErrorCodeNetwork,
			Message: fmt.Sprintf("Failed to create request: %v", err),
			Details: map[string]interface{}{"url": requestURL},
		}
	}

	// Set User-Agent to be respectful
	req.Header.Set("User-Agent", "hatena-bookmark-mcp/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, &types.MCPError{
			Code:    types.ErrorCodeNetwork,
			Message: fmt.Sprintf("Failed to fetch RSS feed: %v", err),
			Details: map[string]interface{}{"url": requestURL},
		}
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Debug("Failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, &types.MCPError{
			Code:    types.ErrorCodeAPI,
			Message: fmt.Sprintf("API returned status %d", resp.StatusCode),
			Details: map[string]interface{}{
				"status_code": resp.StatusCode,
				"url":         requestURL,
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &types.MCPError{
			Code:    types.ErrorCodeNetwork,
			Message: fmt.Sprintf("Failed to read response body: %v", err),
			Details: map[string]interface{}{"url": requestURL},
		}
	}

	return body, nil
}

// getPageOrDefault returns the page number or default value
func (s *BookmarkService) getPageOrDefault(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

// Validation helper functions

func isValidUsername(username string) bool {
	// Username should contain only alphanumeric characters and hyphens
	for _, r := range username {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return len(username) > 0
}

func isValidDateFormat(date string) bool {
	// Check if date is in YYYYMMDD format
	if len(date) != 8 {
		return false
	}
	
	for _, r := range date {
		if r < '0' || r > '9' {
			return false
		}
	}
	
	// Additional validation could be added here to check if it's a valid date
	return true
}

func isValidURL(urlStr string) bool {
	// Basic URL validation
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}