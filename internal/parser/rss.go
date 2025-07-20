package parser

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"hatena-bookmark-mcp/internal/types"
)

// RSSParser handles RSS feed parsing
type RSSParser struct {
	logger *slog.Logger
}

// NewRSSParser creates a new RSS parser instance
func NewRSSParser(logger *slog.Logger) *RSSParser {
	return &RSSParser{
		logger: logger,
	}
}

// ParseRSSFeed parses RSS XML content and returns structured data
func (p *RSSParser) ParseRSSFeed(ctx context.Context, xmlContent []byte) (*types.ParsedRSSData, error) {
	p.logger.Debug("Starting RSS feed parsing", "content_length", len(xmlContent))

	var rss types.RSS
	if err := xml.Unmarshal(xmlContent, &rss); err != nil {
		p.logger.Error("Failed to unmarshal RSS XML", "error", err)
		return nil, &types.MCPError{
			Code:    types.ErrorCodeParsing,
			Message: fmt.Sprintf("Failed to parse RSS XML: %v", err),
			Details: map[string]interface{}{"xml_length": len(xmlContent)},
		}
	}

	bookmarks, err := p.extractBookmarkItems(&rss.Channel)
	if err != nil {
		p.logger.Error("Failed to extract bookmark items", "error", err)
		return nil, err
	}

	p.logger.Info("Successfully parsed RSS feed", 
		"title", rss.Channel.Title,
		"item_count", len(bookmarks))

	return &types.ParsedRSSData{
		Title:     rss.Channel.Title,
		Items:     bookmarks,
		ItemCount: len(bookmarks),
	}, nil
}

// extractBookmarkItems converts RSS items to bookmark items
func (p *RSSParser) extractBookmarkItems(channel *types.Channel) ([]types.BookmarkItem, error) {
	bookmarks := make([]types.BookmarkItem, 0, len(channel.Items))

	for _, item := range channel.Items {
		bookmark, err := p.convertItemToBookmark(item)
		if err != nil {
			p.logger.Warn("Failed to convert RSS item to bookmark", 
				"title", item.Title, 
				"error", err)
			continue
		}
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}

// convertItemToBookmark converts a single RSS item to a bookmark
func (p *RSSParser) convertItemToBookmark(item types.Item) (types.BookmarkItem, error) {
	// Parse the date
	bookmarkedAt, err := p.parseDate(item.PubDate)
	if err != nil {
		p.logger.Warn("Failed to parse date", "pubdate", item.PubDate, "error", err)
		bookmarkedAt = time.Now().Format(time.RFC3339)
	}

	// Extract tags from dc:subject elements
	tags := p.extractTags(item.Subjects)

	// Extract comment from description
	comment := p.extractComment(item.Description)

	return types.BookmarkItem{
		Title:        strings.TrimSpace(item.Title),
		URL:          strings.TrimSpace(item.Link),
		BookmarkedAt: bookmarkedAt,
		Tags:         tags,
		Comment:      comment,
	}, nil
}

// extractTags processes dc:subject elements to extract tag strings
func (p *RSSParser) extractTags(subjects []string) []string {
	tags := make([]string, 0, len(subjects))
	
	for _, subject := range subjects {
		tag := strings.TrimSpace(subject)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

// extractComment extracts user comment from RSS description
func (p *RSSParser) extractComment(description string) string {
	// Hatena Bookmark RSS often includes user comments in the description
	// Try to extract meaningful comment text
	
	// Remove HTML tags if any
	comment := p.stripHTMLTags(description)
	
	// Clean up and trim
	comment = strings.TrimSpace(comment)
	
	// If the comment is too long or seems to be just the article content,
	// it might not be a user comment
	if len(comment) > 500 {
		return ""
	}
	
	return comment
}

// stripHTMLTags removes HTML tags from text
func (p *RSSParser) stripHTMLTags(text string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(text, "")
}

// parseDate converts various date formats to ISO 8601
func (p *RSSParser) parseDate(dateString string) (string, error) {
	if dateString == "" {
		return time.Now().Format(time.RFC3339), nil
	}

	// Common RSS date formats to try
	formats := []string{
		time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC822,      // "02 Jan 06 15:04 MST"
		time.RFC822Z,     // "02 Jan 06 15:04 -0700"
		time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05", // Common alternative format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateString); err == nil {
			return t.Format(time.RFC3339), nil
		}
	}

	p.logger.Warn("Could not parse date, using current time", "original_date", dateString)
	return time.Now().Format(time.RFC3339), fmt.Errorf("could not parse date: %s", dateString)
}