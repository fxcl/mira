package service

import (
	"strings"

	"github.com/pkg/errors"
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"
)

// MenuServiceInterface defines operations for menu management
type MenuServiceInterface interface {
	// Menu CRUD operations
	CreateMenu(param dto.SaveMenu) error
	UpdateMenu(param dto.SaveMenu) error
	DeleteMenu(menuId int) error
	GetMenuList(param dto.MenuListRequest) []dto.MenuListResponse
	GetMenuByMenuId(menuId int) dto.MenuDetailResponse
	GetMenuByMenuName(menuName string) dto.MenuDetailResponse

	// Menu relationship operations
	MenuHasChildren(menuId int) bool
	MenuExistRole(menuId int) bool
	GetPermsByUserId(userId int) []string
	GetMenuIdsByRoleId(roleId int) []int

	// Menu tree and selection operations
	MenuSelect() []dto.SeleteTree
	MenuSeleteToTree(menus []dto.SeleteTree, parentId int) []dto.SeleteTree
	GetMenuMCListByUserId(userId int) []dto.MenuListResponse
	MenusToTree(menus []dto.MenuListResponse, parentId int) []dto.MenuListTreeResponse
	BuildRouterMenus(menus []dto.MenuListTreeResponse) []dto.MenuMetaTreeResponse

	// Helper methods for router building
	GetRouteName(menu dto.MenuListTreeResponse) string
	GetRouteNameOrDefault(name, path string) string
	GetRoutePath(menu dto.MenuListTreeResponse) string
	GetComponent(menu dto.MenuListTreeResponse) string
	IsMenuFrame(menu dto.MenuListTreeResponse) bool
	IsInnerLink(menu dto.MenuListTreeResponse) bool
	IsParentView(menu dto.MenuListTreeResponse) bool
	InnerLinkReplacePach(path string) string
}

// MenuService implements the menu management interface
type MenuService struct{}

// Ensure MenuService implements MenuServiceInterface
var _ MenuServiceInterface = (*MenuService)(nil)

// CreateMenu adds a new menu
func (s *MenuService) CreateMenu(param dto.SaveMenu) error {
	return s.CreateMenuWithErr(param)
}

// CreateMenuWithErr adds a new menu with proper error handling
func (s *MenuService) CreateMenuWithErr(param dto.SaveMenu) error {
	// Input validation
	if param.MenuName == "" {
		return errors.New("menu name cannot be empty")
	}

	err := dal.Gorm.Model(model.SysMenu{}).Create(&model.SysMenu{
		MenuName:  param.MenuName,
		ParentId:  param.ParentId,
		OrderNum:  param.OrderNum,
		Path:      param.Path,
		Component: param.Component,
		Query:     param.Query,
		RouteName: param.RouteName,
		IsFrame:   param.IsFrame,
		IsCache:   param.IsCache,
		MenuType:  param.MenuType,
		Visible:   param.Visible,
		Perms:     param.Perms,
		Icon:      param.Icon,
		Status:    param.Status,
		Remark:    param.Remark,
		CreateBy:  param.CreateBy,
	}).Error
	if err != nil {
		return errors.Wrap(err, "failed to create menu")
	}

	return nil
}

// UpdateMenu updates an existing menu
func (s *MenuService) UpdateMenu(param dto.SaveMenu) error {
	return s.UpdateMenuWithErr(param)
}

// UpdateMenuWithErr updates an existing menu with proper error handling
func (s *MenuService) UpdateMenuWithErr(param dto.SaveMenu) error {
	// Input validation
	if param.MenuId <= 0 {
		return errors.New("invalid menu ID")
	}

	if param.MenuName == "" {
		return errors.New("menu name cannot be empty")
	}

	err := dal.Gorm.Model(model.SysMenu{}).Where("menu_id = ?", param.MenuId).Updates(&model.SysMenu{
		MenuName:  param.MenuName,
		ParentId:  param.ParentId,
		OrderNum:  param.OrderNum,
		Path:      param.Path,
		Component: param.Component,
		Query:     param.Query,
		RouteName: param.RouteName,
		IsFrame:   param.IsFrame,
		IsCache:   param.IsCache,
		MenuType:  param.MenuType,
		Visible:   param.Visible,
		Perms:     param.Perms,
		Icon:      param.Icon,
		Status:    param.Status,
		UpdateBy:  param.UpdateBy,
		Remark:    param.Remark,
	}).Error
	if err != nil {
		return errors.Wrap(err, "failed to update menu")
	}

	return nil
}

