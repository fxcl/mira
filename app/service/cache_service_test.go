package service

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"mira/common/types/redis-key"
)

func TestCacheService_SetAndGet(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	// Test Set and Get
	key := "test_key"
	value := map[string]interface{}{
		"name": "test_user",
		"id":   123,
	}

	err := cache.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	var result map[string]interface{}
	err = cache.Get(ctx, key, &result)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if result["name"] != value["name"] {
		t.Errorf("Expected name %v, got %v", value["name"], result["name"])
	}
	if result["id"] != value["id"] {
		t.Errorf("Expected id %v, got %v", value["id"], result["id"])
	}

	// Clean up
	err = cache.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete cache: %v", err)
	}
}

func TestCacheService_SetMultiple(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	items := []CacheItem{
		{
			Key:        "test_key_1",
			Value:      "value1",
			Expiration: 5 * time.Minute,
		},
		{
			Key:        "test_key_2",
			Value:      map[string]int{"count": 42},
			Expiration: 5 * time.Minute,
		},
	}

	err := cache.SetMultiple(ctx, items)
	assert.NoError(t, err)

	// Verify all items were set
	keys := []string{"test_key_1", "test_key_2"}
	results, err := cache.GetMultiple(ctx, keys)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))

	// Clean up
	err = cache.DeleteMultiple(ctx, keys)
	assert.NoError(t, err)
}

func TestCacheService_Exists(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	key := "test_exists_key"
	value := "test_value"

	// Should not exist initially
	exists, err := cache.Exists(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Set the value
	err = cache.Set(ctx, key, value, 5*time.Minute)
	assert.NoError(t, err)

	// Should exist now
	exists, err = cache.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Clean up
	err = cache.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestCacheService_GetWithFallback(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	key := "test_fallback_key"
	expectedValue := "fallback_value"

	// Define fallback function
	fallback := func() (interface{}, error) {
		return expectedValue, nil
	}

	var result string

	// First call should execute fallback (cache miss)
	err := cache.GetWithFallback(ctx, key, &result, 5*time.Minute, fallback)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)

	// Second call should use cache (cache hit)
	err = cache.GetWithFallback(ctx, key, &result, 5*time.Minute, fallback)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)

	// Clean up
	err = cache.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestCacheService_IncrementDecrement(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	key := "test_counter"

	// Initialize counter
	err := cache.Set(ctx, key, "0", 5*time.Minute)
	assert.NoError(t, err)

	// Increment
	newVal, err := cache.Increment(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), newVal)

	// Increment again
	newVal, err = cache.Increment(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), newVal)

	// Decrement
	newVal, err = cache.Decrement(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), newVal)

	// Clean up
	err = cache.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestCacheService_SetExpireAndGetTTL(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	key := "test_ttl_key"
	value := "test_value"

	// Set with 1 second expiration
	err := cache.Set(ctx, key, value, 1*time.Second)
	assert.NoError(t, err)

	// Check TTL (should be > 0)
	ttl, err := cache.GetTTL(ctx, key)
	assert.NoError(t, err)
	assert.True(t, ttl > 0)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Key should be expired now
	var result string
	err = cache.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestCacheService_InvalidateUserCache(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	userID := 123

	// Set some user cache keys
	keys := []string{
		rediskey.UserProfileKey(userID),
		rediskey.UserPermsKey(userID),
		rediskey.UserRolesKey(userID),
	}

	for _, key := range keys {
		err := cache.Set(ctx, key, "test_value", 5*time.Minute)
		assert.NoError(t, err)
	}

	// Invalidate user cache
	err := cache.InvalidateUserCache(ctx, userID)
	assert.NoError(t, err)

	// Verify all keys are deleted
	for _, key := range keys {
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}

func TestCacheService_Ping(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	err := cache.Ping(ctx)
	assert.NoError(t, err)
}

func TestCacheService_InvalidatePattern(t *testing.T) {
	cache := NewCacheService()
	ctx := context.Background()

	// Set multiple keys with pattern
	keys := []string{
		"test:pattern:key1",
		"test:pattern:key2",
		"test:other:key3",
	}

	for _, key := range keys {
		err := cache.Set(ctx, key, "value", 5*time.Minute)
		assert.NoError(t, err)
	}

	// Invalidate pattern
	err := cache.InvalidatePattern(ctx, "test:pattern:*")
	assert.NoError(t, err)

	// Check which keys were deleted
	exists, err := cache.Exists(ctx, "test:pattern:key1")
	assert.NoError(t, err)
	assert.False(t, exists)

	exists, err = cache.Exists(ctx, "test:pattern:key2")
	assert.NoError(t, err)
	assert.False(t, exists)

	exists, err = cache.Exists(ctx, "test:other:key3")
	assert.NoError(t, err)
	assert.True(t, exists) // Should not be deleted

	// Clean up
	err = cache.Delete(ctx, "test:other:key3")
	assert.NoError(t, err)
}

// Benchmark cache operations
func BenchmarkCacheService_Set(b *testing.B) {
	cache := NewCacheService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "bench_key_" + string(rune(i))
		err := cache.Set(ctx, key, "benchmark_value", 5*time.Minute)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCacheService_Get(b *testing.B) {
	cache := NewCacheService()
	ctx := context.Background()

	// Pre-populate cache
	for i := 0; i < 100; i++ {
		key := "bench_key_" + string(rune(i))
		cache.Set(ctx, key, "benchmark_value", 5*time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "bench_key_" + string(rune(i%100))
		var result string
		err := cache.Get(ctx, key, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCacheService_SetMultiple(b *testing.B) {
	cache := NewCacheService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items := make([]CacheItem, 10)
		for j := 0; j < 10; j++ {
			items[j] = CacheItem{
				Key:        "bench_key_" + string(rune(i*10+j)),
				Value:      "benchmark_value",
				Expiration: 5 * time.Minute,
			}
		}
		err := cache.SetMultiple(ctx, items)
		if err != nil {
			b.Fatal(err)
		}
	}
}