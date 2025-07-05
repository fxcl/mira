package service

import (
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"
)

type RoleService struct{}

// Add role
func (s *RoleService) CreateRole(param dto.SaveRole, menuIds []int) error {
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
		return err
	}

	if len(menuIds) > 0 {
		for _, menuId := range menuIds {
			if err := tx.Model(model.SysRoleMenu{}).Create(&model.SysRoleMenu{
				RoleId: role.RoleId,
				MenuId: menuId,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// Update role
func (s *RoleService) UpdateRole(param dto.SaveRole, menuIds, deptIds []int) error {
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
		return err
	}

	if menuIds != nil {
		if err := tx.Model(model.SysRoleMenu{}).Where("role_id = ?", param.RoleId).Delete(&model.SysRoleMenu{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(menuIds) > 0 {
		for _, menuId := range menuIds {
			if err := tx.Model(model.SysRoleMenu{}).Create(&model.SysRoleMenu{
				RoleId: param.RoleId,
				MenuId: menuId,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if deptIds != nil {
		if err := tx.Model(model.SysRoleDept{}).Where("role_id = ?", param.RoleId).Delete(&model.SysRoleDept{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deptIds) > 0 {
		for _, deptId := range deptIds {
			if err := tx.Model(model.SysRoleDept{}).Create(&model.SysRoleDept{
				RoleId: param.RoleId,
				DeptId: deptId,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// Delete role
func (s *RoleService) DeleteRole(roleIds []int) error {
	tx := dal.Gorm.Begin()

	if err := tx.Model(model.SysRole{}).Where("role_id IN ?", roleIds).Delete(&model.SysRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(model.SysRoleMenu{}).Where("role_id IN ?", roleIds).Delete(&model.SysRoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(model.SysRoleDept{}).Where("role_id IN ?", roleIds).Delete(&model.SysRoleDept{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Get role list
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

	query.Find(&roles)

	return roles, int(count)
}

// Get role details
func (s *RoleService) GetRoleByRoleId(roleId int) dto.RoleDetailResponse {
	var role dto.RoleDetailResponse

	dal.Gorm.Model(model.SysRole{}).Where("role_id = ?", roleId).Last(&role)

	return role
}

// Batch authorize users
func (s *RoleService) AuthUserSelectAll(roleId int, userIds []int) error {
	tx := dal.Gorm.Begin()

	for _, userId := range userIds {
		if err := tx.Model(model.SysUserRole{}).Create(&model.SysUserRole{
			UserId: userId,
			RoleId: roleId,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Batch authorize users
func (s *RoleService) AuthUserDelete(roleId int, userIds []int) error {
	return dal.Gorm.Model(model.SysUserRole{}).Where("role_id = ? AND user_id in ?", roleId, userIds).Delete(&model.SysUserRole{}).Error
}

// Query role list by user ID
func (s *RoleService) GetRoleListByUserId(userId int) []dto.RoleListResponse {
	roles := make([]dto.RoleListResponse, 0)

	dal.Gorm.Model(model.SysRole{}).Select("sys_role.*").
		Joins("JOIN sys_user_role ON sys_role.role_id = sys_user_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS).
		Find(&roles)

	return roles
}

// Query role key by user ID
func (s *RoleService) GetRoleKeysByUserId(userId int) []string {
	roleKeys := make([]string, 0)

	dal.Gorm.Model(model.SysRole{}).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS).
		Pluck("sys_role.role_key", &roleKeys)

	return roleKeys
}

// Query role name by user ID
func (s *RoleService) GetRoleNamesByUserId(userId int) []string {
	var roleNames []string

	dal.Gorm.Model(model.SysRole{}).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS).
		Pluck("sys_role.role_name", &roleNames)

	return roleNames
}

// Query role by role key
func (s *RoleService) GetRoleByRoleName(roleName string) dto.RoleDetailResponse {
	var role dto.RoleDetailResponse

	dal.Gorm.Model(model.SysRole{}).Where("role_name = ?", roleName).Last(&role)

	return role
}

// Query role by role name
func (s *RoleService) GetRoleByRoleKey(roleKey string) dto.RoleDetailResponse {
	var role dto.RoleDetailResponse

	dal.Gorm.Model(model.SysRole{}).Where("role_key = ?", roleKey).Last(&role)

	return role
}
