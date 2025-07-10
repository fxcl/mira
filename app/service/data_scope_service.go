package service

import (
	"strings"

	"mira/app/dto"
	"mira/common/types/constant"

	"gorm.io/gorm"
)

// Data scope constants
const (
	DATA_SCOPE_ALL      = "1" // All data permissions
	DATA_SCOPE_CUSTOM   = "2" // Custom data permissions
	DATA_SCOPE_DEPT     = "3" // Department data permissions
	DATA_SCOPE_DEPT_SUB = "4" // Department and sub-department data permissions
	DATA_SCOPE_PERSONAL = "5" // Personal data only
	SUPER_ADMIN_USER_ID = 1   // Super administrator user ID
)

// DataScopeUserServiceInterface defines the minimum methods required from UserService for data scope
type DataScopeUserServiceInterface interface {
	GetUserByUserId(userId int) dto.UserDetailResponse
}

// DataScopeRoleServiceInterface defines the minimum methods required from RoleService for data scope
type DataScopeRoleServiceInterface interface {
	GetRoleListByUserIdCompat(userId int) []dto.RoleListResponse
}

// DataScopeServiceInterface defines operations for data scope management
type DataScopeServiceInterface interface {
	// GetDataScope returns a function that applies data scope filtering to database queries
	GetDataScope(deptAlias string, userId int, userAlias string) func(*gorm.DB) *gorm.DB
}

// DataScopeService implements the data scope service interface
type DataScopeService struct {
	userService DataScopeUserServiceInterface
	roleService DataScopeRoleServiceInterface
}

// NewDataScopeService creates a new data scope service with dependencies
func NewDataScopeService(userService DataScopeUserServiceInterface, roleService DataScopeRoleServiceInterface) *DataScopeService {
	return &DataScopeService{
		userService: userService,
		roleService: roleService,
	}
}

// For backward compatibility
var defaultDataScopeService *DataScopeService

func init() {
	defaultDataScopeService = &DataScopeService{
		userService: &UserService{},
		roleService: &RoleService{},
	}
}

// GetDataScope gets the data scope.
//
// Call this method in statements that require data permissions.
// To implement data permissions, the dept_id and user_id fields are required.
//
// Parameters:
//   - deptAlias: The alias for the dept table
//   - userId: The ID of the currently authorized user
//   - userAlias: The alias for the user table (optional)
//
// Returns:
//   - func(*gorm.DB) *gorm.DB: A function that applies data scope conditions to a GORM query
//
// Example: dal.Grom.Model(model.User{}).Scopes(dataScopeService.GetDataScope(deptAlias, userId, userAlias)).Find(&[]model.User{})
//
// Data scope: 1-All data permissions; 2-Custom data permissions; 3-Department data permissions;
// 4-Department and sub-department data permissions; 5-Personal data only.
func (s *DataScopeService) GetDataScope(deptAlias string, userId int, userAlias string) func(*gorm.DB) *gorm.DB {
	// Super administrators are not filtered by data permissions
	if userId == SUPER_ADMIN_USER_ID {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	// Set default department alias if not provided
	if deptAlias == "" {
		deptAlias = "sys_dept"
	}

	// Get user information
	user := s.userService.GetUserByUserId(userId)
	if user.UserId == 0 {
		// If user not found, return a scope that returns no data
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("1 = 0") // Return no data if user not found
		}
	}

	// Get the roles of the current user
	roles := s.roleService.GetRoleListByUserIdCompat(user.UserId)
	if len(roles) == 0 {
		// No logging, just continue with empty roles list
	}

	// Collect custom data scope role IDs
	roleIds := s.getCustomDataScopeRoleIds(roles)

	return func(db *gorm.DB) *gorm.DB {
		conditions, args := s.buildDataScopeConditions(roles, user, deptAlias, userAlias, roleIds)
		return s.applyConditions(db, conditions, args)
	}
}

// getCustomDataScopeRoleIds collects role IDs with custom data scope
func (s *DataScopeService) getCustomDataScopeRoleIds(roles []dto.RoleListResponse) []int {
	var roleIds []int
	for _, role := range roles {
		if role.DataScope == DATA_SCOPE_CUSTOM && role.Status == constant.NORMAL_STATUS {
			roleIds = append(roleIds, role.RoleId)
		}
	}
	return roleIds
}

