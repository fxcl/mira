package middleware

import (
	"mira/anima/response"
	"mira/app/security"
	"mira/app/token"
	"mira/common/types/constant"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware for authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authUser := security.GetAuthUser(ctx)
		if authUser == nil {
			response.NewError().SetCode(401).SetMsg("Not logged in").Json(ctx)
			ctx.Abort()
			return
		}

		// If the token is about to expire (less than 20 minutes), refresh it
		if authUser.ExpireTime.Time.Before(time.Now().Add(time.Minute * 20)) {
			token.RefreshToken(ctx, authUser.UserTokenResponse)
		}

		if authUser.Status != constant.NORMAL_STATUS {
			response.NewError().SetCode(601).SetMsg("User is disabled").Json(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
