package service

import (
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"
)

type DeptService struct{}

// Create department
func (s *DeptService) CreateDept(param dto.SaveDept) error {
	return dal.Gorm.Model(model.SysDept{}).Create(&model.SysDept{
		ParentId:  param.ParentId,
		Ancestors: param.Ancestors,
		DeptName:  param.DeptName,
		OrderNum:  param.OrderNum,
		Leader:    param.Leader,
		Phone:     param.Phone,
		Email:     param.Email,
		Status:    param.Status,
		CreateBy:  param.CreateBy,
	}).Error
}

// Update department
func (s *DeptService) UpdateDept(param dto.SaveDept) error {
	return dal.Gorm.Model(model.SysDept{}).Where("dept_id = ?", param.DeptId).Updates(&model.SysDept{
		ParentId:  param.ParentId,
		Ancestors: param.Ancestors,
		DeptName:  param.DeptName,
		OrderNum:  param.OrderNum,
		Leader:    param.Leader,
		Phone:     param.Phone,
		Email:     param.Email,
		Status:    param.Status,
		UpdateBy:  param.UpdateBy,
	}).Error
}

// Delete department
func (s *DeptService) DeleteDept(deptId int) error {
	return dal.Gorm.Model(model.SysDept{}).Where("dept_id = ?", deptId).Delete(&model.SysDept{}).Error
}

// Get department list
func (s *DeptService) GetDeptList(param dto.DeptListRequest, userId int) []dto.DeptListResponse {
	depts := make([]dto.DeptListResponse, 0)

	query := dal.Gorm.Model(model.SysDept{}).Order("order_num, dept_id").Scopes(GetDataScope("sys_dept", userId, ""))

	if param.DeptName != "" {
		query.Where("dept_name LIKE ?", "%"+param.DeptName+"%")
	}

	if param.Status != "" {
		query.Where("status = ?", param.Status)
	}

	query.Find(&depts)

	return depts
}

// Query department information by department ID
func (s *DeptService) GetDeptByDeptId(deptId int) dto.DeptDetailResponse {
	var dept dto.DeptDetailResponse

	dal.Gorm.Model(model.SysDept{}).Where("dept_id = ?", deptId).Last(&dept)

	return dept
}

// Query department information by department name
func (s *DeptService) GetDeptByDeptName(deptName string) dto.DeptDetailResponse {
	var dept dto.DeptDetailResponse

	dal.Gorm.Model(model.SysDept{}).Where("dept_name = ?", deptName).Last(&dept)

	return dept
}

// Get department tree
func (s *DeptService) GetUserDeptTree(userId int) []dto.DeptTreeResponse {
	depts := make([]dto.DeptTreeResponse, 0)

	dal.Gorm.Model(model.SysDept{}).
		Select(
			"dept_id as id",
			"dept_name as label",
			"parent_id",
		).
		Order("order_num, dept_id").
		Where("status = ?", constant.NORMAL_STATUS).
		Scopes(GetDataScope("sys_dept", userId, "")).
		Find(&depts)

	return depts
}

// Get department ID set by role ID
func (s *DeptService) GetDeptIdsByRoleId(roleId int) []int {
	deptIds := make([]int, 0)

	dal.Gorm.Model(model.SysRoleDept{}).
		Joins("JOIN sys_dept ON sys_dept.dept_id = sys_role_dept.dept_id").
		Where("sys_dept.status = ? AND sys_role_dept.role_id = ?", constant.NORMAL_STATUS, roleId).
		Pluck("sys_dept.dept_id", &deptIds)

	return deptIds
}

// Department dropdown tree list
func (s *DeptService) DeptSelect() []dto.SeleteTree {
	depts := make([]dto.SeleteTree, 0)

	dal.Gorm.Model(model.SysDept{}).Order("order_num, dept_id").
		Select("dept_id as id", "dept_name as label", "parent_id").
		Where("status = ?", constant.NORMAL_STATUS).
		Find(&depts)

	return depts
}

// Convert department dropdown list to tree structure
func (s *DeptService) DeptSeleteToTree(depts []dto.SeleteTree, parentId int) []dto.SeleteTree {
	tree := make([]dto.SeleteTree, 0)

	for _, dept := range depts {
		if dept.ParentId == parentId {
			tree = append(tree, dto.SeleteTree{
				Id:       dept.Id,
				Label:    dept.Label,
				ParentId: dept.ParentId,
				Children: s.DeptSeleteToTree(depts, dept.Id),
			})
		}
	}

	return tree
}

// Query whether the department has subordinates
func (s *DeptService) DeptHasChildren(deptId int) bool {
	var count int64

	dal.Gorm.Model(model.SysDept{}).Where("parent_id = ?", deptId).Count(&count)

	return count > 0
}
