package responsewriter

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// Rewrite gin's ResponseWriter to receive the response body
type ResponseWriter struct {
	gin.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	r.Body.Write(b)

	return r.ResponseWriter.Write(b)
}

func (r *ResponseWriter) WriteString(s string) (int, error) {
	r.Body.WriteString(s)

	return r.ResponseWriter.WriteString(s)
}
