package systemcontroller

import (
	"strconv"
	"time"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/utils"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

// ConfigController handles parameter configuration operations.
type ConfigController struct {
	ConfigService *service.ConfigService
}

// NewConfigController creates a new ConfigController.
func NewConfigController(configService *service.ConfigService) *ConfigController {
	return &ConfigController{ConfigService: configService}
}

// List retrieves a paginated list of parameters.
// @Summary Get parameter list
// @Description Retrieves a paginated list of parameters based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.ConfigListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.ConfigListResponse}} "Success"
// @Router /system/config/list [get]
func (c *ConfigController) List(ctx *gin.Context) {
	var param dto.ConfigListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	configs, total := c.ConfigService.GetConfigList(param, true)

	response.NewSuccess().SetPageData(configs, total).Json(ctx)
}

// Detail retrieves the details of a specific parameter.
// @Summary Get parameter details
// @Description Retrieves the details of a parameter by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param configId path int true "Config ID"
// @Success 200 {object} response.Response{data=dto.ConfigDetailResponse} "Success"
// @Router /system/config/{configId} [get]
func (c *ConfigController) Detail(ctx *gin.Context) {
	configId, _ := strconv.Atoi(ctx.Param("configId"))

	config := c.ConfigService.GetConfigByConfigId(configId)

	response.NewSuccess().SetData("data", config).Json(ctx)
}

// Create adds a new parameter.
// @Summary Add parameter
// @Description Adds a new parameter to the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreateConfigRequest true "Parameter data"
// @Success 200 {object} response.Response "Success"
// @Router /system/config [post]
func (c *ConfigController) Create(ctx *gin.Context) {
	var param dto.CreateConfigRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateConfigValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if config := c.ConfigService.GetConfigByConfigKey(param.ConfigKey); config.ConfigId > 0 {
		response.NewError().SetMsg("Failed to add parameter " + param.ConfigName + ", parameter key name already exists").Json(ctx)
		return
	}

	if err := c.ConfigService.CreateConfig(dto.SaveConfig{
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

// Update modifies an existing parameter.
// @Summary Update parameter
// @Description Modifies an existing parameter in the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateConfigRequest true "Parameter data"
// @Success 200 {object} response.Response "Success"
// @Router /system/config [put]
func (c *ConfigController) Update(ctx *gin.Context) {
	var param dto.UpdateConfigRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateConfigValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if config := c.ConfigService.GetConfigByConfigKey(param.ConfigKey); config.ConfigId > 0 && config.ConfigId != param.ConfigId {
		response.NewError().SetMsg("Failed to modify parameter " + param.ConfigName + ", parameter key name already exists").Json(ctx)
		return
	}

	if err := c.ConfigService.UpdateConfig(dto.SaveConfig{
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

// Remove deletes one or more parameters.
// @Summary Delete parameter
// @Description Deletes parameters by their IDs.
// @Tags System
// @Accept json
// @Produce json
// @Param configIds path string true "Config IDs, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /system/config/{configIds} [delete]
func (c *ConfigController) Remove(ctx *gin.Context) {
	configIds, err := utils.StringToIntSlice(ctx.Param("configIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = c.ConfigService.DeleteConfig(configIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// ConfigKey retrieves a configuration value by its key.
// @Summary Get configuration value by key
// @Description Retrieves a configuration value by its unique key.
// @Tags System
// @Accept json
// @Produce json
// @Param configKey path string true "Config Key"
// @Success 200 {object} response.Response{msg=string} "Success"
// @Router /system/config/configKey/{configKey} [get]
func (c *ConfigController) ConfigKey(ctx *gin.Context) {
	configKey := ctx.Param("configKey")

	config := c.ConfigService.GetConfigCacheByConfigKey(configKey)

	response.NewSuccess().SetMsg(config.ConfigValue).Json(ctx)
}

// Export exports parameter data to an Excel file.
// @Summary Export parameters
// @Description Exports parameter data to an Excel file based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.ConfigListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /system/config/export [post]
func (c *ConfigController) Export(ctx *gin.Context) {
	var param dto.ConfigListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.ConfigExportResponse, 0)

	configs, _ := c.ConfigService.GetConfigList(param, false)
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

// RefreshCache refreshes the parameter cache.
// @Summary Refresh parameter cache
// @Description Clears the parameter configuration cache in Redis.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success"
// @Router /system/config/refreshCache [delete]
func (c *ConfigController) RefreshCache(ctx *gin.Context) {
	if err := c.ConfigService.RefreshCache(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
