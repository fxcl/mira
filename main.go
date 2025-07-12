package main

import (
	"log"
	"mira/anima/dal"
	"mira/app/router"
	"mira/config"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	// Load configuration
	if err := config.LoadConfig("application.yaml"); err != nil {
		panic(err)
	}

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := config.Data.Mysql.Username + ":" + config.Data.Mysql.Password + "@tcp(" + config.Data.Mysql.Host + ":" + strconv.Itoa(config.Data.Mysql.Port) + ")/" + config.Data.Mysql.Database + "?charset=" + config.Data.Mysql.Charset + "&parseTime=True&loc=Local"

	// Initialize the data access layer
	dal.InitDal(&dal.Config{
		GormConfig: &dal.GormConfig{
			Dialector: mysql.Open(dsn),
			Opts: &gorm.Config{
				SkipDefaultTransaction: true, // Skip default transaction
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
				Logger: logger.New(log.Default(), logger.Config{
					// LogLevel: logger.Silent, // Do not print logs
					LogLevel:                  logger.Error, // Print error logs
					IgnoreRecordNotFoundError: true,
				}),
			},
			MaxOpenConns: config.Data.Mysql.MaxOpenConns,
			MaxIdleConns: config.Data.Mysql.MaxIdleConns,
		},
		RedisConfig: &dal.RedisConfig{
			Host:     config.Data.Redis.Host,
			Port:     config.Data.Redis.Port,
			Database: config.Data.Redis.Database,
			Password: config.Data.Redis.Password,
		},
	})

	// Set mode
	gin.SetMode(config.Data.Server.Mode)

	// Initialize gin
	server := gin.New()

	// Use recovery middleware
	server.Use(gin.Recovery())

	// Set file resource directory
	// If the front end uses the history routing mode, you need to use nginx proxy
	// Comment out server.Static("/admin", "web/admin")
	// If the front and back ends are not deployed separately, you need to configure the front end to hash routing mode
	// Uncomment server.Static("/admin", "web/admin")
	// And create a web/admin directory in the root directory of the project, and copy the files in the dist after the front-end packaging to this directory
	// server.Static("/admin", "web/admin")
	// Set the upload file directory
	server.Static(config.Data.Ruoyi.UploadPath, config.Data.Ruoyi.UploadPath)

	// Register router
	router.Register(server)

	server.Run(":" + strconv.Itoa(config.Data.Server.Port))
}
