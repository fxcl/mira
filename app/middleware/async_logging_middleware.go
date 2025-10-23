package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"mira/app/service"
	responsewriter "mira/common/response-writer"

	"github.com/gin-gonic/gin"
)

// AsyncLoggingMiddleware provides asynchronous logging with background processing
func AsyncLoggingMiddleware(backgroundService *service.BackgroundTaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for health checks and static files
		if shouldSkipLogging(c) {
			c.Next()
			return
		}

		// Capture request start time
		startTime := time.Now()

		// Read and buffer request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create response writer to capture response
		rw := &responsewriter.ResponseWriter{
			ResponseWriter: c.Writer,
			Body:           bytes.NewBufferString(""),
		}
		c.Writer = rw

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Prepare log data
		logData := prepareLogData(c, bodyBytes, rw, startTime, duration)

		// Submit async logging task
		err := backgroundService.SubmitAuditLogTask(
			logData["action"].(string),
			getUserID(c),
			logData,
		)
		if err != nil {
			// Log error but don't fail the request
			gin.DefaultErrorWriter.Write([]byte("Failed to submit async log: " + err.Error() + "\n"))
		}
	}
}

// shouldSkipLogging determines if logging should be skipped for this request
func shouldSkipLogging(c *gin.Context) bool {
	// Skip health checks
	if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ping" {
		return true
	}

	// Skip static files
	if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/web" {
		return true
	}

	// Skip OPTIONS requests
	if c.Request.Method == "OPTIONS" {
		return true
	}

	return false
}

// prepareLogData prepares comprehensive log data
func prepareLogData(c *gin.Context, bodyBytes []byte, rw *responsewriter.ResponseWriter, startTime time.Time, duration time.Duration) map[string]interface{} {
	// Get user information
	userID := getUserID(c)
	username := getUsername(c)

	// Prepare request data
	requestData := map[string]interface{}{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.RawQuery,
		"user_agent": c.Request.UserAgent(),
		"ip":         c.ClientIP(),
		"headers":    sanitizeHeaders(c.Request.Header),
	}

	// Add request body if present and not too large
	if len(bodyBytes) > 0 && len(bodyBytes) < 1024 {
		requestData["body"] = string(bodyBytes)
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"status_code": c.Writer.Status(),
		"size":        rw.Body.Len(),
	}

	// Add response body if it's JSON and not too large
	if rw.Body.Len() > 0 && rw.Body.Len() < 1024 {
		contentType := c.Writer.Header().Get("Content-Type")
		if isJSONContentType(contentType) {
			var responseBody interface{}
			if err := json.Unmarshal(rw.Body.Bytes(), &responseBody); err == nil {
				responseData["body"] = responseBody
			}
		}
	}

	// Prepare performance data
	performanceData := map[string]interface{}{
		"start_time": startTime,
		"duration":   duration.Milliseconds(),
		"memory":     getMemoryUsage(),
	}

	// Prepare error data if there was an error
	errorData := make(map[string]interface{})
	if len(c.Errors) > 0 {
		lastError := c.Errors.Last()
		errorData = map[string]interface{}{
			"message": lastError.Error(),
			"type":    lastError.Type,
			"meta":    lastError.Meta,
		}
	}

	// Determine action based on request
	action := determineAction(c)

	return map[string]interface{}{
		"action":     action,
		"user_id":    userID,
		"username":   username,
		"request":    requestData,
		"response":   responseData,
		"performance": performanceData,
		"error":      errorData,
		"timestamp":  time.Now(),
	}
}

// sanitizeHeaders removes sensitive headers
func sanitizeHeaders(headers map[string][]string) map[string]string {
	sanitized := make(map[string]string)
	sensitiveHeaders := map[string]bool{
		"authorization":   true,
		"cookie":          true,
		"x-api-key":       true,
		"x-auth-token":    true,
		"x-csrf-token":    true,
		"set-cookie":      true,
	}

	for key, values := range headers {
		lowerKey := key
		if sensitiveHeaders[lowerKey] {
			sanitized[key] = "[REDACTED]"
		} else {
			// Join multiple values with comma
			sanitized[key] = joinStringSlice(values, ", ")
		}
	}

	return sanitized
}

// isJSONContentType checks if content type is JSON
func isJSONContentType(contentType string) bool {
	return contentType == "application/json" ||
		contentType == "application/hal+json" ||
		contentType == "application/vnd.api+json"
}

// joinStringSlice joins a slice of strings with a separator
func joinStringSlice(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}

	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
}

// getMemoryUsage returns current memory usage (simplified)
func getMemoryUsage() map[string]interface{} {
	// In a real implementation, you would use runtime.MemStats
	// For now, return placeholder data
	return map[string]interface{}{
		"allocated":   "N/A",
		"system":      "N/A",
		"gc_cycles":   "N/A",
	}
}

// getUserID extracts user ID from context
func getUserID(c *gin.Context) int {
	// Try to get user ID from common context keys
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int); ok {
			return id
		}
	}

	// Try to get from auth token (implement based on your auth system)
	if token, exists := c.Get("auth_token"); exists {
		// Parse token and extract user ID
		// This is a placeholder implementation
		if authData, ok := token.(map[string]interface{}); ok {
			if id, ok := authData["user_id"].(float64); ok {
				return int(id)
			}
		}
	}

	return 0 // Unknown user
}

// getUsername extracts username from context
func getUsername(c *gin.Context) string {
	// Try to get username from common context keys
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}

	// Try to get from auth token
	if token, exists := c.Get("auth_token"); exists {
		if authData, ok := token.(map[string]interface{}); ok {
			if name, ok := authData["username"].(string); ok {
				return name
			}
		}
	}

	return "unknown"
}

// determineAction determines the action type based on request
func determineAction(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path

	// Determine action based on HTTP method and path pattern
	switch {
	case method == "GET" && path == "/":
		return "page_view"
	case method == "GET" && len(path) > 4 && path[len(path)-4:] == "/list":
		return "list_data"
	case method == "GET" && containsNumericID(path):
		return "get_data"
	case method == "POST":
		return "create_data"
	case method == "PUT" || method == "PATCH":
		return "update_data"
	case method == "DELETE":
		return "delete_data"
	case method == "POST" && (path[len(path)-7:] == "/login" || path[len(path)-9:] == "/logout"):
		return "auth_" + path[len(path)-6:] // login or logout
	default:
		return method + "_" + path
	}
}

// containsNumericID checks if path contains a numeric ID
func containsNumericID(path string) bool {
	// Simple implementation - in production, use regex
	parts := splitPath(path)
	for _, part := range parts {
		if isNumeric(part) {
			return true
		}
	}
	return false
}

// splitPath splits URL path into parts
func splitPath(path string) []string {
	var parts []string
	current := ""

	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// isNumeric checks if string represents a number
func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}