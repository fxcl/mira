package service

import (
	"testing"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
)

func TestMenuService_CreateMenu(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should create menu successfully", func(t *testing.T) {
		// Prepare
		menu := dto.SaveMenu{
			MenuName: "Test Menu",
		}

		// Execute
		err := s.CreateMenu(menu)
		assert.NoError(t, err)

		// Verify
		var result model.SysMenu
		dal.Gorm.First(&result, "menu_name = ?", "Test Menu")
		assert.Equal(t, "Test Menu", result.MenuName)
	})
}

func TestMenuService_UpdateMenu(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should update menu successfully", func(t *testing.T) {
		// Prepare
		menu := &model.SysMenu{
			MenuId:   1,
			MenuName: "Original Menu",
		}
		dal.Gorm.Create(menu)

		updatedMenu := dto.SaveMenu{
			MenuId:   1,
			MenuName: "Updated Menu",
		}

		// Execute
		err := s.UpdateMenu(updatedMenu)
		assert.NoError(t, err)

		// Verify
		var result model.SysMenu
		dal.Gorm.First(&result, 1)
		assert.Equal(t, "Updated Menu", result.MenuName)
	})
}

func TestMenuService_GetMenuList(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return all menus", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysMenu{MenuId: 1, MenuName: "Menu 1"})
		dal.Gorm.Create(&model.SysMenu{MenuId: 2, MenuName: "Menu 2"})

		// Execute
		menus := s.GetMenuList(dto.MenuListRequest{})
		assert.Len(t, menus, 2)
	})
}

func TestMenuService_DeleteMenu(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should delete menu successfully", func(t *testing.T) {
		// Prepare
		menu := &model.SysMenu{
			MenuId: 5,
		}
		dal.Gorm.Create(menu)

		// Execute
		err := s.DeleteMenu(5)
		assert.NoError(t, err)

		// Verify
		var result model.SysMenu
		err = dal.Gorm.First(&result, 5).Error
		assert.Error(t, err, "record not found")
	})
}

func TestMenuService_GetMenuByMenuId(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return menu successfully", func(t *testing.T) {
		// Prepare
		menu := &model.SysMenu{
			MenuId:   1,
			MenuName: "Test Menu",
		}
		dal.Gorm.Create(menu)

		// Execute
		result := s.GetMenuByMenuId(1)
		assert.Equal(t, "Test Menu", result.MenuName)
	})
}

func TestMenuService_GetMenuByMenuName(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return menu successfully", func(t *testing.T) {
		// Prepare
		menu := &model.SysMenu{
			MenuName: "Test Menu",
		}
		dal.Gorm.Create(menu)

		// Execute
		result := s.GetMenuByMenuName("Test Menu")
		assert.Equal(t, "Test Menu", result.MenuName)
	})
}

func TestMenuService_MenuHasChildren(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return true when menu has children", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysMenu{MenuId: 1, ParentId: 0})
		dal.Gorm.Create(&model.SysMenu{MenuId: 2, ParentId: 1})

		// Execute
		hasChildren := s.MenuHasChildren(1)
		assert.True(t, hasChildren)
	})

	t.Run("should return false when menu has no children", func(t *testing.T) {
		dal.Gorm.Exec("DELETE FROM sys_dept")
		// Prepare
		dal.Gorm.Create(&model.SysMenu{MenuId: 3, ParentId: 0})

		// Execute
		hasChildren := s.MenuHasChildren(3)
		assert.False(t, hasChildren)
	})
}

func TestMenuService_MenuExistRole(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return true when menu exist role", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysRoleMenu{MenuId: 1, RoleId: 1})

		// Execute
		exist := s.MenuExistRole(1)
		assert.True(t, exist)
	})
}

func TestMenuService_GetPermsByUserId(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return perms successfully", func(t *testing.T) {
		// Clean up any existing data
		dal.Gorm.Exec("DELETE FROM sys_menu")
		dal.Gorm.Exec("DELETE FROM sys_role")
		dal.Gorm.Exec("DELETE FROM sys_role_menu")
		dal.Gorm.Exec("DELETE FROM sys_user_role")

		// Prepare
		dal.Gorm.Create(&model.SysMenu{MenuId: 1, Perms: "test:perm"})
		dal.Gorm.Create(&model.SysRole{RoleId: 1})
		dal.Gorm.Create(&model.SysRoleMenu{RoleId: 1, MenuId: 1})
		dal.Gorm.Create(&model.SysUserRole{UserId: 2, RoleId: 1})

		// Execute
		perms := s.GetPermsByUserId(2)
		assert.Equal(t, []string{"test:perm"}, perms)
	})
}

