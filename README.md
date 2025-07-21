# Hatena Bookmark MCP Server

A Model Context Protocol (MCP) server for accessing Hatena Bookmark RSS feeds. This server allows AI assistants to retrieve and analyze bookmarks from Hatena Bookmark users with various filtering options.

## Features

- **Get Bookmarks**: Retrieve bookmarks from any Hatena Bookmark user's RSS feed
- **Filtering**: Filter bookmarks by tag, date, or URL
- **Pagination**: Support for paginated results (20 bookmarks per page)
- **Structured Output**: Returns bookmarks in JSON format with metadata
- **Error Handling**: Comprehensive error handling with detailed error messages
- **Logging**: Structured logging with configurable levels

## Installation

### Prerequisites

- Go 1.21 or later
- Git

### Build from Source

```bash
git clone <repository-url>
cd hatena-bookmark-mcp
make build
```

### Run Tests

```bash
make test
```

### Development

```bash
# Install dependencies
make deps

# Run in development mode
make run

# Clean build artifacts
make clean
```

## Usage

### As MCP Server

Add the following configuration to your MCP client settings:

```json
{
  "mcpServers": {
    "hatena-bookmark": {
      "command": "./bin/hatena-bookmark-mcp",
      "args": [],
      "env": {
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

### Available Tools

#### `get_hatena_bookmarks`

Retrieve bookmarks from a Hatena Bookmark user's RSS feed.

**Parameters:**

- `username` (required): Hatena Bookmark username
- `tag` (optional): Filter bookmarks by tag
- `date` (optional): Filter bookmarks by date (YYYYMMDD format)
- `url` (optional): Filter bookmarks by URL
- `page` (optional): Page number for pagination (default: 1)

**Example Usage:**

```json
{
  "name": "get_hatena_bookmarks",
  "arguments": {
    "username": "sample",
    "tag": "programming",
    "page": 1
  }
}
```

**Response Format:**

```json
{
  "user": "sample",
  "page": 1,
  "total_count": 20,
  "filters": {
    "tag": "programming"
  },
  "bookmarks": [
    {
      "title": "Article Title",
      "url": "https://example.com/article",
      "bookmarked_at": "2025-01-20T10:30:00Z",
      "tags": ["programming", "go"],
      "comment": "User comment"
    }
  ]
}
```

## Configuration

### Environment Variables

- `LOG_LEVEL`: Set logging level (`debug`, `info`, `warn`, `error`) - Default: `info`

## API Limitations

### Hatena Bookmark RSS Feed Constraints

- Maximum 20 bookmarks per page
- Public bookmarks only (private bookmarks are not accessible)
- Rate limiting is implemented client-side to respect Hatena's servers

### Input Validation

- **Username**: 1-50 characters, alphanumeric and hyphens only
- **Tag**: Up to 100 characters, no special HTML characters
- **Date**: YYYYMMDD format, valid date range
- **URL**: Valid HTTP/HTTPS URLs only, up to 2000 characters
- **Page**: Positive integers up to 10,000

## Error Handling

The server provides detailed error messages for various scenarios:

- `VALIDATION_ERROR`: Invalid input parameters
- `NETWORK_ERROR`: Network connectivity issues
- `PARSING_ERROR`: RSS feed parsing failures
- `API_ERROR`: Hatena Bookmark API errors

## Development

### Project Structure

```
hatena-bookmark-mcp/
├── cmd/main.go              # Main application entry point
├── internal/
│   ├── service/bookmark.go  # Bookmark service (API interactions)
│   ├── parser/rss.go       # RSS feed parser
│   ├── types/bookmark.go   # Type definitions
│   ├── errors/handler.go   # Error handling utilities
│   └── utils/              # Utility functions
│       ├── validator.go    # Input validation
├── test/                   # Test files
└── Makefile               # Build automation
```

### Key Components

1. **BookmarkService**: Handles HTTP requests to Hatena Bookmark RSS feeds
2. **RSSParser**: Parses RSS XML and converts to structured data
3. **Validator**: Validates input parameters
5. **ErrorHandler**: Centralized error handling and logging

### Adding New Features

1. Define new tool in `cmd/main.go`
2. Add corresponding service methods in `internal/service/`
3. Update type definitions in `internal/types/`
4. Add validation logic in `internal/utils/validator.go`
5. Write tests in `test/`

## License

MIT License

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## Support

For issues and questions:
- Check existing issues in the repository
- Create a new issue with detailed information
- Include log output when reporting bugs
