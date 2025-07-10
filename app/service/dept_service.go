package service

import (
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"

	"github.com/pkg/errors"
)

// DeptServiceInterface defines operations for department management
type DeptServiceInterface interface {
	CreateDept(param dto.SaveDept) error
	UpdateDept(param dto.SaveDept) error
	DeleteDept(deptId int) error
	GetDeptList(param dto.DeptListRequest, userId int) []dto.DeptListResponse
	GetDeptByDeptId(deptId int) dto.DeptDetailResponse
	GetDeptByDeptName(deptName string) dto.DeptDetailResponse
	GetUserDeptTree(userId int) []dto.DeptTreeResponse
	GetDeptIdsByRoleId(roleId int) []int
	DeptSelect() []dto.SeleteTree
	DeptSeleteToTree(depts []dto.SeleteTree, parentId int) []dto.SeleteTree
	DeptHasChildren(deptId int) bool
}

// DeptService implements the department management interface
type DeptService struct{}

// Ensure DeptService implements DeptServiceInterface
var _ DeptServiceInterface = (*DeptService)(nil)

// CreateDept creates a new department
//
// Parameters:
//   - param: Department data transfer object containing all required fields
//
// Returns:
//   - error: Any error that occurred during creation, or nil on success
func (s *DeptService) CreateDept(param dto.SaveDept) error {
	err := dal.Gorm.Model(model.SysDept{}).Create(&model.SysDept{
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
	if err != nil {
		return errors.Wrap(err, "failed to create department")
	}

	return nil
}

// UpdateDept updates an existing department
//
// Parameters:
//   - param: Department data transfer object containing fields to update
//
// Returns:
//   - error: Any error that occurred during update, or nil on success
func (s *DeptService) UpdateDept(param dto.SaveDept) error {
	err := dal.Gorm.Model(model.SysDept{}).Where("dept_id = ?", param.DeptId).Updates(&model.SysDept{
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
	if err != nil {
		return errors.Wrapf(err, "failed to update department with ID %d", param.DeptId)
	}

	return nil
}

// DeleteDept deletes a department by its ID
//
// Parameters:
//   - deptId: Department ID to delete
//
// Returns:
//   - error: Any error that occurred during deletion, or nil on success
func (s *DeptService) DeleteDept(deptId int) error {
	err := dal.Gorm.Model(model.SysDept{}).Where("dept_id = ?", deptId).Delete(&model.SysDept{}).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete department with ID %d", deptId)
	}

	return nil
}

// GetDeptList gets the list of departments based on query parameters
//
// Parameters:
//   - param: Request object containing query conditions
//   - userId: The ID of the currently authorized user (for data scope)
//
// Returns:
//   - []dto.DeptListResponse: List of departments
func (s *DeptService) GetDeptList(param dto.DeptListRequest, userId int) []dto.DeptListResponse {
	depts, err := s.GetDeptListWithErr(param, userId)
	if err != nil {
		// Error already handled in the inner method
	}
	return depts
}

// GetDeptListWithErr gets the list of departments based on query parameters with error reporting
//
// Parameters:
//   - param: Request object containing query conditions
//   - userId: The ID of the currently authorized user (for data scope)
//
// Returns:
//   - []dto.DeptListResponse: List of departments
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DeptService) GetDeptListWithErr(param dto.DeptListRequest, userId int) ([]dto.DeptListResponse, error) {
	depts := make([]dto.DeptListResponse, 0)

	query := dal.Gorm.Model(model.SysDept{}).Order("order_num, dept_id").Scopes(GetDataScope("sys_dept", userId, ""))

	if param.DeptName != "" {
		query = query.Where("dept_name LIKE ?", "%"+param.DeptName+"%")
	}

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if err := query.Find(&depts).Error; err != nil {
		return nil, errors.Wrap(err, "failed to query departments")
	}

	return depts, nil
}

// GetDeptByDeptId gets department details by department ID
//
// Parameters:
//   - deptId: Department ID to look up
//
// Returns:
//   - dto.DeptDetailResponse: Department details, or empty object if not found
func (s *DeptService) GetDeptByDeptId(deptId int) dto.DeptDetailResponse {
	dept, err := s.GetDeptByDeptIdWithErr(deptId)
	if err != nil {
		// Error already handled in the inner method
	}
	return dept
}

// GetDeptByDeptIdWithErr gets department details by department ID with error reporting
//
// Parameters:
//   - deptId: Department ID to look up
//
// Returns:
//   - dto.DeptDetailResponse: Department details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DeptService) GetDeptByDeptIdWithErr(deptId int) (dto.DeptDetailResponse, error) {
	var dept dto.DeptDetailResponse

	if deptId <= 0 {
		return dept, errors.Errorf("invalid department ID: %d", deptId)
	}

	if err := dal.Gorm.Model(model.SysDept{}).Where("dept_id = ?", deptId).Last(&dept).Error; err != nil {
		return dept, errors.Wrapf(err, "failed to get department by ID %d", deptId)
	}

	return dept, nil
}

// GetDeptByDeptName gets department details by department name
//
// Parameters:
//   - deptName: Department name to look up
//
// Returns:
//   - dto.DeptDetailResponse: Department details, or empty object if not found
func (s *DeptService) GetDeptByDeptName(deptName string) dto.DeptDetailResponse {
	dept, err := s.GetDeptByDeptNameWithErr(deptName)
	if err != nil {
		// Error already handled in the inner method
	}
	return dept
}

// GetDeptByDeptNameWithErr gets department details by department name with error reporting
//
// Parameters:
//   - deptName: Department name to look up
//
// Returns:
//   - dto.DeptDetailResponse: Department details
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DeptService) GetDeptByDeptNameWithErr(deptName string) (dto.DeptDetailResponse, error) {
	var dept dto.DeptDetailResponse

	if deptName == "" {
		return dept, errors.New("empty department name provided")
	}

	if err := dal.Gorm.Model(model.SysDept{}).Where("dept_name = ?", deptName).Last(&dept).Error; err != nil {
		return dept, errors.Wrapf(err, "failed to get department by name %s", deptName)
	}

	return dept, nil
}

