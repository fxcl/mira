package security

import (
	"log"

	"mira/anima/dal"
	"mira/config"

	"github.com/go-redis/redismock/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var redisMock redismock.ClientMock

// setup initializes the test environment.
func setup() {
	db, mock := redismock.NewClientMock()
	redisMock = mock
	dal.Redis = db

	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dal.Gorm = gormDB

	// Initialize a minimal config for testing
	config.Data = &config.Config{
		Token: struct {
			Header     string `yaml:"header"`
			Secret     string `yaml:"secret"`
			ExpireTime int    `yaml:"expireTime"`
		}{
			Header:     "Authorization",
			Secret:     "test-secret",
			ExpireTime: 60,
		},
		Ruoyi: struct {
			Name       string `yaml:"name"`
			Version    string `yaml:"version"`
			Copyright  string `yaml:"copyright"`
			Domain     string `yaml:"domain"`
			SSL        bool   `yaml:"ssl"`
			UploadPath string `yaml:"uploadPath"`
		}{
			Name: "test",
		},
	}

	log.Println("Test environment initialized.")
}

// teardown cleans up the test environment.
func teardown() {
	if dal.Gorm != nil {
		db, _ := dal.Gorm.DB()
		db.Close()
	}
	if dal.Redis != nil {
		dal.Redis.Close()
	}
	log.Println("Test environment cleaned up.")
}