// DeleteMenu removes a menu by its ID
func (s *MenuService) DeleteMenu(menuId int) error {
	return s.DeleteMenuWithErr(menuId)
}

// DeleteMenuWithErr removes a menu with proper error handling
func (s *MenuService) DeleteMenuWithErr(menuId int) error {
	// Input validation
	if menuId <= 0 {
		return errors.New("invalid menu ID")
	}

	// Check if menu has children
	if s.MenuHasChildren(menuId) {
		return errors.New("menu has child nodes, cannot delete")
	}

	// Check if menu is assigned to roles
	if s.MenuExistRole(menuId) {
		return errors.New("menu is assigned to roles, cannot delete")
	}

	err := dal.Gorm.Model(model.SysMenu{}).Where("menu_id = ?", menuId).Delete(&model.SysMenu{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete menu")
	}

	return nil
}

// GetMenuList retrieves a list of menus based on search parameters
func (s *MenuService) GetMenuList(param dto.MenuListRequest) []dto.MenuListResponse {
	menus, _ := s.GetMenuListWithErr(param)
	return menus
}

// GetMenuListWithErr retrieves a list of menus with proper error handling
func (s *MenuService) GetMenuListWithErr(param dto.MenuListRequest) ([]dto.MenuListResponse, error) {
	menus := make([]dto.MenuListResponse, 0)

	query := dal.Gorm.Model(model.SysMenu{}).Order("sys_menu.parent_id, sys_menu.order_num, sys_menu.menu_id")

	if param.MenuName != "" {
		query = query.Where("menu_name LIKE ?", "%"+param.MenuName+"%")
	}

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	err := query.Find(&menus).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve menu list")
	}

	return menus, nil
}

// GetMenuByMenuId retrieves a menu by its ID
func (s *MenuService) GetMenuByMenuId(menuId int) dto.MenuDetailResponse {
	menu, _ := s.GetMenuByMenuIdWithErr(menuId)
	return menu
}

// GetMenuByMenuIdWithErr retrieves a menu by its ID with proper error handling
func (s *MenuService) GetMenuByMenuIdWithErr(menuId int) (dto.MenuDetailResponse, error) {
	var menu dto.MenuDetailResponse

	// Input validation
	if menuId <= 0 {
		return menu, errors.New("invalid menu ID")
	}

	err := dal.Gorm.Model(model.SysMenu{}).Where("menu_id = ?", menuId).Last(&menu).Error
	if err != nil {
		return menu, errors.Wrap(err, "failed to retrieve menu by ID")
	}

	return menu, nil
}

// GetMenuByMenuName retrieves a menu by its name
func (s *MenuService) GetMenuByMenuName(menuName string) dto.MenuDetailResponse {
	menu, _ := s.GetMenuByMenuNameWithErr(menuName)
	return menu
}

// GetMenuByMenuNameWithErr retrieves a menu by its name with proper error handling
func (s *MenuService) GetMenuByMenuNameWithErr(menuName string) (dto.MenuDetailResponse, error) {
	var menu dto.MenuDetailResponse

	// Input validation
	if menuName == "" {
		return menu, errors.New("menu name cannot be empty")
	}

	err := dal.Gorm.Model(model.SysMenu{}).Where("menu_name = ?", menuName).Last(&menu).Error
	if err != nil {
		return menu, errors.Wrap(err, "failed to retrieve menu by name")
	}

	return menu, nil
}

// MenuHasChildren checks if a menu has child nodes
func (s *MenuService) MenuHasChildren(menuId int) bool {
	hasChildren, _ := s.MenuHasChildrenWithErr(menuId)
	return hasChildren
}

// MenuHasChildrenWithErr checks if a menu has child nodes with proper error handling
func (s *MenuService) MenuHasChildrenWithErr(menuId int) (bool, error) {
	if menuId <= 0 {
		return false, errors.New("invalid menu ID")
	}

	var count int64
	err := dal.Gorm.Model(model.SysMenu{}).Where("parent_id = ?", menuId).Count(&count).Error
	if err != nil {
		return false, errors.Wrap(err, "failed to check if menu has children")
	}

	return count > 0, nil
}

