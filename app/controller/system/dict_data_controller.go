package systemcontroller

import (
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/types/constant"
	"mira/common/utils"
	"strconv"
	"time"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

type DictDataController struct{}

// Get dictionary data list
func (*DictDataController) List(ctx *gin.Context) {
	var param dto.DictDataListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	dictDatas, total := (&service.DictDataService{}).GetDictDataList(param, true)

	response.NewSuccess().SetPageData(dictDatas, total).Json(ctx)
}

// Get dictionary data details
func (*DictDataController) Detail(ctx *gin.Context) {
	dictCode, _ := strconv.Atoi(ctx.Param("dictCode"))

	dictData := (&service.DictDataService{}).GetDictDataByDictCode(dictCode)

	response.NewSuccess().SetData("data", dictData).Json(ctx)
}

// Add dictionary data
func (*DictDataController) Create(ctx *gin.Context) {
	var param dto.CreateDictDataRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateDictDataValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.DictDataService{}).CreateDictData(dto.SaveDictData{
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

// Update dictionary data
func (*DictDataController) Update(ctx *gin.Context) {
	var param dto.UpdateDictDataRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateDictDataValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.DictDataService{}).UpdateDictData(dto.SaveDictData{
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

// Delete dictionary data
func (*DictDataController) Remove(ctx *gin.Context) {
	dictCodes, err := utils.StringToIntSlice(ctx.Param("dictCodes"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = (&service.DictDataService{}).DeleteDictData(dictCodes); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Query dictionary data by dictionary type
func (*DictDataController) Type(ctx *gin.Context) {
	dictType := ctx.Param("dictType")

	dictDatas := (&service.DictDataService{}).GetDictDataCacheByDictType(dictType)

	for key, dictData := range dictDatas {
		dictDatas[key].Default = dictData.IsDefault == constant.IS_DEFAULT_YES
	}

	response.NewSuccess().SetData("data", dictDatas).Json(ctx)
}

// Data export
func (*DictDataController) Export(ctx *gin.Context) {
	var param dto.DictDataListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.DictDataExportResponse, 0)

	dictDatas, _ := (&service.DictDataService{}).GetDictDataList(param, false)
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
