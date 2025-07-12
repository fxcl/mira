package service

import (
	"testing"
	"time"

	"mira/anima/dal"
	"mira/anima/datetime"
	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
)

func TestOperLogService_CreateSysOperLog(t *testing.T) {
	setup()
	defer teardown()
	s := NewOperLogService()

	t.Run("should create oper log successfully", func(t *testing.T) {
		// Prepare
		log := dto.SaveOperLogRequest{
			Title:    "Test Log",
			OperName: "test_user",
			OperTime: datetime.Datetime{Time: time.Now()},
		}

		// Execute
		err := s.CreateSysOperLogWithErr(log)
		assert.NoError(t, err)

		// Verify
		var result model.SysOperLog
		dal.Gorm.First(&result, "title = ?", "Test Log")
		assert.Equal(t, "Test Log", result.Title)
		assert.Equal(t, "test_user", result.OperName)
	})

	t.Run("should return error when title is empty", func(t *testing.T) {
		// Prepare
		log := dto.SaveOperLogRequest{
			OperName: "test_user",
			OperTime: datetime.Datetime{Time: time.Now()},
		}

		// Execute
		err := s.CreateSysOperLogWithErr(log)
		assert.Error(t, err)
	})
}

func TestOperLogService_DeleteOperLog(t *testing.T) {
	setup()
	defer teardown()
	s := NewOperLogService()

	t.Run("should delete oper logs by ids", func(t *testing.T) {
		// Prepare
		log1 := model.SysOperLog{OperId: 1, Title: "Log 1"}
		log2 := model.SysOperLog{OperId: 2, Title: "Log 2"}
		dal.Gorm.Create(&log1)
		dal.Gorm.Create(&log2)

		// Execute
		err := s.DeleteOperLogWithErr([]int{1})
		assert.NoError(t, err)

		// Verify
		var count int64
		dal.Gorm.Model(&model.SysOperLog{}).Count(&count)
		assert.Equal(t, int64(1), count)

		var result model.SysOperLog
		err = dal.Gorm.First(&result, 1).Error
		assert.Error(t, err, "record not found")
	})

	t.Run("should delete all oper logs", func(t *testing.T) {
		// Prepare
		log1 := model.SysOperLog{OperId: 3, Title: "Log 3"}
		log2 := model.SysOperLog{OperId: 4, Title: "Log 4"}
		dal.Gorm.Create(&log1)
		dal.Gorm.Create(&log2)

		// Execute
		err := s.DeleteOperLogWithErr([]int{})
		assert.NoError(t, err)

		// Verify
		var count int64
		dal.Gorm.Model(&model.SysOperLog{}).Where("oper_id > 0").Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestOperLogService_GetOperLogList(t *testing.T) {
	setup()
	defer teardown()
	s := NewOperLogService()

	t.Run("should get oper log list", func(t *testing.T) {
		// Prepare
		log1 := model.SysOperLog{Title: "Log 1", OperIp: "127.0.0.1", OperName: "user1"}
		log2 := model.SysOperLog{Title: "Log 2", OperIp: "127.0.0.2", OperName: "user2"}
		dal.Gorm.Create(&log1)
		dal.Gorm.Create(&log2)

		// Execute
		params := dto.OperLogListRequest{
			PageRequest: dto.PageRequest{PageNum: 1, PageSize: 10},
		}
		logs, count, err := s.GetOperLogListWithErr(params, true)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Len(t, logs, 2)
	})

	t.Run("should filter by title", func(t *testing.T) {
		// Prepare
		log1 := model.SysOperLog{Title: "Special Log", OperIp: "127.0.0.1", OperName: "user1"}
		log2 := model.SysOperLog{Title: "Another Log", OperIp: "127.0.0.2", OperName: "user2"}
		dal.Gorm.Create(&log1)
		dal.Gorm.Create(&log2)

		// Execute
		params := dto.OperLogListRequest{
			Title:       "Special",
			PageRequest: dto.PageRequest{PageNum: 1, PageSize: 10},
		}
		logs, count, err := s.GetOperLogListWithErr(params, true)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Len(t, logs, 1)
		assert.Equal(t, "Special Log", logs[0].Title)
	})
}
