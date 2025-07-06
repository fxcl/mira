package middleware

import (
	"mira/anima/response"
	"mira/app/security"

	"github.com/gin-gonic/gin"
)

// HasPerm verifies if the user has a specific permission.
func HasPerm(sec *security.Security, perm string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authUserId := security.GetAuthUserId(ctx)
		if authUserId == 1 {
			ctx.Next()
			return
		}

		if hasPerm := sec.HasPerm(authUserId, perm); !hasPerm {
			response.NewError().SetCode(601).SetMsg("Insufficient permissions").Json(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
