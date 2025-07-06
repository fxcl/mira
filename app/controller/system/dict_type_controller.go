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

// DictTypeController handles dictionary type operations.
type DictTypeController struct {
	DictTypeService *service.DictTypeService
}

// NewDictTypeController creates a new DictTypeController.
func NewDictTypeController(dictTypeService *service.DictTypeService) *DictTypeController {
	return &DictTypeController{DictTypeService: dictTypeService}
}

// List retrieves a paginated list of dictionary types.
// @Summary Get dictionary type list
// @Description Retrieves a paginated list of dictionary types based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.DictTypeListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.DictTypeListResponse}} "Success"
// @Router /system/dict/type/list [get]
func (c *DictTypeController) List(ctx *gin.Context) {
	var param dto.DictTypeListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	dictTypes, total := c.DictTypeService.GetDictTypeList(param, true)

	response.NewSuccess().SetPageData(dictTypes, total).Json(ctx)
}

// Detail retrieves the details of a specific dictionary type.
// @Summary Get dictionary type details
// @Description Retrieves the details of a dictionary type by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param dictId path int true "Dictionary ID"
// @Success 200 {object} response.Response{data=dto.DictTypeDetailResponse} "Success"
// @Router /system/dict/type/{dictId} [get]
func (c *DictTypeController) Detail(ctx *gin.Context) {
	dictId, _ := strconv.Atoi(ctx.Param("dictId"))

	dictType := c.DictTypeService.GetDictTypeByDictId(dictId)

	response.NewSuccess().SetData("data", dictType).Json(ctx)
}

// Create adds a new dictionary type.
// @Summary Add dictionary type
// @Description Adds a new dictionary type to the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreateDictTypeRequest true "Dictionary type data"
// @Success 200 {object} response.Response "Success"
// @Router /system/dict/type [post]
func (c *DictTypeController) Create(ctx *gin.Context) {
	var param dto.CreateDictTypeRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateDictTypeValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dictType := c.DictTypeService.GetDcitTypeByDictType(param.DictType); dictType.DictId > 0 {
		response.NewError().SetMsg("Failed to add dictionary " + param.DictName + ", dictionary type already exists").Json(ctx)
		return
	}

	if err := c.DictTypeService.CreateDictType(dto.SaveDictType{
		DictName: param.DictName,
		DictType: param.DictType,
		Status:   param.Status,
		CreateBy: security.GetAuthUserName(ctx),
		Remark:   param.Remark,
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update modifies an existing dictionary type.
// @Summary Update dictionary type
// @Description Modifies an existing dictionary type in the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateDictTypeRequest true "Dictionary type data"
// @Success 200 {object} response.Response "Success"
// @Router /system/dict/type [put]
func (c *DictTypeController) Update(ctx *gin.Context) {
	var param dto.UpdateDictTypeRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateDictTypeValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dictType := c.DictTypeService.GetDcitTypeByDictType(param.DictType); dictType.DictId > 0 && dictType.DictId != param.DictId {
		response.NewError().SetMsg("Failed to modify dictionary " + param.DictName + ", dictionary type already exists").Json(ctx)
		return
	}

	if err := c.DictTypeService.UpdateDictType(dto.SaveDictType{
		DictId:   param.DictId,
		DictName: param.DictName,
		DictType: param.DictType,
		Status:   param.Status,
		UpdateBy: security.GetAuthUserName(ctx),
		Remark:   param.Remark,
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Remove deletes one or more dictionary types.
// @Summary Delete dictionary type
// @Description Deletes dictionary types by their IDs.
// @Tags System
// @Accept json
// @Produce json
// @Param dictIds path string true "Dictionary IDs, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /system/dict/type/{dictIds} [delete]
func (c *DictTypeController) Remove(ctx *gin.Context) {
	dictIds, err := utils.StringToIntSlice(ctx.Param("dictIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = c.DictTypeService.DeleteDictType(dictIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Optionselect retrieves a list of dictionary types for selection.
// @Summary Get dictionary select box list
// @Description Retrieves a list of all active dictionary types for use in dropdowns.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.DictTypeListResponse} "Success"
// @Router /system/dict/type/optionselect [get]
func (c *DictTypeController) Optionselect(ctx *gin.Context) {
	dictTypes, _ := c.DictTypeService.GetDictTypeList(dto.DictTypeListRequest{
		Status: "0",
	}, false)

	response.NewSuccess().SetData("data", dictTypes).Json(ctx)
}

// Export exports dictionary type data to an Excel file.
// @Summary Export dictionary types
// @Description Exports dictionary type data to an Excel file based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.DictTypeListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /system/dict/type/export [post]
func (c *DictTypeController) Export(ctx *gin.Context) {
	var param dto.DictTypeListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.DictTypeExportResponse, 0)

	dictTypes, _ := c.DictTypeService.GetDictTypeList(param, false)
	for _, dictType := range dictTypes {
		list = append(list, dto.DictTypeExportResponse{
			DictId:   dictType.DictId,
			DictName: dictType.DictName,
			DictType: dictType.DictType,
			Status:   dictType.Status,
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("type_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}

// RefreshCache refreshes the dictionary cache.
// @Summary Refresh dictionary cache
// @Description Clears the dictionary cache in Redis.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success"
// @Router /system/dict/type/refreshCache [delete]
func (c *DictTypeController) RefreshCache(ctx *gin.Context) {
	if err := c.DictTypeService.RefreshCache(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
