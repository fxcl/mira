package service

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a unit of work to be executed
type Task interface {
	Execute(ctx context.Context) error
	GetID() string
	GetPriority() int
	GetRetryCount() int
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	TaskID    string
	Success   bool
	Error     error
	Duration  time.Duration
	Retries   int
	Timestamp time.Time
}

// WorkerPool manages a pool of goroutines for concurrent task execution
type WorkerPool struct {
	workers      int
	taskQueue    chan Task
	resultQueue  chan TaskResult
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	running      int32
	stats        *PoolStats
	metrics      *Metrics
}

// PoolStats tracks worker pool statistics
type PoolStats struct {
	TasksSubmitted   int64
	TasksCompleted   int64
	TasksFailed      int64
	TasksRetry       int64
	AverageWaitTime  time.Duration
	AverageExecTime  time.Duration
	ActiveWorkers    int32
	mu               sync.RWMutex
}

// Metrics tracks detailed performance metrics
type Metrics struct {
	TaskCounters     map[string]int64
	ErrorCounters    map[string]int64
	ExecutionTimes   map[string][]time.Duration
	mu               sync.RWMutex
}

// NewWorkerPool creates a new worker pool with specified number of workers
func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workers:     workers,
		taskQueue:   make(chan Task, workers*10), // Buffered queue
		resultQueue: make(chan TaskResult, workers*10),
		ctx:         ctx,
		cancel:      cancel,
		stats:       &PoolStats{},
		metrics:     &Metrics{
			TaskCounters:   make(map[string]int64),
			ErrorCounters:  make(map[string]int64),
			ExecutionTimes: make(map[string][]time.Duration),
		},
	}
}

// Start initializes and starts the worker pool
func (wp *WorkerPool) Start() error {
	if !atomic.CompareAndSwapInt32(&wp.running, 0, 1) {
		return fmt.Errorf("worker pool is already running")
	}

	wp.wg.Add(wp.workers)
	for i := 0; i < wp.workers; i++ {
		go wp.worker(i)
	}

	// Start result processor
	go wp.resultProcessor()

	return nil
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop(timeout time.Duration) error {
	if !atomic.CompareAndSwapInt32(&wp.running, 1, 0) {
		return fmt.Errorf("worker pool is not running")
	}

	// Cancel context to signal workers to stop
	wp.cancel()

	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("worker pool shutdown timeout after %v", timeout)
	}
}

// Submit adds a task to the work queue
func (wp *WorkerPool) Submit(task Task) error {
	if atomic.LoadInt32(&wp.running) != 1 {
		return fmt.Errorf("worker pool is not running")
	}

	atomic.AddInt64(&wp.stats.TasksSubmitted, 1)

	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	default:
		return fmt.Errorf("task queue is full")
	}
}

// SubmitWithPriority adds a task to the queue with priority handling
func (wp *WorkerPool) SubmitWithPriority(task Task) error {
	if atomic.LoadInt32(&wp.running) != 1 {
		return fmt.Errorf("worker pool is not running")
	}

	// For now, use regular submission. In a real implementation,
	// you might want to implement priority queues
	return wp.Submit(task)
}

// worker processes tasks from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	atomic.AddInt32(&wp.stats.ActiveWorkers, 1)
	defer atomic.AddInt32(&wp.stats.ActiveWorkers, -1)

	for {
		select {
		case task := <-wp.taskQueue:
			wp.processTask(task)
		case <-wp.ctx.Done():
			return
		}
	}
}

// processTask executes a single task and handles retries
func (wp *WorkerPool) processTask(task Task) {
	startTime := time.Now()
	maxRetries := 3
	var lastError error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * time.Second
			select {
			case <-time.After(backoff):
			case <-wp.ctx.Done():
				return
			}
		}

		err := task.Execute(wp.ctx)
		if err == nil {
			// Task succeeded
			duration := time.Since(startTime)
			wp.resultQueue <- TaskResult{
				TaskID:    task.GetID(),
				Success:   true,
				Duration:  duration,
				Retries:   attempt,
				Timestamp: time.Now(),
			}

			atomic.AddInt64(&wp.stats.TasksCompleted, 1)
			wp.updateMetrics(task.GetID(), true, duration)
			return
		}

		lastError = err
		if attempt < maxRetries {
			atomic.AddInt64(&wp.stats.TasksRetry, 1)
		}
	}

	// Task failed after all retries
	duration := time.Since(startTime)
	wp.resultQueue <- TaskResult{
		TaskID:    task.GetID(),
		Success:   false,
		Error:     lastError,
		Duration:  duration,
		Retries:   maxRetries,
		Timestamp: time.Now(),
	}

	atomic.AddInt64(&wp.stats.TasksFailed, 1)
	wp.updateMetrics(task.GetID(), false, duration)
}

