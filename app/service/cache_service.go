package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"mira/anima/dal"
	"mira/common/types/redis-key"
)

// CacheService provides high-performance caching with Redis pipelines and TTL management
type CacheService struct{}

// NewCacheService creates a new cache service instance
func NewCacheService() *CacheService {
	return &CacheService{}
}

// CacheItem represents a cached item with metadata
type CacheItem struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
	CreatedAt  time.Time
}

// CacheResult represents the result of a cache operation
type CacheResult struct {
	Success bool
	Error   error
	Data    interface{}
}

// Set stores a value in cache with expiration
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return dal.Redis.Set(ctx, key, jsonValue, expiration).Err()
}

// Get retrieves a value from cache
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := dal.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from cache
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return dal.Redis.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := dal.Redis.Exists(ctx, key).Result()
	return count > 0, err
}

// SetMultiple stores multiple key-value pairs using Redis pipeline for performance
func (c *CacheService) SetMultiple(ctx context.Context, items []CacheItem) error {
	if len(items) == 0 {
		return nil
	}

	pipe := dal.Redis.Pipeline()

	for _, item := range items {
		jsonValue, err := json.Marshal(item.Value)
		if err != nil {
			return fmt.Errorf("failed to marshal cache value for key %s: %w", item.Key, err)
		}
		pipe.Set(ctx, item.Key, jsonValue, item.Expiration)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetMultiple retrieves multiple values from cache using Redis pipeline
func (c *CacheService) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	pipe := dal.Redis.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	results := make(map[string]interface{})
	for key, cmd := range cmds {
		val, err := cmd.Result()
		if err == nil {
			var dest interface{}
			if jsonErr := json.Unmarshal([]byte(val), &dest); jsonErr == nil {
				results[key] = dest
			}
		}
	}

	return results, nil
}

// DeleteMultiple removes multiple keys from cache using Redis pipeline
func (c *CacheService) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	pipe := dal.Redis.Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetWithFallback retrieves a value from cache, or executes the fallback function and caches the result
func (c *CacheService) GetWithFallback(ctx context.Context, key string, dest interface{}, expiration time.Duration, fallback func() (interface{}, error)) error {
	// Try to get from cache first
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil // Cache hit
	}

	// Cache miss, execute fallback
	value, err := fallback()
	if err != nil {
		return fmt.Errorf("fallback function failed: %w", err)
	}

	// Cache the result
	if setErr := c.Set(ctx, key, value, expiration); setErr != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Warning: failed to cache result for key %s: %v\n", key, setErr)
	}

	// Set the destination value
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal fallback result: %w", err)
	}

	return json.Unmarshal(jsonValue, dest)
}

// InvalidatePattern removes all keys matching a pattern
func (c *CacheService) InvalidatePattern(ctx context.Context, pattern string) error {
	keys, err := dal.Redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.DeleteMultiple(ctx, keys)
}

// Increment increments a numeric value in cache
func (c *CacheService) Increment(ctx context.Context, key string) (int64, error) {
	return dal.Redis.Incr(ctx, key).Result()
}

// Decrement decrements a numeric value in cache
func (c *CacheService) Decrement(ctx context.Context, key string) (int64, error) {
	return dal.Redis.Decr(ctx, key).Result()
}

// SetExpire sets expiration time for an existing key
func (c *CacheService) SetExpire(ctx context.Context, key string, expiration time.Duration) error {
	return dal.Redis.Expire(ctx, key, expiration).Err()
}

// GetTTL returns remaining time to live for a key
func (c *CacheService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return dal.Redis.TTL(ctx, key).Result()
}

// Ping checks Redis connection health
func (c *CacheService) Ping(ctx context.Context) error {
	return dal.Redis.Ping(ctx).Err()
}

// FlushAll removes all keys from current database (use with caution)
func (c *CacheService) FlushAll(ctx context.Context) error {
	return dal.Redis.FlushDB(ctx).Err()
}

// CacheStats provides cache statistics
type CacheStats struct {
	TotalKeys   int64
	MemoryUsage int64
	HitRate     float64
}

// GetStats returns cache statistics (basic implementation)
func (c *CacheService) GetStats(ctx context.Context) (*CacheStats, error) {
	info, err := dal.Redis.Info(ctx, "keyspace", "memory").Result()
	if err != nil {
		return nil, err
	}

	// Parse Redis info response (simplified)
	stats := &CacheStats{
		TotalKeys:   0,
		MemoryUsage: 0,
		HitRate:     0.0,
	}

	// This is a simplified implementation
	// In production, you might want to use Redis INFO command parsing library
	_ = info // Placeholder for info parsing

	return stats, nil
}

// InvalidateUserCache removes all cache entries for a specific user
func (c *CacheService) InvalidateUserCache(ctx context.Context, userID int) error {
	cacheKeys := []string{
		rediskey.UserProfileKey(userID),
		rediskey.UserPermsKey(userID),
		rediskey.UserRolesKey(userID),
		rediskey.UserDataScopeKey(userID),
		rediskey.UserSessionKey(userID),
		rediskey.UserAuthTokensKey(userID),
	}

	return c.DeleteMultiple(ctx, cacheKeys)
}