package middleware

import (
	"net/http"
	"time"

	"mira/app/token"
	"mira/common/types/constant"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware for authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenKey, err := token.GetUserTokenKey(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "Not logged in"})
			return
		}

		authUser, err := token.GetAuthUser(ctx.Request.Context(), tokenKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "Not logged in"})
			return
		}

		// If the token is about to expire (less than 20 minutes), refresh it
		if authUser.ExpireTime.Time.Before(time.Now().Add(time.Minute * 20)) {
			token.RefreshToken(ctx.Request.Context(), tokenKey, authUser)
		}

		if authUser.Status == constant.EXCEPTION_STATUS {
			ctx.AbortWithStatusJSON(601, gin.H{"code": 601, "msg": "User is disabled"})
			return
		}

		if authUser.Status != constant.NORMAL_STATUS {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "User status is abnormal"})
			return
		}

		ctx.Next()
	}
}
