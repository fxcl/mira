package service

import (
	"testing"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeptService_CreateDept(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should create dept successfully", func(t *testing.T) {
		param := dto.SaveDept{
			ParentId: 0,
			DeptName: "Test Dept",
			OrderNum: 1,
			Leader:   "Test Leader",
			Phone:    "1234567890",
			Email:    "test@test.com",
			Status:   "0",
			CreateBy: "tester",
		}
		err := s.CreateDept(param)
		assert.NoError(t, err)

		// Verify
		var createdDept model.SysDept
		dal.Gorm.Last(&createdDept)
		assert.Equal(t, "Test Dept", createdDept.DeptName)
	})
}

func TestDeptService_UpdateDept(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should update dept successfully", func(t *testing.T) {
		createParam := dto.SaveDept{
			ParentId: 0,
			DeptName: "Dept to Update",
		}
		s.CreateDept(createParam)
		var createdDept model.SysDept
		dal.Gorm.Last(&createdDept)

		// Execute
		updateParam := dto.SaveDept{
			DeptId:   createdDept.DeptId,
			DeptName: "Updated Dept",
		}
		err := s.UpdateDept(updateParam)
		assert.NoError(t, err)

		// Verify
		var updatedDept model.SysDept
		dal.Gorm.First(&updatedDept, createdDept.DeptId)
		assert.Equal(t, "Updated Dept", updatedDept.DeptName)
	})
}

func TestDeptService_DeleteDept(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should delete dept successfully", func(t *testing.T) {
		createParam := dto.SaveDept{
			DeptName: "Dept to Delete",
		}
		s.CreateDept(createParam)
		var createdDept model.SysDept
		dal.Gorm.Last(&createdDept)

		// Execute
		err := s.DeleteDept(createdDept.DeptId)
		assert.NoError(t, err)

		// Verify
		var deletedDept model.SysDept
		err = dal.Gorm.First(&deletedDept, createdDept.DeptId).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestDeptService_GetDeptByDeptId(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return dept successfully", func(t *testing.T) {
		createParam := dto.SaveDept{
			DeptName: "Dept By Id",
		}
		s.CreateDept(createParam)
		var createdDept model.SysDept
		dal.Gorm.Last(&createdDept)

		// Execute
		dept := s.GetDeptByDeptId(createdDept.DeptId)
		assert.Equal(t, "Dept By Id", dept.DeptName)
	})
}

func TestDeptService_GetDeptByDeptName(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return dept successfully", func(t *testing.T) {
		createParam := dto.SaveDept{
			DeptName: "Dept By Name",
		}
		s.CreateDept(createParam)

		// Execute
		dept := s.GetDeptByDeptName("Dept By Name")
		assert.Equal(t, "Dept By Name", dept.DeptName)
	})
}

func TestDeptService_GetUserDeptTree(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return dept tree successfully", func(t *testing.T) {
		s.CreateDept(dto.SaveDept{DeptId: 1, DeptName: "Parent"})
		s.CreateDept(dto.SaveDept{DeptId: 2, ParentId: 1, DeptName: "Child"})

		// Execute
		tree := s.GetUserDeptTree(1)
		assert.Len(t, tree, 2)
	})
}

func TestDeptService_GetDeptIdsByRoleId(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return dept ids successfully", func(t *testing.T) {
		dal.Gorm.Create(&model.SysDept{DeptId: 1, DeptName: "Dept 1"})
		dal.Gorm.Create(&model.SysRoleDept{RoleId: 1, DeptId: 1})

		// Execute
		deptIds := s.GetDeptIdsByRoleId(1)
		assert.Equal(t, []int{1}, deptIds)
	})
}

func TestDeptService_DeptSelect(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return dept select successfully", func(t *testing.T) {
		s.CreateDept(dto.SaveDept{DeptName: "Dept 1"})
		s.CreateDept(dto.SaveDept{DeptName: "Dept 2"})

		// Execute
		depts := s.DeptSelect()
		assert.Len(t, depts, 2)
	})
}

func TestDeptService_DeptSeleteToTree(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return dept tree successfully", func(t *testing.T) {
		depts := []dto.SeleteTree{
			{Id: 1, Label: "Parent", ParentId: 0},
			{Id: 2, Label: "Child", ParentId: 1},
		}

		// Execute
		tree := s.DeptSeleteToTree(depts, 0)
		assert.Len(t, tree, 1)
		assert.Len(t, tree[0].Children, 1)
		assert.Equal(t, "Child", tree[0].Children[0].Label)
	})
}

func TestDeptService_DeptHasChildren(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return true when dept has children", func(t *testing.T) {
		dal.Gorm.Create(&model.SysDept{DeptId: 1, DeptName: "Parent"})
		dal.Gorm.Create(&model.SysDept{DeptId: 2, ParentId: 1, DeptName: "Child"})

		// Execute
		hasChildren := s.DeptHasChildren(1)
		assert.True(t, hasChildren)
	})

	t.Run("should return false when dept has no children", func(t *testing.T) {
		dal.Gorm.Exec("DELETE FROM sys_dept")
		dal.Gorm.Create(&model.SysDept{DeptId: 1, DeptName: "Parent"})

		// Execute
		hasChildren := s.DeptHasChildren(1)
		assert.False(t, hasChildren)
	})
}

func TestDeptService_GetDeptList(t *testing.T) {
	setup()
	defer teardown()
	s := &DeptService{}

	t.Run("should return all depts", func(t *testing.T) {
		s.CreateDept(dto.SaveDept{DeptName: "Dept 1"})
		s.CreateDept(dto.SaveDept{DeptName: "Dept 2"})

		// Execute
		depts := s.GetDeptList(dto.DeptListRequest{}, 1)
		assert.Len(t, depts, 2)
	})
}
