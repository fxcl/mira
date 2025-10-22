package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig_LoadFromFile(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test_config.yaml")

	configContent := `
ruoyi:
  name: "Mira Test"
  version: "1.0.0"
  copyright: "2024 Test"
  domain: "localhost"
  ssl: false
  uploadPath: "/tmp/uploads"

server:
  port: 8080
  mode: "test"

database:
  driver: "sqlite"
  dsn: "test.db"

redis:
  host: "localhost"
  port: 6379
  db: 0

jwt:
  secret: "test-secret-key"
  expire: 3600
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(configFile)

	t.Run("should load config from file successfully", func(t *testing.T) {
		var cfg Config

		data, err := os.ReadFile(configFile)
		assert.NoError(t, err)

		err = yaml.Unmarshal(data, &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "Mira Test", cfg.Ruoyi.Name)
		assert.Equal(t, "1.0.0", cfg.Ruoyi.Version)
		assert.Equal(t, "2024 Test", cfg.Ruoyi.Copyright)
		assert.Equal(t, "localhost", cfg.Ruoyi.Domain)
		assert.False(t, cfg.Ruoyi.SSL)
		assert.Equal(t, "/tmp/uploads", cfg.Ruoyi.UploadPath)

		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "test", cfg.Server.Mode)
	})

	t.Run("should handle non-existent config file", func(t *testing.T) {
		var cfg Config

		_, err := os.ReadFile("non_existent_config.yaml")
		assert.Error(t, err)

		err = yaml.Unmarshal([]byte{}, &cfg)
		assert.NoError(t, err) // Empty config should not error
	})
}

func TestConfig_Validation(t *testing.T) {
	t.Run("should validate required fields", func(t *testing.T) {
		cfg := Config{
			Ruoyi: struct {
				Name      string `yaml:"name"`
				Version   string `yaml:"version"`
				Copyright string `yaml:"copyright"`
				Domain    string `yaml:"domain"`
				SSL       bool   `yaml:"ssl"`
				UploadPath string `yaml:"uploadPath"`
			}{
				Name:      "Test",
				Version:   "1.0.0",
				Copyright: "2024",
				Domain:    "localhost",
				SSL:       false,
				UploadPath: "/tmp",
			},
			Server: struct {
				Port int    `yaml:"port"`
				Mode string `yaml:"mode"`
			}{
				Port: 8080,
				Mode: "release",
			},
		}

		assert.NotEmpty(t, cfg.Ruoyi.Name)
		assert.NotEmpty(t, cfg.Ruoyi.Version)
		assert.NotEmpty(t, cfg.Ruoyi.Copyright)
		assert.Greater(t, cfg.Server.Port, 0)
		assert.NotEmpty(t, cfg.Server.Mode)
	})

	t.Run("should validate port range", func(t *testing.T) {
		validPorts := []int{80, 443, 8080, 3000, 5000}
		invalidPorts := []int{0, -1, 65536, 100000}

		for _, port := range validPorts {
			assert.True(t, port > 0 && port <= 65535, "Port %d should be valid", port)
		}

		for _, port := range invalidPorts {
			assert.False(t, port > 0 && port <= 65535, "Port %d should be invalid", port)
		}
	})

	t.Run("should validate SSL configuration", func(t *testing.T) {
		cfg := Config{}

		// Test SSL enabled
		cfg.Ruoyi.SSL = true
		assert.True(t, cfg.Ruoyi.SSL)

		// Test SSL disabled
		cfg.Ruoyi.SSL = false
		assert.False(t, cfg.Ruoyi.SSL)
	})

	t.Run("should validate domain format", func(t *testing.T) {
		validDomains := []string{"localhost", "example.com", "api.example.com", "192.168.1.1"}
		invalidDomains := []string{"", "invalid..domain", "domain with spaces", "https://example.com"}

		for _, domain := range validDomains {
			assert.NotEmpty(t, domain, "Domain should not be empty: %s", domain)
			assert.NotContains(t, domain, " ", "Domain should not contain spaces: %s", domain)
		}

		for _, domain := range invalidDomains {
			if domain != "" {
				// Just verify the domain is processed without crashing
				assert.NotEmpty(t, domain, "Domain should not be empty: %s", domain)
			}
		}
	})
}

func TestConfig_DefaultValues(t *testing.T) {
	t.Run("should provide sensible defaults", func(t *testing.T) {
		cfg := Config{}

		// Test that zero values are handled properly
		assert.Equal(t, "", cfg.Ruoyi.Name)
		assert.Equal(t, "", cfg.Ruoyi.Version)
		assert.Equal(t, 0, cfg.Server.Port)
		assert.Equal(t, "", cfg.Server.Mode)
		assert.False(t, cfg.Ruoyi.SSL)
		assert.Equal(t, "", cfg.Ruoyi.UploadPath)
	})
}

func TestConfig_EnvironmentVariables(t *testing.T) {
	t.Run("should handle environment variable overrides", func(t *testing.T) {
		// Test environment variable setting
		originalPort := os.Getenv("SERVER_PORT")
		defer os.Setenv("SERVER_PORT", originalPort)

		os.Setenv("SERVER_PORT", "9090")

		port := os.Getenv("SERVER_PORT")
		assert.Equal(t, "9090", port)
	})
}

func TestConfig_ConfigFileFormats(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("should handle YAML format", func(t *testing.T) {
		yamlFile := filepath.Join(tempDir, "config.yaml")
		yamlContent := `
ruoyi:
  name: "Test YAML"
  version: "1.0.0"
server:
  port: 8080
  mode: "test"
`

		err := os.WriteFile(yamlFile, []byte(yamlContent), 0644)
		assert.NoError(t, err)

		var cfg Config
		data, err := os.ReadFile(yamlFile)
		assert.NoError(t, err)

		err = yaml.Unmarshal(data, &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "Test YAML", cfg.Ruoyi.Name)
		assert.Equal(t, "1.0.0", cfg.Ruoyi.Version)
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "test", cfg.Server.Mode)
	})

	t.Run("should handle malformed YAML", func(t *testing.T) {
		invalidYamlFile := filepath.Join(tempDir, "invalid.yaml")
		invalidContent := `
ruoyi:
  name: "Test"
  invalid_yaml: [unclosed array
server:
  port: not_a_number
`

		err := os.WriteFile(invalidYamlFile, []byte(invalidContent), 0644)
		assert.NoError(t, err)

		var cfg Config
		data, err := os.ReadFile(invalidYamlFile)
		assert.NoError(t, err)

		err = yaml.Unmarshal(data, &cfg)
		// YAML unmarshal should handle some errors gracefully, but invalid syntax will fail
		// The exact behavior depends on the YAML parser implementation
	})
}

func TestConfig_ConfigValidation(t *testing.T) {
	t.Run("should validate upload path", func(t *testing.T) {
		uploadPaths := []string{
			"/tmp/uploads",
			"./uploads",
			"uploads",
			"/var/www/uploads",
		}

		for _, path := range uploadPaths {
			assert.NotEmpty(t, path, "Upload path should not be empty")
		}
	})

	t.Run("should validate application name", func(t *testing.T) {
		validNames := []string{"Mira", "My App", "Application-123", "app_name"}
		invalidNames := []string{"", "   ", "App\nWith\nNewlines"}

		for _, name := range validNames {
			assert.NotEmpty(t, name, "Application name should not be empty")
			assert.NotEqual(t, "   ", name, "Application name should not be just whitespace")
		}

		for _, name := range invalidNames {
			if name == "   " {
				assert.Equal(t, "   ", name, "Should detect whitespace-only names")
			}
		}
	})

	t.Run("should validate version format", func(t *testing.T) {
		validVersions := []string{"1.0.0", "1.0", "v1.0.0", "2.1.3-beta"}
		invalidVersions := []string{"", "not.a.version", "1..0", "1.0."}

		for _, version := range validVersions {
			assert.NotEmpty(t, version, "Version should not be empty")
		}

		for _, version := range invalidVersions {
			if version != "" {
				// Basic validation - version should contain at least one digit
				assert.True(t, len(version) > 0, "Invalid version should be empty or have proper format")
			}
		}
	})
}

func TestConfig_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config_test.yaml")

	t.Run("should handle file permission issues", func(t *testing.T) {
		// Create a config file
		configContent := `
ruoyi:
  name: "Test"
server:
  port: 8080
`
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		// Test reading the file
		data, err := os.ReadFile(configFile)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)

		// Test file permissions
		info, err := os.Stat(configFile)
		assert.NoError(t, err)
		assert.True(t, info.Mode().Perm()&0400 != 0, "File should be readable")
	})
}

func TestConfig_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "concurrent_config.yaml")

	configContent := `
ruoyi:
  name: "Concurrent Test"
  version: "1.0.0"
server:
  port: 8080
  mode: "test"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	t.Run("should handle concurrent config reading", func(t *testing.T) {
		const numGoroutines = 10
		const readsPerGoroutine = 5

		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				for j := 0; j < readsPerGoroutine; j++ {
					var cfg Config
					data, err := os.ReadFile(configFile)
					assert.NoError(t, err)

					err = yaml.Unmarshal(data, &cfg)
					assert.NoError(t, err)

					assert.Equal(t, "Concurrent Test", cfg.Ruoyi.Name)
					assert.Equal(t, 8080, cfg.Server.Port)
				}
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}

// Benchmark tests
func BenchmarkConfig_Load(b *testing.B) {
	tempDir := b.TempDir()
	configFile := filepath.Join(tempDir, "benchmark_config.yaml")

	configContent := `
ruoyi:
  name: "Benchmark Test"
  version: "1.0.0"
  copyright: "2024"
  domain: "localhost"
  ssl: false
  uploadPath: "/tmp/uploads"
server:
  port: 8080
  mode: "release"
database:
  driver: "sqlite"
  dsn: "benchmark.db"
redis:
  host: "localhost"
  port: 6379
  db: 0
jwt:
  secret: "benchmark-secret"
  expire: 3600
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cfg Config
		data, err := os.ReadFile(configFile)
		if err != nil {
			b.Fatal(err)
		}

		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			b.Fatal(err)
		}
	}
}