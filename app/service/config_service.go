package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	rediskey "mira/common/types/redis-key"
	"mira/common/xerrors"
)

// ConfigServiceInterface defines the interface for configuration service, facilitating testing and dependency injection
type ConfigServiceInterface interface {
	CreateConfig(param dto.SaveConfig) error
	UpdateConfig(param dto.SaveConfig) error
	DeleteConfig(configIds []int) error
	GetConfigList(param dto.ConfigListRequest, isPaging bool) ([]dto.ConfigListResponse, int)
	GetConfigByConfigId(configId int) dto.ConfigDetailResponse
	GetConfigByConfigKey(configKey string) dto.ConfigDetailResponse
	GetConfigCacheByConfigKey(configKey string) dto.ConfigDetailResponse
	RefreshCache() error
}

// ConfigService implements the configuration service interface
type ConfigService struct{}

// Ensure ConfigService implements ConfigServiceInterface
var _ ConfigServiceInterface = (*ConfigService)(nil)

// CreateConfig creates a new system configuration parameter
//
// Parameters:
//   - param: Configuration data transfer object containing all required fields
//
// Returns:
//   - error: Any error that occurred during creation, or nil on success
func (s *ConfigService) CreateConfig(param dto.SaveConfig) error {
	if param.ConfigName == "" {
		return xerrors.ErrConfigNameEmpty
	}
	if param.ConfigKey == "" {
		return xerrors.ErrConfigKeyEmpty
	}
	if param.ConfigValue == "" {
		return xerrors.ErrConfigValueEmpty
	}

	err := dal.Gorm.Model(model.SysConfig{}).Create(&model.SysConfig{
		ConfigName:  param.ConfigName,
		ConfigKey:   param.ConfigKey,
		ConfigValue: param.ConfigValue,
		ConfigType:  param.ConfigType,
		CreateBy:    param.CreateBy,
		Remark:      param.Remark,
	}).Error
	if err != nil {
		log.Printf("Failed to create config: %v", err)
		return fmt.Errorf("failed to create config: %w", err)
	}

	// Refresh cache after creating new configuration
	if err := s.RefreshCache(); err != nil {
		log.Printf("Warning: Failed to refresh config cache after creation: %v", err)
	}

	return nil
}

