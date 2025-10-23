package middleware

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// PerformanceMetrics stores performance data
type PerformanceMetrics struct {
	RequestCount    int64         `json:"request_count"`
	ErrorCount      int64         `json:"error_count"`
	AverageResponse time.Duration `json:"average_response_time"`
	MaxResponse     time.Duration `json:"max_response_time"`
	MinResponse     time.Duration `json:"min_response_time"`
	TotalResponse   time.Duration `json:"total_response_time"`
	MemoryUsage     MemoryStats   `json:"memory_usage"`
	LastUpdated     time.Time     `json:"last_updated"`
	mu              sync.RWMutex  `json:"-"`
}

// MemoryStats tracks memory usage
type MemoryStats struct {
	Alloc      uint64 `json:"alloc"`      // Currently allocated memory
	TotalAlloc uint64 `json:"total_alloc"` // Total allocated memory
	Sys        uint64 `json:"sys"`        // System memory
	NumGC      uint32 `json:"num_gc"`     // Number of GC runs
	GCPause    uint64 `json:"gc_pause"`   // GC pause time
}

// RouteMetrics tracks performance per route
type RouteMetrics struct {
	Path            string        `json:"path"`
	Method          string        `json:"method"`
	RequestCount    int64         `json:"request_count"`
	ErrorCount      int64         `json:"error_count"`
	AverageResponse time.Duration `json:"average_response_time"`
	MaxResponse     time.Duration `json:"max_response_time"`
	MinResponse     time.Duration `json:"min_response_time"`
	LastHit         time.Time     `json:"last_hit"`
}

// PerformanceMonitor manages performance monitoring
type PerformanceMonitor struct {
	globalMetrics   PerformanceMetrics
	routeMetrics    map[string]*RouteMetrics
	enabledRoutes   map[string]bool
	slowQueryThreshold time.Duration
	mu             sync.RWMutex
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		routeMetrics:        make(map[string]*RouteMetrics),
		enabledRoutes:       make(map[string]bool),
		slowQueryThreshold:  500 * time.Millisecond,
	}
}

// EnableRouteMonitoring enables monitoring for specific routes
func (pm *PerformanceMonitor) EnableRouteMonitoring(pattern string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabledRoutes[pattern] = true
}

// SetSlowQueryThreshold sets the threshold for slow query detection
func (pm *PerformanceMonitor) SetSlowQueryThreshold(threshold time.Duration) {
	pm.slowQueryThreshold = threshold
}

// PerformanceMonitoringMiddleware provides comprehensive performance monitoring
func (pm *PerformanceMonitor) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		responseTime := time.Since(startTime)

		// Update metrics
		pm.updateMetrics(c, responseTime)

		// Check for slow query
		if responseTime > pm.slowQueryThreshold {
			pm.handleSlowQuery(c, responseTime)
		}

		// Add performance headers
		c.Header("X-Response-Time", responseTime.String())
		c.Header("X-Memory-Usage", fmt.Sprintf("%d", getCurrentMemoryUsage()))
	}
}

// updateMetrics updates performance metrics
func (pm *PerformanceMonitor) updateMetrics(c *gin.Context, responseTime time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Update global metrics
	pm.globalMetrics.RequestCount++
	if c.Writer.Status() >= 400 {
		pm.globalMetrics.ErrorCount++
	}

	pm.globalMetrics.TotalResponse += responseTime
	pm.globalMetrics.AverageResponse = time.Duration(
		int64(pm.globalMetrics.TotalResponse) / pm.globalMetrics.RequestCount)

	if responseTime > pm.globalMetrics.MaxResponse {
		pm.globalMetrics.MaxResponse = responseTime
	}

	if pm.globalMetrics.MinResponse == 0 || responseTime < pm.globalMetrics.MinResponse {
		pm.globalMetrics.MinResponse = responseTime
	}

	pm.globalMetrics.LastUpdated = time.Now()

	// Update memory stats
	pm.updateMemoryStats()

	// Update route-specific metrics
	routeKey := c.Request.Method + ":" + c.FullPath()
	routeMetric, exists := pm.routeMetrics[routeKey]
	if !exists {
		routeMetric = &RouteMetrics{
			Path:   c.FullPath(),
			Method: c.Request.Method,
		}
		pm.routeMetrics[routeKey] = routeMetric
	}

	routeMetric.RequestCount++
	if c.Writer.Status() >= 400 {
		routeMetric.ErrorCount++
	}

	if routeMetric.RequestCount == 1 {
		routeMetric.AverageResponse = responseTime
		routeMetric.MaxResponse = responseTime
		routeMetric.MinResponse = responseTime
	} else {
		routeMetric.AverageResponse = time.Duration(
			(int64(routeMetric.AverageResponse)*(routeMetric.RequestCount-1) + int64(responseTime)) / routeMetric.RequestCount)

		if responseTime > routeMetric.MaxResponse {
			routeMetric.MaxResponse = responseTime
		}

		if responseTime < routeMetric.MinResponse {
			routeMetric.MinResponse = responseTime
		}
	}

	routeMetric.LastHit = time.Now()
}