// resultProcessor processes task results
func (wp *WorkerPool) resultProcessor() {
	for {
		select {
		case result := <-wp.resultQueue:
			wp.processResult(result)
		case <-wp.ctx.Done():
			return
		}
	}
}

// processResult handles a single task result
func (wp *WorkerPool) processResult(result TaskResult) {
	wp.updateStats(result)

	// Log result (in production, you might want to use structured logging)
	if !result.Success {
		fmt.Printf("Task %s failed after %d retries: %v\n",
			result.TaskID, result.Retries, result.Error)
	}
}

// updateStats updates pool statistics
func (wp *WorkerPool) updateStats(result TaskResult) {
	wp.stats.mu.Lock()
	defer wp.stats.mu.Unlock()

	// Update average execution time
	totalCompleted := atomic.LoadInt64(&wp.stats.TasksCompleted)
	if totalCompleted > 0 {
		wp.stats.AverageExecTime = time.Duration(
			(int64(wp.stats.AverageExecTime)*totalCompleted + int64(result.Duration)) / (totalCompleted + 1))
	}
}

// updateMetrics updates detailed metrics
func (wp *WorkerPool) updateMetrics(taskID string, success bool, duration time.Duration) {
	wp.metrics.mu.Lock()
	defer wp.metrics.mu.Unlock()

	taskType := "default"
	if len(taskID) > 10 {
		taskType = taskID[:10] // Use first 10 chars as task type
	}

	// Update counters
	if success {
		wp.metrics.TaskCounters[taskType]++
	} else {
		wp.metrics.ErrorCounters[taskType]++
	}

	// Update execution times
	wp.metrics.ExecutionTimes[taskType] = append(wp.metrics.ExecutionTimes[taskType], duration)

	// Keep only last 100 execution times per task type
	if len(wp.metrics.ExecutionTimes[taskType]) > 100 {
		wp.metrics.ExecutionTimes[taskType] = wp.metrics.ExecutionTimes[taskType][1:]
	}
}

// GetStats returns current pool statistics
func (wp *WorkerPool) GetStats() PoolStats {
	wp.stats.mu.RLock()
	defer wp.stats.mu.RUnlock()

	return *wp.stats
}

// GetMetrics returns detailed metrics
func (wp *WorkerPool) GetMetrics() Metrics {
	wp.metrics.mu.RLock()
	defer wp.metrics.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	metrics := Metrics{
		TaskCounters:   make(map[string]int64),
		ErrorCounters:  make(map[string]int64),
		ExecutionTimes: make(map[string][]time.Duration),
	}

	for k, v := range wp.metrics.TaskCounters {
		metrics.TaskCounters[k] = v
	}
	for k, v := range wp.metrics.ErrorCounters {
		metrics.ErrorCounters[k] = v
	}
	for k, v := range wp.metrics.ExecutionTimes {
		metrics.ExecutionTimes[k] = append([]time.Duration{}, v...)
	}

	return metrics
}

// GetQueueLength returns current queue length
func (wp *WorkerPool) GetQueueLength() int {
	return len(wp.taskQueue)
}

// IsRunning returns true if the worker pool is running
func (wp *WorkerPool) IsRunning() bool {
	return atomic.LoadInt32(&wp.running) == 1
}

// Resize changes the number of workers in the pool
func (wp *WorkerPool) Resize(newWorkers int) error {
	if newWorkers <= 0 {
		return fmt.Errorf("number of workers must be positive")
	}

	if atomic.LoadInt32(&wp.running) != 1 {
		return fmt.Errorf("cannot resize stopped worker pool")
	}

	if newWorkers == wp.workers {
		return nil // No change needed
	}

	// For simplicity, we'll create a new pool with the new size
	// In a production implementation, you might want to add/remove workers dynamically
	wp.Stop(5 * time.Second)

	newPool := NewWorkerPool(newWorkers)
	*wp = *newPool
	return wp.Start()
}

// BatchSubmit submits multiple tasks concurrently
func (wp *WorkerPool) BatchSubmit(tasks []Task) error {
	if len(tasks) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))
	semaphore := make(chan struct{}, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := wp.Submit(t); err != nil {
				errChan <- err
			}
		}(task)
	}

	wg.Wait()
	close(errChan)

	// Return first error if any
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}