package service

import (
	"github.com/pkg/errors"
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
)

// OperLogServiceInterface defines operations for operation log management
type OperLogServiceInterface interface {
	DeleteOperLog(operIds []int) error
	GetOperLogList(param dto.OperLogListRequest, isPaging bool) ([]dto.OperLogListResponse, int)
	CreateSysOperLog(param dto.SaveOperLogRequest) error
}

// OperLogService implements the operation log management interface
type OperLogService struct{}

// Ensure OperLogService implements OperLogServiceInterface
var _ OperLogServiceInterface = (*OperLogService)(nil)

// DeleteOperLog deletes operation logs by IDs or all logs if no IDs are provided
func (s *OperLogService) DeleteOperLog(operIds []int) error {
	return s.DeleteOperLogWithErr(operIds)
}

// DeleteOperLogWithErr deletes operation logs by IDs or all logs if no IDs are provided with error handling
func (s *OperLogService) DeleteOperLogWithErr(operIds []int) error {
	var err error

	if len(operIds) > 0 {
		err = dal.Gorm.Model(model.SysOperLog{}).Where("oper_id IN ?", operIds).Delete(&model.SysOperLog{}).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete operation logs by IDs")
		}
		return nil
	}

	// To solve the "WHERE conditions required" error, add the condition Where("oper_id > ?", 0)
	err = dal.Gorm.Model(model.SysOperLog{}).Where("oper_id > ?", 0).Delete(&model.SysOperLog{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete all operation logs")
	}

	return nil
}

// GetOperLogList retrieves a list of operation logs based on search parameters
func (s *OperLogService) GetOperLogList(param dto.OperLogListRequest, isPaging bool) ([]dto.OperLogListResponse, int) {
	operLogs, count, _ := s.GetOperLogListWithErr(param, isPaging)
	return operLogs, count
}

// GetOperLogListWithErr retrieves a list of operation logs with proper error handling
func (s *OperLogService) GetOperLogListWithErr(param dto.OperLogListRequest, isPaging bool) ([]dto.OperLogListResponse, int, error) {
	var count int64
	operLogs := make([]dto.OperLogListResponse, 0)

	query := dal.Gorm.Model(model.SysOperLog{}).Order(param.OrderByColumn + " " + param.OrderRule)

	if param.OperIp != "" {
		query = query.Where("oper_ip LIKE ?", "%"+param.OperIp+"%")
	}

	if param.Title != "" {
		query = query.Where("title LIKE ?", "%"+param.Title+"%")
	}

	if param.OperName != "" {
		query = query.Where("oper_name LIKE ?", "%"+param.OperName+"%")
	}

	if param.BusinessType != "" {
		query = query.Where("business_type = ?", param.BusinessType)
	}

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if param.BeginTime != "" && param.EndTime != "" {
		query = query.Where("oper_time BETWEEN ? AND ?", param.BeginTime, param.EndTime)
	}

	if isPaging {
		err := query.Count(&count).Error
		if err != nil {
			return nil, 0, errors.Wrap(err, "failed to count operation logs")
		}
		query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	err := query.Find(&operLogs).Error
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to retrieve operation logs")
	}

	return operLogs, int(count), nil
}

// CreateSysOperLog records operation log information asynchronously
// Note: This method starts a goroutine and returns immediately without waiting for completion
func (s *OperLogService) CreateSysOperLog(param dto.SaveOperLogRequest) error {
	// For backward compatibility, we keep the asynchronous behavior
	go func() {
		_ = s.CreateSysOperLogWithErr(param) // Errors are ignored in the asynchronous version
	}()
	return nil
}

// CreateSysOperLogWithErr records operation log information synchronously with proper error handling
func (s *OperLogService) CreateSysOperLogWithErr(param dto.SaveOperLogRequest) error {
	// Input validation
	if param.Title == "" {
		return errors.New("operation title cannot be empty")
	}

	// Create the operation log record
	err := dal.Gorm.Model(model.SysOperLog{}).Create(&model.SysOperLog{
		Title:         param.Title,
		BusinessType:  param.BusinessType,
		Method:        param.Method,
		RequestMethod: param.RequestMethod,
		OperName:      param.OperName,
		DeptName:      param.DeptName,
		OperUrl:       param.OperUrl,
		OperIp:        param.OperIp,
		OperLocation:  param.OperLocation,
		OperParam:     param.OperParam,
		JsonResult:    param.JsonResult,
		Status:        param.Status,
		ErrorMsg:      param.ErrorMsg,
		OperTime:      param.OperTime,
		CostTime:      param.CostTime,
	}).Error
	if err != nil {
		return errors.Wrap(err, "failed to create operation log record")
	}

	return nil
}
