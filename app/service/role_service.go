package service

import (
	"fmt"

	"github.com/pkg/errors"
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"
)

// RoleServiceInterface defines the contract for role management operations
type RoleServiceInterface interface {
	// CreateRole creates a new role with associated menu permissions
	// param: role information
	// menuIds: menu IDs to associate with the role
	// Returns error if operation fails
	CreateRole(param dto.SaveRole, menuIds []int) error

	// UpdateRole updates an existing role with menu and department permissions
	// param: role information with updated values
	// menuIds: menu IDs to associate with the role (nil means no change)
	// deptIds: department IDs to associate with the role (nil means no change)
	// Returns error if operation fails
	UpdateRole(param dto.SaveRole, menuIds, deptIds []int) error

	// DeleteRole removes roles by IDs with their associated permissions
	// roleIds: IDs of roles to delete
	// Returns error if operation fails
	DeleteRole(roleIds []int) error

	// GetRoleList returns a list of roles based on filtering criteria
	// param: filter criteria for roles
	// isPaging: whether to apply pagination
	// Returns role list and total count
	GetRoleList(param dto.RoleListRequest, isPaging bool) ([]dto.RoleListResponse, int)

	// GetRoleByRoleId retrieves detailed role information by role ID
	// roleId: ID of the role to retrieve
	// Returns role details
	GetRoleByRoleId(roleId int) (dto.RoleDetailResponse, error)

	// AuthUserSelectAll assigns a role to multiple users
	// roleId: ID of the role to assign
	// userIds: IDs of users to assign the role to
	// Returns error if operation fails
	AuthUserSelectAll(roleId int, userIds []int) error

	// AuthUserDelete removes role assignment from multiple users
	// roleId: ID of the role to remove
	// userIds: IDs of users to remove the role from
	// Returns error if operation fails
	AuthUserDelete(roleId int, userIds []int) error

	// GetRoleListByUserId returns roles assigned to a specific user
	// userId: ID of the user
	// Returns list of roles assigned to the user and any error encountered
	GetRoleListByUserId(userId int) ([]dto.RoleListResponse, error)

	// GetRoleListByUserIdCompat is a backward compatibility method for DataScopeRoleServiceInterface
	// userId: ID of the user
	// Returns list of roles assigned to the user
	GetRoleListByUserIdCompat(userId int) []dto.RoleListResponse

	// GetRoleKeysByUserId returns role keys for a specific user
	// userId: ID of the user
	// Returns list of role keys assigned to the user
	GetRoleKeysByUserId(userId int) ([]string, error)

	// GetRoleNamesByUserId returns role names for a specific user
	// userId: ID of the user
	// Returns list of role names assigned to the user
	GetRoleNamesByUserId(userId int) ([]string, error)

	// GetRoleByRoleName retrieves role details by role name
	// roleName: name of the role to retrieve
	// Returns role details and error if not found
	GetRoleByRoleName(roleName string) (dto.RoleDetailResponse, error)

	// GetRoleByRoleKey retrieves role details by role key
	// roleKey: key of the role to retrieve
	// Returns role details and error if not found
	GetRoleByRoleKey(roleKey string) (dto.RoleDetailResponse, error)
}

// RoleService implements RoleServiceInterface for role management
type RoleService struct{}

// NewRoleService creates a new instance of RoleService
func NewRoleService() RoleServiceInterface {
	return &RoleService{}
}

// GetRoleListByUserIdCompat is a backward compatibility method for DataScopeRoleServiceInterface
func (s *RoleService) GetRoleListByUserIdCompat(userId int) []dto.RoleListResponse {
	roles, _ := s.GetRoleListByUserId(userId)
	return roles
}