// buildDataScopeConditions builds SQL conditions for data scope filtering
func (s *DataScopeService) buildDataScopeConditions(
	roles []dto.RoleListResponse,
	user dto.UserDetailResponse,
	deptAlias string,
	userAlias string,
	roleIds []int,
) ([]string, []interface{}) {
	var sqlCondition []string
	var sqlArg []interface{}

	for _, role := range roles {
		// All data permissions
		if role.DataScope == DATA_SCOPE_ALL {
			return nil, nil // No conditions needed for all data
		}

		// Add appropriate condition based on role data scope
		s.addDataScopeCondition(
			role.DataScope,
			role.RoleId,
			user,
			deptAlias,
			userAlias,
			roleIds,
			&sqlCondition,
			&sqlArg,
		)
	}

	return sqlCondition, sqlArg
}

// addDataScopeCondition adds a specific condition based on data scope type
func (s *DataScopeService) addDataScopeCondition(
	dataScope string,
	roleId int,
	user dto.UserDetailResponse,
	deptAlias string,
	userAlias string,
	roleIds []int,
	sqlCondition *[]string,
	sqlArg *[]interface{},
) {
	switch dataScope {
	case DATA_SCOPE_CUSTOM:
		s.addCustomScopeCondition(roleId, deptAlias, roleIds, sqlCondition, sqlArg)
	case DATA_SCOPE_DEPT:
		s.addDeptScopeCondition(user.DeptId, deptAlias, sqlCondition, sqlArg)
	case DATA_SCOPE_DEPT_SUB:
		s.addDeptSubScopeCondition(user.DeptId, deptAlias, sqlCondition, sqlArg)
	case DATA_SCOPE_PERSONAL:
		s.addPersonalScopeCondition(user.UserId, deptAlias, userAlias, sqlCondition, sqlArg)
	}
}

// addCustomScopeCondition adds conditions for custom data scope
func (s *DataScopeService) addCustomScopeCondition(
	roleId int,
	deptAlias string,
	roleIds []int,
	sqlCondition *[]string,
	sqlArg *[]interface{},
) {
	if len(roleIds) > 0 {
		*sqlCondition = append(*sqlCondition, deptAlias+".dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id IN (?))")
		*sqlArg = append(*sqlArg, roleIds)
	} else {
		*sqlCondition = append(*sqlCondition, deptAlias+".dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id = ?)")
		*sqlArg = append(*sqlArg, roleId)
	}
}

// addDeptScopeCondition adds conditions for department data scope
func (s *DataScopeService) addDeptScopeCondition(
	deptId int,
	deptAlias string,
	sqlCondition *[]string,
	sqlArg *[]interface{},
) {
	*sqlCondition = append(*sqlCondition, deptAlias+".dept_id = ?")
	*sqlArg = append(*sqlArg, deptId)
}

// addDeptSubScopeCondition adds conditions for department and sub-department data scope
func (s *DataScopeService) addDeptSubScopeCondition(
	deptId int,
	deptAlias string,
	sqlCondition *[]string,
	sqlArg *[]interface{},
) {
	*sqlCondition = append(*sqlCondition, deptAlias+".dept_id IN (SELECT dept_id FROM sys_dept WHERE dept_id = ? OR find_in_set(?, ancestors))")
	*sqlArg = append(*sqlArg, deptId, deptId)
}

// addPersonalScopeCondition adds conditions for personal data scope
func (s *DataScopeService) addPersonalScopeCondition(
	userId int,
	deptAlias string,
	userAlias string,
	sqlCondition *[]string,
	sqlArg *[]interface{},
) {
	if userAlias != "" {
		*sqlCondition = append(*sqlCondition, userAlias+".user_id = ?")
		*sqlArg = append(*sqlArg, userId)
	} else {
		// If data permission is for personal data only and there is no userAlias, do not query any data
		*sqlCondition = append(*sqlCondition, deptAlias+".dept_id = ?")
		*sqlArg = append(*sqlArg, 0)
	}
}

// applyConditions applies the built conditions to the database query
func (s *DataScopeService) applyConditions(db *gorm.DB, conditions []string, args []interface{}) *gorm.DB {
	if conditions == nil {
		return db // All data permissions case
	}

	if len(conditions) > 0 {
		return db.Where(strings.Join(conditions, " OR "), args...)
	}

	return db
}

// For backward compatibility, allow the function to be called directly
func GetDataScope(deptAlias string, userId int, userAlias string) func(*gorm.DB) *gorm.DB {
	return defaultDataScopeService.GetDataScope(deptAlias, userId, userAlias)
}
