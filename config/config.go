package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Project-related configuration
	Ruoyi struct {
		// Name
		Name string `yaml:"name"`
		// Version
		Version string `yaml:"version"`
		// Copyright year
		Copyright string `yaml:"copyright"`
		// Domain name
		Domain string `yaml:"domain"`
		// Enable SSL
		SSL bool `yaml:"ssl"`
		// File upload path
		UploadPath string `yaml:"uploadPath"`
	} `yaml:"ruoyi"`

	// Development environment configuration
	Server struct {
		// Port
		Port int `yaml:"port"`
		// Mode, optional values: debug, test, release
		Mode string `yaml:"mode"`
	} `yaml:"server"`

	// Database configuration
	Mysql struct {
		Host string `yaml:"host"`
		// Port, default is 3306
		Port int `yaml:"port"`
		// Database name
		Database string `yaml:"database"`
		// Username
		Username string `yaml:"username"`
		// Password
		Password string `yaml:"password"`
		// Charset
		Charset string `yaml:"charset"`
		// Maximum number of idle connections in the connection pool
		MaxIdleConns int `yaml:"maxIdleConns"`
		// Maximum number of open connections in the connection pool
		MaxOpenConns int `yaml:"maxOpenConns"`
	} `yaml:"mysql"`

	// Redis configuration
	Redis struct {
		Host string `yaml:"host"`
		// Port, default is 6379
		Port int `yaml:"port"`
		// Database index
		Database int `yaml:"database"`
		// Password
		Password string `yaml:"password"`
	} `yaml:"redis"`

	// Token configuration
	Token struct {
		// Custom token identifier
		Header string `yaml:"header"`
		// Token secret key
		Secret string `yaml:"secret"`
		// Token validity period (default 30 minutes)
		ExpireTime int `yaml:"expireTime"`
	} `yaml:"token"`

	// User configuration
	User struct {
		Password struct {
			// Maximum password error attempts
			MaxRetryCount int `yaml:"maxRetryCount"`
			// Password lock time (default 10 minutes)
			LockTime int `yaml:"lockTime"`
		} `yaml:"password"`
	} `yaml:"user"`
}

var Data *Config

// GetConfig returns the global configuration object.
func GetConfig() *Config {
	return Data
}

// LoadConfig loads configuration from the given path.
func LoadConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, &Data)
}
