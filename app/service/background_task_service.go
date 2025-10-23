package service

import (
	"context"
	"fmt"
	"time"
)

// TaskType represents different types of background tasks
type TaskType string

const (
	TaskTypeSendEmail     TaskType = "send_email"
	TaskTypeProcessData   TaskType = "process_data"
	TaskTypeGenerateReport TaskType = "generate_report"
	TaskTypeCleanupCache  TaskType = "cleanup_cache"
	TaskTypeAuditLog      TaskType = "audit_log"
)

// BaseTask provides common functionality for all background tasks
type BaseTask struct {
	ID        string
	Type      TaskType
	Priority  int
	RetryCount int
	CreatedAt time.Time
}

// BackgroundTask represents a generic background task
type BackgroundTask struct {
	BaseTask
	Data     map[string]interface{}
	Handler  func(ctx context.Context, data map[string]interface{}) error
}

// Execute implements the Task interface
func (bt *BackgroundTask) Execute(ctx context.Context) error {
	if bt.Handler == nil {
		return fmt.Errorf("task handler not defined for task %s", bt.ID)
	}
	return bt.Handler(ctx, bt.Data)
}

// GetID returns the task ID
func (bt *BackgroundTask) GetID() string {
	return bt.ID
}

// GetPriority returns the task priority
func (bt *BackgroundTask) GetPriority() int {
	return bt.Priority
}

// GetRetryCount returns the current retry count
func (bt *BackgroundTask) GetRetryCount() int {
	return bt.RetryCount
}

// BackgroundTaskService manages background task execution
type BackgroundTaskService struct {
	workerPool *WorkerPool
	cache      *CacheService
}

// NewBackgroundTaskService creates a new background task service
func NewBackgroundTaskService(workers int) *BackgroundTaskService {
	return &BackgroundTaskService{
		workerPool: NewWorkerPool(workers),
		cache:      NewCacheService(),
	}
}

// Start starts the background task service
func (bts *BackgroundTaskService) Start() error {
	return bts.workerPool.Start()
}

// Stop stops the background task service
func (bts *BackgroundTaskService) Stop(timeout time.Duration) error {
	return bts.workerPool.Stop(timeout)
}

// SubmitTask submits a background task for execution
func (bts *BackgroundTaskService) SubmitTask(taskType TaskType, data map[string]interface{}, priority int) error {
	task := &BackgroundTask{
		BaseTask: BaseTask{
			ID:        fmt.Sprintf("%s_%d", taskType, time.Now().UnixNano()),
			Type:      taskType,
			Priority:  priority,
			RetryCount: 0,
			CreatedAt: time.Now(),
		},
		Data: data,
	}

	// Set handler based on task type
	switch taskType {
	case TaskTypeSendEmail:
		task.Handler = bts.sendEmailHandler
	case TaskTypeProcessData:
		task.Handler = bts.processDataHandler
	case TaskTypeGenerateReport:
		task.Handler = bts.generateReportHandler
	case TaskTypeCleanupCache:
		task.Handler = bts.cleanupCacheHandler
	case TaskTypeAuditLog:
		task.Handler = bts.auditLogHandler
	default:
		return fmt.Errorf("unknown task type: %s", taskType)
	}

	return bts.workerPool.Submit(task)
}

// SubmitEmailTask submits an email sending task
func (bts *BackgroundTaskService) SubmitEmailTask(to, subject, body string) error {
	data := map[string]interface{}{
		"to":      to,
		"subject": subject,
		"body":    body,
	}
	return bts.SubmitTask(TaskTypeSendEmail, data, 1)
}

// SubmitDataProcessingTask submits a data processing task
func (bts *BackgroundTaskService) SubmitDataProcessingTask(dataType string, data interface{}) error {
	taskData := map[string]interface{}{
		"type": dataType,
		"data": data,
	}
	return bts.SubmitTask(TaskTypeProcessData, taskData, 2)
}

// SubmitReportGenerationTask submits a report generation task
func (bts *BackgroundTaskService) SubmitReportGenerationTask(reportType string, params map[string]interface{}) error {
	taskData := map[string]interface{}{
		"type":   reportType,
		"params": params,
	}
	return bts.SubmitTask(TaskTypeGenerateReport, taskData, 3)
}

// SubmitCacheCleanupTask submits a cache cleanup task
func (bts *BackgroundTaskService) SubmitCacheCleanupTask(pattern string) error {
	taskData := map[string]interface{}{
		"pattern": pattern,
	}
	return bts.SubmitTask(TaskTypeCleanupCache, taskData, 0) // Low priority
}

// SubmitAuditLogTask submits an audit logging task
func (bts *BackgroundTaskService) SubmitAuditLogTask(action string, userID int, details map[string]interface{}) error {
	taskData := map[string]interface{}{
		"action":  action,
		"user_id": userID,
		"details": details,
	}
	return bts.SubmitTask(TaskTypeAuditLog, taskData, 1)
}