// CreateRole creates a new role with associated menu permissions
func (s *RoleService) CreateRole(param dto.SaveRole, menuIds []int) error {
	// Validate input parameters
	if param.RoleName == "" {
		return errors.New("role name cannot be empty")
	}
	if param.RoleKey == "" {
		return errors.New("role key cannot be empty")
	}

	tx := dal.Gorm.Begin()

	role := model.SysRole{
		RoleName:          param.RoleName,
		RoleKey:           param.RoleKey,
		RoleSort:          param.RoleSort,
		MenuCheckStrictly: param.MenuCheckStrictly,
		DeptCheckStrictly: param.DeptCheckStrictly,
		Status:            param.Status,
		CreateBy:          param.CreateBy,
		Remark:            param.Remark,
	}

	if err := tx.Model(model.SysRole{}).Create(&role).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to create role")
	}

	if len(menuIds) > 0 {
		for _, menuId := range menuIds {
			if err := tx.Model(model.SysRoleMenu{}).Create(&model.SysRoleMenu{
				RoleId: role.RoleId,
				MenuId: menuId,
			}).Error; err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "failed to associate role with menu ID %d", menuId)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// UpdateRole updates an existing role with menu and department permissions
func (s *RoleService) UpdateRole(param dto.SaveRole, menuIds, deptIds []int) error {
	// Validate input parameters
	if param.RoleId <= 0 {
		return errors.New("invalid role ID")
	}
	if param.RoleName == "" {
		return errors.New("role name cannot be empty")
	}
	if param.RoleKey == "" {
		return errors.New("role key cannot be empty")
	}

	tx := dal.Gorm.Begin()

	if err := tx.Model(model.SysRole{}).Where("role_id = ?", param.RoleId).Updates(&model.SysRole{
		RoleName:          param.RoleName,
		RoleKey:           param.RoleKey,
		RoleSort:          param.RoleSort,
		DataScope:         param.DataScope,
		MenuCheckStrictly: param.MenuCheckStrictly,
		DeptCheckStrictly: param.DeptCheckStrictly,
		Status:            param.Status,
		UpdateBy:          param.UpdateBy,
		Remark:            param.Remark,
	}).Error; err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed to update role with ID %d", param.RoleId)
	}

	if menuIds != nil {
		if err := tx.Model(model.SysRoleMenu{}).Where("role_id = ?", param.RoleId).Delete(&model.SysRoleMenu{}).Error; err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "failed to delete menus for role ID %d", param.RoleId)
		}

		if len(menuIds) > 0 {
			for _, menuId := range menuIds {
				if err := tx.Model(model.SysRoleMenu{}).Create(&model.SysRoleMenu{
					RoleId: param.RoleId,
					MenuId: menuId,
				}).Error; err != nil {
					tx.Rollback()
					return errors.Wrapf(err, "failed to associate role ID %d with menu ID %d", param.RoleId, menuId)
				}
			}
		}
	}

	if deptIds != nil {
		if err := tx.Model(model.SysRoleDept{}).Where("role_id = ?", param.RoleId).Delete(&model.SysRoleDept{}).Error; err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "failed to delete departments for role ID %d", param.RoleId)
		}

		if len(deptIds) > 0 {
			for _, deptId := range deptIds {
				if err := tx.Model(model.SysRoleDept{}).Create(&model.SysRoleDept{
					RoleId: param.RoleId,
					DeptId: deptId,
				}).Error; err != nil {
					tx.Rollback()
					return errors.Wrapf(err, "failed to associate role ID %d with department ID %d", param.RoleId, deptId)
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// DeleteRole removes roles by IDs with their associated permissions
func (s *RoleService) DeleteRole(roleIds []int) error {
	// Validate input parameters
	if len(roleIds) == 0 {
		return errors.New("role IDs cannot be empty")
	}

	tx := dal.Gorm.Begin()

	if err := tx.Model(model.SysRole{}).Where("role_id IN ?", roleIds).Delete(&model.SysRole{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete roles")
	}

	if err := tx.Model(model.SysRoleMenu{}).Where("role_id IN ?", roleIds).Delete(&model.SysRoleMenu{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete role menus")
	}

	if err := tx.Model(model.SysRoleDept{}).Where("role_id IN ?", roleIds).Delete(&model.SysRoleDept{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete role departments")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// GetRoleList returns a list of roles based on filtering criteria
func (s *RoleService) GetRoleList(param dto.RoleListRequest, isPaging bool) ([]dto.RoleListResponse, int) {
	var count int64
	roles := make([]dto.RoleListResponse, 0)

	query := dal.Gorm.Model(model.SysRole{}).Order("role_sort, role_id")

	if param.RoleName != "" {
		query.Where("role_name LIKE ?", "%"+param.RoleName+"%")
	}

	if param.RoleKey != "" {
		query.Where("role_key LIKE ?", "%"+param.RoleKey+"%")
	}

	if param.Status != "" {
		query.Where("status = ?", param.Status)
	}

	if param.BeginTime != "" && param.EndTime != "" {
		query = query.Where("sys_user.create_time BETWEEN ? AND ?", param.BeginTime, param.EndTime)
	}

	if isPaging {
		query.Count(&count).Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	if err := query.Find(&roles).Error; err != nil {
		// Just continue execution, but error is not returned to caller
	}

	return roles, int(count)
}

// GetRoleByRoleId retrieves detailed role information by role ID
func (s *RoleService) GetRoleByRoleId(roleId int) (dto.RoleDetailResponse, error) {
	var role dto.RoleDetailResponse

	if roleId <= 0 {
		return role, errors.New("invalid role ID")
	}

	result := dal.Gorm.Model(model.SysRole{}).Where("role_id = ?", roleId).Last(&role)

	if result.Error != nil {
		return role, errors.Wrapf(result.Error, "failed to fetch role with ID %d", roleId)
	}

	if result.RowsAffected == 0 {
		return role, errors.Errorf("role with ID %d not found", roleId)
	}

	return role, nil
}

// AuthUserSelectAll assigns a role to multiple users
func (s *RoleService) AuthUserSelectAll(roleId int, userIds []int) error {
	// Validate input parameters
	if roleId <= 0 {
		return errors.New("invalid role ID")
	}
	if len(userIds) == 0 {
		return errors.New("user IDs cannot be empty")
	}

	tx := dal.Gorm.Begin()

	for _, userId := range userIds {
		if err := tx.Model(model.SysUserRole{}).Create(&model.SysUserRole{
			UserId: userId,
			RoleId: roleId,
		}).Error; err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "failed to assign role ID %d to user ID %d", roleId, userId)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// AuthUserDelete removes role assignment from multiple users
func (s *RoleService) AuthUserDelete(roleId int, userIds []int) error {
	// Validate input parameters
	if roleId <= 0 {
		return errors.New("invalid role ID")
	}
	if len(userIds) == 0 {
		return errors.New("user IDs cannot be empty")
	}

	if err := dal.Gorm.Model(model.SysUserRole{}).Where("role_id = ? AND user_id in ?", roleId, userIds).Delete(&model.SysUserRole{}).Error; err != nil {
		return errors.Wrapf(err, "failed to remove role ID %d from users", roleId)
	}

	return nil
}

// GetRoleListByUserId returns roles assigned to a specific user
func (s *RoleService) GetRoleListByUserId(userId int) ([]dto.RoleListResponse, error) {
	roles := make([]dto.RoleListResponse, 0)

	if userId <= 0 {
		return roles, errors.New("invalid user ID")
	}

	if err := dal.Gorm.Model(model.SysRole{}).Select("sys_role.*").
		Joins("JOIN sys_user_role ON sys_role.role_id = sys_user_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS).
		Find(&roles).Error; err != nil {
		return roles, errors.Wrapf(err, "failed to fetch roles for user ID %d", userId)
	}

	return roles, nil
}

// GetRoleKeysByUserId returns role keys for a specific user
func (s *RoleService) GetRoleKeysByUserId(userId int) ([]string, error) {
	roleKeys := make([]string, 0)

	if userId <= 0 {
		return roleKeys, errors.New("invalid user ID")
	}

	if err := dal.Gorm.Model(model.SysRole{}).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS).
		Pluck("sys_role.role_key", &roleKeys).Error; err != nil {
		return roleKeys, errors.Wrapf(err, "failed to fetch role keys for user ID %d", userId)
	}

	return roleKeys, nil
}

// GetRoleNamesByUserId returns role names for a specific user
func (s *RoleService) GetRoleNamesByUserId(userId int) ([]string, error) {
	var roleNames []string

	if userId <= 0 {
		return roleNames, errors.New("invalid user ID")
	}

	if err := dal.Gorm.Model(model.SysRole{}).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS).
		Pluck("sys_role.role_name", &roleNames).Error; err != nil {
		return roleNames, errors.Wrapf(err, "failed to fetch role names for user ID %d", userId)
	}

	return roleNames, nil
}

// GetRoleByRoleName retrieves role details by role name
func (s *RoleService) GetRoleByRoleName(roleName string) (dto.RoleDetailResponse, error) {
	var role dto.RoleDetailResponse

	if roleName == "" {
		return role, errors.New("role name cannot be empty")
	}

	result := dal.Gorm.Model(model.SysRole{}).Where("role_name = ?", roleName).Last(&role)

	if result.Error != nil {
		return role, errors.Wrapf(result.Error, "failed to fetch role with name %s", roleName)
	}

	if result.RowsAffected == 0 {
		return role, fmt.Errorf("role with name %s not found", roleName)
	}

	return role, nil
}

// GetRoleByRoleKey retrieves role details by role key
func (s *RoleService) GetRoleByRoleKey(roleKey string) (dto.RoleDetailResponse, error) {
	var role dto.RoleDetailResponse

	if roleKey == "" {
		return role, errors.New("role key cannot be empty")
	}

	result := dal.Gorm.Model(model.SysRole{}).Where("role_key = ?", roleKey).Last(&role)

	if result.Error != nil {
		return role, errors.Wrapf(result.Error, "failed to fetch role with key %s", roleKey)
	}

	if result.RowsAffected == 0 {
		return role, fmt.Errorf("role with key %s not found", roleKey)
	}

	return role, nil
}