// MenuExistRole checks if a menu is assigned to any roles
func (s *MenuService) MenuExistRole(menuId int) bool {
	exists, _ := s.MenuExistRoleWithErr(menuId)
	return exists
}

// MenuExistRoleWithErr checks if a menu is assigned to any roles with proper error handling
func (s *MenuService) MenuExistRoleWithErr(menuId int) (bool, error) {
	if menuId <= 0 {
		return false, errors.New("invalid menu ID")
	}

	var count int64
	err := dal.Gorm.Model(model.SysRoleMenu{}).Where("menu_id = ?", menuId).Count(&count).Error
	if err != nil {
		return false, errors.Wrap(err, "failed to check if menu is assigned to roles")
	}

	return count > 0, nil
}

// GetPermsByUserId retrieves menu permissions for a user
func (s *MenuService) GetPermsByUserId(userId int) []string {
	perms, _ := s.GetPermsByUserIdWithErr(userId)
	return perms
}

// GetPermsByUserIdWithErr retrieves menu permissions for a user with proper error handling
func (s *MenuService) GetPermsByUserIdWithErr(userId int) ([]string, error) {
	if userId <= 0 {
		return nil, errors.New("invalid user ID")
	}

	perms := make([]string, 0)

	// Super administrator has all permissions
	if userId == 1 {
		perms = append(perms, "*:*:*")
		return perms, nil
	}

	err := dal.Gorm.Model(model.SysMenu{}).
		Joins("JOIN sys_role_menu ON sys_menu.menu_id = sys_role_menu.menu_id").
		Joins("JOIN sys_role ON sys_role_menu.role_id = sys_role.role_id").
		Joins("JOIN sys_user_role ON sys_role.role_id = sys_user_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_menu.status = ?", userId, constant.NORMAL_STATUS).
		Pluck("sys_menu.perms", &perms).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve user permissions")
	}

	return perms, nil
}

// GetMenuIdsByRoleId retrieves menu IDs assigned to a role
func (s *MenuService) GetMenuIdsByRoleId(roleId int) []int {
	menuIds, _ := s.GetMenuIdsByRoleIdWithErr(roleId)
	return menuIds
}

// GetMenuIdsByRoleIdWithErr retrieves menu IDs assigned to a role with proper error handling
func (s *MenuService) GetMenuIdsByRoleIdWithErr(roleId int) ([]int, error) {
	if roleId <= 0 {
		return nil, errors.New("invalid role ID")
	}

	menuIds := make([]int, 0)

	err := dal.Gorm.Model(model.SysRoleMenu{}).
		Joins("JOIN sys_menu ON sys_menu.menu_id = sys_role_menu.menu_id").
		Where("sys_menu.status = ? AND sys_role_menu.role_id = ?", constant.NORMAL_STATUS, roleId).
		Pluck("sys_menu.menu_id", &menuIds).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve menu IDs for role")
	}

	return menuIds, nil
}

// MenuSelect retrieves a list of menus for dropdown selection
func (s *MenuService) MenuSelect() []dto.SeleteTree {
	menus, _ := s.MenuSelectWithErr()
	return menus
}

// MenuSelectWithErr retrieves a list of menus for dropdown selection with proper error handling
func (s *MenuService) MenuSelectWithErr() ([]dto.SeleteTree, error) {
	menus := make([]dto.SeleteTree, 0)

	err := dal.Gorm.Model(model.SysMenu{}).Order("order_num, menu_id").
		Select("menu_id as id", "menu_name as label", "parent_id").
		Where("status = ?", constant.NORMAL_STATUS).
		Find(&menus).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve menu select list")
	}

	return menus, nil
}

// MenuSeleteToTree converts a flat menu list to a hierarchical tree structure
func (s *MenuService) MenuSeleteToTree(menus []dto.SeleteTree, parentId int) []dto.SeleteTree {
	tree, _ := s.MenuSeleteToTreeWithErr(menus, parentId)
	return tree
}

