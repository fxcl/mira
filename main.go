package main

import (
	"context"
	"log"
	"mira/anima/dal"
	"mira/app/router"
	"mira/config"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	// Initialize the data access layer with performance optimizations
	dal.InitDal(&dal.Config{
		GormConfig: &dal.GormConfig{
			Dialector:         mysql.Open(dsn),
			Opts: &gorm.Config{
				SkipDefaultTransaction: true, // Skip default transaction
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
				Logger: logger.New(log.Default(), logger.Config{
					LogLevel:                  logger.Error, // Print error logs
					IgnoreRecordNotFoundError: true,
				}),
			},
			MaxOpenConns:      config.Data.Mysql.MaxOpenConns,
			MaxIdleConns:      config.Data.Mysql.MaxIdleConns,
			ConnMaxLifetime:   30 * time.Minute, // Optimized connection lifetime
			ConnMaxIdleTime:   5 * time.Minute,  // Optimized idle time
			PrepareStmt:       true,              // Enable prepared statements for performance
			DisableForeignKeyConstraintWhenMigrating: true, // Disable FK constraints for bulk operations
		},
		RedisConfig: &dal.RedisConfig{
			Host:               config.Data.Redis.Host,
			Port:               config.Data.Redis.Port,
			Database:           config.Data.Redis.Database,
			Password:           config.Data.Redis.Password,
			PoolSize:           50,  // Optimized pool size
			MinIdleConns:       10,  // Minimum idle connections
			MaxRetries:         3,   // Retry attempts
			DialTimeout:        5 * time.Second,
			ReadTimeout:        3 * time.Second,
			WriteTimeout:       3 * time.Second,
			PoolTimeout:        4 * time.Second,
			IdleTimeout:        5 * time.Minute,
			IdleCheckFrequency: time.Minute,
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

	// Create optimized HTTP server with performance settings
	srv := &http.Server{
		Addr:           ":" + strconv.Itoa(config.Data.Server.Port),
		Handler:        server,
		ReadTimeout:    10 * time.Second,  // Max time to read request
		WriteTimeout:   10 * time.Second,  // Max time to write response
		IdleTimeout:    60 * time.Second,  // Max time for keep-alive connections
		ReadHeaderTimeout: 5 * time.Second, // Max time to read headers
		MaxHeaderBytes: 1 << 20,          // 1MB max header size
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %d", config.Data.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
