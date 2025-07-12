package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mira/app/dto"
	"mira/app/service"
	"mira/app/token"
	"mira/common/types/constant"

	"mira/anima/dal"
	"mira/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var redisMock redismock.ClientMock

func setup() {
	db, mock := redismock.NewClientMock()
	redisMock = mock
	dal.Redis = db
	config.Data = &config.Config{
		Token: struct {
			Header     string `yaml:"header"`
			Secret     string `yaml:"secret"`
			ExpireTime int    `yaml:"expireTime"`
		}{
			Header:     "Authorization",
			Secret:     "your-secret-key",
			ExpireTime: 30,
		},
	}
}

func teardown() {
	dal.Redis.Close()
}

// MockOperLogService is a mock type for the OperLogService
type MockOperLogService struct {
	mock.Mock
}

var _ service.OperLogServiceInterface = (*MockOperLogService)(nil)

// DeleteOperLog is a mock method
func (m *MockOperLogService) DeleteOperLog(operIds []int) error {
	args := m.Called(operIds)
	return args.Error(0)
}

// GetOperLogList is a mock method
func (m *MockOperLogService) GetOperLogList(param dto.OperLogListRequest, isPaging bool) ([]dto.OperLogListResponse, int) {
	args := m.Called(param, isPaging)
	return args.Get(0).([]dto.OperLogListResponse), args.Int(1)
}

// CreateSysOperLog is a mock method
func (m *MockOperLogService) CreateSysOperLog(param dto.SaveOperLogRequest) error {
	args := m.Called(param)
	return args.Error(0)
}

func TestOperLogMiddleware(t *testing.T) {
	setup()
	defer teardown()
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("should log operation successfully", func(t *testing.T) {
		// Create a new mock service
		mockService := new(MockOperLogService)

		// Setup expectations
		mockService.On("CreateSysOperLog", mock.AnythingOfType("dto.SaveOperLogRequest")).Return(nil)

		// Create a new Gin engine
		r := gin.New()
		r.Use(OperLogMiddleware(mockService, "Test Title", constant.REQUEST_BUSINESS_TYPE_INSERT, func(c *gin.Context) *token.UserTokenResponse {
			return &token.UserTokenResponse{
				UserTokenResponse: dto.UserTokenResponse{
					UserId:   1,
					NickName: "testuser",
					DeptName: "testdept",
				},
			}
		}))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// Create a request to pass to our handler.
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"key":"value"}`))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to record the response.
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "CreateSysOperLog", mock.AnythingOfType("dto.SaveOperLogRequest"))

		// You can add more specific assertions on the captured argument
		// For example, check if the title is correct
		// Note: This requires a bit more setup to capture the argument.
		// For simplicity, we are just checking if the method was called.
	})

	t.Run("should log operation with error status when response code is not 200", func(t *testing.T) {
		// Create a new mock service
		mockService := new(MockOperLogService)

		// Setup expectations to capture the argument
		var capturedRequest dto.SaveOperLogRequest
		mockService.On("CreateSysOperLog", mock.AnythingOfType("dto.SaveOperLogRequest")).
			Run(func(args mock.Arguments) {
				capturedRequest = args.Get(0).(dto.SaveOperLogRequest)
			}).
			Return(nil)

		// Create a new Gin engine
		r := gin.New()
		r.Use(OperLogMiddleware(mockService, "Test Error", constant.REQUEST_BUSINESS_TYPE_OTHER, func(c *gin.Context) *token.UserTokenResponse {
			return nil
		}))
		r.POST("/test_error", func(c *gin.Context) {
			// Simulate an error response
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Bad Request"})
		})

		// Create a request
		req := httptest.NewRequest(http.MethodPost, "/test_error", nil)
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "CreateSysOperLog", mock.AnythingOfType("dto.SaveOperLogRequest"))

		// Assert on the captured request for error case
		assert.Equal(t, "Test Error", capturedRequest.Title)
		assert.Equal(t, constant.EXCEPTION_STATUS, capturedRequest.Status)
		assert.Equal(t, "Bad Request", capturedRequest.ErrorMsg)
	})

	t.Run("should prioritize form body parameters over query parameters", func(t *testing.T) {
		// Create a new mock service
		mockService := new(MockOperLogService)

		// Setup expectations to capture the argument
		var capturedRequest dto.SaveOperLogRequest
		mockService.On("CreateSysOperLog", mock.AnythingOfType("dto.SaveOperLogRequest")).
			Run(func(args mock.Arguments) {
				capturedRequest = args.Get(0).(dto.SaveOperLogRequest)
			}).
			Return(nil)

		// Create a new Gin engine
		r := gin.New()
		r.Use(OperLogMiddleware(mockService, "Test Collision", constant.REQUEST_BUSINESS_TYPE_INSERT, func(c *gin.Context) *token.UserTokenResponse {
			return &token.UserTokenResponse{
				UserTokenResponse: dto.UserTokenResponse{
					UserId:   1,
					NickName: "collosionuser",
					DeptName: "collisiondept",
				},
			}
		}))
		r.POST("/test_collision", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// Create a request with form data and a query parameter with the same key
		formData := "key=from_form"
		req := httptest.NewRequest(http.MethodPost, "/test_collision?key=from_query", bytes.NewBufferString(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Create a ResponseRecorder to record the response.
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "CreateSysOperLog", mock.AnythingOfType("dto.SaveOperLogRequest"))

		// Assert on the captured request for form data
		var loggedParams map[string]interface{}
		err := json.Unmarshal([]byte(capturedRequest.OperParam), &loggedParams)
		assert.NoError(t, err)

		// The form value should be kept, not the query value.
		assert.Equal(t, []interface{}{"from_form"}, loggedParams["key"])
	})
}
