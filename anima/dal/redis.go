package dal

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host               string
	Port               int
	Database           int
	Password           string
	PoolSize           int
	MinIdleConns       int
	MaxRetries         int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

var Redis *redis.Client

func initRedis(config *RedisConfig) {
	// Set default values for performance
	if config.PoolSize == 0 {
		config.PoolSize = 50 // Connection pool size
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 10 // Minimum idle connections
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}
	if config.PoolTimeout == 0 {
		config.PoolTimeout = 4 * time.Second
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 5 * time.Minute
	}
	if config.IdleCheckFrequency == 0 {
		config.IdleCheckFrequency = time.Minute
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:               config.Host + ":" + strconv.Itoa(config.Port),
		Password:           config.Password,
		DB:                 config.Database,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
		// Enable connection pooling
		MaxConnAge:         30 * time.Minute,
	})

	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
