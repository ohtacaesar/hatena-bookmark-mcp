package types

// GetHatenaBookmarksParams represents the parameters for the get_hatena_bookmarks tool
type GetHatenaBookmarksParams struct {
	Username string `json:"username"`           // Required: Hatena Bookmark username
	Tag      string `json:"tag,omitempty"`      // Optional: Filtering tag
	Date     string `json:"date,omitempty"`     // Optional: Date filter (YYYYMMDD)
	URL      string `json:"url,omitempty"`      // Optional: URL filter
	Page     int    `json:"page,omitempty"`     // Optional: Page number (default: 1)
}

// GetHatenaBookmarksResponse represents the response from the get_hatena_bookmarks tool
type GetHatenaBookmarksResponse struct {
	User       string          `json:"user"`
	Page       int             `json:"page"`
	TotalCount int             `json:"total_count"`
	Filters    *FilterParams   `json:"filters,omitempty"`
	Bookmarks  []BookmarkItem  `json:"bookmarks"`
}

// FilterParams represents the applied filters
type FilterParams struct {
	Tag  string `json:"tag,omitempty"`
	Date string `json:"date,omitempty"`
	URL  string `json:"url,omitempty"`
}

// BookmarkItem represents a single bookmark entry
type BookmarkItem struct {
	Title        string   `json:"title"`
	URL          string   `json:"url"`
	BookmarkedAt string   `json:"bookmarked_at"` // ISO 8601 format
	Tags         []string `json:"tags"`
	Comment      string   `json:"comment,omitempty"`
}

// RSS XML structure for parsing Hatena Bookmark RSS feeds
type RSS struct {
	XMLName string  `xml:"rss"`
	Version string  `xml:"version,attr"`
	Channel Channel `xml:"channel"`
}

// Channel represents the RSS channel
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents a single RSS item (bookmark)
type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Subjects    []string `xml:"http://purl.org/dc/elements/1.1/ subject"`
}

// ParsedRSSData represents the intermediate parsed RSS data
type ParsedRSSData struct {
	Title     string
	Items     []BookmarkItem
	ItemCount int
}

// Error types for better error handling
type ErrorCode string

const (
	ErrorCodeValidation ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNetwork    ErrorCode = "NETWORK_ERROR"
	ErrorCodeParsing    ErrorCode = "PARSING_ERROR"
	ErrorCodeAPI        ErrorCode = "API_ERROR"
)

// MCPError represents an error response for MCP
type MCPError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *MCPError) Error() string {
	return e.Message
}