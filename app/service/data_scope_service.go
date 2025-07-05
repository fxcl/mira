package service

import (
	"mira/common/types/constant"
	"strings"

	"gorm.io/gorm"
)

// GetDataScope gets the data scope.
//
// Call this method in statements that require data permissions.
// To implement data permissions, the dept_id and user_id fields are required.
//
// Required fields: deptAlias, userId, where deptAlias is the alias for the dept table, and userId is the ID of the currently authorized user.
//
// Example: dal.Grom.Model(model.User{}).Scopes(GetDataScope(deptAlias, userId, userAlias)...).Find(&[]model.User{})
//
// Data scope: 1-All data permissions; 2-Custom data permissions; 3-Department data permissions; 4-Department and sub-department data permissions; 5-Personal data only.
func GetDataScope(deptAlias string, userId int, userAlias string) func(*gorm.DB) *gorm.DB {
	// Super administrators are not filtered by data permissions
	if userId == 1 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	if deptAlias == "" {
		deptAlias = "sys_dept"
	}

	// Get user information
	user := (&UserService{}).GetUserByUserId(userId)

	var roleIds []int

	// Get the roles of the current user
	roles := (&RoleService{}).GetRoleListByUserId(user.UserId)
	for _, role := range roles {
		if role.DataScope == "2" && role.Status == constant.NORMAL_STATUS {
			roleIds = append(roleIds, role.RoleId)
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		var sqlCondition []string
		var sqlArg []interface{}

		for _, role := range roles {

			// All data permissions
			if role.DataScope == "1" {
				return db
			}

			// Custom data permissions
			if role.DataScope == "2" {
				if len(roleIds) > 0 {
					sqlCondition = append(sqlCondition, deptAlias+".dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id IN (?))")
					sqlArg = append(sqlArg, roleIds)
				} else {
					sqlCondition = append(sqlCondition, deptAlias+".dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id = ?)")
					sqlArg = append(sqlArg, role.RoleId)
				}
			}

			// Department data permissions
			if role.DataScope == "3" {
				sqlCondition = append(sqlCondition, deptAlias+".dept_id = ?")
				sqlArg = append(sqlArg, user.DeptId)
			}

			// Department and sub-department data permissions
			if role.DataScope == "4" {
				sqlCondition = append(sqlCondition, deptAlias+".dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = ? OR find_in_set(?, ancestors) )")
				sqlArg = append(sqlArg, user.DeptId, user.DeptId)
			}

			// Personal data only
			if role.DataScope == "5" {
				if userAlias != "" {
					sqlCondition = append(sqlCondition, userAlias+".user_id = ?")
					sqlArg = append(sqlArg, user.UserId)
				} else {
					// If data permission is for personal data only and there is no userAlias, do not query any data
					sqlCondition = append(sqlCondition, deptAlias+".dept_id = ?")
					sqlArg = append(sqlArg, 0)
				}
			}
		}

		if len(sqlCondition) > 0 {
			return db.Where(strings.Join(sqlCondition, " OR "), sqlArg...)
		}

		return db
	}
}
