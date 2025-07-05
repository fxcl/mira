package validator

import (
	"errors"
	"mira/app/dto"
	"mira/common/utils"
)

// CreateRoleValidator validates the request to create a role.
func CreateRoleValidator(param dto.CreateRoleRequest) error {
	if param.RoleName == "" {
		return errors.New("please enter the role name")
	}

	if param.RoleKey == "" {
		return errors.New("please enter the permission string")
	}

	return nil
}

// UpdateRoleValidator validates the request to update a role.
func UpdateRoleValidator(param dto.UpdateRoleRequest) error {
	if param.RoleId <= 0 {
		return errors.New("parameter error")
	}

	if param.RoleName == "" {
		return errors.New("please enter the role name")
	}

	if param.RoleKey == "" {
		return errors.New("please enter the permission string")
	}

	return nil
}

// RemoveRoleValidator validates the request to remove a role.
func RemoveRoleValidator(roleIds []int, roleId int, roleName string) error {
	if utils.Contains(roleIds, 1) {
		return errors.New("the super administrator cannot be deleted")
	}

	if utils.Contains(roleIds, roleId) {
		return errors.New("the " + roleName + " role cannot be deleted")
	}

	return nil
}

// ChangeRoleStatusValidator validates the request to change the role status.
func ChangeRoleStatusValidator(param dto.UpdateRoleRequest) error {
	if param.RoleId <= 0 {
		return errors.New("parameter error")
	}

	if param.Status == "" {
		return errors.New("please select a status")
	}

	return nil
}