// GetUserDeptTree gets the department tree structure for a user
//
// Parameters:
//   - userId: The ID of the user to get the department tree for
//
// Returns:
//   - []dto.DeptTreeResponse: List of department tree nodes
func (s *DeptService) GetUserDeptTree(userId int) []dto.DeptTreeResponse {
	depts, err := s.GetUserDeptTreeWithErr(userId)
	if err != nil {
		// Error already handled in the inner method
	}
	return depts
}

// GetUserDeptTreeWithErr gets the department tree structure for a user with error reporting
//
// Parameters:
//   - userId: The ID of the user to get the department tree for
//
// Returns:
//   - []dto.DeptTreeResponse: List of department tree nodes
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DeptService) GetUserDeptTreeWithErr(userId int) ([]dto.DeptTreeResponse, error) {
	depts := make([]dto.DeptTreeResponse, 0)

	if err := dal.Gorm.Model(model.SysDept{}).
		Select(
			"dept_id as id",
			"dept_name as label",
			"parent_id",
		).
		Order("order_num, dept_id").
		Where("status = ?", constant.NORMAL_STATUS).
		Scopes(GetDataScope("sys_dept", userId, "")).
		Find(&depts).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to get department tree for user ID %d", userId)
	}

	return depts, nil
}

// GetDeptIdsByRoleId gets department IDs associated with a role
//
// Parameters:
//   - roleId: The ID of the role to get department IDs for
//
// Returns:
//   - []int: List of department IDs
func (s *DeptService) GetDeptIdsByRoleId(roleId int) []int {
	deptIds, err := s.GetDeptIdsByRoleIdWithErr(roleId)
	if err != nil {
		// Error already handled in the inner method
	}
	return deptIds
}

// GetDeptIdsByRoleIdWithErr gets department IDs associated with a role with error reporting
//
// Parameters:
//   - roleId: The ID of the role to get department IDs for
//
// Returns:
//   - []int: List of department IDs
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DeptService) GetDeptIdsByRoleIdWithErr(roleId int) ([]int, error) {
	deptIds := make([]int, 0)

	if roleId <= 0 {
		return deptIds, errors.Errorf("invalid role ID: %d", roleId)
	}

	if err := dal.Gorm.Model(model.SysRoleDept{}).
		Joins("JOIN sys_dept ON sys_dept.dept_id = sys_role_dept.dept_id").
		Where("sys_dept.status = ? AND sys_role_dept.role_id = ?", constant.NORMAL_STATUS, roleId).
		Pluck("sys_dept.dept_id", &deptIds).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to get department IDs for role ID %d", roleId)
	}

	return deptIds, nil
}

// DeptSelect gets department dropdown list data
//
// Returns:
//   - []dto.SeleteTree: List of department select options
func (s *DeptService) DeptSelect() []dto.SeleteTree {
	depts, err := s.DeptSelectWithErr()
	if err != nil {
		// Error already handled in the inner method
	}
	return depts
}

// DeptSelectWithErr gets department dropdown list data with error reporting
//
// Returns:
//   - []dto.SeleteTree: List of department select options
//   - error: Any error that occurred during retrieval, or nil on success
func (s *DeptService) DeptSelectWithErr() ([]dto.SeleteTree, error) {
	depts := make([]dto.SeleteTree, 0)

	if err := dal.Gorm.Model(model.SysDept{}).Order("order_num, dept_id").
		Select("dept_id as id", "dept_name as label", "parent_id").
		Where("status = ?", constant.NORMAL_STATUS).
		Find(&depts).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get department select list")
	}

	return depts, nil
}

// DeptSeleteToTree converts department dropdown list to tree structure
//
// Parameters:
//   - depts: List of department select options
//   - parentId: Parent ID to start building tree from (typically 0 for root)
//
// Returns:
//   - []dto.SeleteTree: Hierarchical tree of departments
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

// DeptHasChildren checks if a department has subordinate departments
//
// Parameters:
//   - deptId: Department ID to check
//
// Returns:
//   - bool: true if the department has subordinates, false otherwise
func (s *DeptService) DeptHasChildren(deptId int) bool {
	hasChildren, err := s.DeptHasChildrenWithErr(deptId)
	if err != nil {
		// Error already handled in the inner method
	}
	return hasChildren
}

// DeptHasChildrenWithErr checks if a department has subordinate departments with error reporting
//
// Parameters:
//   - deptId: Department ID to check
//
// Returns:
//   - bool: true if the department has subordinates, false otherwise
//   - error: Any error that occurred during the check, or nil on success
func (s *DeptService) DeptHasChildrenWithErr(deptId int) (bool, error) {
	var count int64

	if deptId <= 0 {
		return false, errors.Errorf("invalid department ID: %d", deptId)
	}

	if err := dal.Gorm.Model(model.SysDept{}).Where("parent_id = ?", deptId).Count(&count).Error; err != nil {
		return false, errors.Wrapf(err, "failed to check if department ID %d has children", deptId)
	}

	return count > 0, nil
}
