package validator

import (
	"errors"

	"mira/app/dto"
	"mira/common/utils"
	"mira/common/xerrors"
)

// CreateRoleValidator validates the request to create a role.
func CreateRoleValidator(param dto.CreateRoleRequest) error {
	if param.RoleName == "" {
		return xerrors.ErrRoleNameEmpty
	}

	if param.RoleKey == "" {
		return xerrors.ErrRoleKeyEmpty
	}

	return nil
}

// UpdateRoleValidator validates the request to update a role.
func UpdateRoleValidator(param dto.UpdateRoleRequest) error {
	if param.RoleId <= 0 {
		return xerrors.ErrParam
	}

	if param.RoleName == "" {
		return xerrors.ErrRoleNameEmpty
	}

	if param.RoleKey == "" {
		return xerrors.ErrRoleKeyEmpty
	}

	return nil
}

// RemoveRoleValidator validates the request to remove a role.
func RemoveRoleValidator(roleIds []int, roleId int, roleName string) error {
	if utils.Contains(roleIds, 1) {
		return xerrors.ErrRoleSuperAdminDelete
	}

	if utils.Contains(roleIds, roleId) {
		return errors.New("the " + roleName + " role cannot be deleted")
	}

	return nil
}

// ChangeRoleStatusValidator validates the request to change the role status.
func ChangeRoleStatusValidator(param dto.UpdateRoleRequest) error {
	if param.RoleId <= 0 {
		return xerrors.ErrParam
	}

	if param.Status == "" {
		return xerrors.ErrRoleStatusEmpty
	}

	return nil
}
