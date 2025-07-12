package service

import (
	"testing"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
)

func TestRoleService_CreateRole(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should create role successfully", func(t *testing.T) {
		// Prepare
		role := dto.SaveRole{
			RoleName: "Test Role",
			RoleKey:  "test_role",
			RoleSort: 1,
			Status:   "0",
			Remark:   "Test Remark",
			CreateBy: "test_user",
		}
		menuIds := []int{1, 2}

		// Execute
		err := s.CreateRole(role, menuIds)
		assert.NoError(t, err)

		// Verify
		var result model.SysRole
		dal.Gorm.First(&result, "role_name = ?", "Test Role")
		assert.Equal(t, "test_role", result.RoleKey)

		var roleMenus []model.SysRoleMenu
		dal.Gorm.Find(&roleMenus, "role_id = ?", result.RoleId)
		assert.Len(t, roleMenus, 2)
	})
}

func TestRoleService_UpdateRole(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should update role successfully", func(t *testing.T) {
		// Prepare
		role := model.SysRole{RoleId: 1, RoleName: "Old Role", RoleKey: "old_key"}
		dal.Gorm.Create(&role)
		update := dto.SaveRole{
			RoleId:   1,
			RoleName: "New Role",
			RoleKey:  "new_key",
		}
		menuIds := []int{3, 4}
		deptIds := []int{1, 2}

		// Execute
		err := s.UpdateRole(update, menuIds, deptIds)
		assert.NoError(t, err)

		// Verify
		var result model.SysRole
		dal.Gorm.First(&result, 1)
		assert.Equal(t, "New Role", result.RoleName)
		assert.Equal(t, "new_key", result.RoleKey)

		var roleMenus []model.SysRoleMenu
		dal.Gorm.Find(&roleMenus, "role_id = ?", 1)
		assert.Len(t, roleMenus, 2)

		var roleDepts []model.SysRoleDept
		dal.Gorm.Find(&roleDepts, "role_id = ?", 1)
		assert.Len(t, roleDepts, 2)
	})
}

func TestRoleService_DeleteRole(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should delete role successfully", func(t *testing.T) {
		// Prepare
		role := model.SysRole{RoleId: 1, RoleName: "Test Role", RoleKey: "test_key"}
		dal.Gorm.Create(&role)

		// Execute
		err := s.DeleteRole([]int{1})
		assert.NoError(t, err)

		// Verify
		var result model.SysRole
		err = dal.Gorm.First(&result, 1).Error
		assert.Error(t, err, "record not found")
	})
}

func TestRoleService_GetRoleList(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should get role list", func(t *testing.T) {
		// Prepare
		role1 := model.SysRole{RoleId: 1, RoleName: "Role 1", RoleKey: "key1"}
		role2 := model.SysRole{RoleId: 2, RoleName: "Role 2", RoleKey: "key2"}
		dal.Gorm.Create(&role1)
		dal.Gorm.Create(&role2)

		// Execute
		params := dto.RoleListRequest{
			PageRequest: dto.PageRequest{PageNum: 1, PageSize: 10},
		}
		roles, count := s.GetRoleList(params, true)

		// Verify
		assert.Equal(t, 2, count)
		assert.Len(t, roles, 2)
	})
}

func TestRoleService_GetRoleByRoleId(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should get role by role id", func(t *testing.T) {
		// Prepare
		role := model.SysRole{RoleId: 1, RoleName: "Test Role", RoleKey: "test_key"}
		dal.Gorm.Create(&role)

		// Execute
		result, err := s.GetRoleByRoleId(1)
		assert.NoError(t, err)
		assert.Equal(t, "Test Role", result.RoleName)
	})
}

func TestRoleService_GetRoleByRoleName(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should get role by role name", func(t *testing.T) {
		// Prepare
		role := model.SysRole{RoleId: 1, RoleName: "Test Role", RoleKey: "test_key"}
		dal.Gorm.Create(&role)

		// Execute
		result, err := s.GetRoleByRoleName("Test Role")
		assert.NoError(t, err)
		assert.Equal(t, "test_key", result.RoleKey)
	})
}

func TestRoleService_GetRoleByRoleKey(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should get role by role key", func(t *testing.T) {
		// Prepare
		role := model.SysRole{RoleId: 1, RoleName: "Test Role", RoleKey: "test_key"}
		dal.Gorm.Create(&role)

		// Execute
		result, err := s.GetRoleByRoleKey("test_key")
		assert.NoError(t, err)
		assert.Equal(t, "Test Role", result.RoleName)
	})
}

func TestRoleService_GetRoleKeysByUserId(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should get role keys by user id", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysRole{RoleId: 1, RoleKey: "key1", Status: "0"})
		dal.Gorm.Create(&model.SysRole{RoleId: 2, RoleKey: "key2", Status: "0"})
		dal.Gorm.Create(&model.SysUserRole{UserId: 1, RoleId: 1})
		dal.Gorm.Create(&model.SysUserRole{UserId: 1, RoleId: 2})

		// Execute
		roleKeys, err := s.GetRoleKeysByUserId(1)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"key1", "key2"}, roleKeys)
	})
}

func TestRoleService_GetRoleNamesByUserId(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should get role names by user id", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysRole{RoleId: 1, RoleName: "Role 1", Status: "0"})
		dal.Gorm.Create(&model.SysRole{RoleId: 2, RoleName: "Role 2", Status: "0"})
		dal.Gorm.Create(&model.SysUserRole{UserId: 1, RoleId: 1})
		dal.Gorm.Create(&model.SysUserRole{UserId: 1, RoleId: 2})

		// Execute
		roleNames, err := s.GetRoleNamesByUserId(1)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"Role 1", "Role 2"}, roleNames)
	})
}

func TestRoleService_AuthUserSelectAll(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should auth user select all", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysRole{RoleId: 1})
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysUser{UserId: 2})

		// Execute
		err := s.AuthUserSelectAll(1, []int{1, 2})
		assert.NoError(t, err)

		// Verify
		var userRoles []model.SysUserRole
		dal.Gorm.Find(&userRoles, "role_id = ?", 1)
		assert.Len(t, userRoles, 2)
	})
}

func TestRoleService_AuthUserDelete(t *testing.T) {
	setup()
	defer teardown()
	s := NewRoleService()

	t.Run("should auth user delete", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysRole{RoleId: 1})
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysUser{UserId: 2})
		dal.Gorm.Create(&model.SysUserRole{RoleId: 1, UserId: 1})
		dal.Gorm.Create(&model.SysUserRole{RoleId: 1, UserId: 2})

		// Execute
		err := s.AuthUserDelete(1, []int{1})
		assert.NoError(t, err)

		// Verify
		var userRoles []model.SysUserRole
		dal.Gorm.Find(&userRoles, "role_id = ?", 1)
		assert.Len(t, userRoles, 1)
		assert.Equal(t, 2, userRoles[0].UserId)
	})
}
