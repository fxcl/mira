package security

import (
	"mira/app/service"
	"mira/app/token"

	"github.com/gin-gonic/gin"
)

// Security provides methods for security checks like permission and role verification.
type Security struct {
	UserService *service.UserService
}

// NewSecurity creates a new Security instance.
func NewSecurity(userService *service.UserService) *Security {
	return &Security{UserService: userService}
}

// Get user id
func GetAuthUserId(ctx *gin.Context) int {
	authUser, err := token.GetAuthUser(ctx)
	if err != nil {
		return 0
	}
	return authUser.UserId
}

// Get department id
func GetAuthDeptId(ctx *gin.Context) int {
	authUser, err := token.GetAuthUser(ctx)
	if err != nil {
		return 0
	}
	return authUser.DeptId
}

// Get user account
func GetAuthUserName(ctx *gin.Context) string {
	authUser, err := token.GetAuthUser(ctx)
	if err != nil {
		return ""
	}
	return authUser.UserName
}

// Get user
func GetAuthUser(ctx *gin.Context) *token.UserTokenResponse {
	authUser, err := token.GetAuthUser(ctx)
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
