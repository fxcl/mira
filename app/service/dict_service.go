package service

import (
	"context"
	"encoding/json"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"

	"github.com/pkg/errors"

	rediskey "mira/common/types/redis-key"
)

// DictTypeServiceInterface defines operations for dictionary type management
type DictTypeServiceInterface interface {
	CreateDictType(param dto.SaveDictType) error
	UpdateDictType(param dto.SaveDictType) error
	DeleteDictType(dictIds []int) error
	GetDictTypeList(param dto.DictTypeListRequest, isPaging bool) ([]dto.DictTypeListResponse, int)
	GetDictTypeByDictId(dictId int) dto.DictTypeDetailResponse
	GetDcitTypeByDictType(dictType string) dto.DictTypeDetailResponse
	RefreshCache() error
}

// DictTypeService implements the dictionary type management interface
type DictTypeService struct{}

// Ensure DictTypeService implements DictTypeServiceInterface
var _ DictTypeServiceInterface = (*DictTypeService)(nil)

// CreateDictType creates a new dictionary type
//
// Parameters:
//   - param: Dictionary type data transfer object containing all required fields
//
// Returns:
//   - error: Any error that occurred during creation, or nil on success
func (s *DictTypeService) CreateDictType(param dto.SaveDictType) error {
	err := dal.Gorm.Model(model.SysDictType{}).Create(&model.SysDictType{
		DictName: param.DictName,
		DictType: param.DictType,
		Status:   param.Status,
		Remark:   param.Remark,
		CreateBy: param.CreateBy,
	}).Error
	if err != nil {
		return errors.Wrap(err, "failed to create dictionary type")
	}

	return nil
}

// UpdateDictType updates an existing dictionary type
//
// Parameters:
//   - param: Dictionary type data transfer object containing fields to update
//
// Returns:
//   - error: Any error that occurred during update, or nil on success
func (s *DictTypeService) UpdateDictType(param dto.SaveDictType) error {
	err := dal.Gorm.Model(model.SysDictType{}).Where("dict_id = ?", param.DictId).Updates(&model.SysDictType{
		DictName: param.DictName,
		DictType: param.DictType,
		Status:   param.Status,
		Remark:   param.Remark,
		UpdateBy: param.UpdateBy,
	}).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update dictionary type with ID %d", param.DictId)
	}

	return nil
}

// DeleteDictType deletes dictionary types by their IDs
//
// Parameters:
//   - dictIds: Array of dictionary type IDs to delete
//
// Returns:
//   - error: Any error that occurred during deletion, or nil on success
func (s *DictTypeService) DeleteDictType(dictIds []int) error {
	if len(dictIds) == 0 {
		return errors.New("no dictionary type IDs provided for deletion")
	}

	err := dal.Gorm.Model(model.SysDictType{}).Where("dict_id IN ?", dictIds).Delete(&model.SysDictType{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete dictionary types")
	}

	return nil
}

// GetDictTypeList gets the list of dictionary types based on query parameters
//
// Parameters:
//   - param: Request object containing query conditions
//   - isPaging: Whether pagination is needed
//
// Returns:
//   - []dto.DictTypeListResponse: List of dictionary types
//   - int: Total record count if isPaging is true; otherwise 0
func (s *DictTypeService) GetDictTypeList(param dto.DictTypeListRequest, isPaging bool) ([]dto.DictTypeListResponse, int) {
	dictTypes, count, err := s.GetDictTypeListWithErr(param, isPaging)
	if err != nil {
		// Error already handled in the inner method
	}
	return dictTypes, count
}

// GetDictTypeListWithErr gets the list of dictionary types based on query parameters with error reporting
//
// Parameters:
//   - param: Request object containing query conditions
//   - isPaging: Whether pagination is needed
//
// Returns:
//   - []dto.DictTypeListResponse: List of dictionary types
//   - int: Total record count if isPaging is true; otherwise 0
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DictTypeService) GetDictTypeListWithErr(param dto.DictTypeListRequest, isPaging bool) ([]dto.DictTypeListResponse, int, error) {
	var count int64
	dictTypes := make([]dto.DictTypeListResponse, 0)

	query := dal.Gorm.Model(model.SysDictType{}).Order("dict_id")

	if param.DictName != "" {
		query = query.Where("dict_name LIKE ?", "%"+param.DictName+"%")
	}

	if param.DictType != "" {
		query = query.Where("dict_type LIKE ?", "%"+param.DictType+"%")
	}

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if param.BeginTime != "" && param.EndTime != "" {
		query = query.Where("create_time BETWEEN ? AND ?", param.BeginTime, param.EndTime)
	}

	if isPaging {
		if err := query.Count(&count).Error; err != nil {
			return nil, 0, errors.Wrap(err, "failed to count dictionary types")
		}
		query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	if err := query.Find(&dictTypes).Error; err != nil {
		return nil, 0, errors.Wrap(err, "failed to query dictionary types")
	}

	return dictTypes, int(count), nil
}

// GetDictTypeByDictId gets dictionary type details by ID
//
// Parameters:
//   - dictId: Dictionary type ID to look up
//
// Returns:
//   - dto.DictTypeDetailResponse: Dictionary type details, or empty object if not found
func (s *DictTypeService) GetDictTypeByDictId(dictId int) dto.DictTypeDetailResponse {
	dictType, err := s.GetDictTypeByDictIdWithErr(dictId)
	if err != nil {
		// Error already handled in the inner method
	}
	return dictType
}

// GetDictTypeByDictIdWithErr gets dictionary type details by ID with error reporting
//
// Parameters:
//   - dictId: Dictionary type ID to look up
//
// Returns:
//   - dto.DictTypeDetailResponse: Dictionary type details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DictTypeService) GetDictTypeByDictIdWithErr(dictId int) (dto.DictTypeDetailResponse, error) {
	var dictType dto.DictTypeDetailResponse

	if dictId <= 0 {
		return dictType, errors.Errorf("invalid dictionary type ID: %d", dictId)
	}

	if err := dal.Gorm.Model(model.SysDictType{}).Where("dict_id = ?", dictId).Last(&dictType).Error; err != nil {
		return dictType, errors.Wrapf(err, "failed to get dictionary type by ID %d", dictId)
	}

	return dictType, nil
}

