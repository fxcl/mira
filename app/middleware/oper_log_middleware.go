package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"mira/anima/datetime"
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/service"
	"mira/app/token"
	ipaddress "mira/common/ip-address"
	responsewriter "mira/common/response-writer"
	"mira/common/types/constant"

	"github.com/gin-gonic/gin"
)

// OperLogMiddleware is a middleware for operation logs.
// title: title of the operation module
// businessType: operation type, constant.REQUEST_BUSINESS_TYPE_*
func OperLogMiddleware(operLogService service.OperLogServiceInterface, title string, businessType int, getAuthUser func(ctx *gin.Context) *token.UserTokenResponse) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var operName, deptName string

		if authUser := getAuthUser(ctx); authUser != nil {
			operName = authUser.NickName
			deptName = authUser.DeptName
		}

		// Record the request start time to calculate the request duration.
		requestStartTime := time.Now()

		// Because the request body will be consumed after reading, to avoid EOF errors,
		// the request body needs to be cached and reassigned to ctx.Request.Body after each use.
		bodyBytes, _ := ctx.GetRawData()
		// Reassign the cached request body to ctx.Request.Body for ctx.ShouldBind to use below.
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		rw := &responsewriter.ResponseWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
		}

		param := make(map[string]interface{}, 0)

		// Convert query parameters to a map.
		for key, value := range ctx.Request.URL.Query() {
			param[key] = value
		}

		// Then, parse the form data to overwrite any query parameters with the same key.
		// This handles application/x-www-form-urlencoded.
		if err := ctx.Request.ParseForm(); err == nil {
			for key, value := range ctx.Request.PostForm {
				param[key] = value
			}
		}

		// For JSON data, unmarshal it into the map.
		// This will also overwrite any query parameters with the same key.
		if ctx.ContentType() == "application/json" {
			var jsonParams map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &jsonParams); err == nil {
				for key, value := range jsonParams {
					param[key] = value
				}
			}
		}

		// Because the request body may have been consumed,
		// the cached request body needs to be reassigned to ctx.Request.Body.
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		operParam, _ := json.Marshal(&param)

		ipInfo, err := ipaddress.GetAddress(ctx.ClientIP(), ctx.Request.UserAgent())
		if err != nil {
			ipInfo = &ipaddress.IpAddress{}
		}

		sysOperLog := dto.SaveOperLogRequest{
			Title:         title,
			BusinessType:  businessType,
			Method:        ctx.HandlerName(),
			RequestMethod: ctx.Request.Method,
			OperName:      operName,
			DeptName:      deptName,
			OperUrl:       ctx.Request.URL.Path,
			OperIp:        ipInfo.Ip,
			OperLocation:  ipInfo.Addr,
			OperParam:     string(operParam),
			JsonResult:    "",
			Status:        constant.NORMAL_STATUS,
			ErrorMsg:      "",
			OperTime:      datetime.Datetime{Time: time.Now()},
			CostTime:      0,
		}

		ctx.Writer = rw

		ctx.Next()

		sysOperLog.JsonResult = rw.Body.String()

		// Parse the response
		var body response.Response
		json.Unmarshal(rw.Body.Bytes(), &body)

		if body.Code != 200 {
			sysOperLog.Status = constant.EXCEPTION_STATUS
			sysOperLog.ErrorMsg = body.Msg
		}

		duration := time.Since(requestStartTime)
		sysOperLog.CostTime = int(duration.Milliseconds())

		operLogService.CreateSysOperLog(sysOperLog)
	}
}
