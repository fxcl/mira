package responsewriter

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	writer := &ResponseWriter{
		ResponseWriter: c.Writer,
		Body:           bytes.NewBufferString(""),
	}
	c.Writer = writer

	expectedBody := "test response"
	expectedStatusCode := http.StatusOK

	c.String(expectedStatusCode, expectedBody)

	assert.Equal(t, expectedStatusCode, writer.StatusCode, "status code should match")
	assert.Equal(t, expectedBody, writer.Body.String(), "body should match")
	assert.Equal(t, expectedStatusCode, w.Code, "underlying writer status code should match")
	assert.Equal(t, expectedBody, w.Body.String(), "underlying writer body should match")
}
