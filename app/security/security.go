package security

import (
	"mira/app/service"
	"mira/app/token"

	"github.com/gin-gonic/gin"
)

// Security provides methods for security checks like permission and role verification.
type Security struct {
	UserService service.UserServiceInterface
}

// SecurityInterface defines the methods for security checks.
type SecurityInterface interface {
	HasPerm(userId int, perm string) bool
}

// NewSecurity creates a new Security instance.
func NewSecurity(userService service.UserServiceInterface) *Security {
	return &Security{UserService: userService}
}

// Get user id
func GetAuthUserId(ctx *gin.Context) int {
	val, ok := ctx.Get(token.UserTokenKey)
	if !ok {
		return 0
	}
	return val.(*token.UserTokenResponse).UserId
}

// Get department id
func GetAuthDeptId(ctx *gin.Context) int {
	tokenKey, err := token.GetUserTokenKey(ctx)
	if err != nil {
		return 0
	}
	authUser, err := token.GetAuthUser(ctx.Request.Context(), tokenKey)
	if err != nil {
		return 0
	}
	return authUser.DeptId
}

// Get user account
func GetAuthUserName(ctx *gin.Context) string {
	tokenKey, err := token.GetUserTokenKey(ctx)
	if err != nil {
		return ""
	}
	authUser, err := token.GetAuthUser(ctx.Request.Context(), tokenKey)
	if err != nil {
		return ""
	}
	return authUser.UserName
}

// Get user
func GetAuthUser(ctx *gin.Context) *token.UserTokenResponse {
	tokenKey, err := token.GetUserTokenKey(ctx)
	if err != nil {
		return nil
	}
	authUser, err := token.GetAuthUser(ctx.Request.Context(), tokenKey)
	if err != nil {
		return nil
	}
	return authUser
}

// HasPerm checks if the user has a specific permission.
func (s *Security) HasPerm(userId int, perm string) bool {
	return s.UserService.UserHasPerms(userId, []string{perm})
}

// LacksPerm checks if the user does not have a specific permission.
func (s *Security) LacksPerm(userId int, perm string) bool {
	return !s.UserService.UserHasPerms(userId, []string{perm})
}

// HasAnyPerms checks if the user has any of the given permissions.
func (s *Security) HasAnyPerms(userId int, perms []string) bool {
	return s.UserService.UserHasPerms(userId, perms)
}

// HasRole checks if the user has a specific role.
func (s *Security) HasRole(userId int, roleKey string) bool {
	return s.UserService.UserHasRoles(userId, []string{roleKey})
}

// LacksRole checks if the user does not have a specific role.
func (s *Security) LacksRole(userId int, roleKey string) bool {
	return !s.UserService.UserHasRoles(userId, []string{roleKey})
}

// HasAnyRoles checks if the user has any of the given roles.
func (s *Security) HasAnyRoles(userId int, roleKeys []string) bool {
	return s.UserService.UserHasRoles(userId, roleKeys)
}
