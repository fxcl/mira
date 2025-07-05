package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mira/anima/datetime"
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/service"
	ipaddress "mira/common/ip-address"
	responsewriter "mira/common/response-writer"
	"mira/common/types/constant"
	"time"

	"github.com/gin-gonic/gin"
)

// LogininforMiddleware records login information.
func LogininforMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Because the request body will be consumed after reading, to avoid EOF errors,
		// the request body needs to be cached and reassigned to ctx.Request.Body after each use.
		bodyBytes, _ := ctx.GetRawData()
		// Reassign the cached request body to ctx.Request.Body for ctx.ShouldBind to use below.
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		rw := &responsewriter.ResponseWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
		}

		var param dto.LoginRequest

		if err := ctx.ShouldBind(&param); err != nil {
			response.NewError().SetCode(400).SetMsg(err.Error()).Json(ctx)
			ctx.Abort()
			return
		}

		// Because the request body will be consumed after ctx.ShouldBind,
		// the cached request body needs to be reassigned to ctx.Request.Body.
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		ipaddress := ipaddress.GetAddress(ctx.ClientIP(), ctx.Request.UserAgent())

		logininfor := dto.SaveLogininforRequest{
			UserName:      param.Username,
			Ipaddr:        ipaddress.Ip,
			LoginLocation: ipaddress.Addr,
			Browser:       ipaddress.Browser,
			Os:            ipaddress.Os,
			Status:        constant.NORMAL_STATUS,
			LoginTime:     datetime.Datetime{Time: time.Now()},
		}

		ctx.Writer = rw

		ctx.Next()

		// Parse the response
		var body response.Response
		json.Unmarshal(rw.Body.Bytes(), &body)

		if body.Code != 200 {
			logininfor.Status = constant.EXCEPTION_STATUS
		}
		logininfor.Msg = body.Msg

		(&service.LogininforService{}).CreateSysLogininfor(logininfor)
	}
}
