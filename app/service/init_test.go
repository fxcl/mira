package service

import (
	"log"

	"mira/anima/dal"
	"mira/app/model"
	"mira/config"

	"github.com/glebarez/sqlite"
	"github.com/go-redis/redismock/v8"
	"gorm.io/gorm"
)

var redisMock redismock.ClientMock

// setup initializes the test environment, including the database connection
func setup() {
	db, mock := redismock.NewClientMock()
	redisMock = mock
	dal.Redis = db

	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dal.Gorm = gormDB

	// Auto-migrate the schema for the SysConfig model
	dal.Gorm.AutoMigrate(&model.SysConfig{})
	dal.Gorm.AutoMigrate(&model.SysDept{})
	dal.Gorm.AutoMigrate(&model.SysRoleDept{})
	dal.Gorm.AutoMigrate(&model.SysDictType{})
	dal.Gorm.AutoMigrate(&model.SysDictData{})
	dal.Gorm.AutoMigrate(&model.SysLogininfor{})
	dal.Gorm.AutoMigrate(&model.SysMenu{})
	dal.Gorm.AutoMigrate(&model.SysRoleMenu{})
	dal.Gorm.AutoMigrate(&model.SysUserRole{})
	dal.Gorm.AutoMigrate(&model.SysRole{})
	dal.Gorm.AutoMigrate(&model.SysOperLog{})
	dal.Gorm.AutoMigrate(&model.SysPost{})
	dal.Gorm.AutoMigrate(&model.SysUser{})
	dal.Gorm.AutoMigrate(&model.SysUserPost{})

	// Initialize a minimal config for testing
	config.Data = &config.Config{
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

// teardown cleans up the test environment
func teardown() {
	if dal.Gorm != nil {
		dal.Gorm.Exec("DELETE FROM sys_config")
		dal.Gorm.Exec("DELETE FROM sys_dept")
		dal.Gorm.Exec("DELETE FROM sys_role_dept")
		dal.Gorm.Exec("DELETE FROM sys_dict_type")
		dal.Gorm.Exec("DELETE FROM sys_dict_data")
		dal.Gorm.Exec("DELETE FROM sys_logininfor")
		dal.Gorm.Exec("DELETE FROM sys_menu")
		dal.Gorm.Exec("DELETE FROM sys_role_menu")
		dal.Gorm.Exec("DELETE FROM sys_user_role")
		dal.Gorm.Exec("DELETE FROM sys_role")
		dal.Gorm.Exec("DELETE FROM sys_oper_log")
		dal.Gorm.Exec("DELETE FROM sys_post")
		dal.Gorm.Exec("DELETE FROM sys_user")
		dal.Gorm.Exec("DELETE FROM sys_user_post")
		db, _ := dal.Gorm.DB()
		db.Close()
	}
	if dal.Redis != nil {
		dal.Redis.Close()
	}
	log.Println("Test environment cleaned up.")
}
