package monitorcontroller

import (
	"strconv"
	"time"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/service"
	"mira/common/utils"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

// OperlogController handles operation log related operations.
type OperlogController struct {
	OperLogService *service.OperLogService
}

// NewOperlogController creates a new OperlogController.
func NewOperlogController(operLogService *service.OperLogService) *OperlogController {
	return &OperlogController{OperLogService: operLogService}
}

// List retrieves a paginated list of operation logs.
// @Summary Get operation log list
// @Description Retrieves a paginated list of operation logs based on query parameters.
// @Tags Monitor
// @Accept json
// @Produce json
// @Param query body dto.OperLogListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.OperLogListResponse}} "Success"
// @Router /monitor/operlog/list [get]
func (c *OperlogController) List(ctx *gin.Context) {
	var param dto.OperLogListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	param.OrderRule, param.OrderByColumn = utils.ParseSort(param.IsAsc, param.OrderByColumn, "operTime")

	operLogs, total := c.OperLogService.GetOperLogList(param, true)

	response.NewSuccess().SetPageData(operLogs, total).Json(ctx)
}

// Remove deletes one or more operation logs.
// @Summary Delete operation log
// @Description Deletes operation logs by their IDs.
// @Tags Monitor
// @Accept json
// @Produce json
// @Param operIds path string true "Operation log IDs, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /monitor/operlog/{operIds} [delete]
func (c *OperlogController) Remove(ctx *gin.Context) {
	operIds, err := utils.StringToIntSlice(ctx.Param("operIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = c.OperLogService.DeleteOperLog(operIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Clean clears all operation logs.
// @Summary Clean operation log
// @Description Clears all operation logs from the system.
// @Tags Monitor
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success"
// @Router /monitor/operlog/clean [delete]
func (c *OperlogController) Clean(ctx *gin.Context) {
	if err := c.OperLogService.DeleteOperLog(nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Export exports operation logs to an Excel file.
// @Summary Export operation log
// @Description Exports operation logs to an Excel file based on query parameters.
// @Tags Monitor
// @Accept json
// @Produce json
// @Param query body dto.OperLogListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /monitor/operlog/export [post]
func (c *OperlogController) Export(ctx *gin.Context) {
	var param dto.OperLogListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	param.OrderRule, param.OrderByColumn = utils.ParseSort(param.IsAsc, param.OrderByColumn, "operTime")

	list := make([]dto.OperLogExportResponse, 0)

	operLogs, _ := c.OperLogService.GetOperLogList(param, false)
	for _, operLog := range operLogs {
		list = append(list, dto.OperLogExportResponse{
			OperId:        operLog.OperId,
			Title:         operLog.Title,
			BusinessType:  operLog.BusinessType,
			Method:        operLog.Method,
			RequestMethod: operLog.RequestMethod,
			OperName:      operLog.OperName,
			DeptName:      operLog.DeptName,
			OperUrl:       operLog.OperUrl,
			OperIp:        operLog.OperIp,
			OperLocation:  operLog.OperLocation,
			OperParam:     operLog.OperParam,
			JsonResult:    operLog.JsonResult,
			Status:        operLog.Status,
			ErrorMsg:      operLog.ErrorMsg,
			OperTime:      operLog.OperTime.Format("2006-01-02 15:04:05"),
			CostTime:      strconv.Itoa(operLog.CostTime) + "ms",
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("operlog_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}
