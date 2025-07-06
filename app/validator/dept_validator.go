package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreateDeptValidator validates the request to create a department.
func CreateDeptValidator(param dto.CreateDeptRequest) error {
	if param.ParentId <= 0 {
		return xerrors.ErrParentDeptEmpty
	}

	if param.DeptName == "" {
		return xerrors.ErrDeptNameEmpty
	}

	return nil
}

// UpdateDeptValidator validates the request to update a department.
func UpdateDeptValidator(param dto.UpdateDeptRequest) error {
	if param.DeptId <= 0 {
		return xerrors.ErrParam
	}

	if param.DeptId != 100 && param.ParentId <= 0 {
		return xerrors.ErrParentDeptEmpty
	}

	if param.DeptName == "" {
		return xerrors.ErrDeptNameEmpty
	}

	if param.DeptId == param.ParentId {
		return xerrors.ErrDeptParentSelf
	}

	return nil
}
