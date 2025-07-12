package service

import (
	"encoding/json"
	"testing"
	"time"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	rediskey "mira/common/types/redis-key"
	"mira/common/xerrors"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestConfigService_CreateConfig(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should return error when config name is empty", func(t *testing.T) {
		param := dto.SaveConfig{
			ConfigName:  "",
			ConfigKey:   "some-key",
			ConfigValue: "some-value",
			ConfigType:  "Y",
			CreateBy:    "test",
		}
		err := s.CreateConfig(param)
		assert.Error(t, err)
		assert.Equal(t, xerrors.ErrConfigNameEmpty, err)
	})

	t.Run("should return error when config key is empty", func(t *testing.T) {
		param := dto.SaveConfig{
			ConfigName:  "some-name",
			ConfigKey:   "",
			ConfigValue: "some-value",
			ConfigType:  "Y",
			CreateBy:    "test",
		}
		err := s.CreateConfig(param)
		assert.Error(t, err)
		assert.Equal(t, xerrors.ErrConfigKeyEmpty, err)
	})

	t.Run("should return error when config value is empty", func(t *testing.T) {
		param := dto.SaveConfig{
			ConfigName:  "some-name",
			ConfigKey:   "some-key",
			ConfigValue: "",
			ConfigType:  "Y",
			CreateBy:    "test",
		}
		err := s.CreateConfig(param)
		assert.Error(t, err)
		assert.Equal(t, xerrors.ErrConfigValueEmpty, err)
	})

	t.Run("should create config successfully", func(t *testing.T) {
		param := dto.SaveConfig{
			ConfigName:  "test-config",
			ConfigKey:   "test-key",
			ConfigValue: "test-value",
			ConfigType:  "Y",
			CreateBy:    "tester",
		}

		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)

		err := s.CreateConfig(param)
		assert.NoError(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestConfigService_UpdateConfig(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should return error when config id is invalid", func(t *testing.T) {
		param := dto.SaveConfig{
			ConfigId: 0,
		}
		err := s.UpdateConfig(param)
		assert.Error(t, err)
		assert.Equal(t, xerrors.ErrParam, err)
	})

	t.Run("should update config successfully", func(t *testing.T) {
		// First, create a config to update
		createParam := dto.SaveConfig{
			ConfigName:  "config-to-update",
			ConfigKey:   "key-to-update",
			ConfigValue: "value-to-update",
			ConfigType:  "Y",
			CreateBy:    "tester",
		}
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		err := s.CreateConfig(createParam)
		assert.NoError(t, err)

		// Now, update the config
		updateParam := dto.SaveConfig{
			ConfigId:    1, // Assuming this is the first record
			ConfigName:  "updated-config",
			ConfigKey:   "updated-key",
			ConfigValue: "updated-value",
			UpdateBy:    "tester",
		}
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		err = s.UpdateConfig(updateParam)
		assert.NoError(t, err)

		// Verify the update
		var updatedConfig model.SysConfig
		err = dal.Gorm.First(&updatedConfig, 1).Error
		assert.NoError(t, err)
		assert.Equal(t, "updated-config", updatedConfig.ConfigName)
		assert.Equal(t, "updated-key", updatedConfig.ConfigKey)
		assert.Equal(t, "updated-value", updatedConfig.ConfigValue)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestConfigService_DeleteConfig(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should return error when config ids is empty", func(t *testing.T) {
		err := s.DeleteConfig([]int{})
		assert.Error(t, err)
		assert.Equal(t, xerrors.ErrParam, err)
	})

	t.Run("should delete config successfully", func(t *testing.T) {
		// First, create a config to delete
		createParam := dto.SaveConfig{
			ConfigName:  "config-to-delete",
			ConfigKey:   "key-to-delete",
			ConfigValue: "value-to-delete",
			ConfigType:  "Y",
			CreateBy:    "tester",
		}
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		err := s.CreateConfig(createParam)
		assert.NoError(t, err)

		// Now, delete the config
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		err = s.DeleteConfig([]int{1})
		assert.NoError(t, err)

		// Verify the deletion
		var deletedConfig model.SysConfig
		err = dal.Gorm.First(&deletedConfig, 1).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestConfigService_GetConfigList(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should return all configs when no params are given", func(t *testing.T) {
		// Create some configs
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		s.CreateConfig(dto.SaveConfig{ConfigName: "c1", ConfigKey: "k1", ConfigValue: "v1"})
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		s.CreateConfig(dto.SaveConfig{ConfigName: "c2", ConfigKey: "k2", ConfigValue: "v2"})

		configs, count := s.GetConfigList(dto.ConfigListRequest{}, false)
		assert.Len(t, configs, 2)
		assert.Equal(t, 0, count)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestConfigService_GetConfigByConfigId(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should return empty config when id is invalid", func(t *testing.T) {
		config := s.GetConfigByConfigId(0)
		assert.Empty(t, config)
	})

	t.Run("should return config successfully", func(t *testing.T) {
		// Create a config
		createParam := dto.SaveConfig{
			ConfigName:  "test-config-by-id",
			ConfigKey:   "test-key-by-id",
			ConfigValue: "test-value-by-id",
		}
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		err := s.CreateConfig(createParam)
		assert.NoError(t, err)

		// Get the created config to find its ID
		var createdConfig model.SysConfig
		dal.Gorm.Last(&createdConfig)

		// Get the config
		config := s.GetConfigByConfigId(createdConfig.ConfigId)
		assert.NotEmpty(t, config)
		assert.Equal(t, "test-config-by-id", config.ConfigName)
	})
}

func TestConfigService_GetConfigByConfigKey(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should return empty config when key is empty", func(t *testing.T) {
		config := s.GetConfigByConfigKey("")
		assert.Empty(t, config)
	})

	t.Run("should return config successfully", func(t *testing.T) {
		// Create a config
		createParam := dto.SaveConfig{
			ConfigName:  "test-config-by-key",
			ConfigKey:   "test-key-by-key",
			ConfigValue: "test-value-by-key",
		}
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		s.CreateConfig(createParam)

		// Get the config
		config := s.GetConfigByConfigKey("test-key-by-key")
		assert.NotEmpty(t, config)
		assert.Equal(t, "test-config-by-key", config.ConfigName)
	})
}

func TestConfigService_GetConfigCacheByConfigKey(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should return config from cache", func(t *testing.T) {
		// Mock cache
		cachedConfig := `{"ConfigId":1,"ConfigName":"cached-config","ConfigKey":"cached-key","ConfigValue":"cached-value"}`
		redisMock.ExpectHGet(rediskey.SysConfigKey(), "cached-key").SetVal(cachedConfig)

		// Get the config
		config := s.GetConfigCacheByConfigKey("cached-key")
		assert.NotEmpty(t, config)
		assert.Equal(t, "cached-config", config.ConfigName)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return config from db and set cache", func(t *testing.T) {
		// Create a config
		createParam := dto.SaveConfig{
			ConfigName:  "db-config",
			ConfigKey:   "db-key",
			ConfigValue: "db-value",
		}
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		s.CreateConfig(createParam)

		// Get the created config to find its ID
		var createdConfig dto.ConfigDetailResponse
		dal.Gorm.Model(&model.SysConfig{}).Last(&createdConfig)

		// Marshal the created config to JSON
		configBytes, err := json.Marshal(&createdConfig)
		assert.NoError(t, err)

		// Mock cache miss and set
		redisMock.ExpectHGet(rediskey.SysConfigKey(), "db-key").RedisNil()
		redisMock.ExpectHSet(rediskey.SysConfigKey(), "db-key", string(configBytes)).SetVal(1)
		redisMock.ExpectExpire(rediskey.SysConfigKey(), 24*time.Hour).SetVal(true)

		// Get the config
		config := s.GetConfigCacheByConfigKey("db-key")
		assert.NotEmpty(t, config)
		assert.Equal(t, "db-config", config.ConfigName)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestConfigService_RefreshCache(t *testing.T) {
	setup()
	defer teardown()
	s := &ConfigService{}

	t.Run("should refresh cache successfully", func(t *testing.T) {
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetVal(1)
		err := s.RefreshCache()
		assert.NoError(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return error when redis fails", func(t *testing.T) {
		redisMock.ExpectDel(rediskey.SysConfigKey()).SetErr(gorm.ErrInvalidDB)
		err := s.RefreshCache()
		assert.Error(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}
