package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hatena-bookmark-mcp/internal/service"
	"hatena-bookmark-mcp/internal/types"
)

const (
	ServerName    = "hatena-bookmark-mcp"
	ServerVersion = "1.0.0"
)

// GetHatenaBookmarksParams represents the parameters for the tool
type GetHatenaBookmarksParams struct {
	Username string `json:"username"`
	Tag      string `json:"tag,omitempty"`
	Date     string `json:"date,omitempty"`
	URL      string `json:"url,omitempty"`
	Page     int    `json:"page,omitempty"`
}

func main() {
	// Initialize logger
	logger := initLogger()
	logger.Info("Starting Hatena Bookmark MCP Server", "version", ServerVersion)

	// Initialize services
	bookmarkService := service.NewBookmarkService(logger)

	// Create MCP server with implementation
	server := mcp.NewServer(&mcp.Implementation{
		Name:    ServerName,
		Version: ServerVersion,
	}, nil)

	// Register the get_hatena_bookmarks tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_hatena_bookmarks",
		Description: "Retrieve bookmarks from Hatena Bookmark RSS feed for a specified user with optional filtering",
	}, func(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[GetHatenaBookmarksParams]) (*mcp.CallToolResultFor[interface{}], error) {
		return handleGetBookmarks(ctx, params.Arguments, bookmarkService, logger)
	})

	logger.Info("Registered MCP tools", "tool_count", 1)

	// Start server with stdio transport
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// initLogger initializes the structured logger
func initLogger() *slog.Logger {
	// Get log level from environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	
	var level slog.Level
	switch logLevel {
	case "debug", "DEBUG":
		level = slog.LevelDebug
	case "warn", "WARN":
		level = slog.LevelWarn
	case "error", "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create logger with JSON handler for structured output
	opts := &slog.HandlerOptions{
		Level: level,
	}
	
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// handleGetBookmarks handles the get_hatena_bookmarks tool call
func handleGetBookmarks(
	ctx context.Context,
	arguments GetHatenaBookmarksParams,
	bookmarkService *service.BookmarkService,
	logger *slog.Logger,
) (*mcp.CallToolResultFor[interface{}], error) {
	logger.Debug("Handling get_hatena_bookmarks request", "arguments", arguments)

	// Convert to internal types
	params := types.GetHatenaBookmarksParams{
		Username: arguments.Username,
		Tag:      arguments.Tag,
		Date:     arguments.Date,
		URL:      arguments.URL,
		Page:     arguments.Page,
	}

	// Get bookmarks from service
	result, err := bookmarkService.GetBookmarks(ctx, params)
	if err != nil {
		logger.Error("Failed to get bookmarks", "error", err, "params", params)
		
		// Check if it's an MCP error
		if mcpErr, ok := err.(*types.MCPError); ok {
			return &mcp.CallToolResultFor[interface{}]{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: mcpErr.Message},
				},
			}, nil
		}
		
		// Generic error
		return &mcp.CallToolResultFor[interface{}]{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "An unexpected error occurred while fetching bookmarks"},
			},
		}, nil
	}

	logger.Info("Successfully retrieved bookmarks", 
		"username", params.Username,
		"bookmark_count", len(result.Bookmarks))

	return createSuccessResult(result), nil
}

// createSuccessResult creates a successful MCP tool result
func createSuccessResult(result *types.GetHatenaBookmarksResponse) *mcp.CallToolResultFor[interface{}] {
	// Convert result to JSON for display
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	
	return &mcp.CallToolResultFor[interface{}]{
		IsError: false,
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(resultJSON)},
		},
	}
}