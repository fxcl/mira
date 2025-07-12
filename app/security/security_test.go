package security

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"mira/anima/datetime"
	"mira/app/dto"
	"mira/app/service"
	"mira/app/token"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// MockUserService is a mock implementation of the UserService for testing.
type MockUserService struct{}

func (m *MockUserService) CreateUser(param dto.SaveUser, roleIds, postIds []int) error {
	return nil
}

func (m *MockUserService) UpdateUser(param dto.SaveUser, roleIds, postIds []int) error {
	return nil
}

func (m *MockUserService) DeleteUser(userIds []int) error {
	return nil
}

func (m *MockUserService) AddAuthRole(userId int, roleIds []int) error {
	return nil
}

func (m *MockUserService) GetUserList(param dto.UserListRequest, userId int, isPaging bool) ([]dto.UserListResponse, int) {
	return nil, 0
}

func (m *MockUserService) GetUserByUserId(userId int) dto.UserDetailResponse {
	return dto.UserDetailResponse{}
}

func (m *MockUserService) GetUserByUsername(userName string) dto.UserTokenResponse {
	return dto.UserTokenResponse{}
}

func (m *MockUserService) GetUserByEmail(email string) dto.UserTokenResponse {
	return dto.UserTokenResponse{}
}

func (m *MockUserService) GetUserByPhonenumber(phonenumber string) dto.UserTokenResponse {
	return dto.UserTokenResponse{}
}

func (m *MockUserService) DeptListToTree(depts []dto.DeptTreeResponse, parentId int) []dto.DeptTreeResponse {
	return nil
}

func (m *MockUserService) GetUserListByRoleId(param dto.RoleAuthUserAllocatedListRequest, userId int, isAllocation bool) ([]dto.UserListResponse, int) {
	return nil, 0
}

func (m *MockUserService) UserHasDeptByDeptId(deptId int) bool {
	return false
}

func (m *MockUserService) UserHasPerms(userId int, perms []string) bool {
	if userId == 123 && len(perms) > 0 && perms[0] == "system:user:list" {
		return true
	}
	return false
}

func (m *MockUserService) UserHasRoles(userId int, roleKeys []string) bool {
	if userId == 123 && len(roleKeys) > 0 && roleKeys[0] == "admin" {
		return true
	}
	return false
}

var _ service.UserServiceInterface = (*MockUserService)(nil)

func TestGetAuthUserId(t *testing.T) {
	setup()
	defer teardown()

	// Create a test user and token
	testUser := &token.UserTokenResponse{
		UserTokenResponse: dto.UserTokenResponse{
			UserId:   123,
			DeptId:   1,
			UserName: "testuser",
		},
		ExpireTime: datetime.Datetime{Time: time.Now().Add(time.Hour)},
	}
	userBytes, _ := json.Marshal(testUser)

	claims := &token.SysUserClaim{
		Uuid: "test-uuid",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mira",
		},
	}
	jwtToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))

	t.Run("should return user ID when token is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

		redisMock.ExpectGet("test:user:token:test-uuid").SetVal(string(userBytes))

		userId := GetAuthUserId(ctx)
		assert.Equal(t, 123, userId)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return 0 when token is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		userId := GetAuthUserId(ctx)
		assert.Equal(t, 0, userId)
	})
}

func TestGetAuthDeptId(t *testing.T) {
	setup()
	defer teardown()

	// Create a test user and token
	testUser := &token.UserTokenResponse{
		UserTokenResponse: dto.UserTokenResponse{
			UserId:   123,
			DeptId:   456,
			UserName: "testuser",
		},
		ExpireTime: datetime.Datetime{Time: time.Now().Add(time.Hour)},
	}
	userBytes, _ := json.Marshal(testUser)

	claims := &token.SysUserClaim{
		Uuid: "test-uuid-dept",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mira",
		},
	}
	jwtToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))

	t.Run("should return dept ID when token is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

		redisMock.ExpectGet("test:user:token:test-uuid-dept").SetVal(string(userBytes))

		deptId := GetAuthDeptId(ctx)
		assert.Equal(t, 456, deptId)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return 0 when token is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		deptId := GetAuthDeptId(ctx)
		assert.Equal(t, 0, deptId)
	})
}

func TestGetAuthUserName(t *testing.T) {
	setup()
	defer teardown()

	// Create a test user and token
	testUser := &token.UserTokenResponse{
		UserTokenResponse: dto.UserTokenResponse{
			UserId:   123,
			DeptId:   456,
			UserName: "testuser",
		},
		ExpireTime: datetime.Datetime{Time: time.Now().Add(time.Hour)},
	}
	userBytes, _ := json.Marshal(testUser)

	claims := &token.SysUserClaim{
		Uuid: "test-uuid-username",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mira",
		},
	}
	jwtToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))

	t.Run("should return username when token is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

		redisMock.ExpectGet("test:user:token:test-uuid-username").SetVal(string(userBytes))

		userName := GetAuthUserName(ctx)
		assert.Equal(t, "testuser", userName)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return empty string when token is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		userName := GetAuthUserName(ctx)
		assert.Equal(t, "", userName)
	})
}

func TestGetAuthUser(t *testing.T) {
	setup()
	defer teardown()

	// Create a test user and token
	testUser := &token.UserTokenResponse{
		UserTokenResponse: dto.UserTokenResponse{
			UserId:   123,
			DeptId:   456,
			UserName: "testuser",
		},
		ExpireTime: datetime.Datetime{Time: time.Now().Add(time.Hour)},
	}
	userBytes, _ := json.Marshal(testUser)

	claims := &token.SysUserClaim{
		Uuid: "test-uuid-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mira",
		},
	}
	jwtToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))

	t.Run("should return user when token is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

		redisMock.ExpectGet("test:user:token:test-uuid-user").SetVal(string(userBytes))

		authUser := GetAuthUser(ctx)
		assert.NotNil(t, authUser)
		assert.Equal(t, 123, authUser.UserId)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("should return nil when token is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		authUser := GetAuthUser(ctx)
		assert.Nil(t, authUser)
	})
}

func TestSecurity_HasPerm(t *testing.T) {
	s := &Security{UserService: &MockUserService{}}

	t.Run("should return true when user has permission", func(t *testing.T) {
		hasPerm := s.HasPerm(123, "system:user:list")
		assert.True(t, hasPerm)
	})

	t.Run("should return false when user does not have permission", func(t *testing.T) {
		hasPerm := s.HasPerm(123, "system:user:delete")
		assert.False(t, hasPerm)
	})

	t.Run("should return false for non-existent user", func(t *testing.T) {
		hasPerm := s.HasPerm(456, "system:user:list")
		assert.False(t, hasPerm)
	})
}
