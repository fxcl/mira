package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should handle pre-flight OPTIONS request", func(t *testing.T) {
		r := gin.New()
		r.Use(Cors())
		r.OPTIONS("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE, UPDATE", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Origin, X-Requested-With, Content-Type, Accept, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("should handle GET request", func(t *testing.T) {
		r := gin.New()
		r.Use(Cors())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}
