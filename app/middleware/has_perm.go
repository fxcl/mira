package middleware

import (
	"mira/anima/response"
	"mira/app/security"

	"github.com/gin-gonic/gin"
)

// HasPerm verifies if the user has a specific permission.
//
// This is to implement the @PreAuthorize("@ss.hasPermi('system:user:list')") annotation.
//
// Usage: api.GET("/system/user/deptTree", middleware.HasPerm("system:user:list"), (&systemcontroller.UserController{}).DeptTree)
func HasPerm(perm string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authUserId := security.GetAuthUserId(ctx)
		if authUserId == 1 {
			ctx.Next()
			return
		}

		if hasPerm := security.HasPerm(authUserId, perm); !hasPerm {
			response.NewError().SetCode(601).SetMsg("Insufficient permissions").Json(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
