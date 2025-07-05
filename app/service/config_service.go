package service

import (
	"context"
	"encoding/json"
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	rediskey "mira/common/types/redis-key"
)

type ConfigService struct{}

// Create configuration parameter
func (s *ConfigService) CreateConfig(param dto.SaveConfig) error {
	return dal.Gorm.Model(model.SysConfig{}).Create(&model.SysConfig{
		ConfigName:  param.ConfigName,
		ConfigKey:   param.ConfigKey,
		ConfigValue: param.ConfigValue,
		ConfigType:  param.ConfigType,
		CreateBy:    param.CreateBy,
		Remark:      param.Remark,
	}).Error
}

// Update configuration parameter
func (s *ConfigService) UpdateConfig(param dto.SaveConfig) error {
	return dal.Gorm.Model(model.SysConfig{}).Where("config_id = ?", param.ConfigId).Updates(&model.SysConfig{
		ConfigName:  param.ConfigName,
		ConfigKey:   param.ConfigKey,
		ConfigValue: param.ConfigValue,
		ConfigType:  param.ConfigType,
		UpdateBy:    param.UpdateBy,
		Remark:      param.Remark,
	}).Error
}

// Delete configuration parameters
func (s *ConfigService) DeleteConfig(configIds []int) error {
	return dal.Gorm.Model(model.SysConfig{}).Where("config_id IN ?", configIds).Delete(&model.SysConfig{}).Error
}

// Get configuration parameter list
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
		query.Count(&count).Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	query.Find(&configs)

	return configs, int(count)
}

// Get configuration parameter details
func (s *ConfigService) GetConfigByConfigId(configId int) dto.ConfigDetailResponse {
	var config dto.ConfigDetailResponse

	dal.Gorm.Model(model.SysConfig{}).Where("config_id = ?", configId).Last(&config)

	return config
}

// Get configuration parameter value by config key
func (s *ConfigService) GetConfigByConfigKey(configKey string) dto.ConfigDetailResponse {
	var config dto.ConfigDetailResponse

	dal.Gorm.Model(model.SysConfig{}).Where("config_key = ?", configKey).Last(&config)

	return config
}

// Get configuration by config key
func (s *ConfigService) GetConfigCacheByConfigKey(configKey string) dto.ConfigDetailResponse {
	var config dto.ConfigDetailResponse

	// If the cache is not empty, do not read from the database to reduce database pressure
	if configCache, _ := dal.Redis.HGet(context.Background(), rediskey.SysConfigKey, configKey).Result(); configCache != "" {
		if err := json.Unmarshal([]byte(configCache), &config); err == nil {
			return config
		}
	}

	// Read the configuration from the database and record it to the cache
	config = s.GetConfigByConfigKey(configKey)
	if config.ConfigId > 0 {
		configBytes, _ := json.Marshal(&config)
		dal.Redis.HSet(context.Background(), rediskey.SysConfigKey, configKey, string(configBytes)).Result()
	}

	return config
}
