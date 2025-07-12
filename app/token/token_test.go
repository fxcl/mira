package token

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"mira/anima/dal"
	"mira/anima/datetime"
	"mira/app/dto"
	"mira/config"

	rediskey "mira/common/types/redis-key"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := config.LoadConfig("../../application-example.yaml"); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	os.Exit(m.Run())
}

func TestGenerateToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	dal.Redis = db

	now := time.Now()
	user := dto.UserTokenResponse{
		UserId:   1,
		UserName: "test",
	}

	claims := &SysUserClaim{
		Uuid: "test-uuid",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(config.Data.Token.ExpireTime))),
			Issuer:    "mira",
		},
	}

	// Mock successful Redis SET
	expectedRedisValue := &UserTokenResponse{
		UserTokenResponse: user,
		ExpireTime:        datetime.Datetime{Time: claims.ExpiresAt.Time},
	}
	mock.ExpectSet(rediskey.UserTokenKey()+claims.Uuid, expectedRedisValue, time.Minute*time.Duration(config.Data.Token.ExpireTime)).SetVal("OK")

	_, err := GenerateToken(claims, user)
	assert.NoError(t, err)

	// Mock Redis SET failure
	mock.ExpectSet(rediskey.UserTokenKey()+claims.Uuid, expectedRedisValue, time.Minute*time.Duration(config.Data.Token.ExpireTime)).SetErr(fmt.Errorf("redis error"))
	_, err = GenerateToken(claims, user)
	assert.Error(t, err)
}

func TestGetUserTokenKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("no_authorization_header", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)

		_, err := GetUserTokenKey(ctx)
		assert.EqualError(t, err, ErrPleaseLoginFirst.Error())
	})

	t.Run("invalid_authorization_header", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", "invalid")

		_, err := GetUserTokenKey(ctx)
		assert.EqualError(t, err, ErrAuthFormat.Error())
	})

	t.Run("valid_token", func(t *testing.T) {
		config.Data.Token.Secret = "test-secret"
		claims := &SysUserClaim{
			Uuid: "test-uuid",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Data.Token.Secret))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", "Bearer "+token)

		tokenKey, err := GetUserTokenKey(ctx)
		assert.NoError(t, err)
		assert.Equal(t, rediskey.UserTokenKey()+"test-uuid", tokenKey)
	})
}

func TestRefreshToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	dal.Redis = db

	user := &UserTokenResponse{
		UserTokenResponse: dto.UserTokenResponse{
			UserId:   1,
			UserName: "test",
		},
		ExpireTime: datetime.Datetime{Time: time.Now().Add(time.Minute * time.Duration(config.Data.Token.ExpireTime))},
	}

	mock.ExpectSet("test-key", user, time.Minute*time.Duration(config.Data.Token.ExpireTime)).SetVal("OK")

	err := RefreshToken(context.Background(), "test-key", user)
	assert.NoError(t, err)
}

func TestGetAuthUser(t *testing.T) {
	db, mock := redismock.NewClientMock()
	dal.Redis = db

	user := &UserTokenResponse{
		UserTokenResponse: dto.UserTokenResponse{
			UserId:   1,
			UserName: "test",
		},
	}
	userBytes, _ := user.MarshalBinary()
	mock.ExpectGet("test-key").SetVal(string(userBytes))

	authUser, err := GetAuthUser(context.Background(), "test-key")
	assert.NoError(t, err)
	assert.Equal(t, user.UserId, authUser.UserId)
}

func TestDeleteToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	dal.Redis = db

	mock.ExpectDel("test-key").SetVal(1)

	err := DeleteToken(context.Background(), "test-key")
	assert.NoError(t, err)
}