// UpdateConfig updates system configuration parameter
//
// Parameters:
//   - param: Configuration data transfer object containing fields to update
//
// Returns:
//   - error: Any error that occurred during update, or nil on success
func (s *ConfigService) UpdateConfig(param dto.SaveConfig) error {
	if param.ConfigId <= 0 {
		return xerrors.ErrParam
	}

	err := dal.Gorm.Model(model.SysConfig{}).Where("config_id = ?", param.ConfigId).Updates(&model.SysConfig{
		ConfigName:  param.ConfigName,
		ConfigKey:   param.ConfigKey,
		ConfigValue: param.ConfigValue,
		ConfigType:  param.ConfigType,
		UpdateBy:    param.UpdateBy,
		Remark:      param.Remark,
	}).Error
	if err != nil {
		log.Printf("Failed to update config: %v", err)
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Refresh cache after updating configuration
	if err := s.RefreshCache(); err != nil {
		log.Printf("Warning: Failed to refresh config cache after update: %v", err)
	}

	return nil
}

// DeleteConfig deletes system configuration parameters
//
// Parameters:
//   - configIds: Array of configuration IDs to delete
//
// Returns:
//   - error: Any error that occurred during deletion, or nil on success
func (s *ConfigService) DeleteConfig(configIds []int) error {
	if len(configIds) == 0 {
		return xerrors.ErrParam
	}

	err := dal.Gorm.Model(model.SysConfig{}).Where("config_id IN ?", configIds).Delete(&model.SysConfig{}).Error
	if err != nil {
		log.Printf("Failed to delete configs: %v", err)
		return fmt.Errorf("failed to delete configs: %w", err)
	}

	// Refresh cache after deleting configurations
	if err := s.RefreshCache(); err != nil {
		log.Printf("Warning: Failed to refresh config cache after deletion: %v", err)
	}

	return nil
}

// GetConfigList gets the list of configuration parameters
//
// Parameters:
//   - param: Request object containing query conditions
//   - isPaging: Whether pagination is needed
//
// Returns:
//   - []dto.ConfigListResponse: List of configurations
//   - int: Total record count if isPaging is true; otherwise 0
func (s *ConfigService) GetConfigList(param dto.ConfigListRequest, isPaging bool) ([]dto.ConfigListResponse, int) {
	var count int64
	configs := make([]dto.ConfigListResponse, 0)

	query := dal.Gorm.Model(model.SysConfig{}).Order("config_id")

	if param.ConfigName != "" {
		query = query.Where("config_name LIKE ?", "%"+param.ConfigName+"%")
	}

	if param.ConfigKey != "" {
		query = query.Where("config_key LIKE ?", "%"+param.ConfigKey+"%")
	}

	if param.ConfigType != "" {
		query = query.Where("config_type = ?", param.ConfigType)
	}

	if param.BeginTime != "" && param.EndTime != "" {
		query = query.Where("create_time BETWEEN ? AND ?", param.BeginTime, param.EndTime)
	}

	if isPaging {
		if err := query.Count(&count).Error; err != nil {
			log.Printf("Failed to count configs: %v", err)
		}
		query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	if err := query.Find(&configs).Error; err != nil {
		log.Printf("Failed to query configs: %v", err)
	}

	return configs, int(count)
}

// GetConfigByConfigId gets configuration parameter details by config ID
//
// Parameters:
//   - configId: Configuration ID
//
// Returns:
//   - dto.ConfigDetailResponse: Configuration details, or empty object if not found
func (s *ConfigService) GetConfigByConfigId(configId int) dto.ConfigDetailResponse {
	var config dto.ConfigDetailResponse

	if configId <= 0 {
		return config
	}

	if err := dal.Gorm.Model(model.SysConfig{}).Where("config_id = ?", configId).Last(&config).Error; err != nil {
		log.Printf("Failed to get config by ID %d: %v", configId, err)
	}

	return config
}

// GetConfigByConfigKey gets configuration parameter details by config key
//
// Parameters:
//   - configKey: Configuration key
//
// Returns:
//   - dto.ConfigDetailResponse: Configuration details, or empty object if not found
func (s *ConfigService) GetConfigByConfigKey(configKey string) dto.ConfigDetailResponse {
	var config dto.ConfigDetailResponse

	if configKey == "" {
		return config
	}

	if err := dal.Gorm.Model(model.SysConfig{}).Where("config_key = ?", configKey).Last(&config).Error; err != nil {
		log.Printf("Failed to get config by key %s: %v", configKey, err)
	}

	return config
}

// GetConfigCacheByConfigKey gets configuration parameter by config key (prioritizing cache)
//
// Parameters:
//   - configKey: Configuration key
//
// Returns:
//   - dto.ConfigDetailResponse: Configuration details, or empty object if not found
//
// Caching Strategy:
//   - Prioritizes fetching from Redis cache
//   - Falls back to database query on cache miss or error
//   - Writes query results to cache with 24-hour expiration time
//   - Logs errors but does not interrupt flow (graceful degradation)
func (s *ConfigService) GetConfigCacheByConfigKey(configKey string) dto.ConfigDetailResponse {
	var config dto.ConfigDetailResponse
	ctx := context.Background()

	// If cache is not empty, avoid reading from database to reduce database pressure
	configCache, err := dal.Redis.HGet(ctx, rediskey.SysConfigKey(), configKey).Result()
	if err != nil {
		// Log error but continue execution (fallback to database)
		log.Printf("Redis error when getting config for key %s: %v", configKey, err)
	} else if configCache != "" {
		if err := json.Unmarshal([]byte(configCache), &config); err != nil {
			// Log deserialization error
			log.Printf("Failed to unmarshal config for key %s: %v", configKey, err)
		} else {
			return config
		}
	}

	// Read configuration from database and record to cache
	config = s.GetConfigByConfigKey(configKey)
	if config.ConfigId > 0 {
		configBytes, err := json.Marshal(&config)
		if err != nil {
			log.Printf("Failed to marshal config for key %s: %v", configKey, err)
			return config
		}

		// Set cache
		_, err = dal.Redis.HSet(ctx, rediskey.SysConfigKey(), configKey, string(configBytes)).Result()
		if err != nil {
			// Log error but don't affect return value
			log.Printf("Failed to set config cache for key %s: %v", configKey, err)
		} else {
			// Set cache expiration time (if not already set)
			dal.Redis.Expire(ctx, rediskey.SysConfigKey(), 24*time.Hour)
		}
	}

	return config
}

// RefreshCache refreshes the configuration cache
//
// Returns:
//   - error: Any error that occurred during refresh, or nil on success
func (s *ConfigService) RefreshCache() error {
	ctx := context.Background()
	err := dal.Redis.Del(ctx, rediskey.SysConfigKey()).Err()
	if err != nil {
		log.Printf("Failed to refresh config cache: %v", err)
		return fmt.Errorf("failed to refresh config cache: %w", err)
	}
	return nil
}