// GetDcitTypeByDictType gets dictionary type details by type code
//
// Parameters:
//   - dictType: Dictionary type code to look up
//
// Returns:
//   - dto.DictTypeDetailResponse: Dictionary type details, or empty object if not found
func (s *DictTypeService) GetDcitTypeByDictType(dictType string) dto.DictTypeDetailResponse {
	dictTypeResult, err := s.GetDcitTypeByDictTypeWithErr(dictType)
	if err != nil {
		// Error already handled in the inner method
	}
	return dictTypeResult
}

// GetDcitTypeByDictTypeWithErr gets dictionary type details by type code with error reporting
//
// Parameters:
//   - dictType: Dictionary type code to look up
//
// Returns:
//   - dto.DictTypeDetailResponse: Dictionary type details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DictTypeService) GetDcitTypeByDictTypeWithErr(dictType string) (dto.DictTypeDetailResponse, error) {
	var dictTypeResult dto.DictTypeDetailResponse

	if dictType == "" {
		return dictTypeResult, errors.New("empty dictionary type provided")
	}

	if err := dal.Gorm.Model(model.SysDictType{}).Where("dict_type = ?", dictType).Last(&dictTypeResult).Error; err != nil {
		return dictTypeResult, errors.Wrapf(err, "failed to get dictionary type by type %s", dictType)
	}

	return dictTypeResult, nil
}

// RefreshCache refreshes the dictionary cache
//
// Returns:
//   - error: Any error that occurred during refresh, or nil on success
func (s *DictTypeService) RefreshCache() error {
	err := dal.Redis.Del(context.Background(), rediskey.SysDictKey()).Err()
	if err != nil {
		return errors.Wrap(err, "failed to refresh dictionary cache")
	}
	return nil
}

// DictDataServiceInterface defines operations for dictionary data management
type DictDataServiceInterface interface {
	CreateDictData(param dto.SaveDictData) error
	UpdateDictData(param dto.SaveDictData) error
	DeleteDictData(dictCodes []int) error
	GetDictDataList(param dto.DictDataListRequest, isPaging bool) ([]dto.DictDataListResponse, int)
	GetDictDataByDictCode(dictCode int) dto.DictDataDetailResponse
	GetDictDataByDictType(dictType string) []dto.DictDataListResponse
	GetDictDataCacheByDictType(dictType string) []dto.DictDataListResponse
}

