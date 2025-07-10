package service

import (
	"github.com/pkg/errors"
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"
	"mira/common/xerrors"
)

// UserServiceInterface defines operations for user management
type UserServiceInterface interface {
	CreateUser(param dto.SaveUser, roleIds, postIds []int) error
	UpdateUser(param dto.SaveUser, roleIds, postIds []int) error
	DeleteUser(userIds []int) error
	AddAuthRole(userId int, roleIds []int) error
	GetUserList(param dto.UserListRequest, userId int, isPaging bool) ([]dto.UserListResponse, int)
	GetUserByUserId(userId int) dto.UserDetailResponse
	GetUserByUsername(userName string) dto.UserTokenResponse
	GetUserByEmail(email string) dto.UserTokenResponse
	GetUserByPhonenumber(phonenumber string) dto.UserTokenResponse
	DeptListToTree(depts []dto.DeptTreeResponse, parentId int) []dto.DeptTreeResponse
	GetUserListByRoleId(param dto.RoleAuthUserAllocatedListRequest, userId int, isAllocation bool) ([]dto.UserListResponse, int)
	UserHasDeptByDeptId(deptId int) bool
	UserHasPerms(userId int, perms []string) bool
	UserHasRoles(userId int, roles []string) bool
}

// UserService implements the user management interface
type UserService struct{}

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)

