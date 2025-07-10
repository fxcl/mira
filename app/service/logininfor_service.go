package service

import (
	"context"

	"github.com/pkg/errors"
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	rediskey "mira/common/types/redis-key"
)

// LogininforServiceInterface defines operations for login information management
type LogininforServiceInterface interface {
	DeleteLogininfor(infoIds []int) error
	GetLogininforList(param dto.LogininforListRequest, isPaging bool) ([]dto.LogininforListResponse, int)
	Unlock(userName string) error
	CreateSysLogininfor(param dto.SaveLogininforRequest) error
}

// LogininforService implements the login information management interface
type LogininforService struct{}

// Ensure LogininforService implements LogininforServiceInterface
var _ LogininforServiceInterface = (*LogininforService)(nil)

// DeleteLogininfor deletes login logs by IDs or all logs if no IDs are provided
func (s *LogininforService) DeleteLogininfor(infoIds []int) error {
	return s.DeleteLogininforWithErr(infoIds)
}

// DeleteLogininforWithErr deletes login logs by IDs or all logs if no IDs are provided with error handling
func (s *LogininforService) DeleteLogininforWithErr(infoIds []int) error {
	var err error

	if len(infoIds) > 0 {
		err = dal.Gorm.Model(model.SysLogininfor{}).Where("info_id IN ?", infoIds).Delete(&model.SysLogininfor{}).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete login logs by IDs")
		}
		return nil
	}

	// To solve the "WHERE conditions required" error, add the condition Where("info_id > ?", 0)
	err = dal.Gorm.Model(model.SysLogininfor{}).Where("info_id > ?", 0).Delete(&model.SysLogininfor{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete all login logs")
	}

	return nil
}

// GetLogininforList retrieves a list of login information records based on search parameters
func (s *LogininforService) GetLogininforList(param dto.LogininforListRequest, isPaging bool) ([]dto.LogininforListResponse, int) {
	logininfos, count, _ := s.GetLogininforListWithErr(param, isPaging)
	return logininfos, count
}

// GetLogininforListWithErr retrieves a list of login information records with error handling
func (s *LogininforService) GetLogininforListWithErr(param dto.LogininforListRequest, isPaging bool) ([]dto.LogininforListResponse, int, error) {
	var count int64
	logininfos := make([]dto.LogininforListResponse, 0)

	query := dal.Gorm.Model(model.SysLogininfor{}).Order(param.OrderByColumn + " " + param.OrderRule)

	if param.Ipaddr != "" {
		query = query.Where("ipaddr LIKE ?", "%"+param.Ipaddr+"%")
	}

	if param.UserName != "" {
		query = query.Where("user_name LIKE ?", "%"+param.UserName+"%")
	}

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if param.BeginTime != "" && param.EndTime != "" {
		query = query.Where("login_time BETWEEN ? AND ?", param.BeginTime, param.EndTime)
	}

	if isPaging {
		err := query.Count(&count).Error
		if err != nil {
			return logininfos, 0, errors.Wrap(err, "failed to count login records")
		}
		query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	err := query.Find(&logininfos).Error
	if err != nil {
		return logininfos, 0, errors.Wrap(err, "failed to retrieve login records")
	}

	return logininfos, int(count), nil
}

// Unlock removes the login error count cache for a user, effectively unlocking their account
func (s *LogininforService) Unlock(userName string) error {
	return s.UnlockWithErr(userName)
}

// UnlockWithErr removes the login error count cache for a user with proper error handling
func (s *LogininforService) UnlockWithErr(userName string) error {
	if userName == "" {
		return errors.New("username cannot be empty")
	}

	_, err := dal.Redis.Del(context.Background(), rediskey.LoginPasswordErrorKey+userName).Result()
	if err != nil {
		return errors.Wrap(err, "failed to delete login error cache for user")
	}

	return nil
}

// CreateSysLogininfor records login information asynchronously
// Note: This method starts a goroutine and returns immediately without waiting for completion
func (s *LogininforService) CreateSysLogininfor(param dto.SaveLogininforRequest) error {
	// For backward compatibility, we keep the asynchronous behavior
	go func() {
		_ = s.CreateSysLogininforWithErr(param) // Errors are ignored in the asynchronous version
	}()
	return nil
}

// CreateSysLogininforWithErr records login information synchronously with proper error handling
func (s *LogininforService) CreateSysLogininforWithErr(param dto.SaveLogininforRequest) error {
	// Input validation
	if param.UserName == "" {
		return errors.New("username cannot be empty")
	}

	// Create the login record
	err := dal.Gorm.Model(model.SysLogininfor{}).Create(&model.SysLogininfor{
		UserName:      param.UserName,
		Ipaddr:        param.Ipaddr,
		LoginLocation: param.LoginLocation,
		Browser:       param.Browser,
		Os:            param.Os,
		Status:        param.Status,
		Msg:           param.Msg,
		LoginTime:     param.LoginTime,
	}).Error

	if err != nil {
		return errors.Wrap(err, "failed to create login record")
	}

	return nil
}
