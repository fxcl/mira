package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mira/anima/datetime"
	"mira/app/dto"
	"mira/app/token"
	"mira/common/types/constant"
	rediskey "mira/common/types/redis-key"
	"mira/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// Helper function to manually sign a token for testing purposes
func signTestToken(claims *token.SysUserClaim) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Data.Token.Secret))
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("No authenticated user", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		r.Use(AuthMiddleware())
		r.GET("/", func(c *gin.Context) {
			t.Error("Next handler was called unexpectedly")
		})

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"code":401,"msg":"Not logged in"}`, w.Body.String())
	})

	t.Run("User is authenticated", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		claims := token.GetClaims()
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		redisKey := rediskey.UserTokenKey() + claims.Uuid

		userResponse := &token.UserTokenResponse{
			UserTokenResponse: dto.UserTokenResponse{
				UserId:   1,
				UserName: "testuser",
				Status:   constant.NORMAL_STATUS,
			},
			ExpireTime: datetime.Datetime{Time: time.Now().Add(time.Hour)},
		}
		userResponseBytes, _ := userResponse.MarshalBinary()

		// Expect only the GET call, as we are not testing token generation
		redisMock.ExpectGet(redisKey).SetVal(string(userResponseBytes))

		// Manually sign the token to bypass GenerateToken's SET call
		authToken, err := signTestToken(claims)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authToken)

		nextCalled := false
		r.Use(AuthMiddleware())
		r.GET("/", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.True(t, nextCalled, "ctx.Next() was not called")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("User is disabled", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		claims := token.GetClaims()
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		redisKey := rediskey.UserTokenKey() + claims.Uuid

		userResponse := &token.UserTokenResponse{
			UserTokenResponse: dto.UserTokenResponse{
				UserId:   1,
				UserName: "disableduser",
				Status:   constant.EXCEPTION_STATUS,
			},
			ExpireTime: datetime.Datetime{Time: time.Now().Add(time.Hour)},
		}
		userResponseBytes, _ := userResponse.MarshalBinary()

		// Expect only the GET call
		redisMock.ExpectGet(redisKey).SetVal(string(userResponseBytes))

		// Manually sign the token
		authToken, err := signTestToken(claims)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authToken)

		r.Use(AuthMiddleware())
		r.GET("/", func(c *gin.Context) {
			t.Error("Next handler was called unexpectedly")
		})

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, 601, w.Code)
		assert.JSONEq(t, `{"code":601,"msg":"User is disabled"}`, w.Body.String())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("Token is refreshed", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		claims := token.GetClaims()
		// Use a fixed time truncated to seconds to avoid precision issues
		expireTime := time.Now().Add(10 * time.Minute).Truncate(time.Second)
		claims.ExpiresAt = jwt.NewNumericDate(expireTime)
		redisKey := rediskey.UserTokenKey() + claims.Uuid

		userResponse := &token.UserTokenResponse{
			UserTokenResponse: dto.UserTokenResponse{
				UserId:   1,
				UserName: "testuser",
				Status:   constant.NORMAL_STATUS,
			},
			ExpireTime: datetime.Datetime{Time: expireTime},
		}
		userResponseBytes, _ := userResponse.MarshalBinary()

		// Expect Redis calls
		redisMock.ExpectGet(redisKey).SetVal(string(userResponseBytes))
		// The object passed to RefreshToken will be the unmarshaled version of userResponse,
		// which should be identical now that we've truncated the time.
		redisMock.ExpectSet(redisKey, userResponse, time.Minute*time.Duration(config.Data.Token.ExpireTime)).SetVal("ok")

		// Manually sign the token
		authToken, err := signTestToken(claims)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authToken)

		nextCalled := false
		r.Use(AuthMiddleware())
		r.GET("/", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.True(t, nextCalled, "ctx.Next() was not called")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}