func TestMenuService_GetMenuIdsByRoleId(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return menu ids successfully", func(t *testing.T) {
		// Prepare
		dal.Gorm.Exec("DELETE FROM sys_menu")
		dal.Gorm.Exec("DELETE FROM sys_role_menu")
		dal.Gorm.Create(&model.SysMenu{MenuId: 1, Status: "0"})
		dal.Gorm.Create(&model.SysMenu{MenuId: 2, Status: "0"})
		dal.Gorm.Create(&model.SysRoleMenu{RoleId: 1, MenuId: 1})
		dal.Gorm.Create(&model.SysRoleMenu{RoleId: 1, MenuId: 2})

		// Execute
		menuIds := s.GetMenuIdsByRoleId(1)
		assert.ElementsMatch(t, []int{1, 2}, menuIds)
	})
}

func TestMenuService_MenuSelect(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return menu select successfully", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysMenu{MenuId: 1, MenuName: "Menu 1", Status: "0"})
		dal.Gorm.Create(&model.SysMenu{MenuId: 2, MenuName: "Menu 2", Status: "0"})

		// Execute
		menus := s.MenuSelect()
		assert.Len(t, menus, 2)
	})
}

func TestMenuService_MenuSeleteToTree(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return menu tree successfully", func(t *testing.T) {
		// Prepare
		menus := []dto.SeleteTree{
			{Id: 1, Label: "Parent", ParentId: 0},
			{Id: 2, Label: "Child", ParentId: 1},
		}

		// Execute
		tree := s.MenuSeleteToTree(menus, 0)
		assert.Len(t, tree, 1)
		assert.Len(t, tree[0].Children, 1)
		assert.Equal(t, "Child", tree[0].Children[0].Label)
	})
}

func TestMenuService_GetMenuMCListByUserId(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return menu list successfully", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysMenu{MenuId: 1, MenuName: "Menu 1", MenuType: "M", Status: "0"})
		dal.Gorm.Create(&model.SysMenu{MenuId: 2, MenuName: "Menu 2", MenuType: "C", Status: "0"})
		dal.Gorm.Create(&model.SysMenu{MenuId: 3, MenuName: "Menu 3", MenuType: "F", Status: "0"})
		dal.Gorm.Create(&model.SysRole{RoleId: 1, Status: "0"})
		dal.Gorm.Create(&model.SysRoleMenu{RoleId: 1, MenuId: 1})
		dal.Gorm.Create(&model.SysRoleMenu{RoleId: 1, MenuId: 2})
		dal.Gorm.Create(&model.SysUserRole{UserId: 2, RoleId: 1})

		// Execute
		menus := s.GetMenuMCListByUserId(2)
		assert.Len(t, menus, 2)
	})
}

func TestMenuService_MenusToTree(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should return menu tree successfully", func(t *testing.T) {
		// Prepare
		menus := []dto.MenuListResponse{
			{MenuId: 1, MenuName: "Parent", ParentId: 0},
			{MenuId: 2, MenuName: "Child", ParentId: 1},
		}

		// Execute
		tree := s.MenusToTree(menus, 0)
		assert.Len(t, tree, 1)
		assert.Len(t, tree[0].Children, 1)
		assert.Equal(t, "Child", tree[0].Children[0].MenuName)
	})
}

func TestMenuService_BuildRouterMenus(t *testing.T) {
	setup()
	defer teardown()
	s := NewMenuService()

	t.Run("should build router menus successfully", func(t *testing.T) {
		// Prepare
		menus := []dto.MenuListTreeResponse{
			{
				MenuListResponse: dto.MenuListResponse{MenuId: 1, MenuName: "Parent", MenuType: "M", Path: "/parent"},
				Children: []dto.MenuListTreeResponse{
					{
						MenuListResponse: dto.MenuListResponse{MenuId: 2, MenuName: "Child", MenuType: "C", Path: "child", Component: "ChildComponent"},
					},
				},
			},
		}

		// Execute
		routers := s.BuildRouterMenus(menus)
		assert.Len(t, routers, 1)
		assert.Equal(t, "Parent", routers[0].Name)
		assert.Len(t, routers[0].Children, 1)
		assert.Equal(t, "Child", routers[0].Children[0].Name)
	})
}
