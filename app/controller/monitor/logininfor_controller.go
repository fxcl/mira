package monitorcontroller

import (
	"mira/anima/dal"
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/service"
	"mira/common/utils"
	"regexp"
	"strings"
	"time"

	rediskey "mira/common/types/redis-key"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

type LogininforController struct{}

// Login log list
func (*LogininforController) List(ctx *gin.Context) {
	var param dto.LogininforListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	// The default sorting rule is descending (DESC)
	param.OrderRule = "DESC"
	if strings.HasPrefix(param.IsAsc, "asc") {
		param.OrderRule = "ASC"
	}

	// Sort field camel case to snake case
	if param.OrderByColumn == "" {
		param.OrderByColumn = "loginTime"
	}
	param.OrderByColumn = strings.ToLower(regexp.MustCompile("([A-Z])").ReplaceAllString(param.OrderByColumn, "_${1}"))

	logininfors, total := (&service.LogininforService{}).GetLogininforList(param, true)

	response.NewSuccess().SetPageData(logininfors, total).Json(ctx)
}

// Delete login log
func (*LogininforController) Remove(ctx *gin.Context) {
	infoIds, err := utils.StringToIntSlice(ctx.Param("infoIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = (&service.LogininforService{}).DeleteLogininfor(infoIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Clean login log
func (*LogininforController) Clean(ctx *gin.Context) {
	if err := (&service.LogininforService{}).DeleteLogininfor(nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Account unlock (delete the 10-minute cache for login error count limit)
func (*LogininforController) Unlock(ctx *gin.Context) {
	if _, err := dal.Redis.Del(ctx.Request.Context(), rediskey.LoginPasswordErrorKey+ctx.Param("userName")).Result(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Data export
func (*LogininforController) Export(ctx *gin.Context) {
	var param dto.LogininforListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	// The default sorting rule is descending (DESC)
	param.OrderRule = "DESC"
	if strings.HasPrefix(param.IsAsc, "asc") {
		param.OrderRule = "ASC"
	}

	// Sort field camel case to snake case
	if param.OrderByColumn == "" {
		param.OrderByColumn = "loginTime"
	}
	param.OrderByColumn = strings.ToLower(regexp.MustCompile("([A-Z])").ReplaceAllString(param.OrderByColumn, "_${1}"))

	list := make([]dto.LogininforExportResponse, 0)

	logininfors, _ := (&service.LogininforService{}).GetLogininforList(param, false)
	for _, logininfor := range logininfors {
		list = append(list, dto.LogininforExportResponse{
			InfoId:        logininfor.InfoId,
			UserName:      logininfor.UserName,
			Status:        logininfor.Status,
			Ipaddr:        logininfor.Ipaddr,
			LoginLocation: logininfor.LoginLocation,
			Browser:       logininfor.Browser,
			Os:            logininfor.Os,
			Msg:           logininfor.Msg,
			LoginTime:     logininfor.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("logininfor_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}
