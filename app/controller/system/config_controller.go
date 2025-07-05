package systemcontroller

import (
	"mira/anima/dal"
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/utils"
	"strconv"
	"time"

	rediskey "mira/common/types/redis-key"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

type ConfigController struct{}

// Parameter list
func (*ConfigController) List(ctx *gin.Context) {
	var param dto.ConfigListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	configs, total := (&service.ConfigService{}).GetConfigList(param, true)

	response.NewSuccess().SetPageData(configs, total).Json(ctx)
}

// Parameter details
func (*ConfigController) Detail(ctx *gin.Context) {
	configId, _ := strconv.Atoi(ctx.Param("configId"))

	config := (&service.ConfigService{}).GetConfigByConfigId(configId)

	response.NewSuccess().SetData("data", config).Json(ctx)
}

// Add parameter
func (*ConfigController) Create(ctx *gin.Context) {
	var param dto.CreateConfigRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateConfigValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if config := (&service.ConfigService{}).GetConfigByConfigKey(param.ConfigKey); config.ConfigId > 0 {
		response.NewError().SetMsg("Failed to add parameter " + param.ConfigName + ", parameter key name already exists").Json(ctx)
		return
	}

	if err := (&service.ConfigService{}).CreateConfig(dto.SaveConfig{
		ConfigName:  param.ConfigName,
		ConfigKey:   param.ConfigKey,
		ConfigValue: param.ConfigValue,
		ConfigType:  param.ConfigType,
		Remark:      param.Remark,
		CreateBy:    security.GetAuthUserName(ctx),
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update parameter
func (*ConfigController) Update(ctx *gin.Context) {
	var param dto.UpdateConfigRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateConfigValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if config := (&service.ConfigService{}).GetConfigByConfigKey(param.ConfigKey); config.ConfigId > 0 && config.ConfigId != param.ConfigId {
		response.NewError().SetMsg("Failed to modify parameter " + param.ConfigName + ", parameter key name already exists").Json(ctx)
		return
	}

	if err := (&service.ConfigService{}).UpdateConfig(dto.SaveConfig{
		ConfigId:    param.ConfigId,
		ConfigName:  param.ConfigName,
		ConfigKey:   param.ConfigKey,
		ConfigValue: param.ConfigValue,
		ConfigType:  param.ConfigType,
		Remark:      param.Remark,
		UpdateBy:    security.GetAuthUserName(ctx),
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Delete parameter
func (*ConfigController) Remove(ctx *gin.Context) {
	configIds, err := utils.StringToIntSlice(ctx.Param("configIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = (&service.ConfigService{}).DeleteConfig(configIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Get configuration value by configuration key
func (*ConfigController) ConfigKey(ctx *gin.Context) {
	configKey := ctx.Param("configKey")

	config := (&service.ConfigService{}).GetConfigCacheByConfigKey(configKey)

	response.NewSuccess().SetMsg(config.ConfigValue).Json(ctx)
}

// Data export
func (*ConfigController) Export(ctx *gin.Context) {
	var param dto.ConfigListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.ConfigExportResponse, 0)

	configs, _ := (&service.ConfigService{}).GetConfigList(param, false)
	for _, config := range configs {
		list = append(list, dto.ConfigExportResponse{
			ConfigId:    config.ConfigId,
			ConfigName:  config.ConfigName,
			ConfigKey:   config.ConfigKey,
			ConfigValue: config.ConfigValue,
			ConfigType:  config.ConfigType,
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("config_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}

// Refresh cache
func (*ConfigController) RefreshCache(ctx *gin.Context) {
	if err := dal.Redis.Del(ctx.Request.Context(), rediskey.SysConfigKey).Err(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
