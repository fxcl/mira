package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mira/app/dto"
	"mira/app/security"
	"mira/app/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSecurity is a mock type for the Security
type MockSecurity struct {
	mock.Mock
}

var _ security.SecurityInterface = (*MockSecurity)(nil)

// HasPerm is a mock method
func (m *MockSecurity) HasPerm(userId int, perm string) bool {
	args := m.Called(userId, perm)
	return args.Bool(0)
}

func TestHasPerm(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should allow access when user has permission", func(t *testing.T) {
		mockSecurity := new(MockSecurity)
		mockSecurity.On("HasPerm", 2, "test:perm").Return(true)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(token.UserTokenKey, &token.UserTokenResponse{
				UserTokenResponse: dto.UserTokenResponse{
					UserId: 2,
				},
			})
			c.Next()
		})
		r.Use(HasPerm(mockSecurity, "test:perm"))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSecurity.AssertExpectations(t)
	})

	t.Run("should deny access when user does not have permission", func(t *testing.T) {
		mockSecurity := new(MockSecurity)
		mockSecurity.On("HasPerm", 2, "test:perm").Return(false)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(token.UserTokenKey, &token.UserTokenResponse{
				UserTokenResponse: dto.UserTokenResponse{
					UserId: 2,
				},
			})
			c.Next()
		})
		r.Use(HasPerm(mockSecurity, "test:perm"))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSecurity.AssertExpectations(t)
	})

	t.Run("should allow access for admin user", func(t *testing.T) {
		mockSecurity := new(MockSecurity)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(token.UserTokenKey, &token.UserTokenResponse{
				UserTokenResponse: dto.UserTokenResponse{
					UserId: 1,
				},
			})
			c.Next()
		})
		r.Use(HasPerm(mockSecurity, "any:perm"))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSecurity.AssertNotCalled(t, "HasPerm", mock.Anything, mock.Anything)
	})
}
