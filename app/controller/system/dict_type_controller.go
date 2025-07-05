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

type DictTypeController struct{}

// Dictionary type list
func (*DictTypeController) List(ctx *gin.Context) {
	var param dto.DictTypeListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	dictTypes, total := (&service.DictTypeService{}).GetDictTypeList(param, true)

	response.NewSuccess().SetPageData(dictTypes, total).Json(ctx)
}

// Dictionary type details
func (*DictTypeController) Detail(ctx *gin.Context) {
	dictId, _ := strconv.Atoi(ctx.Param("dictId"))

	dictType := (&service.DictTypeService{}).GetDictTypeByDictId(dictId)

	response.NewSuccess().SetData("data", dictType).Json(ctx)
}

// Add dictionary type
func (*DictTypeController) Create(ctx *gin.Context) {
	var param dto.CreateDictTypeRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateDictTypeValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dictType := (&service.DictTypeService{}).GetDcitTypeByDictType(param.DictType); dictType.DictId > 0 {
		response.NewError().SetMsg("Failed to add dictionary " + param.DictName + ", dictionary type already exists").Json(ctx)
		return
	}

	if err := (&service.DictTypeService{}).CreateDictType(dto.SaveDictType{
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

// Update dictionary type
func (*DictTypeController) Update(ctx *gin.Context) {
	var param dto.UpdateDictTypeRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateDictTypeValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dictType := (&service.DictTypeService{}).GetDcitTypeByDictType(param.DictType); dictType.DictId > 0 && dictType.DictId != param.DictId {
		response.NewError().SetMsg("Failed to modify dictionary " + param.DictName + ", dictionary type already exists").Json(ctx)
		return
	}

	if err := (&service.DictTypeService{}).UpdateDictType(dto.SaveDictType{
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

// Delete dictionary type
func (*DictTypeController) Remove(ctx *gin.Context) {
	dictIds, err := utils.StringToIntSlice(ctx.Param("dictIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = (&service.DictTypeService{}).DeleteDictType(dictIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Get dictionary selection box list
func (*DictTypeController) Optionselect(ctx *gin.Context) {
	dictTypes, _ := (&service.DictTypeService{}).GetDictTypeList(dto.DictTypeListRequest{
		Status: "0",
	}, false)

	response.NewSuccess().SetData("data", dictTypes).Json(ctx)
}

// Data export
func (*DictTypeController) Export(ctx *gin.Context) {
	var param dto.DictTypeListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.DictTypeExportResponse, 0)

	dictTypes, _ := (&service.DictTypeService{}).GetDictTypeList(param, false)
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

// Refresh cache
func (*DictTypeController) RefreshCache(ctx *gin.Context) {
	if err := dal.Redis.Del(ctx.Request.Context(), rediskey.SysDictKey).Err(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
