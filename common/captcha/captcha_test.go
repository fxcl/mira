package captcha

import (
	"testing"

	"github.com/mojocn/base64Captcha"
	"github.com/stretchr/testify/assert"
)

func TestNewCaptcha(t *testing.T) {
	captcha := NewCaptcha()

	assert.NotNil(t, captcha)
	assert.NotNil(t, captcha.captcha)
}

func TestCaptcha_Generate_Basic(t *testing.T) {
	t.Run("should create captcha driver successfully", func(t *testing.T) {
		driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
		assert.NotNil(t, driver)
	})

	t.Run("should create captcha with driver", func(t *testing.T) {
		driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
		store := &MockStore{}
		captcha := base64Captcha.NewCaptcha(driver, store)
		assert.NotNil(t, captcha)
	})
}

// MockStore for testing without Redis
type MockStore struct{}

func (m *MockStore) Set(id string, value string) error {
	// Mock implementation - does nothing
	return nil
}

func (m *MockStore) Get(id string, clear bool) string {
	// Mock implementation - returns empty
	return ""
}

func (m *MockStore) Verify(id, answer string, clear bool) bool {
	// Mock implementation - always returns false for testing
	return false
}

func TestCaptcha_WithMockStore(t *testing.T) {
	driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
	store := &MockStore{}
	captcha := base64Captcha.NewCaptcha(driver, store)

	t.Run("should generate captcha with mock store", func(t *testing.T) {
		id, b64s, _, err := captcha.Generate()

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
		assert.NotEmpty(t, b64s)
		assert.True(t, len(id) > 0)
		assert.True(t, len(b64s) > 0)
	})

	t.Run("should verify captcha with mock store", func(t *testing.T) {
		id, _, _, err := captcha.Generate()
		assert.NoError(t, err)

		result := captcha.Verify(id, "1234", true)
		assert.False(t, result) // Expected to return false due to mock implementation
	})

	t.Run("should handle non-existent captcha ID", func(t *testing.T) {
		nonExistentId := "non-existent-id"

		result := captcha.Verify(nonExistentId, "1234", true)
		assert.False(t, result)
	})
}

func TestCaptcha_DifferentParameters(t *testing.T) {
	t.Run("should work with different driver parameters", func(t *testing.T) {
		driverConfigs := []struct {
			height   int
			width    int
			length   int
			noise    float64
			showLine int
		}{
			{40, 100, 4, 0.7, 1},
			{50, 120, 5, 0.8, 2},
			{60, 140, 6, 0.9, 3},
		}

		for _, config := range driverConfigs {
			driver := base64Captcha.NewDriverDigit(config.height, config.width, config.length, config.noise, config.showLine)
			store := &MockStore{}
			captcha := base64Captcha.NewCaptcha(driver, store)

			id, b64s, _, err := captcha.Generate()

			assert.NoError(t, err, "Captcha generation should not fail for config: %+v", config)
			assert.NotEmpty(t, id, "Captcha ID should not be empty for config: %+v", config)
			assert.NotEmpty(t, b64s, "Captcha image should not be empty for config: %+v", config)
		}
	})
}

func TestCaptcha_ConcurrentGeneration(t *testing.T) {
	driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
	store := &MockStore{}
	captcha := base64Captcha.NewCaptcha(driver, store)

	t.Run("should handle concurrent generation", func(t *testing.T) {
		const numGoroutines = 5
		const captchasPerGoroutine = 3

		results := make(chan struct {
			id   string
			b64s string
		}, numGoroutines*captchasPerGoroutine)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				for j := 0; j < captchasPerGoroutine; j++ {
					id, b64s, _, err := captcha.Generate()
					if err == nil {
						results <- struct {
							id   string
							b64s string
						}{id: id, b64s: b64s}
					}
				}
			}()
		}

		// Collect results
		generatedCaptchas := make(map[string]bool)
		for i := 0; i < numGoroutines*captchasPerGoroutine; i++ {
			result := <-results

			// Verify all captchas are unique
			assert.False(t, generatedCaptchas[result.id], "Captcha ID should be unique: %s", result.id)
			generatedCaptchas[result.id] = true

			// Verify basic format
			assert.NotEmpty(t, result.id)
			assert.NotEmpty(t, result.b64s)
		}
	})
}

func TestCaptcha_Performance(t *testing.T) {
	driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
	store := &MockStore{}
	captcha := base64Captcha.NewCaptcha(driver, store)

	t.Run("should generate captchas efficiently", func(t *testing.T) {
		const numCaptchas = 10

		for i := 0; i < numCaptchas; i++ {
			id, b64s, _, err := captcha.Generate()
			assert.NoError(t, err)
			assert.NotEmpty(t, id)
			assert.NotEmpty(t, b64s)
		}
	})
}

func TestCaptcha_EdgeCases(t *testing.T) {
	t.Run("should handle multiple rapid generations", func(t *testing.T) {
		driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
		store := &MockStore{}
		captcha := base64Captcha.NewCaptcha(driver, store)

		ids := make([]string, 5)

		for i := 0; i < 5; i++ {
			id, _, _, err := captcha.Generate()
			assert.NoError(t, err)
			ids[i] = id
		}

		// Verify all IDs are unique
		idSet := make(map[string]bool)
		for _, id := range ids {
			assert.False(t, idSet[id], "Duplicate captcha ID found: %s", id)
			idSet[id] = true
		}
	})
}

// Benchmark tests
func BenchmarkCaptcha_Generate(b *testing.B) {
	driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
	store := &MockStore{}
	captcha := base64Captcha.NewCaptcha(driver, store)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		captcha.Generate()
	}
}

func BenchmarkCaptcha_Verify(b *testing.B) {
	driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)
	store := &MockStore{}
	captcha := base64Captcha.NewCaptcha(driver, store)

	// Pre-generate a captcha ID for benchmarking
	id, _, _, err := captcha.Generate()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		captcha.Verify(id, "1234", true)
	}
}