package validator

import (
	"errors"

	"mira/app/dto"
	"mira/common/utils"
	"mira/common/xerrors"
)

// CreateRoleValidator validates the request to create a role.
func CreateRoleValidator(param dto.CreateRoleRequest) error {
	switch {
	case param.RoleName == "":
		return xerrors.ErrRoleNameEmpty
	case param.RoleKey == "":
		return xerrors.ErrRoleKeyEmpty
	default:
		return nil
	}
}

// UpdateRoleValidator validates the request to update a role.
func UpdateRoleValidator(param dto.UpdateRoleRequest) error {
	switch {
	case param.RoleId <= 0:
		return xerrors.ErrParam
	case param.RoleName == "":
		return xerrors.ErrRoleNameEmpty
	case param.RoleKey == "":
		return xerrors.ErrRoleKeyEmpty
	default:
		return nil
	}
}

// RemoveRoleValidator validates the request to remove a role.
func RemoveRoleValidator(roleIds []int, roleId int, roleName string) error {
	switch {
	case utils.Contains(roleIds, 1):
		return xerrors.ErrRoleSuperAdminDelete
	case utils.Contains(roleIds, roleId):
		return errors.New("the " + roleName + " role cannot be deleted")
	default:
		return nil
	}
}

// ChangeRoleStatusValidator validates the request to change the role status.
func ChangeRoleStatusValidator(param dto.UpdateRoleRequest) error {
	switch {
	case param.RoleId <= 0:
		return xerrors.ErrParam
	case param.Status == "":
		return xerrors.ErrRoleStatusEmpty
	default:
		return nil
	}
}