// DictDataService implements the dictionary data management interface
type DictDataService struct{}

// Ensure DictDataService implements DictDataServiceInterface
var _ DictDataServiceInterface = (*DictDataService)(nil)

// CreateDictData creates a new dictionary data entry
//
// Parameters:
//   - param: Dictionary data transfer object containing all required fields
//
// Returns:
//   - error: Any error that occurred during creation, or nil on success
func (s *DictDataService) CreateDictData(param dto.SaveDictData) error {
	err := dal.Gorm.Model(model.SysDictData{}).Create(&model.SysDictData{
		DictSort:  param.DictSort,
		DictLabel: param.DictLabel,
		DictValue: param.DictValue,
		DictType:  param.DictType,
		CssClass:  param.CssClass,
		ListClass: param.ListClass,
		IsDefault: param.IsDefault,
		Status:    param.Status,
		Remark:    param.Remark,
		CreateBy:  param.CreateBy,
	}).Error
	if err != nil {
		return errors.Wrap(err, "failed to create dictionary data")
	}

	return nil
}

// UpdateDictData updates an existing dictionary data entry
//
// Parameters:
//   - param: Dictionary data transfer object containing fields to update
//
// Returns:
//   - error: Any error that occurred during update, or nil on success
func (s *DictDataService) UpdateDictData(param dto.SaveDictData) error {
	err := dal.Gorm.Model(model.SysDictData{}).Where("dict_code = ?", param.DictCode).Updates(&model.SysDictData{
		DictSort:  param.DictSort,
		DictLabel: param.DictLabel,
		DictValue: param.DictValue,
		DictType:  param.DictType,
		CssClass:  param.CssClass,
		ListClass: param.ListClass,
		IsDefault: param.IsDefault,
		Status:    param.Status,
		Remark:    param.Remark,
		UpdateBy:  param.UpdateBy,
	}).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update dictionary data with code %d", param.DictCode)
	}

	return nil
}

// DeleteDictData deletes dictionary data entries by their codes
//
// Parameters:
//   - dictCodes: Array of dictionary data codes to delete
//
// Returns:
//   - error: Any error that occurred during deletion, or nil on success
func (s *DictDataService) DeleteDictData(dictCodes []int) error {
	if len(dictCodes) == 0 {
		return errors.New("no dictionary data codes provided for deletion")
	}

	err := dal.Gorm.Model(model.SysDictData{}).Where("dict_code IN ?", dictCodes).Delete(&model.SysDictData{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete dictionary data")
	}

	return nil
}

// GetDictDataList gets the list of dictionary data based on query parameters
//
// Parameters:
//   - param: Request object containing query conditions
//   - isPaging: Whether pagination is needed
//
// Returns:
//   - []dto.DictDataListResponse: List of dictionary data
//   - int: Total record count if isPaging is true; otherwise 0
func (s *DictDataService) GetDictDataList(param dto.DictDataListRequest, isPaging bool) ([]dto.DictDataListResponse, int) {
	dictDatas, count, err := s.GetDictDataListWithErr(param, isPaging)
	if err != nil {
		// Error already handled in the inner method
	}
	return dictDatas, count
}

// GetDictDataListWithErr gets the list of dictionary data based on query parameters with error reporting
//
// Parameters:
//   - param: Request object containing query conditions
//   - isPaging: Whether pagination is needed
//
// Returns:
//   - []dto.DictDataListResponse: List of dictionary data
//   - int: Total record count if isPaging is true; otherwise 0
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DictDataService) GetDictDataListWithErr(param dto.DictDataListRequest, isPaging bool) ([]dto.DictDataListResponse, int, error) {
	var count int64
	dictDatas := make([]dto.DictDataListResponse, 0)

	query := dal.Gorm.Model(model.SysDictData{}).Order("dict_code")

	if param.DictLabel != "" {
		query = query.Where("dict_label LIKE ?", "%"+param.DictLabel+"%")
	}

	if param.DictType != "" {
		query = query.Where("dict_type LIKE ?", "%"+param.DictType+"%")
	}

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if isPaging {
		if err := query.Count(&count).Error; err != nil {
			return nil, 0, errors.Wrap(err, "failed to count dictionary data")
		}
		query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	if err := query.Find(&dictDatas).Error; err != nil {
		return nil, 0, errors.Wrap(err, "failed to query dictionary data")
	}

	return dictDatas, int(count), nil
}

// GetDictDataByDictCode gets dictionary data details by code
//
// Parameters:
//   - dictCode: Dictionary data code to look up
//
// Returns:
//   - dto.DictDataDetailResponse: Dictionary data details, or empty object if not found
func (s *DictDataService) GetDictDataByDictCode(dictCode int) dto.DictDataDetailResponse {
	dictData, err := s.GetDictDataByDictCodeWithErr(dictCode)
	if err != nil {
		// Error already handled in the inner method
	}
	return dictData
}

// GetDictDataByDictCodeWithErr gets dictionary data details by code with error reporting
//
// Parameters:
//   - dictCode: Dictionary data code to look up
//
// Returns:
//   - dto.DictDataDetailResponse: Dictionary data details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DictDataService) GetDictDataByDictCodeWithErr(dictCode int) (dto.DictDataDetailResponse, error) {
	var dictData dto.DictDataDetailResponse

	if dictCode <= 0 {
		return dictData, errors.Errorf("invalid dictionary data code: %d", dictCode)
	}

	if err := dal.Gorm.Model(model.SysDictData{}).Where("dict_code = ?", dictCode).Last(&dictData).Error; err != nil {
		return dictData, errors.Wrapf(err, "failed to get dictionary data by code %d", dictCode)
	}

	return dictData, nil
}

