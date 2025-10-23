package tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mira/app/middleware"
	"mira/app/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// PerformanceIntegrationTestSuite tests all performance optimizations together
type PerformanceIntegrationTestSuite struct {
	suite.Suite
	router           *gin.Engine
	cacheService     *service.CacheService
	backgroundService *service.BackgroundTaskService
	perfMonitor      *middleware.PerformanceMonitor
}

func (suite *PerformanceIntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Initialize services
	suite.cacheService = service.NewCacheService()
	suite.backgroundService = service.NewBackgroundTaskService(5)
	suite.perfMonitor = middleware.NewPerformanceMonitor()

	// Start background service
	err := suite.backgroundService.Start()
	suite.Require().NoError(err)

	// Setup router with all performance middleware
	suite.router = gin.New()

	// Add performance monitoring middleware
	suite.router.Use(suite.perfMonitor.Middleware())
	suite.router.Use(middleware.ContextMiddleware(30 * time.Second))
	suite.router.Use(middleware.MemoryOptimizationMiddleware())
	suite.router.Use(middleware.WithRequestID())
	suite.router.Use(middleware.WithRequestTimestamp())

	// Add async logging middleware
	suite.router.Use(middleware.AsyncLoggingMiddleware(suite.backgroundService))

	// Add performance test routes
	suite.setupTestRoutes()

	// Add metrics endpoint
	suite.router.GET("/metrics", middleware.MetricsHandler())
}

func (suite *PerformanceIntegrationTestSuite) TearDownSuite() {
	if suite.backgroundService != nil {
		err := suite.backgroundService.Stop(5 * time.Second)
		suite.NoError(err)
	}
}