// CreateUser creates a new system user with assigned roles and posts
//
// Parameters:
//   - param: User data transfer object containing all required fields
//   - roleIds: List of role IDs to assign to the user
//   - postIds: List of post IDs to assign to the user
//
// Returns:
//   - error: Any error that occurred during creation, or nil on success
func (s *UserService) CreateUser(param dto.SaveUser, roleIds, postIds []int) error {
	// Validate parameters
	if param.UserName == "" {
		return xerrors.ErrUserNameEmpty
	}
	if param.Password == "" {
		return xerrors.ErrUserPasswordEmpty
	}
	if param.NickName == "" {
		return xerrors.ErrUserNicknameEmpty
	}

	tx := dal.Gorm.Begin()

	user := model.SysUser{
		DeptId:      param.DeptId,
		UserName:    param.UserName,
		NickName:    param.NickName,
		UserType:    param.UserType,
		Email:       param.Email,
		Phonenumber: param.Phonenumber,
		Sex:         param.Sex,
		Avatar:      param.Avatar,
		Password:    param.Password,
		LoginIP:     param.LoginIP,
		LoginDate:   param.LoginDate,
		Status:      param.Status,
		CreateBy:    param.CreateBy,
		Remark:      param.Remark,
	}

	if err := tx.Model(model.SysUser{}).Create(&user).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to create user")
	}

	if len(roleIds) > 0 {
		for _, roleId := range roleIds {
			if err := tx.Model(model.SysUserRole{}).Create(&model.SysUserRole{
				UserId: user.UserId,
				RoleId: roleId,
			}).Error; err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "failed to assign role ID %d", roleId)
			}
		}
	}

	if len(postIds) > 0 {
		for _, postId := range postIds {
			if err := tx.Model(model.SysUserPost{}).Create(&model.SysUserPost{
				UserId: user.UserId,
				PostId: postId,
			}).Error; err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "failed to assign post ID %d", postId)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// UpdateUser updates an existing system user with assigned roles and posts
//
// Parameters:
//   - param: User data transfer object containing fields to update
//   - roleIds: List of role IDs to assign to the user
//   - postIds: List of post IDs to assign to the user
//
// Returns:
//   - error: Any error that occurred during update, or nil on success
func (s *UserService) UpdateUser(param dto.SaveUser, roleIds, postIds []int) error {
	if param.UserId <= 0 {
		return xerrors.ErrParam
	}

	tx := dal.Gorm.Begin()

	if err := tx.Model(model.SysUser{}).Where("user_id = ?", param.UserId).Updates(&model.SysUser{
		DeptId:      param.DeptId,
		NickName:    param.NickName,
		UserType:    param.UserType,
		Email:       param.Email,
		Phonenumber: param.Phonenumber,
		Sex:         param.Sex,
		Avatar:      param.Avatar,
		Password:    param.Password,
		LoginIP:     param.LoginIP,
		LoginDate:   param.LoginDate,
		Status:      param.Status,
		UpdateBy:    param.UpdateBy,
		Remark:      param.Remark,
	}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to update user")
	}

	if roleIds != nil {
		if err := tx.Model(model.SysUserRole{}).Where("user_id = ?", param.UserId).Delete(&model.SysUserRole{}).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "failed to delete existing user roles")
		}
	}
	if len(roleIds) > 0 {
		for _, roleId := range roleIds {
			if err := tx.Model(model.SysUserRole{}).Create(&model.SysUserRole{
				UserId: param.UserId,
				RoleId: roleId,
			}).Error; err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "failed to assign role ID %d", roleId)
			}
		}
	}

	if postIds != nil {
		if err := tx.Model(model.SysUserPost{}).Where("user_id = ?", param.UserId).Delete(&model.SysUserPost{}).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "failed to delete existing user posts")
		}
	}
	if len(postIds) > 0 {
		for _, postId := range postIds {
			if err := tx.Model(model.SysUserPost{}).Create(&model.SysUserPost{
				UserId: param.UserId,
				PostId: postId,
			}).Error; err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "failed to assign post ID %d", postId)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// DeleteUser deletes users by their IDs
//
// Parameters:
//   - userIds: Array of user IDs to delete
//
// Returns:
//   - error: Any error that occurred during deletion, or nil on success
func (s *UserService) DeleteUser(userIds []int) error {
	if len(userIds) == 0 {
		return xerrors.ErrParam
	}

	// Check for super admin deletion
	for _, userId := range userIds {
		if userId == 1 {
			return xerrors.ErrUserSuperAdminDelete
		}
	}

	tx := dal.Gorm.Begin()

	if err := tx.Model(model.SysUser{}).Where("user_id IN ?", userIds).Delete(&model.SysUser{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete users")
	}

	if err := tx.Model(model.SysUserRole{}).Where("user_id IN ?", userIds).Delete(&model.SysUserRole{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete user roles")
	}

	if err := tx.Model(model.SysUserPost{}).Where("user_id IN ?", userIds).Delete(&model.SysUserPost{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete user posts")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// AddAuthRole assigns roles to a user
//
// Parameters:
//   - userId: User ID to assign roles to
//   - roleIds: List of role IDs to assign
//
// Returns:
//   - error: Any error that occurred during role assignment, or nil on success
func (s *UserService) AddAuthRole(userId int, roleIds []int) error {
	if userId <= 0 {
		return xerrors.ErrParam
	}

	tx := dal.Gorm.Begin()

	// Clean up user roles
	if err := tx.Model(model.SysUserRole{}).Where("user_id = ?", userId).Delete(&model.SysUserRole{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete existing user roles")
	}

	// Re-insert assigned roles
	if len(roleIds) > 0 {
		for _, roleId := range roleIds {
			if err := tx.Model(model.SysUserRole{}).Create(&model.SysUserRole{
				UserId: userId,
				RoleId: roleId,
			}).Error; err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "failed to assign role ID %d", roleId)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// GetUserList gets the list of users based on query parameters
//
// Parameters:
//   - param: Request object containing query conditions
//   - userId: The ID of the currently authorized user (for data scope)
//   - isPaging: Whether pagination is needed
//
// Returns:
//   - []dto.UserListResponse: List of users
//   - int: Total record count if isPaging is true; otherwise 0
func (s *UserService) GetUserList(param dto.UserListRequest, userId int, isPaging bool) ([]dto.UserListResponse, int) {
	users, count, err := s.GetUserListWithErr(param, userId, isPaging)
	if err != nil {
		// Error already handled in the inner method
	}
	return users, count
}

// GetUserListWithErr gets the list of users based on query parameters with error reporting
//
// Parameters:
//   - param: Request object containing query conditions
//   - userId: The ID of the currently authorized user (for data scope)
//   - isPaging: Whether pagination is needed
//
// Returns:
//   - []dto.UserListResponse: List of users
//   - int: Total record count if isPaging is true; otherwise 0
//   - error: Any error that occurred during retrieval, or nil on success
func (s *UserService) GetUserListWithErr(param dto.UserListRequest, userId int, isPaging bool) ([]dto.UserListResponse, int, error) {
	var count int64
	users := make([]dto.UserListResponse, 0)

	query := dal.Gorm.Model(model.SysUser{}).
		Select("sys_user.*", "sys_dept.dept_name", "sys_dept.leader").
		Joins("LEFT JOIN sys_dept ON sys_user.dept_id = sys_dept.dept_id").
		Scopes(GetDataScope("sys_dept", userId, "sys_user"))

	if param.UserName != "" {
		query = query.Where("sys_user.user_name LIKE ?", "%"+param.UserName+"%")
	}

	if param.Phonenumber != "" {
		query = query.Where("sys_user.phonenumber LIKE ?", "%"+param.Phonenumber+"%")
	}

	if param.Status != "" {
		query = query.Where("sys_user.status = ?", param.Status)
	}

	if param.DeptId != 0 {
		query = query.Where("sys_user.dept_id = ?", param.DeptId)
	}

	if param.BeginTime != "" && param.EndTime != "" {
		query = query.Where("sys_user.create_time BETWEEN ? AND ?", param.BeginTime, param.EndTime)
	}

	if isPaging {
		if err := query.Count(&count).Error; err != nil {
			return nil, 0, errors.Wrap(err, "failed to count users")
		}
		query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, errors.Wrap(err, "failed to query users")
	}

	return users, int(count), nil
}

// GetUserByUserId gets user details by user ID
//
// Parameters:
//   - userId: User ID to look up
//
// Returns:
//   - dto.UserDetailResponse: User details, or empty object if not found
func (s *UserService) GetUserByUserId(userId int) dto.UserDetailResponse {
	user, err := s.GetUserByUserIdWithErr(userId)
	if err != nil {
		// Error already logged in the inner method
	}
	return user
}

// GetUserByUserIdWithErr gets user details by user ID with error reporting
//
// Parameters:
//   - userId: User ID to look up
//
// Returns:
//   - dto.UserDetailResponse: User details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *UserService) GetUserByUserIdWithErr(userId int) (dto.UserDetailResponse, error) {
	var user dto.UserDetailResponse

	if userId <= 0 {
		return user, errors.Errorf("invalid user ID: %d", userId)
	}

	if err := dal.Gorm.Model(model.SysUser{}).Where("user_id = ?", userId).Last(&user).Error; err != nil {
		return user, errors.Wrapf(err, "failed to get user by ID %d", userId)
	}

	return user, nil
}

// GetUserByUsername gets user details by username
//
// Parameters:
//   - userName: Username to look up
//
// Returns:
//   - dto.UserTokenResponse: User details, or empty object if not found
func (s *UserService) GetUserByUsername(userName string) dto.UserTokenResponse {
	user, err := s.GetUserByUsernameWithErr(userName)
	if err != nil {
		// Error already handled in the inner method
	}
	return user
}

// GetUserByUsernameWithErr gets user details by username with error reporting
//
// Parameters:
//   - userName: Username to look up
//
// Returns:
//   - dto.UserTokenResponse: User details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *UserService) GetUserByUsernameWithErr(userName string) (dto.UserTokenResponse, error) {
	var user dto.UserTokenResponse

	if userName == "" {
		return user, errors.New("empty username provided")
	}

	if err := dal.Gorm.Model(model.SysUser{}).
		Select(
			"sys_user.user_id",
			"sys_user.dept_id",
			"sys_user.user_name",
			"sys_user.nick_name",
			"sys_user.user_type",
			"sys_user.password",
			"sys_user.status",
			"sys_dept.dept_name",
		).
		Joins("LEFT JOIN sys_dept ON sys_user.dept_id = sys_dept.dept_id").
		Where("sys_user.user_name = ?", userName).
		Last(&user).Error; err != nil {
		return user, errors.Wrapf(err, "failed to get user by username %s", userName)
	}

	return user, nil
}

// GetUserByEmail gets user details by email
//
// Parameters:
//   - email: Email address to look up
//
// Returns:
//   - dto.UserTokenResponse: User details, or empty object if not found
func (s *UserService) GetUserByEmail(email string) dto.UserTokenResponse {
	user, err := s.GetUserByEmailWithErr(email)
	if err != nil {
		// Error already handled in the inner method
	}
	return user
}

// GetUserByEmailWithErr gets user details by email with error reporting
//
// Parameters:
//   - email: Email address to look up
//
// Returns:
//   - dto.UserTokenResponse: User details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *UserService) GetUserByEmailWithErr(email string) (dto.UserTokenResponse, error) {
	var user dto.UserTokenResponse

	if email == "" {
		return user, errors.New("empty email provided")
	}

	if err := dal.Gorm.Model(model.SysUser{}).
		Select(
			"sys_user.user_id",
			"sys_user.dept_id",
			"sys_user.user_name",
			"sys_user.nick_name",
			"sys_user.user_type",
			"sys_user.password",
			"sys_user.status",
			"sys_dept.dept_name",
		).
		Joins("LEFT JOIN sys_dept ON sys_user.dept_id = sys_dept.dept_id").
		Where("sys_user.email = ?", email).
		Last(&user).Error; err != nil {
		return user, errors.Wrapf(err, "failed to get user by email %s", email)
	}

	return user, nil
}

// GetUserByPhonenumber gets user details by phone number
//
// Parameters:
//   - phonenumber: Phone number to look up
//
// Returns:
//   - dto.UserTokenResponse: User details, or empty object if not found
func (s *UserService) GetUserByPhonenumber(phonenumber string) dto.UserTokenResponse {
	user, err := s.GetUserByPhonenumberWithErr(phonenumber)
	if err != nil {
		// Error already handled in the inner method
	}
	return user
}

// GetUserByPhonenumberWithErr gets user details by phone number with error reporting
//
// Parameters:
//   - phonenumber: Phone number to look up
//
// Returns:
//   - dto.UserTokenResponse: User details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *UserService) GetUserByPhonenumberWithErr(phonenumber string) (dto.UserTokenResponse, error) {
	var user dto.UserTokenResponse

	if phonenumber == "" {
		return user, errors.New("empty phone number provided")
	}

	if err := dal.Gorm.Model(model.SysUser{}).
		Select(
			"sys_user.user_id",
			"sys_user.dept_id",
			"sys_user.user_name",
			"sys_user.nick_name",
			"sys_user.user_type",
			"sys_user.password",
			"sys_user.status",
			"sys_dept.dept_name",
		).
		Joins("LEFT JOIN sys_dept ON sys_user.dept_id = sys_dept.dept_id").
		Where("sys_user.phonenumber = ?", phonenumber).
		Last(&user).Error; err != nil {
		return user, errors.Wrapf(err, "failed to get user by phone number %s", phonenumber)
	}

	return user, nil
}

// DeptListToTree converts a flat department list to a hierarchical tree structure
//
// Parameters:
//   - depts: List of departments
//   - parentId: Parent ID to start building tree from (typically 0 for root)
//
// Returns:
//   - []dto.DeptTreeResponse: Hierarchical tree of departments
func (s *UserService) DeptListToTree(depts []dto.DeptTreeResponse, parentId int) []dto.DeptTreeResponse {
	tree := make([]dto.DeptTreeResponse, 0)

	// Build tree structure
	for _, dept := range depts {
		if dept.ParentId == parentId {
			dept.Children = s.DeptListToTree(depts, dept.Id)
			tree = append(tree, dept)
		}
	}

	return tree
}

// GetUserListByRoleId gets users who have been assigned or not assigned to a specific role
//
// Parameters:
//   - param: Request object containing query conditions
//   - userId: The ID of the currently authorized user (for data scope)
//   - isAllocation: true for allocated users, false for unallocated users
//
// Returns:
//   - []dto.UserListResponse: List of users
//   - int: Total record count
func (s *UserService) GetUserListByRoleId(param dto.RoleAuthUserAllocatedListRequest, userId int, isAllocation bool) ([]dto.UserListResponse, int) {
	users, count, err := s.GetUserListByRoleIdWithErr(param, userId, isAllocation)
	if err != nil {
		// Error already handled in the inner method
	}
	return users, count
}

// GetUserListByRoleIdWithErr gets users who have been assigned or not assigned to a specific role with error reporting
//
// Parameters:
//   - param: Request object containing query conditions
//   - userId: The ID of the currently authorized user (for data scope)
//   - isAllocation: true for allocated users, false for unallocated users
//
// Returns:
//   - []dto.UserListResponse: List of users
//   - int: Total record count
//   - error: Any error that occurred during retrieval, or nil on success
func (s *UserService) GetUserListByRoleIdWithErr(param dto.RoleAuthUserAllocatedListRequest, userId int, isAllocation bool) ([]dto.UserListResponse, int, error) {
	var count int64
	users := make([]dto.UserListResponse, 0)

	query := dal.Gorm.Model(model.SysUser{}).
		Select("sys_user.*", "sys_dept.dept_name", "sys_dept.leader").
		Joins("LEFT JOIN sys_dept ON sys_user.dept_id = sys_dept.dept_id").
		Scopes(GetDataScope("sys_dept", userId, "sys_user"))

	if isAllocation {
		query.Joins("JOIN sys_user_role ON sys_user_role.user_id = sys_user.user_id").
			Where("sys_user_role.role_id = ?", param.RoleId)
	} else {
		query.Joins("LEFT JOIN sys_user_role ON sys_user_role.user_id = sys_user.user_id").
			Where("sys_user_role.user_id IS NULL")
	}

	if param.UserName != "" {
		query = query.Where("sys_user.user_name LIKE ?", "%"+param.UserName+"%")
	}

	if param.Phonenumber != "" {
		query = query.Where("sys_user.phonenumber LIKE ?", "%"+param.Phonenumber+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed to count users for role ID %d", param.RoleId)
	}

	query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed to query users for role ID %d", param.RoleId)
	}

	return users, int(count), nil
}

// UserHasDeptByDeptId checks if any user belongs to a specific department
//
// Parameters:
//   - deptId: Department ID to check
//
// Returns:
//   - bool: true if at least one user belongs to the department, false otherwise
func (s *UserService) UserHasDeptByDeptId(deptId int) bool {
	hasUsers, err := s.UserHasDeptByDeptIdWithErr(deptId)
	if err != nil {
		// Error already handled in the inner method
	}
	return hasUsers
}

// UserHasDeptByDeptIdWithErr checks if any user belongs to a specific department with error reporting
//
// Parameters:
//   - deptId: Department ID to check
//
// Returns:
//   - bool: true if at least one user belongs to the department, false otherwise
//   - error: Any error that occurred during the check, or nil on success
func (s *UserService) UserHasDeptByDeptIdWithErr(deptId int) (bool, error) {
	var count int64

	if deptId <= 0 {
		return false, errors.Errorf("invalid department ID: %d", deptId)
	}

	if err := dal.Gorm.Model(model.SysUser{}).Where("dept_id = ?", deptId).Count(&count).Error; err != nil {
		return false, errors.Wrapf(err, "failed to check if users exist for department ID %d", deptId)
	}

	return count > 0, nil
}

// UserHasPerms checks if a user has specific permissions
//
// Parameters:
//   - userId: User ID to check
//   - perms: List of permission strings to check
//
// Returns:
//   - bool: true if the user has at least one of the specified permissions, false otherwise
func (s *UserService) UserHasPerms(userId int, perms []string) bool {
	hasPerms, err := s.UserHasPermsWithErr(userId, perms)
	if err != nil {
		// Error already handled in the inner method
	}
	return hasPerms
}

// UserHasPermsWithErr checks if a user has specific permissions with error reporting
//
// Parameters:
//   - userId: User ID to check
//   - perms: List of permission strings to check
//
// Returns:
//   - bool: true if the user has at least one of the specified permissions, false otherwise
//   - error: Any error that occurred during the check, or nil on success
func (s *UserService) UserHasPermsWithErr(userId int, perms []string) (bool, error) {
	var count int64

	if userId <= 0 || len(perms) == 0 {
		return false, nil
	}

	if err := dal.Gorm.Model(model.SysUserRole{}).
		Joins("JOIN sys_role ON sys_user_role.role_id = sys_role.role_id AND sys_role.status = ?", constant.NORMAL_STATUS).
		Joins("JOIN sys_role_menu ON sys_role_menu.role_id = sys_role.role_id").
		Joins("JOIN sys_menu ON sys_menu.menu_id = sys_role_menu.menu_id AND sys_menu.status = ?", constant.NORMAL_STATUS).
		Where("sys_role.delete_time IS NULL AND sys_menu.delete_time IS NULL").
		Where("sys_user_role.user_id = ? AND sys_menu.perms IN ?", userId, perms).
		Count(&count).Error; err != nil {
		return false, errors.Wrapf(err, "failed to check if user ID %d has permissions %v", userId, perms)
	}

	return count > 0, nil
}

// UserHasRoles checks if a user has specific roles
//
// Parameters:
//   - userId: User ID to check
//   - roles: List of role keys to check
//
// Returns:
//   - bool: true if the user has at least one of the specified roles, false otherwise
func (s *UserService) UserHasRoles(userId int, roles []string) bool {
	hasRoles, err := s.UserHasRolesWithErr(userId, roles)
	if err != nil {
		// Error already handled in the inner method
	}
	return hasRoles
}

// UserHasRolesWithErr checks if a user has specific roles with error reporting
//
// Parameters:
//   - userId: User ID to check
//   - roles: List of role keys to check
//
// Returns:
//   - bool: true if the user has at least one of the specified roles, false otherwise
//   - error: Any error that occurred during the check, or nil on success
func (s *UserService) UserHasRolesWithErr(userId int, roles []string) (bool, error) {
	var count int64

	if userId <= 0 || len(roles) == 0 {
		return false, nil
	}

	if err := dal.Gorm.Model(model.SysUserRole{}).
		Joins("JOIN sys_role ON sys_user_role.role_id = sys_role.role_id AND sys_role.status = ?", constant.NORMAL_STATUS).
		Where("sys_role.delete_time IS NULL").
		Where("sys_user_role.user_id = ? AND sys_role.role_key IN ?", userId, roles).
		Count(&count).Error; err != nil {
		return false, errors.Wrapf(err, "failed to check if user ID %d has roles %v", userId, roles)
	}

	return count > 0, nil
}
