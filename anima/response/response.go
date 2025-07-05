package response

import (
	"github.com/gin-gonic/gin"
)

// Response
type Response struct {
	Status int
	Code   int
	Msg    string
	Data   map[string]interface{}
}

// NewSuccess initializes a successful response.
func NewSuccess() *Response {
	return &Response{
		Status: 200,
		Code:   200,
		Msg:    "Success",
		Data:   make(map[string]interface{}),
	}
}

// NewError initializes a failed response.
func NewError() *Response {
	return &Response{
		Status: 200,
		Code:   500,
		Msg:    "Failed",
		Data:   make(map[string]interface{}),
	}
}

// SetStatus sets the status code.
func (r *Response) SetStatus(status int) *Response {
	r.Status = status

	return r
}

// SetCode sets the response code.
func (r *Response) SetCode(code int) *Response {
	r.Code = code

	return r
}

// SetMsg sets the response message.
func (r *Response) SetMsg(msg string) *Response {
	r.Msg = msg

	return r
}

// SetData sets the response data.
func (r *Response) SetData(key string, value interface{}) *Response {
	if key == "code" || key == "msg" {
		return r
	}

	r.Data[key] = value

	return r
}

// SetPageData sets the paginated response data.
func (r *Response) SetPageData(rows interface{}, total int) *Response {
	r.Data["rows"] = rows
	r.Data["total"] = total

	return r
}

// SetDataMap sets the response data.
func (r *Response) SetDataMap(data map[string]interface{}) *Response {
	for key, value := range data {
		if key == "code" || key == "msg" {
			continue
		}

		r.Data[key] = value
	}

	return r
}

// Json serializes and returns.
func (r *Response) Json(ctx *gin.Context) {
	response := gin.H{
		"code": r.Code,
		"msg":  r.Msg,
	}

	for key, value := range r.Data {
		response[key] = value
	}

	ctx.JSON(r.Status, response)
}
