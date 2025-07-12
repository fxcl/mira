package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreateDeptValidator validates the request to create a department.
func CreateDeptValidator(param dto.CreateDeptRequest) error {
	switch {
	case param.ParentId <= 0:
		return xerrors.ErrParentDeptEmpty
	case param.DeptName == "":
		return xerrors.ErrDeptNameEmpty
	default:
		return nil
	}
}

// UpdateDeptValidator validates the request to update a department.
func UpdateDeptValidator(param dto.UpdateDeptRequest) error {
	switch {
	case param.DeptId <= 0:
		return xerrors.ErrParam
	case param.DeptId != 100 && param.ParentId <= 0:
		return xerrors.ErrParentDeptEmpty
	case param.DeptName == "":
		return xerrors.ErrDeptNameEmpty
	case param.DeptId == param.ParentId:
		return xerrors.ErrDeptParentSelf
	default:
		return nil
	}
}