// MenuSeleteToTreeWithErr converts a flat menu list to a hierarchical tree structure with proper error handling
func (s *MenuService) MenuSeleteToTreeWithErr(menus []dto.SeleteTree, parentId int) ([]dto.SeleteTree, error) {
	if menus == nil {
		return nil, errors.New("menu list cannot be nil")
	}

	tree := make([]dto.SeleteTree, 0)

	for _, menu := range menus {
		if menu.ParentId == parentId {
			children, err := s.MenuSeleteToTreeWithErr(menus, menu.Id)
			if err != nil {
				return nil, err
			}

			tree = append(tree, dto.SeleteTree{
				Id:       menu.Id,
				Label:    menu.Label,
				ParentId: menu.ParentId,
				Children: children,
			})
		}
	}

	return tree, nil
}

// GetMenuMCListByUserId retrieves menu permissions (M-directory; C-menu) for a user
func (s *MenuService) GetMenuMCListByUserId(userId int) []dto.MenuListResponse {
	menus, _ := s.GetMenuMCListByUserIdWithErr(userId)
	return menus
}

// GetMenuMCListByUserIdWithErr retrieves menu permissions for a user with proper error handling
func (s *MenuService) GetMenuMCListByUserIdWithErr(userId int) ([]dto.MenuListResponse, error) {
	if userId <= 0 {
		return nil, errors.New("invalid user ID")
	}

	menus := make([]dto.MenuListResponse, 0)

	query := dal.Gorm.Model(model.SysMenu{}).
		Distinct("sys_menu.*").
		Order("sys_menu.parent_id, sys_menu.order_num").
		Joins("LEFT JOIN sys_role_menu ON sys_menu.menu_id = sys_role_menu.menu_id").
		Joins("LEFT JOIN sys_role ON sys_role_menu.role_id = sys_role.role_id").
		Joins("LEFT JOIN sys_user_role ON sys_role.role_id = sys_user_role.role_id").
		Where("sys_menu.status = ? AND sys_menu.menu_type IN ?", constant.NORMAL_STATUS, []string{"M", "C"})

	if userId > 1 {
		query = query.Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS)
	}

	err := query.Find(&menus).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve user menu permissions")
	}

	return menus, nil
}

// MenusToTree converts a flat menu list to a hierarchical tree structure
func (s *MenuService) MenusToTree(menus []dto.MenuListResponse, parentId int) []dto.MenuListTreeResponse {
	tree, _ := s.MenusToTreeWithErr(menus, parentId)
	return tree
}

// MenusToTreeWithErr converts a flat menu list to a hierarchical tree structure with proper error handling
func (s *MenuService) MenusToTreeWithErr(menus []dto.MenuListResponse, parentId int) ([]dto.MenuListTreeResponse, error) {
	if menus == nil {
		return nil, errors.New("menu list cannot be nil")
	}

	tree := make([]dto.MenuListTreeResponse, 0)

	for _, menu := range menus {
		if menu.ParentId == parentId {
			children, err := s.MenusToTreeWithErr(menus, menu.MenuId)
			if err != nil {
				return nil, err
			}

			tree = append(tree, dto.MenuListTreeResponse{
				MenuListResponse: menu,
				Children:         children,
			})
		}
	}

	return tree, nil
}

// BuildRouterMenus builds the menu structure required for front-end routing
func (s *MenuService) BuildRouterMenus(menus []dto.MenuListTreeResponse) []dto.MenuMetaTreeResponse {
	routers, _ := s.BuildRouterMenusWithErr(menus)
	return routers
}

