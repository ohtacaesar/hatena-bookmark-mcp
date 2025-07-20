package errors

import (
	"log/slog"

	"hatena-bookmark-mcp/internal/types"
)

// ErrorHandler handles error processing and logging
type ErrorHandler struct {
	logger *slog.Logger
}

// NewErrorHandler creates a new error handler instance
func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleNetworkError processes network-related errors
func (h *ErrorHandler) HandleNetworkError(err error) *types.MCPError {
	h.logger.Error("Network error occurred", "error", err)
	
	return &types.MCPError{
		Code:    types.ErrorCodeNetwork,
		Message: "Network request failed",
		Details: map[string]interface{}{
			"original_error": err.Error(),
		},
	}
}

// HandleParsingError processes RSS parsing errors
func (h *ErrorHandler) HandleParsingError(err error) *types.MCPError {
	h.logger.Error("RSS parsing error occurred", "error", err)
	
	return &types.MCPError{
		Code:    types.ErrorCodeParsing,
		Message: "Failed to parse RSS feed",
		Details: map[string]interface{}{
			"original_error": err.Error(),
		},
	}
}

// HandleValidationError processes parameter validation errors
func (h *ErrorHandler) HandleValidationError(message string) *types.MCPError {
	h.logger.Warn("Validation error", "message", message)
	
	return &types.MCPError{
		Code:    types.ErrorCodeValidation,
		Message: message,
	}
}

// HandleAPIError processes API-related errors
func (h *ErrorHandler) HandleAPIError(statusCode int, message string) *types.MCPError {
	h.logger.Error("API error occurred", 
		"status_code", statusCode, 
		"message", message)
	
	return &types.MCPError{
		Code:    types.ErrorCodeAPI,
		Message: message,
		Details: map[string]interface{}{
			"status_code": statusCode,
		},
	}
}

// LogInfo logs informational messages
func (h *ErrorHandler) LogInfo(message string, args ...interface{}) {
	h.logger.Info(message, args...)
}

// LogDebug logs debug messages
func (h *ErrorHandler) LogDebug(message string, args ...interface{}) {
	h.logger.Debug(message, args...)
}

// LogWarn logs warning messages
func (h *ErrorHandler) LogWarn(message string, args ...interface{}) {
	h.logger.Warn(message, args...)
}