// GetDictDataByDictType gets dictionary data by dictionary type
//
// Parameters:
//   - dictType: Dictionary type to look up
//
// Returns:
//   - []dto.DictDataListResponse: List of dictionary data
func (s *DictDataService) GetDictDataByDictType(dictType string) []dto.DictDataListResponse {
	dictDatas, err := s.GetDictDataByDictTypeWithErr(dictType)
	if err != nil {
		// Error already handled in the inner method
	}
	return dictDatas
}

// GetDictDataByDictTypeWithErr gets dictionary data by dictionary type with error reporting
//
// Parameters:
//   - dictType: Dictionary type to look up
//
// Returns:
//   - []dto.DictDataListResponse: List of dictionary data
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DictDataService) GetDictDataByDictTypeWithErr(dictType string) ([]dto.DictDataListResponse, error) {
	dictDatas := make([]dto.DictDataListResponse, 0)

	if dictType == "" {
		return dictDatas, errors.New("empty dictionary type provided")
	}

	if err := dal.Gorm.Model(model.SysDictData{}).Where("status = ? AND dict_type = ?", constant.NORMAL_STATUS, dictType).Find(&dictDatas).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to get dictionary data by type %s", dictType)
	}

	return dictDatas, nil
}

// GetDictDataCacheByDictType gets dictionary data by dictionary type with caching
//
// Parameters:
//   - dictType: Dictionary type to look up
//
// Returns:
//   - []dto.DictDataListResponse: List of dictionary data
func (s *DictDataService) GetDictDataCacheByDictType(dictType string) []dto.DictDataListResponse {
	dictDatas, err := s.GetDictDataCacheByDictTypeWithErr(dictType)
	if err != nil {
		// Error already handled in the inner method
	}
	return dictDatas
}

// GetDictDataCacheByDictTypeWithErr gets dictionary data by dictionary type with caching and error reporting
//
// Parameters:
//   - dictType: Dictionary type to look up
//
// Returns:
//   - []dto.DictDataListResponse: List of dictionary data
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DictDataService) GetDictDataCacheByDictTypeWithErr(dictType string) ([]dto.DictDataListResponse, error) {
	dictDatas := make([]dto.DictDataListResponse, 0)

	if dictType == "" {
		return dictDatas, errors.New("empty dictionary type provided")
	}

	// Try to get from cache first
	cache, err := dal.Redis.HGet(context.Background(), rediskey.SysDictKey(), dictType).Result()
	if err == nil && cache != "" {
		err = json.Unmarshal([]byte(cache), &dictDatas)
		if err == nil {
			return dictDatas, nil
		}
	}

	// Get from DB if cache fails
	dictDatas, err = s.GetDictDataByDictTypeWithErr(dictType)
	if err != nil {
		return nil, err
	}

	// Set cache
	cacheBytes, err := json.Marshal(dictDatas)
	if err == nil {
		dal.Redis.HSet(context.Background(), rediskey.SysDictKey(), dictType, string(cacheBytes))
	}

	return dictDatas, nil
}
