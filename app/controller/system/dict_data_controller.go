package systemcontroller

import (
	"strconv"
	"time"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/types/constant"
	"mira/common/utils"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

// DictDataController handles dictionary data operations.
type DictDataController struct {
	DictDataService *service.DictDataService
}

// NewDictDataController creates a new DictDataController.
func NewDictDataController(dictDataService *service.DictDataService) *DictDataController {
	return &DictDataController{DictDataService: dictDataService}
}

// List retrieves a paginated list of dictionary data.
// @Summary Get dictionary data list
// @Description Retrieves a paginated list of dictionary data based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.DictDataListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.DictDataListResponse}} "Success"
// @Router /system/dict/data/list [get]
func (c *DictDataController) List(ctx *gin.Context) {
	var param dto.DictDataListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	dictDatas, total := c.DictDataService.GetDictDataList(param, true)

	response.NewSuccess().SetPageData(dictDatas, total).Json(ctx)
}

// Detail retrieves the details of a specific dictionary data item.
// @Summary Get dictionary data details
// @Description Retrieves the details of a dictionary data item by its code.
// @Tags System
// @Accept json
// @Produce json
// @Param dictCode path int true "Dictionary Code"
// @Success 200 {object} response.Response{data=dto.DictDataDetailResponse} "Success"
// @Router /system/dict/data/{dictCode} [get]
func (c *DictDataController) Detail(ctx *gin.Context) {
	dictCode, _ := strconv.Atoi(ctx.Param("dictCode"))

	dictData := c.DictDataService.GetDictDataByDictCode(dictCode)

	response.NewSuccess().SetData("data", dictData).Json(ctx)
}

// Create adds new dictionary data.
// @Summary Add dictionary data
// @Description Adds new dictionary data to the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreateDictDataRequest true "Dictionary data"
// @Success 200 {object} response.Response "Success"
// @Router /system/dict/data [post]
func (c *DictDataController) Create(ctx *gin.Context) {
	var param dto.CreateDictDataRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateDictDataValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := c.DictDataService.CreateDictData(dto.SaveDictData{
		DictSort:  param.DictSort,
		DictLabel: param.DictLabel,
		DictValue: param.DictValue,
		DictType:  param.DictType,
		CssClass:  param.CssClass,
		ListClass: param.ListClass,
		IsDefault: param.IsDefault,
		Status:    param.Status,
		CreateBy:  security.GetAuthUserName(ctx),
		Remark:    param.Remark,
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update modifies existing dictionary data.
// @Summary Update dictionary data
// @Description Modifies existing dictionary data in the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateDictDataRequest true "Dictionary data"
// @Success 200 {object} response.Response "Success"
// @Router /system/dict/data [put]
func (c *DictDataController) Update(ctx *gin.Context) {
	var param dto.UpdateDictDataRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateDictDataValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := c.DictDataService.UpdateDictData(dto.SaveDictData{
		DictCode:  param.DictCode,
		DictSort:  param.DictSort,
		DictLabel: param.DictLabel,
		DictValue: param.DictValue,
		DictType:  param.DictType,
		CssClass:  param.CssClass,
		ListClass: param.ListClass,
		IsDefault: param.IsDefault,
		Status:    param.Status,
		UpdateBy:  security.GetAuthUserName(ctx),
		Remark:    param.Remark,
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Remove deletes one or more dictionary data items.
// @Summary Delete dictionary data
// @Description Deletes dictionary data items by their codes.
// @Tags System
// @Accept json
// @Produce json
// @Param dictCodes path string true "Dictionary Codes, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /system/dict/data/{dictCodes} [delete]
func (c *DictDataController) Remove(ctx *gin.Context) {
	dictCodes, err := utils.StringToIntSlice(ctx.Param("dictCodes"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = c.DictDataService.DeleteDictData(dictCodes); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Type retrieves dictionary data by dictionary type.
// @Summary Get dictionary data by type
// @Description Retrieves a list of dictionary data items for a given dictionary type.
// @Tags System
// @Accept json
// @Produce json
// @Param dictType path string true "Dictionary Type"
// @Success 200 {object} response.Response{data=[]dto.DictDataCacheResponse} "Success"
// @Router /system/dict/data/type/{dictType} [get]
func (c *DictDataController) Type(ctx *gin.Context) {
	dictType := ctx.Param("dictType")

	dictDatas := c.DictDataService.GetDictDataCacheByDictType(dictType)

	for key, dictData := range dictDatas {
		dictDatas[key].Default = dictData.IsDefault == constant.IS_DEFAULT_YES
	}

	response.NewSuccess().SetData("data", dictDatas).Json(ctx)
}

// Export exports dictionary data to an Excel file.
// @Summary Export dictionary data
// @Description Exports dictionary data to an Excel file based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.DictDataListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /system/dict/data/export [post]
func (c *DictDataController) Export(ctx *gin.Context) {
	var param dto.DictDataListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.DictDataExportResponse, 0)

	dictDatas, _ := c.DictDataService.GetDictDataList(param, false)
	for _, dictData := range dictDatas {
		list = append(list, dto.DictDataExportResponse{
			DictCode:  dictData.DictCode,
			DictSort:  dictData.DictSort,
			DictLabel: dictData.DictLabel,
			DictValue: dictData.DictValue,
			DictType:  dictData.DictType,
			IsDefault: dictData.IsDefault,
			Status:    dictData.Status,
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("data_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}
