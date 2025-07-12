package service

import (
	"encoding/json"
	"testing"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	rediskey "mira/common/types/redis-key"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDictTypeService_CreateDictType(t *testing.T) {
	setup()
	defer teardown()
	s := &DictTypeService{}

	t.Run("should create dict type successfully", func(t *testing.T) {
		param := dto.SaveDictType{
			DictName: "Test Dict Type",
			DictType: "test_dict_type",
			Status:   "0",
			CreateBy: "tester",
		}
		err := s.CreateDictType(param)
		assert.NoError(t, err)

		// Verify
		var createdDictType model.SysDictType
		dal.Gorm.Last(&createdDictType)
		assert.Equal(t, "Test Dict Type", createdDictType.DictName)
	})
}

func TestDictTypeService_UpdateDictType(t *testing.T) {
	setup()
	defer teardown()
	s := &DictTypeService{}

	t.Run("should update dict type successfully", func(t *testing.T) {
		createParam := dto.SaveDictType{
			DictName: "Dict Type to Update",
			DictType: "type_to_update",
		}
		s.CreateDictType(createParam)
		var createdDictType model.SysDictType
		dal.Gorm.Last(&createdDictType)

		// Execute
		updateParam := dto.SaveDictType{
			DictId:   createdDictType.DictId,
			DictName: "Updated Dict Type",
		}
		err := s.UpdateDictType(updateParam)
		assert.NoError(t, err)

		// Verify
		var updatedDictType model.SysDictType
		dal.Gorm.First(&updatedDictType, createdDictType.DictId)
		assert.Equal(t, "Updated Dict Type", updatedDictType.DictName)
	})
}

func TestDictTypeService_DeleteDictType(t *testing.T) {
	setup()
	defer teardown()
	s := &DictTypeService{}

	t.Run("should delete dict type successfully", func(t *testing.T) {
		createParam := dto.SaveDictType{
			DictName: "Dict Type to Delete",
			DictType: "type_to_delete",
		}
		s.CreateDictType(createParam)
		var createdDictType model.SysDictType
		dal.Gorm.Last(&createdDictType)

		// Execute
		err := s.DeleteDictType([]int{createdDictType.DictId})
		assert.NoError(t, err)

		// Verify
		var deletedDictType model.SysDictType
		err = dal.Gorm.First(&deletedDictType, createdDictType.DictId).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestDictTypeService_GetDictTypeList(t *testing.T) {
	setup()
	defer teardown()
	s := &DictTypeService{}

	t.Run("should return all dict types", func(t *testing.T) {
		s.CreateDictType(dto.SaveDictType{DictName: "Dict Type 1", DictType: "type1"})
		s.CreateDictType(dto.SaveDictType{DictName: "Dict Type 2", DictType: "type2"})

		// Execute
		dictTypes, count := s.GetDictTypeList(dto.DictTypeListRequest{}, false)
		assert.Len(t, dictTypes, 2)
		assert.Equal(t, 0, count)
	})
}

func TestDictTypeService_GetDictTypeByDictId(t *testing.T) {
	setup()
	defer teardown()
	s := &DictTypeService{}

	t.Run("should return dict type successfully", func(t *testing.T) {
		createParam := dto.SaveDictType{
			DictName: "Dict Type By Id",
			DictType: "type_by_id",
		}
		s.CreateDictType(createParam)
		var createdDictType model.SysDictType
		dal.Gorm.Last(&createdDictType)

		// Execute
		dictType := s.GetDictTypeByDictId(createdDictType.DictId)
		assert.Equal(t, "Dict Type By Id", dictType.DictName)
	})
}

func TestDictTypeService_GetDcitTypeByDictType(t *testing.T) {
	setup()
	defer teardown()
	s := &DictTypeService{}

	t.Run("should return dict type successfully", func(t *testing.T) {
		createParam := dto.SaveDictType{
			DictName: "Dict Type By Type",
			DictType: "type_by_type",
		}
		s.CreateDictType(createParam)

		// Execute
		dictType := s.GetDcitTypeByDictType("type_by_type")
		assert.Equal(t, "Dict Type By Type", dictType.DictName)
	})
}

func TestDictTypeService_RefreshCache(t *testing.T) {
	setup()
	defer teardown()
	s := &DictTypeService{}

	t.Run("should refresh cache successfully", func(t *testing.T) {
		redisMock.ExpectDel(rediskey.SysDictKey()).SetVal(1)
		err := s.RefreshCache()
		assert.NoError(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return error when redis fails", func(t *testing.T) {
		redisMock.ExpectDel(rediskey.SysDictKey()).SetErr(gorm.ErrInvalidDB)
		err := s.RefreshCache()
		assert.Error(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestDictDataService_CreateDictData(t *testing.T) {
	setup()
	defer teardown()
	s := &DictDataService{}

	t.Run("should create dict data successfully", func(t *testing.T) {
		param := dto.SaveDictData{
			DictLabel: "Test Label",
			DictValue: "test_value",
			DictType:  "test_type",
			Status:    "0",
			CreateBy:  "tester",
		}
		err := s.CreateDictData(param)
		assert.NoError(t, err)

		// Verify
		var createdDictData model.SysDictData
		dal.Gorm.Last(&createdDictData)
		assert.Equal(t, "Test Label", createdDictData.DictLabel)
	})
}

func TestDictDataService_UpdateDictData(t *testing.T) {
	setup()
	defer teardown()
	s := &DictDataService{}

	t.Run("should update dict data successfully", func(t *testing.T) {
		createParam := dto.SaveDictData{
			DictLabel: "Label to Update",
			DictValue: "value_to_update",
			DictType:  "type_to_update",
		}
		s.CreateDictData(createParam)
		var createdDictData model.SysDictData
		dal.Gorm.Last(&createdDictData)

		// Execute
		updateParam := dto.SaveDictData{
			DictCode:  createdDictData.DictCode,
			DictLabel: "Updated Label",
		}
		err := s.UpdateDictData(updateParam)
		assert.NoError(t, err)

		// Verify
		var updatedDictData model.SysDictData
		dal.Gorm.First(&updatedDictData, createdDictData.DictCode)
		assert.Equal(t, "Updated Label", updatedDictData.DictLabel)
	})
}

func TestDictDataService_DeleteDictData(t *testing.T) {
	setup()
	defer teardown()
	s := &DictDataService{}

	t.Run("should delete dict data successfully", func(t *testing.T) {
		createParam := dto.SaveDictData{
			DictLabel: "Label to Delete",
			DictValue: "value_to_delete",
			DictType:  "type_to_delete",
		}
		s.CreateDictData(createParam)
		var createdDictData model.SysDictData
		dal.Gorm.Last(&createdDictData)

		// Execute
		err := s.DeleteDictData([]int{createdDictData.DictCode})
		assert.NoError(t, err)

		// Verify
		var deletedDictData model.SysDictData
		err = dal.Gorm.First(&deletedDictData, createdDictData.DictCode).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestDictDataService_GetDictDataList(t *testing.T) {
	setup()
	defer teardown()
	s := &DictDataService{}

	t.Run("should return all dict data", func(t *testing.T) {
		s.CreateDictData(dto.SaveDictData{DictLabel: "Label 1", DictValue: "value1", DictType: "type1"})
		s.CreateDictData(dto.SaveDictData{DictLabel: "Label 2", DictValue: "value2", DictType: "type2"})

		// Execute
		dictDatas, count := s.GetDictDataList(dto.DictDataListRequest{}, false)
		assert.Len(t, dictDatas, 2)
		assert.Equal(t, 0, count)
	})
}

func TestDictDataService_GetDictDataByDictCode(t *testing.T) {
	setup()
	defer teardown()
	s := &DictDataService{}

	t.Run("should return dict data successfully", func(t *testing.T) {
		createParam := dto.SaveDictData{
			DictLabel: "Label By Code",
			DictValue: "value_by_code",
			DictType:  "type_by_code",
		}
		s.CreateDictData(createParam)
		var createdDictData model.SysDictData
		dal.Gorm.Last(&createdDictData)

		// Execute
		dictData := s.GetDictDataByDictCode(createdDictData.DictCode)
		assert.Equal(t, "Label By Code", dictData.DictLabel)
	})
}

func TestDictDataService_GetDictDataByDictType(t *testing.T) {
	setup()
	defer teardown()
	s := &DictDataService{}

	t.Run("should return dict data successfully", func(t *testing.T) {
		s.CreateDictData(dto.SaveDictData{DictLabel: "Label 1", DictValue: "value1", DictType: "test_type"})
		s.CreateDictData(dto.SaveDictData{DictLabel: "Label 2", DictValue: "value2", DictType: "test_type"})

		// Execute
		dictDatas := s.GetDictDataByDictType("test_type")
		assert.Len(t, dictDatas, 2)
	})
}

func TestDictDataService_GetDictDataCacheByDictType(t *testing.T) {
	setup()
	defer teardown()
	s := &DictDataService{}

	t.Run("should return dict data from cache", func(t *testing.T) {
		// Mock cache
		cachedData := `[{"DictLabel":"Label 1","DictValue":"value1","DictType":"test_type"}]`
		redisMock.ExpectHGet(rediskey.SysDictKey(), "test_type").SetVal(cachedData)

		// Execute
		dictDatas := s.GetDictDataCacheByDictType("test_type")
		assert.Len(t, dictDatas, 1)
		assert.Equal(t, "Label 1", dictDatas[0].DictLabel)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return dict data from db and set cache", func(t *testing.T) {
		s.CreateDictData(dto.SaveDictData{DictLabel: "Label 1", DictValue: "value1", DictType: "test_type"})
		var createdDictData []dto.DictDataListResponse
		dal.Gorm.Model(&model.SysDictData{}).Find(&createdDictData)
		createdDictDataBytes, _ := json.Marshal(createdDictData)

		// Mock cache miss and set
		redisMock.ExpectHGet(rediskey.SysDictKey(), "test_type").RedisNil()
		redisMock.ExpectHSet(rediskey.SysDictKey(), "test_type", string(createdDictDataBytes)).SetVal(1)

		// Execute
		dictDatas := s.GetDictDataCacheByDictType("test_type")
		assert.Len(t, dictDatas, 1)
		assert.Equal(t, "Label 1", dictDatas[0].DictLabel)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}