// updateMemoryStats updates memory usage statistics
func (pm *PerformanceMonitor) updateMemoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	pm.globalMetrics.MemoryUsage = MemoryStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
		GCPause:    m.PauseTotalNs,
	}
}

// handleSlowQuery handles slow query detection and logging
func (pm *PerformanceMonitor) handleSlowQuery(c *gin.Context, responseTime time.Duration) {
	// Log slow query
	slowQueryLog := map[string]interface{}{
		"timestamp":     time.Now(),
		"method":        c.Request.Method,
		"path":          c.FullPath(),
		"response_time": responseTime.String(),
		"status_code":   c.Writer.Status(),
		"user_agent":    c.Request.UserAgent(),
		"ip":           c.ClientIP(),
		"request_id":   GetRequestID(c),
	}

	logData, _ := json.Marshal(slowQueryLog)
	gin.DefaultErrorWriter.Write([]byte("SLOW QUERY: " + string(logData) + "\n"))
}

// GetGlobalMetrics returns global performance metrics
func (pm *PerformanceMonitor) GetGlobalMetrics() PerformanceMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// Return a copy to prevent concurrent access issues
	return PerformanceMetrics{
		RequestCount:    pm.globalMetrics.RequestCount,
		ErrorCount:      pm.globalMetrics.ErrorCount,
		AverageResponse: pm.globalMetrics.AverageResponse,
		MaxResponse:     pm.globalMetrics.MaxResponse,
		MinResponse:     pm.globalMetrics.MinResponse,
		TotalResponse:   pm.globalMetrics.TotalResponse,
		MemoryUsage:     pm.globalMetrics.MemoryUsage,
		LastUpdated:     pm.globalMetrics.LastUpdated,
	}
}

// GetRouteMetrics returns metrics for all routes
func (pm *PerformanceMonitor) GetRouteMetrics() map[string]RouteMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]RouteMetrics)
	for key, metric := range pm.routeMetrics {
		result[key] = RouteMetrics{
			Path:            metric.Path,
			Method:          metric.Method,
			RequestCount:    metric.RequestCount,
			ErrorCount:      metric.ErrorCount,
			AverageResponse: metric.AverageResponse,
			MaxResponse:     metric.MaxResponse,
			MinResponse:     metric.MinResponse,
			LastHit:         metric.LastHit,
		}
	}

	return result
}

// GetMetricsJSON returns metrics as JSON
func (pm *PerformanceMonitor) GetMetricsJSON() (string, error) {
	metrics := struct {
		Global PerformanceMetrics           `json:"global"`
		Routes  map[string]RouteMetrics    `json:"routes"`
		Summary struct {
			TotalRoutes    int     `json:"total_routes"`
			ErrorRate      float64 `json:"error_rate"`
			AvgResponse    string  `json:"avg_response"`
			MemoryUsageMB  float64 `json:"memory_usage_mb"`
		} `json:"summary"`
	}{
		Global: pm.GetGlobalMetrics(),
		Routes: pm.GetRouteMetrics(),
	}

	// Calculate summary
	metrics.Summary.TotalRoutes = len(metrics.Routes)
	if metrics.Global.RequestCount > 0 {
		metrics.Summary.ErrorRate = float64(metrics.Global.ErrorCount) / float64(metrics.Global.RequestCount) * 100
	}
	metrics.Summary.AvgResponse = metrics.Global.AverageResponse.String()
	metrics.Summary.MemoryUsageMB = float64(metrics.Global.MemoryUsage.Alloc) / 1024 / 1024

	jsonData, err := json.MarshalIndent(metrics, "", "  ")
	return string(jsonData), err
}

// ResetMetrics resets all performance metrics
func (pm *PerformanceMonitor) ResetMetrics() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.globalMetrics = PerformanceMetrics{}
	pm.routeMetrics = make(map[string]*RouteMetrics)
}

// getCurrentMemoryUsage returns current memory usage in bytes
func getCurrentMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// MemoryOptimizationMiddleware provides memory optimization features
func MemoryOptimizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Force garbage collection before request if memory usage is high
		if getCurrentMemoryUsage() > 100*1024*1024 { // 100MB threshold
			runtime.GC()
		}

		c.Next()

		// Force garbage collection after request if needed
		if getCurrentMemoryUsage() > 200*1024*1024 { // 200MB threshold
			go runtime.GC()
		}
	}
}