// Task handlers
func (bts *BackgroundTaskService) sendEmailHandler(ctx context.Context, data map[string]interface{}) error {
	to, ok := data["to"].(string)
	if !ok {
		return fmt.Errorf("email recipient not specified")
	}

	subject, _ := data["subject"].(string)
	body, _ := data["body"].(string)

	// Simulate email sending (in production, integrate with actual email service)
	fmt.Printf("Sending email to %s: %s\n%s\n", to, subject, body)
	time.Sleep(100 * time.Millisecond) // Simulate network delay

	// Cache sent email record
	cacheKey := fmt.Sprintf("email_sent:%s:%d", to, time.Now().Unix())
	bts.cache.Set(ctx, cacheKey, map[string]interface{}{
		"to":       to,
		"subject":  subject,
		"sent_at":  time.Now(),
	}, 24*time.Hour)

	return nil
}

func (bts *BackgroundTaskService) processDataHandler(ctx context.Context, data map[string]interface{}) error {
	dataType, ok := data["type"].(string)
	if !ok {
		return fmt.Errorf("data type not specified")
	}

	// Simulate data processing
	fmt.Printf("Processing data of type: %s\n", dataType)
	time.Sleep(500 * time.Millisecond) // Simulate processing time

	// Cache processing result
	cacheKey := fmt.Sprintf("processed_data:%s:%d", dataType, time.Now().Unix())
	bts.cache.Set(ctx, cacheKey, map[string]interface{}{
		"type":       dataType,
		"processed":  true,
		"processed_at": time.Now(),
	}, 1*time.Hour)

	return nil
}

func (bts *BackgroundTaskService) generateReportHandler(ctx context.Context, data map[string]interface{}) error {
	reportType, ok := data["type"].(string)
	if !ok {
		return fmt.Errorf("report type not specified")
	}

	params, _ := data["params"].(map[string]interface{})

	// Simulate report generation
	fmt.Printf("Generating report of type: %s with params: %v\n", reportType, params)
	time.Sleep(2 * time.Second) // Simulate report generation time

	// Cache report result
	cacheKey := fmt.Sprintf("report:%s:%d", reportType, time.Now().Unix())
	bts.cache.Set(ctx, cacheKey, map[string]interface{}{
		"type":        reportType,
		"params":      params,
		"generated":   true,
		"generated_at": time.Now(),
	}, 6*time.Hour)

	return nil
}

func (bts *BackgroundTaskService) cleanupCacheHandler(ctx context.Context, data map[string]interface{}) error {
	pattern, ok := data["pattern"].(string)
	if !ok {
		return fmt.Errorf("cleanup pattern not specified")
	}

	// Simulate cache cleanup
	fmt.Printf("Cleaning up cache with pattern: %s\n", pattern)
	time.Sleep(200 * time.Millisecond) // Simulate cleanup time

	// Perform actual cache cleanup
	err := bts.cache.InvalidatePattern(ctx, pattern)
	if err != nil {
		return fmt.Errorf("failed to cleanup cache pattern %s: %w", pattern, err)
	}

	return nil
}

func (bts *BackgroundTaskService) auditLogHandler(ctx context.Context, data map[string]interface{}) error {
	action, ok := data["action"].(string)
	if !ok {
		return fmt.Errorf("audit action not specified")
	}

	userID, _ := data["user_id"].(int)
	details, _ := data["details"].(map[string]interface{})

	// Simulate audit logging
	fmt.Printf("Audit log: User %d performed action %s with details %v\n", userID, action, details)
	time.Sleep(50 * time.Millisecond) // Simulate logging time

	// Cache audit record
	cacheKey := fmt.Sprintf("audit_log:%s:%d:%d", action, userID, time.Now().Unix())
	bts.cache.Set(ctx, cacheKey, map[string]interface{}{
		"action":      action,
		"user_id":     userID,
		"details":     details,
		"logged_at":   time.Now(),
	}, 30*24*time.Hour) // Keep for 30 days

	return nil
}

// GetStats returns background task service statistics
func (bts *BackgroundTaskService) GetStats() PoolStats {
	return bts.workerPool.GetStats()
}

// GetMetrics returns detailed metrics
func (bts *BackgroundTaskService) GetMetrics() Metrics {
	return bts.workerPool.GetMetrics()
}

// GetQueueLength returns current queue length
func (bts *BackgroundTaskService) GetQueueLength() int {
	return bts.workerPool.GetQueueLength()
}

// IsRunning returns true if the service is running
func (bts *BackgroundTaskService) IsRunning() bool {
	return bts.workerPool.IsRunning()
}

// ResizeWorkerPool changes the number of workers
func (bts *BackgroundTaskService) ResizeWorkerPool(workers int) error {
	return bts.workerPool.Resize(workers)
}