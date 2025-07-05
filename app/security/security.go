package security

import (
	"mira/app/service"
	"mira/app/token"

	"github.com/gin-gonic/gin"
)

// Get user id
//
// For example: security.GetAuthUserId(ctx)
func GetAuthUserId(ctx *gin.Context) int {
	authUser, err := token.GetAuhtUser(ctx)
	if err != nil {
		return 0
	}

	return authUser.UserId
}

// Get department id
//
// For example: security.GetAuthDeptId(ctx)
func GetAuthDeptId(ctx *gin.Context) int {
	authUser, err := token.GetAuhtUser(ctx)
	if err != nil {
		return 0
	}

	return authUser.DeptId
}

// Get user account
//
// For example: security.GetAuthUserName(ctx)
func GetAuthUserName(ctx *gin.Context) string {
	authUser, err := token.GetAuhtUser(ctx)
	if err != nil {
		return ""
	}

	return authUser.UserName
}

// Get user
//
// For example: security.GetAuthUser(ctx)
func GetAuthUser(ctx *gin.Context) *token.UserTokenResponse {
	authUser, err := token.GetAuhtUser(ctx)
	if err != nil {
		return nil
	}

	return authUser
}

// Check if the user has a certain permission, equivalent to @PreAuthorize("@ss.hasPermi('system:user:list')")
//
// For example: if HasPerm(security.GetAuthUserId(ctx), "system:user:list") { ... }
func HasPerm(userId int, perm string) bool {
	return (&service.UserService{}).UserHasPerms(userId, []string{perm})
}

// Check if the user does not have a certain permission, the logic is opposite to HasPerm, equivalent to @PreAuthorize("@ss.lacksPermi('system:user:list')")
//
// For example: if LacksPerm(security.GetAuthUserId(ctx), "system:user:list") { ... }
func LacksPerm(userId int, perm string) bool {
	return !(&service.UserService{}).UserHasPerms(userId, []string{perm})
}

// Check if the user has any of the following permissions, equivalent to @PreAuthorize("@ss.hasAnyPermi('system:user:add, system:user:edit')")
//
// For example: if HasAnyPerms(security.GetAuthUserId(ctx), []string{"system:user:add", "system:user:edit"}) { ... }
func HasAnyPerms(userId int, perms []string) bool {
	return (&service.UserService{}).UserHasPerms(userId, perms)
}

// Check if the user has a certain role, equivalent to @PreAuthorize("@ss.hasRole('user')")
//
// For example: if HasRole(security.GetAuthUserId(ctx), "user") { ... }
func HasRole(userId int, roleKey string) bool {
	return (&service.UserService{}).UserHasRoles(userId, []string{roleKey})
}

// Check if the user does not have a certain role, the logic is opposite to HasRole, equivalent to @PreAuthorize("@ss.lacksRole('user')")
//
// For example: if LacksRole(security.GetAuthUserId(ctx), "user") { ... }
func LacksRole(userId int, roleKey string) bool {
	return !(&service.UserService{}).UserHasRoles(userId, []string{roleKey})
}

// Check if the user has any of the following roles, equivalent to @PreAuthorize("@ss.hasAnyRoles('user, admin')")
//
// For example: if HasAnyRoles(security.GetAuthUserId(ctx), []string{"user", "admin"}) { ... }
func HasAnyRoles(userId int, roleKey []string) bool {
	return (&service.UserService{}).UserHasPerms(userId, roleKey)
}
