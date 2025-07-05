package dal

import (
	"time"

	"gorm.io/gorm"
)

type GormConfig struct {
	Dialector    gorm.Dialector
	Opts         gorm.Option
	MaxOpenConns int
	MaxIdleConns int
}

var Gorm *gorm.DB

func initGorm(config *GormConfig) {
	var err error

	Gorm, err = gorm.Open(config.Dialector, config.Opts)
	if err != nil {
		panic(err)
	}

	sqlDB, err := Gorm.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}
}
