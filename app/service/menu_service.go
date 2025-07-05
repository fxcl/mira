package service

import (
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"
	"strings"
)

type MenuService struct{}

// Add menu
func (s *MenuService) CreateMenu(param dto.SaveMenu) error {
	return dal.Gorm.Model(model.SysMenu{}).Create(&model.SysMenu{
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
}

// Update menu
func (s *MenuService) UpdateMenu(param dto.SaveMenu) error {
	return dal.Gorm.Model(model.SysMenu{}).Where("menu_id = ?", param.MenuId).Updates(&model.SysMenu{
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
}

// Delete menu
func (s *MenuService) DeleteMenu(menuId int) error {
	return dal.Gorm.Model(model.SysMenu{}).Where("menu_id = ?", menuId).Delete(&model.SysMenu{}).Error
}

// Menu list
func (s *MenuService) GetMenuList(param dto.MenuListRequest) []dto.MenuListResponse {
	menus := make([]dto.MenuListResponse, 0)

	query := dal.Gorm.Model(model.SysMenu{}).Order("sys_menu.parent_id, sys_menu.order_num, sys_menu.menu_id")

	if param.MenuName != "" {
		query.Where("menu_name LIKE ?", "%"+param.MenuName+"%")
	}

	if param.Status != "" {
		query.Where("status = ?", param.Status)
	}

	query.Find(&menus)

	return menus
}

// Query menu by menu ID
func (s *MenuService) GetMenuByMenuId(menuId int) dto.MenuDetailResponse {
	var menu dto.MenuDetailResponse

	dal.Gorm.Model(model.SysMenu{}).Where("menu_id = ?", menuId).Last(&menu)

	return menu
}

// Query menu by menu name
func (s *MenuService) GetMenuByMenuName(menuName string) dto.MenuDetailResponse {
	var menu dto.MenuDetailResponse

	dal.Gorm.Model(model.SysMenu{}).Where("menu_name = ?", menuName).Last(&menu)

	return menu
}

// Query whether there are sub-menus
func (s *MenuService) MenuHasChildren(menuId int) bool {
	var count int64

	dal.Gorm.Model(model.SysMenu{}).Where("parent_id = ?", menuId).Count(&count)

	return count > 0
}

// Query whether the menu has been assigned to permissions
func (s *MenuService) MenuExistRole(menuId int) bool {
	var count int64

	dal.Gorm.Model(model.SysRoleMenu{}).Where("menu_id = ?", menuId).Count(&count)

	return count > 0
}

// Query menu permission perms by user ID
func (s *MenuService) GetPermsByUserId(userId int) []string {
	perms := make([]string, 0)

	// Super administrator has all permissions
	if userId == 1 {
		perms = append(perms, "*:*:*")
	} else {
		dal.Gorm.Model(model.SysMenu{}).
			Joins("JOIN sys_role_menu ON sys_menu.menu_id = sys_role_menu.menu_id").
			Joins("JOIN sys_role ON sys_role_menu.role_id = sys_role.role_id").
			Joins("JOIN sys_user_role ON sys_role.role_id = sys_user_role.role_id").
			Where("sys_user_role.user_id = ? AND sys_menu.status = ?", userId, constant.NORMAL_STATUS).
			Pluck("sys_menu.perms", &perms)
	}

	return perms
}

// Query the set of menu IDs owned by role ID
func (s *MenuService) GetMenuIdsByRoleId(roleId int) []int {
	menuIds := make([]int, 0)

	dal.Gorm.Model(model.SysRoleMenu{}).
		Joins("JOIN sys_menu ON sys_menu.menu_id = sys_role_menu.menu_id").
		Where("sys_menu.status = ? AND sys_role_menu.role_id = ?", constant.NORMAL_STATUS, roleId).
		Pluck("sys_menu.menu_id", &menuIds)

	return menuIds
}

// Menu drop-down tree list
func (s *MenuService) MenuSelect() []dto.SeleteTree {
	menus := make([]dto.SeleteTree, 0)

	dal.Gorm.Model(model.SysMenu{}).Order("order_num, menu_id").
		Select("menu_id as id", "menu_name as label", "parent_id").
		Where("status = ?", constant.NORMAL_STATUS).
		Find(&menus)

	return menus
}

// Convert menu drop-down list to tree structure
func (s *MenuService) MenuSeleteToTree(menus []dto.SeleteTree, parentId int) []dto.SeleteTree {
	tree := make([]dto.SeleteTree, 0)

	for _, menu := range menus {
		if menu.ParentId == parentId {
			tree = append(tree, dto.SeleteTree{
				Id:       menu.Id,
				Label:    menu.Label,
				ParentId: menu.ParentId,
				Children: s.MenuSeleteToTree(menus, menu.Id),
			})
		}
	}

	return tree
}

// Query the menu permissions owned by the user ID (M-directory; C-menu; F-button)
func (s *MenuService) GetMenuMCListByUserId(userId int) []dto.MenuListResponse {
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

	query.Find(&menus)

	return menus
}

// Convert menu permission list to tree structure
func (s *MenuService) MenusToTree(menus []dto.MenuListResponse, parentId int) []dto.MenuListTreeResponse {
	tree := make([]dto.MenuListTreeResponse, 0)

	for _, menu := range menus {
		if menu.ParentId == parentId {
			tree = append(tree, dto.MenuListTreeResponse{
				MenuListResponse: menu,
				Children:         s.MenusToTree(menus, menu.MenuId),
			})
		}
	}

	return tree
}

// Build the menu required for the front-end routing
func (s *MenuService) BuildRouterMenus(menus []dto.MenuListTreeResponse) []dto.MenuMetaTreeResponse {
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
			router.Children = s.BuildRouterMenus(menu.Children)
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

	return routers
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