// GCPool manages a pool of reusable objects to reduce GC pressure
type GCPool struct {
	// Add pools for frequently used objects here
	// For example: byte buffer pools, JSON encoder pools, etc.
}

// NewGCPool creates a new GC pool
func NewGCPool() *GCPool {
	return &GCPool{}
}

// Global performance monitor instance
var globalPerformanceMonitor *PerformanceMonitor
var performanceMonitorOnce sync.Once

// GetGlobalPerformanceMonitor returns the global performance monitor instance
func GetGlobalPerformanceMonitor() *PerformanceMonitor {
	performanceMonitorOnce.Do(func() {
		globalPerformanceMonitor = NewPerformanceMonitor()
	})
	return globalPerformanceMonitor
}

// PerformanceMiddleware returns a configured performance monitoring middleware
func PerformanceMiddleware() gin.HandlerFunc {
	return GetGlobalPerformanceMonitor().Middleware()
}

// MetricsHandler provides an HTTP endpoint for metrics
func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pm := GetGlobalPerformanceMonitor()

		acceptHeader := c.GetHeader("Accept")
		if acceptHeader == "application/json" {
			metricsJSON, err := pm.GetMetricsJSON()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.Header("Content-Type", "application/json")
			c.String(200, metricsJSON)
		} else {
			// Return HTML format
			html := generateMetricsHTML(pm)
			c.Header("Content-Type", "text/html")
			c.String(200, html)
		}
	}
}

// generateMetricsHTML generates an HTML dashboard for metrics
func generateMetricsHTML(pm *PerformanceMonitor) string {
	globalMetrics := pm.GetGlobalMetrics()
	routeMetrics := pm.GetRouteMetrics()

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Performance Metrics Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .metric { margin: 10px 0; padding: 10px; border: 1px solid #ddd; }
        .metric h3 { margin: 0 0 10px 0; color: #333; }
        table { width: 100%; border-collapse: collapse; margin: 10px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .error { color: red; }
        .slow { color: orange; }
    </style>
    <script>
        setTimeout(function(){ location.reload(); }, 5000); // Auto refresh every 5 seconds
    </script>
</head>
<body>
    <h1>Performance Metrics Dashboard</h1>

    <div class="metric">
        <h3>Global Metrics</h3>
        <p><strong>Total Requests:</strong> ` + fmt.Sprintf("%d", globalMetrics.RequestCount) + `</p>
        <p><strong>Error Count:</strong> ` + fmt.Sprintf("%d", globalMetrics.ErrorCount) + `</p>
        <p><strong>Average Response Time:</strong> ` + globalMetrics.AverageResponse.String() + `</p>
        <p><strong>Max Response Time:</strong> ` + globalMetrics.MaxResponse.String() + `</p>
        <p><strong>Min Response Time:</strong> ` + globalMetrics.MinResponse.String() + `</p>
        <p><strong>Memory Usage:</strong> ` + fmt.Sprintf("%.2f MB", float64(globalMetrics.MemoryUsage.Alloc)/1024/1024) + `</p>
        <p><strong>Last Updated:</strong> ` + globalMetrics.LastUpdated.Format("2006-01-02 15:04:05") + `</p>
    </div>

    <div class="metric">
        <h3>Route Metrics</h3>
        <table>
            <tr>
                <th>Method</th>
                <th>Path</th>
                <th>Requests</th>
                <th>Errors</th>
                <th>Avg Response</th>
                <th>Max Response</th>
                <th>Last Hit</th>
            </tr>`

	for _, metric := range routeMetrics {
		errorRate := "0%"
		if metric.RequestCount > 0 {
			errorRate = fmt.Sprintf("%.2f%%", float64(metric.ErrorCount)/float64(metric.RequestCount)*100)
		}

		html += `
            <tr>
                <td>` + metric.Method + `</td>
                <td>` + metric.Path + `</td>
                <td>` + fmt.Sprintf("%d", metric.RequestCount) + `</td>
                <td class="error">` + fmt.Sprintf("%d (%s)", metric.ErrorCount, errorRate) + `</td>
                <td>` + metric.AverageResponse.String() + `</td>
                <td>` + metric.MaxResponse.String() + `</td>
                <td>` + metric.LastHit.Format("15:04:05") + `</td>
            </tr>`
	}

	html += `
        </table>
    </div>

    <p><em>Page auto-refreshes every 5 seconds</em></p>
</body>
</html>`

	return html
}