// BuildRouterMenusWithErr builds the menu structure for front-end routing with proper error handling
func (s *MenuService) BuildRouterMenusWithErr(menus []dto.MenuListTreeResponse) ([]dto.MenuMetaTreeResponse, error) {
	if menus == nil {
		return nil, errors.New("menu tree list cannot be nil")
	}

	routers := make([]dto.MenuMetaTreeResponse, 0)

	for _, menu := range menus {
		router := dto.MenuMetaTreeResponse{
			Name:      s.GetRouteName(menu),
			Path:      s.GetRoutePath(menu),
			Component: s.GetComponent(menu),
			Hidden:    menu.Visible == "1",
			Meta: dto.MenuMetaResponse{
				Title:   menu.MenuName,
				Icon:    menu.Icon,
				NoCache: menu.IsCache == 1,
			},
		}

		if len(menu.Children) > 0 && menu.MenuType == constant.MENU_TYPE_DIRECTORY {
			router.AlwaysShow = true
			router.Redirect = "noRedirect"

			children, err := s.BuildRouterMenusWithErr(menu.Children)
			if err != nil {
				return nil, err
			}
			router.Children = children
		} else if s.IsMenuFrame(menu) {
			children := dto.MenuMetaTreeResponse{
				Path:      menu.Path,
				Component: menu.Component,
				Name:      s.GetRouteNameOrDefault(menu.RouteName, menu.Path),
				Meta: dto.MenuMetaResponse{
					Title:   menu.MenuName,
					Icon:    menu.Icon,
					NoCache: menu.IsCache == 1,
				},
				Query: menu.Query,
			}
			router.Children = append(router.Children, children)
		} else if menu.ParentId == 0 && s.IsInnerLink(menu) {
			router.Meta = dto.MenuMetaResponse{
				Title: menu.MenuName,
				Icon:  menu.Icon,
			}
			router.Path = "/"
			children := dto.MenuMetaTreeResponse{
				Path:      s.InnerLinkReplacePach(menu.Path),
				Component: constant.INNER_LINK_COMPONENT,
				Name:      s.GetRouteNameOrDefault(menu.RouteName, menu.Path),
				Meta: dto.MenuMetaResponse{
					Title: menu.MenuName,
					Icon:  menu.Icon,
					Link:  menu.Path,
				},
			}
			router.Children = append(router.Children, children)
		}

		routers = append(routers, router)
	}

	return routers, nil
}

// Get route name
func (s *MenuService) GetRouteName(menu dto.MenuListTreeResponse) string {
	if s.IsMenuFrame(menu) {
		return ""
	}

	return s.GetRouteNameOrDefault(menu.RouteName, menu.Path)
}

// Get the route name, if the route name is not configured, take the route address
func (s *MenuService) GetRouteNameOrDefault(name, path string) string {
	if name == "" {
		name = path
	}

	return strings.ToUpper(string(name[0])) + name[1:]
}

// Get route address
func (s *MenuService) GetRoutePath(menu dto.MenuListTreeResponse) string {
	routePath := menu.Path

	// Inner chain opens the external network mode
	if menu.ParentId != 0 && !s.IsInnerLink(menu) {
		routePath = s.InnerLinkReplacePach(routePath)
	}

	// Not an external link and is a first-level directory (type is directory)
	if menu.ParentId == 0 && menu.MenuType == constant.MENU_TYPE_DIRECTORY && menu.IsFrame == constant.MENU_NO_FRAME {
		routePath = "/" + routePath
	} else if s.IsMenuFrame(menu) {
		// Not an external link and is a first-level directory (type is menu)
		routePath = "/"
	}

	return routePath
}

// Get component information
func (s *MenuService) GetComponent(menu dto.MenuListTreeResponse) string {
	component := constant.LAYOUT_COMPONENT

	if menu.Component != "" && !s.IsMenuFrame(menu) {
		component = menu.Component
	} else if menu.Component == "" && menu.ParentId != 0 && s.IsInnerLink(menu) {
		component = constant.INNER_LINK_COMPONENT
	} else if menu.Component == "" && s.IsParentView(menu) {
		component = constant.PARENT_VIEW_COMPONENT
	}

	return component
}

// Whether it is an internal jump of the menu
func (s *MenuService) IsMenuFrame(menu dto.MenuListTreeResponse) bool {
	return menu.ParentId == 0 && constant.MENU_TYPE_MENU == menu.MenuType && menu.IsFrame == constant.MENU_NO_FRAME
}

// Whether it is an inner chain component
func (s *MenuService) IsInnerLink(menu dto.MenuListTreeResponse) bool {
	return menu.IsFrame == constant.MENU_NO_FRAME && strings.HasPrefix(menu.Path, "http")
}

// Whether it is a parent_view component
func (s *MenuService) IsParentView(menu dto.MenuListTreeResponse) bool {
	return menu.ParentId != 0 && menu.MenuType == constant.MENU_TYPE_DIRECTORY
}

// Inner chain domain name special character replacement
func (s *MenuService) InnerLinkReplacePach(path string) string {
	// Remove http:// and https://
	path = strings.ReplaceAll(path, "http://", "")
	path = strings.ReplaceAll(path, "https://", "")
	path = strings.ReplaceAll(path, "www.", "")

	// Replace . with /
	path = strings.ReplaceAll(path, ".", "/")

	// Replace : with /
	path = strings.ReplaceAll(path, ":", "/")

	return path
}
