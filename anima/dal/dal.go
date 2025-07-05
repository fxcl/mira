package dal

import "sync"

type Config struct {
	GormConfig  *GormConfig
	RedisConfig *RedisConfig
}

var once sync.Once

// Initialize the data access layer
func InitDal(config *Config) {
	once.Do(func() {
		// Initialize the database
		if config.GormConfig != nil {
			initGorm(config.GormConfig)
		}

		// Initialize Redis
		if config.RedisConfig != nil {
			initRedis(config.RedisConfig)
		}
	})
}
