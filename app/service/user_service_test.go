package service

import (
	"testing"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should create user successfully", func(t *testing.T) {
		// Prepare
		user := dto.SaveUser{
			UserName: "test_user",
			Password: "password",
			NickName: "Test User",
		}
		roleIds := []int{1}
		postIds := []int{1}

		// Execute
		err := s.CreateUser(user, roleIds, postIds)
		assert.NoError(t, err)

		// Verify
		var result model.SysUser
		dal.Gorm.First(&result, "user_name = ?", "test_user")
		assert.Equal(t, "Test User", result.NickName)

		var userRoles []model.SysUserRole
		dal.Gorm.Find(&userRoles, "user_id = ?", result.UserId)
		assert.Len(t, userRoles, 1)

		var userPosts []model.SysUserPost
		dal.Gorm.Find(&userPosts, "user_id = ?", result.UserId)
		assert.Len(t, userPosts, 1)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should update user successfully", func(t *testing.T) {
		// Prepare
		user := model.SysUser{UserId: 1, UserName: "old_user", NickName: "Old User"}
		dal.Gorm.Create(&user)
		update := dto.SaveUser{
			UserId:   1,
			NickName: "New User",
		}
		roleIds := []int{2}
		postIds := []int{2}

		// Execute
		err := s.UpdateUser(update, roleIds, postIds)
		assert.NoError(t, err)

		// Verify
		var result model.SysUser
		dal.Gorm.First(&result, 1)
		assert.Equal(t, "New User", result.NickName)

		var userRoles []model.SysUserRole
		dal.Gorm.Find(&userRoles, "user_id = ?", 1)
		assert.Len(t, userRoles, 1)
		assert.Equal(t, 2, userRoles[0].RoleId)

		var userPosts []model.SysUserPost
		dal.Gorm.Find(&userPosts, "user_id = ?", 1)
		assert.Len(t, userPosts, 1)
		assert.Equal(t, 2, userPosts[0].PostId)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should delete user successfully", func(t *testing.T) {
		// Prepare
		user := model.SysUser{UserId: 2, UserName: "test_user"}
		dal.Gorm.Create(&user)

		// Execute
		err := s.DeleteUser([]int{2})
		assert.NoError(t, err)

		// Verify
		var result model.SysUser
		err = dal.Gorm.First(&result, 2).Error
		assert.Error(t, err, "record not found")
	})

	t.Run("should not delete super admin", func(t *testing.T) {
		// Execute
		err := s.DeleteUser([]int{1})
		assert.Error(t, err)
	})
}

func TestUserService_GetUserList(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should get user list", func(t *testing.T) {
		// Prepare
		user1 := model.SysUser{UserId: 2, UserName: "user1"}
		user2 := model.SysUser{UserId: 3, UserName: "user2"}
		dal.Gorm.Create(&user1)
		dal.Gorm.Create(&user2)

		// Execute
		params := dto.UserListRequest{
			PageRequest: dto.PageRequest{PageNum: 1, PageSize: 10},
		}
		users, count, err := s.GetUserListWithErr(params, 1, true)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Len(t, users, 2)
	})
}

func TestUserService_GetUserByUserId(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should get user by user id", func(t *testing.T) {
		// Prepare
		user := model.SysUser{UserId: 1, UserName: "test_user", NickName: "Test User"}
		dal.Gorm.Create(&user)

		// Execute
		result, err := s.GetUserByUserIdWithErr(1)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", result.NickName)
	})
}

func TestUserService_GetUserByUsername(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should get user by username", func(t *testing.T) {
		// Prepare
		user := model.SysUser{UserId: 1, UserName: "test_user", NickName: "Test User"}
		dal.Gorm.Create(&user)

		// Execute
		result, err := s.GetUserByUsernameWithErr("test_user")
		assert.NoError(t, err)
		assert.Equal(t, "Test User", result.NickName)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should get user by email", func(t *testing.T) {
		// Prepare
		user := model.SysUser{UserId: 1, UserName: "test_user", Email: "test@example.com"}
		dal.Gorm.Create(&user)

		// Execute
		result, err := s.GetUserByEmailWithErr("test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "test_user", result.UserName)
	})
}

func TestUserService_GetUserByPhonenumber(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should get user by phonenumber", func(t *testing.T) {
		// Prepare
		user := model.SysUser{UserId: 1, UserName: "test_user", Phonenumber: "1234567890"}
		dal.Gorm.Create(&user)

		// Execute
		result, err := s.GetUserByPhonenumberWithErr("1234567890")
		assert.NoError(t, err)
		assert.Equal(t, "test_user", result.UserName)
	})
}

func TestUserService_UserHasPerms(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should return true when user has perms", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysRole{RoleId: 1, Status: "0"})
		dal.Gorm.Create(&model.SysMenu{MenuId: 1, Perms: "test:perm"})
		dal.Gorm.Create(&model.SysUserRole{UserId: 1, RoleId: 1})
		dal.Gorm.Create(&model.SysRoleMenu{RoleId: 1, MenuId: 1})

		// Execute
		hasPerms, err := s.UserHasPermsWithErr(1, []string{"test:perm"})
		assert.NoError(t, err)
		assert.True(t, hasPerms)
	})
}

func TestUserService_UserHasRoles(t *testing.T) {
	setup()
	defer teardown()
	s := &UserService{}

	t.Run("should return true when user has roles", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysRole{RoleId: 1, RoleKey: "test_role", Status: "0"})
		dal.Gorm.Create(&model.SysUserRole{UserId: 1, RoleId: 1})

		// Execute
		hasRoles, err := s.UserHasRolesWithErr(1, []string{"test_role"})
		assert.NoError(t, err)
		assert.True(t, hasRoles)
	})
}
