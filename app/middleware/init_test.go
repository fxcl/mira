package middleware

import (
	"os"
	"testing"

	"mira/anima/dal"
	"mira/config"

	"github.com/go-redis/redismock/v8"
)

// redisMock can be used by other tests in this package
var redisMock redismock.ClientMock

func TestMain(m *testing.M) {
	setupTest()
	code := m.Run()
	teardownTest()
	os.Exit(code)
}

// setupTest is a package-level setup function
func setupTest() {
	db, mock := redismock.NewClientMock()
	redisMock = mock
	dal.Redis = db
	config.Data = &config.Config{
		Token: struct {
			Header     string `yaml:"header"`
			Secret     string `yaml:"secret"`
			ExpireTime int    `yaml:"expireTime"`
		}{
			Header:     "Authorization",
			Secret:     "your-secret-key",
			ExpireTime: 30,
		},
	}
}

// teardownTest is a package-level teardown function
func teardownTest() {
	if dal.Redis != nil {
		dal.Redis.Close()
	}
}