func (suite *PerformanceIntegrationTestSuite) setupTestRoutes() {
	// Fast endpoint
	suite.router.GET("/fast", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "fast response"})
	})

	// Slow endpoint
	suite.router.GET("/slow", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		c.JSON(200, gin.H{"message": "slow response"})
	})

	// Cache test endpoint
	suite.router.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		ctx := c.Request.Context()

		var value string
		err := suite.cacheService.GetWithFallback(ctx, key, &value, 5*time.Minute, func() (interface{}, error) {
			return "cached_value_" + key, nil
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"key": key, "value": value})
	})

	// Background task test endpoint
	suite.router.POST("/task", func(c *gin.Context) {
		var request struct {
			TaskType string                 `json:"task_type"`
			Data     map[string]interface{} `json:"data"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var taskType service.TaskType
		switch request.TaskType {
		case "send_email":
			taskType = service.TaskTypeSendEmail
		case "process_data":
			taskType = service.TaskTypeProcessData
		case "generate_report":
			taskType = service.TaskTypeGenerateReport
		default:
			c.JSON(400, gin.H{"error": "invalid task type"})
			return
		}

		err := suite.backgroundService.SubmitTask(taskType, request.Data, 1)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "task submitted"})
	})

	// Data scope test endpoint
	suite.router.GET("/data-scope/:userId", func(c *gin.Context) {
		userID := 0
		fmt.Sscanf(c.Param("userId"), "%d", &userID)

		// Simulate data scope query
		c.JSON(200, gin.H{
			"user_id": userID,
			"data_scope": "custom",
			"accessible_depts": []int{1, 2, 3},
			"accessible_users": []int{1, 2, 3, 4, 5},
		})
	})

	// Memory test endpoint
	suite.router.GET("/memory-test", func(c *gin.Context) {
		// Allocate some memory
		data := make([][]byte, 1000)
		for i := range data {
			data[i] = make([]byte, 1024)
		}

		c.JSON(200, gin.H{
			"message": "memory allocated",
			"size_mb": len(data) * len(data[0]) / 1024 / 1024,
		})
	})
}

func (suite *PerformanceIntegrationTestSuite) TestCachePerformance() {
	ctx := context.Background()

	// Test cache set/get performance
	start := time.Now()
	key := fmt.Sprintf("test_key_%d", time.Now().UnixNano())

	// First call (cache miss)
	req := httptest.NewRequest("GET", "/cache/"+key, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	firstCallDuration := time.Since(start)
	suite.Equal(200, w.Code)

	// Second call (cache hit)
	start = time.Now()
	req = httptest.NewRequest("GET", "/cache/"+key, nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	secondCallDuration := time.Since(start)
	suite.Equal(200, w.Code)

	// Cache hit should be faster than cache miss
	suite.Less(secondCallDuration, firstCallDuration)

	fmt.Printf("Cache miss: %v, Cache hit: %v\n", firstCallDuration, secondCallDuration)
}

func (suite *PerformanceIntegrationTestSuite) TestBackgroundTaskPerformance() {
	// Test background task submission
	taskData := map[string]interface{}{
		"to":      "test@example.com",
		"subject": "Test Email",
		"body":    "This is a test email",
	}

	start := time.Now()
	req := httptest.NewRequest("POST", "/task", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = httptest.NewBodyFromString(fmt.Sprintf(`{
		"task_type": "send_email",
		"data": %s
	}`, fmt.Sprintf(`{"to":"%s","subject":"%s","body":"%s"}`, taskData["to"], taskData["subject"], taskData["body"])))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	submissionDuration := time.Since(start)
	suite.Equal(200, w.Code)

	// Task submission should be fast (non-blocking)
	suite.Less(submissionDuration, 50*time.Millisecond)

	// Wait a bit for background processing
	time.Sleep(200 * time.Millisecond)

	// Check background service stats
	stats := suite.backgroundService.GetStats()
	suite.Greater(stats.TasksSubmitted, int64(0))

	fmt.Printf("Task submission duration: %v\n", submissionDuration)
}

func (suite *PerformanceIntegrationTestSuite) TestPerformanceMonitoring() {
	// Make multiple requests to generate metrics
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/fast", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		suite.Equal(200, w.Code)
	}

	// Test slow request
	req := httptest.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)

	// Get metrics
	req = httptest.NewRequest("GET", "/metrics", nil)
	req.Header.Set("Accept", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)

	// Check metrics endpoint works
	suite.Contains(w.Body.String(), "global")
	suite.Contains(w.Body.String(), "routes")

	fmt.Printf("Metrics response length: %d\n", len(w.Body.String()))
}

func (suite *PerformanceIntegrationTestSuite) TestContextMiddleware() {
	req := httptest.NewRequest("GET", "/fast", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
	suite.NotEmpty(w.Header().Get("X-Request-ID"))
}

func (suite *PerformanceIntegrationTestSuite) TestMemoryOptimization() {
	// Test memory usage endpoint
	req := httptest.NewRequest("GET", "/memory-test", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	// Multiple requests should trigger memory optimization
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/memory-test", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		suite.Equal(200, w.Code)
	}
}

func (suite *PerformanceIntegrationTestSuite) TestConcurrentRequests() {
	const numRequests = 100
	const numGoroutines = 10

	requestChan := make(chan int, numRequests)
	resultChan := make(chan time.Duration, numRequests)

	// Generate requests
	for i := 0; i < numRequests; i++ {
		requestChan <- i
	}
	close(requestChan)

	// Start goroutines
	start := time.Now()
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for range requestChan {
				reqStart := time.Now()
				req := httptest.NewRequest("GET", "/fast", nil)
				w := httptest.NewRecorder()
				suite.router.ServeHTTP(w, req)
				reqDuration := time.Since(reqStart)
				resultChan <- reqDuration
			}
		}()
	}

	// Collect results
	var totalDuration time.Duration
	var maxDuration time.Duration
	var minDuration time.Duration = time.Hour

	for i := 0; i < numRequests; i++ {
		duration := <-resultChan
		totalDuration += duration
		if duration > maxDuration {
			maxDuration = duration
		}
		if duration < minDuration {
			minDuration = duration
		}
	}

	totalTime := time.Since(start)
	avgDuration := totalDuration / time.Duration(numRequests)

	fmt.Printf("Concurrent test - Total: %v, Avg: %v, Min: %v, Max: %v\n",
		totalTime, avgDuration, minDuration, maxDuration)

	// Performance assertions
	suite.Less(avgDuration, 10*time.Millisecond)
	suite.Less(maxDuration, 100*time.Millisecond)
}

func (suite *PerformanceIntegrationTestSuite) TestCacheInvalidation() {
	ctx := context.Background()
	key := "invalidation_test"

	// Set cache value
	err := suite.cacheService.Set(ctx, key, "test_value", 1*time.Minute)
	suite.NoError(err)

	// Verify value exists
	var value string
	err = suite.cacheService.Get(ctx, key, &value)
	suite.NoError(err)
	suite.Equal("test_value", value)

	// Invalidate cache
	err = suite.cacheService.Delete(ctx, key)
	suite.NoError(err)

	// Verify value is gone
	err = suite.cacheService.Get(ctx, key, &value)
	suite.Error(err)
}

func (suite *PerformanceIntegrationTestSuite) TestWorkerPoolStats() {
	// Submit some tasks to generate stats
	for i := 0; i < 5; i++ {
		taskData := map[string]interface{}{
			"type": fmt.Sprintf("test_type_%d", i),
			"data": fmt.Sprintf("test_data_%d", i),
		}

		err := suite.backgroundService.SubmitTask(service.TaskTypeProcessData, taskData, 1)
		suite.NoError(err)
	}

	// Wait for processing
	time.Sleep(1 * time.Second)

	// Check stats
	stats := suite.backgroundService.GetStats()
	suite.Greater(stats.TasksSubmitted, int64(0))

	metrics := suite.backgroundService.GetMetrics()
	suite.NotEmpty(metrics.TaskCounters)

	fmt.Printf("Worker pool stats - Submitted: %d, Completed: %d, Failed: %d\n",
		stats.TasksSubmitted, stats.TasksCompleted, stats.TasksFailed)
}

// Benchmark tests
func BenchmarkCacheOperations(b *testing.B) {
	cache := service.NewCacheService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)

		cache.Set(ctx, key, value, 5*time.Minute)

		var retrieved string
		cache.Get(ctx, key, &retrieved)
	}
}

func BenchmarkBackgroundTaskSubmission(b *testing.B) {
	bgService := service.NewBackgroundTaskService(10)
	err := bgService.Start()
	if err != nil {
		b.Fatal(err)
	}
	defer bgService.Stop(5 * time.Second)

	taskData := map[string]interface{}{
		"test": "benchmark",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bgService.SubmitTask(service.TaskTypeProcessData, taskData, 1)
	}
}

func BenchmarkConcurrentRequests(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ContextMiddleware(30 * time.Second))
	router.GET("/bench", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "benchmark"})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/bench", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

// Test runner
func TestPerformanceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(PerformanceIntegrationTestSuite))
}

// Helper function to create request body
func httptest.NewBodyFromString(s string) *bytes.Buffer {
	return bytes.NewBufferString(s)
}