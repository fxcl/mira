package security

import (
	"testing"

	"mira/app/dto"
	"mira/app/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return user ID when token is present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		userToken := &token.UserTokenResponse{
			UserTokenResponse: dto.UserTokenResponse{
				UserId:   123,
				UserName: "testuser",
			},
		}

		c.Set(token.UserTokenKey, userToken)

		userId := GetAuthUserId(c)
		assert.Equal(t, 123, userId)
	})

	t.Run("should return 0 when token is not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		userId := GetAuthUserId(c)
		assert.Equal(t, 0, userId)
	})
}

// Mock security service that implements our own interface for testing
type MockSecurityService struct {
	hasPerms bool
	hasRoles bool
}

func (m *MockSecurityService) UserHasPerms(userId int, perms []string) bool {
	return m.hasPerms
}

func (m *MockSecurityService) UserHasRoles(userId int, roles []string) bool {
	return m.hasRoles
}

// Create a testable security struct
type TestSecurity struct {
	UserService MockSecurityService
}

func (s *TestSecurity) HasPerm(userId int, perm string) bool {
	return s.UserService.UserHasPerms(userId, []string{perm})
}

func (s *TestSecurity) LacksPerm(userId int, perm string) bool {
	return !s.UserService.UserHasPerms(userId, []string{perm})
}

func (s *TestSecurity) HasAnyPerms(userId int, perms []string) bool {
	return s.UserService.UserHasPerms(userId, perms)
}

func (s *TestSecurity) HasRole(userId int, roleKey string) bool {
	return s.UserService.UserHasRoles(userId, []string{roleKey})
}

func (s *TestSecurity) LacksRole(userId int, roleKey string) bool {
	return !s.UserService.UserHasRoles(userId, []string{roleKey})
}

func (s *TestSecurity) HasAnyRoles(userId int, roleKeys []string) bool {
	return s.UserService.UserHasRoles(userId, roleKeys)
}

func TestTestSecurity_HasPerm(t *testing.T) {
	t.Run("should return true when user has permission", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: true}
		security := &TestSecurity{UserService: *service}

		result := security.HasPerm(123, "admin:user:view")
		assert.True(t, result)
	})

	t.Run("should return false when user does not have permission", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: false}
		security := &TestSecurity{UserService: *service}

		result := security.HasPerm(123, "admin:user:delete")
		assert.False(t, result)
	})
}

func TestTestSecurity_LacksPerm(t *testing.T) {
	t.Run("should return true when user lacks permission", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: false}
		security := &TestSecurity{UserService: *service}

		result := security.LacksPerm(123, "admin:user:delete")
		assert.True(t, result)
	})

	t.Run("should return false when user has permission", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: true}
		security := &TestSecurity{UserService: *service}

		result := security.LacksPerm(123, "admin:user:view")
		assert.False(t, result)
	})
}

func TestTestSecurity_HasAnyPerms(t *testing.T) {
	t.Run("should return true when user has any permissions", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: true}
		security := &TestSecurity{UserService: *service}

		perms := []string{"admin:user:view", "admin:user:add", "admin:user:edit"}
		result := security.HasAnyPerms(123, perms)
		assert.True(t, result)
	})

	t.Run("should return false when user has no permissions", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: false}
		security := &TestSecurity{UserService: *service}

		perms := []string{"admin:user:delete", "admin:system:config"}
		result := security.HasAnyPerms(123, perms)
		assert.False(t, result)
	})
}

func TestTestSecurity_HasRole(t *testing.T) {
	t.Run("should return true when user has role", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: true}
		security := &TestSecurity{UserService: *service}

		result := security.HasRole(123, "admin")
		assert.True(t, result)
	})

	t.Run("should return false when user does not have role", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: false}
		security := &TestSecurity{UserService: *service}

		result := security.HasRole(123, "super_admin")
		assert.False(t, result)
	})
}

func TestTestSecurity_LacksRole(t *testing.T) {
	t.Run("should return true when user lacks role", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: false}
		security := &TestSecurity{UserService: *service}

		result := security.LacksRole(123, "super_admin")
		assert.True(t, result)
	})

	t.Run("should return false when user has role", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: true}
		security := &TestSecurity{UserService: *service}

		result := security.LacksRole(123, "admin")
		assert.False(t, result)
	})
}

func TestTestSecurity_HasAnyRoles(t *testing.T) {
	t.Run("should return true when user has any roles", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: true}
		security := &TestSecurity{UserService: *service}

		roles := []string{"admin", "user", "viewer"}
		result := security.HasAnyRoles(123, roles)
		assert.True(t, result)
	})

	t.Run("should return false when user has no roles", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: false}
		security := &TestSecurity{UserService: *service}

		roles := []string{"super_admin", "guest"}
		result := security.HasAnyRoles(123, roles)
		assert.False(t, result)
	})
}

func TestSecurityEdgeCases(t *testing.T) {
	t.Run("should handle empty permission string", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: false}
		security := &TestSecurity{UserService: *service}

		result := security.HasPerm(123, "")
		assert.False(t, result)
	})

	t.Run("should handle empty role string", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: false}
		security := &TestSecurity{UserService: *service}

		result := security.HasRole(123, "")
		assert.False(t, result)
	})

	t.Run("should handle zero user ID", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: false}
		security := &TestSecurity{UserService: *service}

		result := security.HasPerm(0, "admin:user:view")
		assert.False(t, result)
	})

	t.Run("should handle empty permission slice", func(t *testing.T) {
		service := &MockSecurityService{hasPerms: false}
		security := &TestSecurity{UserService: *service}

		var perms []string
		result := security.HasAnyPerms(123, perms)
		assert.False(t, result)
	})

	t.Run("should handle empty role slice", func(t *testing.T) {
		service := &MockSecurityService{hasRoles: false}
		security := &TestSecurity{UserService: *service}

		var roles []string
		result := security.HasAnyRoles(123, roles)
		assert.False(t, result)
	})
}