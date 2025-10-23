package dal

import (
	"time"

	"gorm.io/gorm"
)

type GormConfig struct {
	Dialector         gorm.Dialector
	Opts              gorm.Option
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	ConnMaxIdleTime   time.Duration
	PrepareStmt       bool
	DisableForeignKeyConstraintWhenMigrating bool
}

var Gorm *gorm.DB

func initGorm(config *GormConfig) {
	var err error

	// Set default values for performance
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 30 * time.Minute // Reduced from 1 hour
	}
	if config.ConnMaxIdleTime == 0 {
		config.ConnMaxIdleTime = 5 * time.Minute
	}

	// Create optimized GORM config
	gormConfig := &gorm.Config{
		SkipDefaultTransaction:                   true,
		PrepareStmt:                               config.PrepareStmt,
		DisableForeignKeyConstraintWhenMigrating:  config.DisableForeignKeyConstraintWhenMigrating,
	}

	// Merge with provided options
	if config.Opts != nil {
		if opts, ok := config.Opts.(*gorm.Config); ok {
			if opts.Logger != nil {
				gormConfig.Logger = opts.Logger
			}
			if opts.NamingStrategy != nil {
				gormConfig.NamingStrategy = opts.NamingStrategy
			}
		}
	}

	Gorm, err = gorm.Open(config.Dialector, gormConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := Gorm.DB()
	if err != nil {
		panic(err)
	}

	// Optimize connection pool settings
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}
}
