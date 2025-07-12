package service

import (
	"testing"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// MockUserService is a mock implementation of DataScopeUserServiceInterface for testing
type MockUserService struct {
	User dto.UserDetailResponse
}

func (m *MockUserService) GetUserByUserId(userId int) dto.UserDetailResponse {
	return m.User
}

// MockRoleService is a mock implementation of DataScopeRoleServiceInterface for testing
type MockRoleService struct {
	Roles []dto.RoleListResponse
}

func (m *MockRoleService) GetRoleListByUserIdCompat(userId int) []dto.RoleListResponse {
	return m.Roles
}

func TestDataScopeService_GetDataScope(t *testing.T) {
	setup()
	defer teardown()
	t.Run("should return no-op scope for super admin", func(t *testing.T) {
		// Setup
		userService := &MockUserService{}
		roleService := &MockRoleService{}
		dataScopeService := NewDataScopeService(userService, roleService)

		// Execute
		scope := dataScopeService.GetDataScope("sys_dept", SUPER_ADMIN_USER_ID, "sys_user")
		db := scope(dal.Gorm)

		// Assert
		assert.NotNil(t, db)
		assert.Empty(t, db.Statement.Clauses["where"].Expression)
	})

	t.Run("should return no-op scope for user with all data scope", func(t *testing.T) {
		// Setup
		userService := &MockUserService{
			User: dto.UserDetailResponse{UserId: 2, DeptId: 1},
		}
		roleService := &MockRoleService{
			Roles: []dto.RoleListResponse{
				{RoleId: 1, DataScope: DATA_SCOPE_ALL},
			},
		}
		dataScopeService := NewDataScopeService(userService, roleService)

		// Execute
		scope := dataScopeService.GetDataScope("sys_dept", 2, "sys_user")
		db := scope(dal.Gorm)

		// Assert
		assert.NotNil(t, db)
		assert.Empty(t, db.Statement.Clauses["where"].Expression)
	})

	t.Run("should return dept scope", func(t *testing.T) {
		// Setup
		userService := &MockUserService{
			User: dto.UserDetailResponse{UserId: 2, DeptId: 101},
		}
		roleService := &MockRoleService{
			Roles: []dto.RoleListResponse{
				{RoleId: 2, DataScope: DATA_SCOPE_DEPT},
			},
		}
		dataScopeService := NewDataScopeService(userService, roleService)

		// Execute
		scope := dataScopeService.GetDataScope("d", 2, "u")
		db := scope(dal.Gorm.Session(&gorm.Session{DryRun: true}))
		tx := db.Find(&model.SysDept{})

		// Assert
		assert.Equal(t, "SELECT * FROM `sys_dept` WHERE d.dept_id = ? AND `sys_dept`.`delete_time` IS NULL", tx.Statement.SQL.String())
		assert.Equal(t, []interface{}{101}, tx.Statement.Vars)
	})

	t.Run("should return dept and sub-dept scope", func(t *testing.T) {
		// Setup
		userService := &MockUserService{
			User: dto.UserDetailResponse{UserId: 2, DeptId: 101},
		}
		roleService := &MockRoleService{
			Roles: []dto.RoleListResponse{
				{RoleId: 2, DataScope: DATA_SCOPE_DEPT_SUB},
			},
		}
		dataScopeService := NewDataScopeService(userService, roleService)

		// Execute
		scope := dataScopeService.GetDataScope("d", 2, "u")
		db := scope(dal.Gorm.Session(&gorm.Session{DryRun: true}))
		tx := db.Find(&model.SysDept{})

		// Assert
		assert.Equal(t, "SELECT * FROM `sys_dept` WHERE (d.dept_id IN (SELECT dept_id FROM sys_dept WHERE dept_id = ? OR find_in_set(?, ancestors))) AND `sys_dept`.`delete_time` IS NULL", tx.Statement.SQL.String())
		assert.Equal(t, []interface{}{101, 101}, tx.Statement.Vars)
	})

	t.Run("should return personal scope", func(t *testing.T) {
		// Setup
		userService := &MockUserService{
			User: dto.UserDetailResponse{UserId: 2, DeptId: 101},
		}
		roleService := &MockRoleService{
			Roles: []dto.RoleListResponse{
				{RoleId: 2, DataScope: DATA_SCOPE_PERSONAL},
			},
		}
		dataScopeService := NewDataScopeService(userService, roleService)

		// Execute
		scope := dataScopeService.GetDataScope("d", 2, "u")
		db := scope(dal.Gorm.Session(&gorm.Session{DryRun: true}))
		tx := db.Find(&model.SysUser{})

		// Assert
		assert.Equal(t, "SELECT * FROM `sys_user` WHERE u.user_id = ? AND `sys_user`.`delete_time` IS NULL", tx.Statement.SQL.String())
		assert.Equal(t, []interface{}{2}, tx.Statement.Vars)
	})

	t.Run("should return custom scope", func(t *testing.T) {
		// Setup
		userService := &MockUserService{
			User: dto.UserDetailResponse{UserId: 2, DeptId: 101},
		}
		roleService := &MockRoleService{
			Roles: []dto.RoleListResponse{
				{RoleId: 2, DataScope: DATA_SCOPE_CUSTOM, Status: "0"},
				{RoleId: 3, DataScope: DATA_SCOPE_CUSTOM, Status: "0"},
			},
		}
		dataScopeService := NewDataScopeService(userService, roleService)

		// Execute
		scope := dataScopeService.GetDataScope("d", 2, "u")
		db := scope(dal.Gorm.Session(&gorm.Session{DryRun: true}))
		tx := db.Find(&model.SysDept{})

		// Assert
		assert.Equal(t, "SELECT * FROM `sys_dept` WHERE (d.dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id IN (?,?)) OR d.dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id IN (?,?))) AND `sys_dept`.`delete_time` IS NULL", tx.Statement.SQL.String())
		assert.Equal(t, []interface{}{2, 3, 2, 3}, tx.Statement.Vars)
	})
}
