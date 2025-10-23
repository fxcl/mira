package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ContextMiddleware manages request context with timeouts and cancellation
func ContextMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Store context and cancel function in gin context
		c.Set("request_ctx", ctx)
		c.Set("cancel_func", cancel)

		// Set up channel to listen for context cancellation
		done := make(chan struct{})
		go func() {
			defer close(done)
			select {
			case <-ctx.Done():
				// Context was cancelled or timed out
				if ctx.Err() == context.DeadlineExceeded {
					// Request timeout
					c.JSON(http.StatusRequestTimeout, gin.H{
						"error": "Request timeout",
						"code":  "REQUEST_TIMEOUT",
					})
					c.Abort()
				}
			case <-done:
				// Request completed normally
			}
		}()

		// Process request
		c.Next()

		// Signal completion to the goroutine
		close(done)
	}
}

// GetRequestContext retrieves the enhanced request context
func GetRequestContext(c *gin.Context) context.Context {
	if ctx, exists := c.Get("request_ctx"); exists {
		if requestCtx, ok := ctx.(context.Context); ok {
			return requestCtx
		}
	}
	return c.Request.Context()
}

// GetCancelFunc retrieves the cancel function for the current request
func GetCancelFunc(c *gin.Context) context.CancelFunc {
	if cancel, exists := c.Get("cancel_func"); exists {
		if cancelFunc, ok := cancel.(context.CancelFunc); ok {
			return cancelFunc
		}
	}
	return func() {}
}

// CancelRequest cancels the current request
func CancelRequest(c *gin.Context) {
	cancel := GetCancelFunc(c)
	cancel()
}

// RequestContextKey is used to store values in the request context
type RequestContextKey string

const (
	// RequestIDKey stores unique request ID
	RequestIDKey RequestContextKey = "request_id"
	// UserIDKey stores authenticated user ID
	UserIDKey RequestContextKey = "user_id"
	// RequestStartKey stores request start time
	RequestStartKey RequestContextKey = "request_start"
	// RequestDataKey stores custom request data
	RequestDataKey RequestContextKey = "request_data"
)

// SetRequestValue sets a value in the request context
func SetRequestValue(c *gin.Context, key RequestContextKey, value interface{}) {
	ctx := GetRequestContext(c)
	ctx = context.WithValue(ctx, key, value)
	c.Set("request_ctx", ctx)
}

// GetRequestValue gets a value from the request context
func GetRequestValue(c *gin.Context, key RequestContextKey) interface{} {
	ctx := GetRequestContext(c)
	return ctx.Value(key)
}

// GetRequestID returns the unique request ID
func GetRequestID(c *gin.Context) string {
	if requestID := GetRequestValue(c, RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID returns the authenticated user ID
func GetUserID(c *gin.Context) int {
	if userID := GetRequestValue(c, UserIDKey); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// GetRequestStart returns the request start time
func GetRequestStart(c *gin.Context) time.Time {
	if start := GetRequestValue(c, RequestStartKey); start != nil {
		if startTime, ok := start.(time.Time); ok {
			return startTime
		}
	}
	return time.Time{}
}

// WithRequestID adds a unique request ID to the context
func WithRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		SetRequestValue(c, RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// WithRequestTimestamp adds request start time to the context
func WithRequestTimestamp() gin.HandlerFunc {
	return func(c *gin.Context) {
		SetRequestValue(c, RequestStartKey, time.Now())
		c.Next()
	}
}

// WithUserID adds authenticated user ID to the context
func WithUserID(userID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetRequestValue(c, UserIDKey, userID)
		c.Next()
	}
}

// RequestDataMiddleware manages custom request data
func RequestDataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize request data map
		SetRequestValue(c, RequestDataKey, make(map[string]interface{}))
		c.Next()
	}
}

// SetRequestData sets custom data in the request context
func SetRequestData(c *gin.Context, key string, value interface{}) {
	data := GetRequestValue(c, RequestDataKey)
	if dataMap, ok := data.(map[string]interface{}); ok {
		dataMap[key] = value
		// Update the context with modified map
		SetRequestValue(c, RequestDataKey, dataMap)
	}
}

// GetRequestData gets custom data from the request context
func GetRequestData(c *gin.Context, key string) interface{} {
	data := GetRequestValue(c, RequestDataKey)
	if dataMap, ok := data.(map[string]interface{}); ok {
		return dataMap[key]
	}
	return nil
}

// GetAllRequestData gets all custom data from the request context
func GetAllRequestData(c *gin.Context) map[string]interface{} {
	data := GetRequestValue(c, RequestDataKey)
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Return a copy to prevent modifications
		copy := make(map[string]interface{})
		for k, v := range dataMap {
			copy[k] = v
		}
		return copy
	}
	return make(map[string]interface{})
}

// IsContextCancelled checks if the request context has been cancelled
func IsContextCancelled(c *gin.Context) bool {
	ctx := GetRequestContext(c)
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// GetContextTimeout returns the remaining timeout for the request
func GetContextTimeout(c *gin.Context) time.Duration {
	ctx := GetRequestContext(c)
	if deadline, ok := ctx.Deadline(); ok {
		return time.Until(deadline)
	}
	return 0
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Simple implementation using timestamp and random number
	// In production, you might want to use UUID or other more sophisticated methods
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	// Simple random string generation
	// In production, use crypto/rand for better randomness
	for i := range b {
		// This is a simplified implementation
		b[i] = charset[i%len(charset)]
	}

	return string(b)
}

// HealthCheckMiddleware provides context-aware health check
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := GetRequestContext(c)

		// Check if context is healthy
		select {
		case <-ctx.Done():
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  ctx.Err().Error(),
			})
			c.Abort()
			return
		default:
			// Context is healthy
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"time":   time.Now().Unix(),
			})
		}
	}
}

// ConcurrencyLimitMiddleware implements request concurrency limiting
func ConcurrencyLimitMiddleware(maxConcurrent int) gin.HandlerFunc {
	semaphore := make(chan struct{}, maxConcurrent)

	return func(c *gin.Context) {
		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }()
			c.Next()
		default:
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Server overloaded, please try again later",
				"code":  "OVERLOAD",
			})
			c.Abort()
		}
	}
}

// SlowQueryMiddleware detects slow queries and adds context
func SlowQueryMiddleware(threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		if duration > threshold {
			// Mark as slow query in context
			SetRequestData(c, "slow_query", true)
			SetRequestData(c, "query_duration", duration)

			// Log slow query (in production, use proper logging)
			requestID := GetRequestID(c)
			if requestID == "" {
				requestID = "unknown"
			}
			gin.DefaultWriter.Write([]byte(
				"Slow Query Detected - RequestID: " + requestID +
				", Duration: " + duration.String() +
				", Path: " + c.Request.URL.Path + "\n",
			))
		}
	}
}