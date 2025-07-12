package token

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"mira/anima/dal"
	"mira/anima/datetime"
	"mira/app/dto"
	"mira/common/uuid"
	"mira/config"

	rediskey "mira/common/types/redis-key"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	UserTokenKey = "user_token"
)

var (
	ErrPleaseLoginFirst      = errors.New("please log in first")
	ErrAuthFormat            = errors.New("authorization format error")
	ErrTokenFormat           = errors.New("token format error")
	ErrTokenNotValidYet      = errors.New("token not yet valid")
	ErrTokenValidationFailed = errors.New("token validation failed")
)

// SysUserClaim represents the authorization claims.
type SysUserClaim struct {
	Uuid string `json:"uuid"`
	jwt.RegisteredClaims
}

// GetClaims gets the authorization claims.
func GetClaims() *SysUserClaim {
	uuid, _ := uuid.New()

	return &SysUserClaim{
		Uuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()), // Issued at
			NotBefore: jwt.NewNumericDate(time.Now()), // Effective at
			Issuer:    "mira",                         // Issuer
		},
	}
}

// GenerateToken generates a token.
func GenerateToken(claims *SysUserClaim, user dto.UserTokenResponse) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Data.Token.Secret))
	if err != nil {
		return "", err
	}

	var expireTime time.Time
	if claims.ExpiresAt != nil {
		expireTime = claims.ExpiresAt.Time
	}
	err = dal.Redis.Set(context.Background(), rediskey.UserTokenKey()+claims.Uuid, &UserTokenResponse{
		UserTokenResponse: user,
		ExpireTime:        datetime.Datetime{Time: expireTime},
	}, time.Minute*time.Duration(config.Data.Token.ExpireTime)).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

// RefreshToken refreshes the token.
func RefreshToken(ctx context.Context, tokenKey string, user *UserTokenResponse) error {
	return dal.Redis.Set(ctx, tokenKey, user, time.Minute*time.Duration(config.Data.Token.ExpireTime)).Err()
}

// GetAuthUser parses the token.
func GetAuthUser(ctx context.Context, tokenKey string) (*UserTokenResponse, error) {
	var user UserTokenResponse
	if err := dal.Redis.Get(ctx, tokenKey).Scan(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteToken deletes the token.
func DeleteToken(ctx context.Context, tokenKey string) error {
	return dal.Redis.Del(ctx, tokenKey).Err()
}

// GetUserTokenKey gets the redis key for the authorized user.
func GetUserTokenKey(ctx *gin.Context) (string, error) {
	authorization := ctx.GetHeader(config.Data.Token.Header)
	if authorization == "" {
		return "", ErrPleaseLoginFirst
	}

	tokenSplit := strings.Split(authorization, " ")
	if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
		return "", ErrAuthFormat
	}

	token, err := jwt.ParseWithClaims(tokenSplit[1], &SysUserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Data.Token.Secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", ErrTokenFormat
			}
			if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", ErrTokenNotValidYet
			}
			return "", ErrTokenValidationFailed
		}
		return "", err
	}

	if claims, ok := token.Claims.(*SysUserClaim); ok && token.Valid {
		return rediskey.UserTokenKey() + claims.Uuid, nil
	}

	return "", ErrTokenValidationFailed
}

type UserTokenResponse struct {
	dto.UserTokenResponse
	ExpireTime datetime.Datetime `json:"expireTime"`
}

// MarshalBinary serializes dto.UserTokenResponse for redis read/write.
func (u UserTokenResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary deserializes dto.UserTokenResponse for redis read/write.
func (u *UserTokenResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
