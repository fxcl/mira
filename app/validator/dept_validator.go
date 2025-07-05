package validator

import (
	"errors"
	"mira/app/dto"
)

// CreateDeptValidator validates the request to create a department.
func CreateDeptValidator(param dto.CreateDeptRequest) error {
	if param.ParentId <= 0 {
		return errors.New("please select the parent department")
	}

	if param.DeptName == "" {
		return errors.New("please enter the department name")
	}

	return nil
}

// UpdateDeptValidator validates the request to update a department.
func UpdateDeptValidator(param dto.UpdateDeptRequest) error {
	if param.DeptId <= 0 {
		return errors.New("parameter error")
	}

	if param.DeptId != 100 && param.ParentId <= 0 {
		return errors.New("please select the parent department")
	}

	if param.DeptName == "" {
		return errors.New("please enter the department name")
	}

	if param.DeptId == param.ParentId {
		return errors.New("failed to modify menu " + param.DeptName + ", the parent department cannot be itself")
	}

	return nil
}
