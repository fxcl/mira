package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mira/app/dto"
	"mira/app/service"
	"mira/common/types/constant"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogininforService is a mock type for the LogininforService
type MockLogininforService struct {
	mock.Mock
}

var _ service.LogininforServiceInterface = (*MockLogininforService)(nil)

// CreateSysLogininfor is a mock method
func (m *MockLogininforService) CreateSysLogininfor(param dto.SaveLogininforRequest) error {
	args := m.Called(param)
	return args.Error(0)
}

// GetLogininforList is a mock method
func (m *MockLogininforService) GetLogininforList(param dto.LogininforListRequest, isPaging bool) ([]dto.LogininforListResponse, int) {
	args := m.Called(param, isPaging)
	return args.Get(0).([]dto.LogininforListResponse), args.Int(1)
}

// DeleteLogininfor is a mock method
func (m *MockLogininforService) DeleteLogininfor(logininforIds []int) error {
	args := m.Called(logininforIds)
	return args.Error(0)
}

// Unlock is a mock method
func (m *MockLogininforService) Unlock(userName string) error {
	args := m.Called(userName)
	return args.Error(0)
}

func TestLogininforMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("should log login information successfully", func(t *testing.T) {
		// Create a new mock service
		mockService := new(MockLogininforService)

		// Setup expectations
		var capturedRequest dto.SaveLogininforRequest
		mockService.On("CreateSysLogininfor", mock.AnythingOfType("dto.SaveLogininforRequest")).
			Run(func(args mock.Arguments) {
				capturedRequest = args.Get(0).(dto.SaveLogininforRequest)
			}).
			Return(nil)

		// Create a new Gin engine
		r := gin.New()
		r.Use(LogininforMiddleware(mockService))
		r.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Login successful"})
		})

		// Create a request to pass to our handler.
		loginReq := dto.LoginRequest{Username: "testuser", Password: "password"}
		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to record the response.
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "CreateSysLogininfor", mock.AnythingOfType("dto.SaveLogininforRequest"))

		// Assert on the captured request
		assert.Equal(t, "testuser", capturedRequest.UserName)
		assert.Equal(t, constant.NORMAL_STATUS, capturedRequest.Status)
		assert.Equal(t, "Login successful", capturedRequest.Msg)
	})

	t.Run("should log login information with error status when login fails", func(t *testing.T) {
		// Create a new mock service
		mockService := new(MockLogininforService)

		// Setup expectations
		var capturedRequest dto.SaveLogininforRequest
		mockService.On("CreateSysLogininfor", mock.AnythingOfType("dto.SaveLogininforRequest")).
			Run(func(args mock.Arguments) {
				capturedRequest = args.Get(0).(dto.SaveLogininforRequest)
			}).
			Return(nil)

		// Create a new Gin engine
		r := gin.New()
		r.Use(LogininforMiddleware(mockService))
		r.POST("/login_fail", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid credentials"})
		})

		// Create a request to pass to our handler.
		loginReq := dto.LoginRequest{Username: "wronguser", Password: "wrongpassword"}
		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/login_fail", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to record the response.
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "CreateSysLogininfor", mock.AnythingOfType("dto.SaveLogininforRequest"))

		// Assert on the captured request
		assert.Equal(t, "wronguser", capturedRequest.UserName)
		assert.Equal(t, constant.EXCEPTION_STATUS, capturedRequest.Status)
		assert.Equal(t, "Invalid credentials", capturedRequest.Msg)
	})
}
