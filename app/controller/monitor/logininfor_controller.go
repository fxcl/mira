package monitorcontroller

import (
	"time"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/service"
	"mira/common/utils"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

// LogininforController handles login log related operations.
type LogininforController struct {
	LogininforService *service.LogininforService
}

// NewLogininforController creates a new LogininforController.
func NewLogininforController(logininforService *service.LogininforService) *LogininforController {
	return &LogininforController{LogininforService: logininforService}
}

// List retrieves a paginated list of login logs.
// @Summary Get login log list
// @Description Retrieves a paginated list of login logs based on query parameters.
// @Tags Monitor
// @Accept json
// @Produce json
// @Param query body dto.LogininforListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.LogininforListResponse}} "Success"
// @Router /monitor/logininfor/list [get]
func (c *LogininforController) List(ctx *gin.Context) {
	var param dto.LogininforListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	param.OrderRule, param.OrderByColumn = utils.ParseSort(param.IsAsc, param.OrderByColumn, "loginTime")

	logininfors, total := c.LogininforService.GetLogininforList(param, true)

	response.NewSuccess().SetPageData(logininfors, total).Json(ctx)
}

// Remove deletes one or more login logs.
// @Summary Delete login log
// @Description Deletes login logs by their IDs.
// @Tags Monitor
// @Accept json
// @Produce json
// @Param infoIds path string true "Login log IDs, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /monitor/logininfor/{infoIds} [delete]
func (c *LogininforController) Remove(ctx *gin.Context) {
	infoIds, err := utils.StringToIntSlice(ctx.Param("infoIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = c.LogininforService.DeleteLogininfor(infoIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Clean clears all login logs.
// @Summary Clean login log
// @Description Clears all login logs from the system.
// @Tags Monitor
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success"
// @Router /monitor/logininfor/clean [delete]
func (c *LogininforController) Clean(ctx *gin.Context) {
	if err := c.LogininforService.DeleteLogininfor(nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Unlock unlocks a user account.
// @Summary Unlock user account
// @Description Unlocks a user account by deleting the login failure cache.
// @Tags Monitor
// @Accept json
// @Produce json
// @Param userName path string true "Username to unlock"
// @Success 200 {object} response.Response "Success"
// @Router /monitor/logininfor/unlock/{userName} [get]
func (c *LogininforController) Unlock(ctx *gin.Context) {
	if err := c.LogininforService.Unlock(ctx.Param("userName")); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Export exports login logs to an Excel file.
// @Summary Export login log
// @Description Exports login logs to an Excel file based on query parameters.
// @Tags Monitor
// @Accept json
// @Produce json
// @Param query body dto.LogininforListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /monitor/logininfor/export [post]
func (c *LogininforController) Export(ctx *gin.Context) {
	var param dto.LogininforListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	param.OrderRule, param.OrderByColumn = utils.ParseSort(param.IsAsc, param.OrderByColumn, "loginTime")

	list := make([]dto.LogininforExportResponse, 0)

	logininfors, _ := c.LogininforService.GetLogininforList(param, false)
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
