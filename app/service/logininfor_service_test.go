package service

import (
	"testing"
	"time"

	"mira/anima/dal"
	"mira/anima/datetime"
	rediskey "mira/common/types/redis-key"

	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLogininforService_CreateSysLogininfor(t *testing.T) {
	setup()
	defer teardown()
	s := &LogininforService{}

	t.Run("should create logininfor successfully", func(t *testing.T) {
		param := dto.SaveLogininforRequest{
			UserName:      "testuser",
			Ipaddr:        "127.0.0.1",
			LoginLocation: "Test Location",
			Browser:       "Test Browser",
			Os:            "Test OS",
			Status:        "0",
			Msg:           "Login successful",
			LoginTime:     datetime.Datetime{Time: time.Now()},
		}
		err := s.CreateSysLogininfor(param)
		assert.NoError(t, err)

		// Verify
		// Since the creation is async, we'll wait a bit and then check the db
		time.Sleep(100 * time.Millisecond)
		var createdLogininfor model.SysLogininfor
		dal.Gorm.Last(&createdLogininfor)
		assert.Equal(t, "testuser", createdLogininfor.UserName)
	})
}

func TestLogininforService_DeleteLogininfor(t *testing.T) {
	setup()
	defer teardown()
	s := &LogininforService{}

	t.Run("should delete logininfor successfully", func(t *testing.T) {
		s.CreateSysLogininfor(dto.SaveLogininforRequest{UserName: "testuser"})
		time.Sleep(100 * time.Millisecond)
		var createdLogininfor model.SysLogininfor
		dal.Gorm.Last(&createdLogininfor)

		// Execute
		err := s.DeleteLogininfor([]int{createdLogininfor.InfoId})
		assert.NoError(t, err)

		// Verify
		var deletedLogininfor model.SysLogininfor
		err = dal.Gorm.First(&deletedLogininfor, createdLogininfor.InfoId).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestLogininforService_GetLogininforList(t *testing.T) {
	setup()
	defer teardown()
	s := &LogininforService{}

	t.Run("should return all logininfors", func(t *testing.T) {
		s.CreateSysLogininfor(dto.SaveLogininforRequest{UserName: "user1"})
		s.CreateSysLogininfor(dto.SaveLogininforRequest{UserName: "user2"})
		time.Sleep(100 * time.Millisecond)

		// Execute
		logininfors, count := s.GetLogininforList(dto.LogininforListRequest{
			OrderByColumn: "info_id",
			OrderRule:     "desc",
		}, false)
		assert.Len(t, logininfors, 2)
		assert.Equal(t, 0, count)
	})
}

func TestLogininforService_Unlock(t *testing.T) {
	setup()
	defer teardown()
	s := &LogininforService{}

	t.Run("should unlock user successfully", func(t *testing.T) {
		// Setup
		redisMock.ExpectDel(rediskey.LoginPasswordErrorKey() + "testuser").SetVal(1)

		// Execute
		err := s.Unlock("testuser")
		assert.NoError(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}
