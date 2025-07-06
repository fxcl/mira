package systemcontroller

import (
	"strconv"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"

	"github.com/gin-gonic/gin"
)

// MenuController handles menu-related operations.
type MenuController struct {
	MenuService *service.MenuService
}

// NewMenuController creates a new MenuController.
func NewMenuController(menuService *service.MenuService) *MenuController {
	return &MenuController{MenuService: menuService}
}

// List retrieves the menu list.
// @Summary Get menu list
// @Description Retrieves the menu list based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.MenuListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=[]dto.MenuListResponse} "Success"
// @Router /system/menu/list [get]
func (c *MenuController) List(ctx *gin.Context) {
	var param dto.MenuListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	menus := c.MenuService.GetMenuList(param)

	response.NewSuccess().SetData("data", menus).Json(ctx)
}

// Detail retrieves the details of a specific menu.
// @Summary Get menu details
// @Description Retrieves the details of a menu by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param menuId path int true "Menu ID"
// @Success 200 {object} response.Response{data=dto.MenuDetailResponse} "Success"
// @Router /system/menu/{menuId} [get]
func (c *MenuController) Detail(ctx *gin.Context) {
	menuId, _ := strconv.Atoi(ctx.Param("menuId"))

	menu := c.MenuService.GetMenuByMenuId(menuId)

	response.NewSuccess().SetData("data", menu).Json(ctx)
}

// Treeselect retrieves the menu tree for selection.
// @Summary Get menu drop-down tree list
// @Description Retrieves the menu tree structure for use in dropdowns.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.MenuTreeSelectResponse} "Success"
// @Router /system/menu/treeselect [get]
func (c *MenuController) Treeselect(ctx *gin.Context) {
	menus := c.MenuService.MenuSelect()

	tree := c.MenuService.MenuSeleteToTree(menus, 0)

	response.NewSuccess().SetData("data", tree).Json(ctx)
}

// RoleMenuTreeselect retrieves the menu tree for a specific role.
// @Summary Load corresponding role menu list tree
// @Description Retrieves the menu tree and the menu IDs associated with a specific role.
// @Tags System
// @Accept json
// @Produce json
// @Param roleId path int true "Role ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "Success, returns 'menus' and 'checkedKeys'"
// @Router /system/menu/roleMenuTreeselect/{roleId} [get]
func (c *MenuController) RoleMenuTreeselect(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))
	roleHasMenuIds := c.MenuService.GetMenuIdsByRoleId(roleId)

	menus := c.MenuService.MenuSelect()
	tree := c.MenuService.MenuSeleteToTree(menus, 0)

	response.NewSuccess().SetData("menus", tree).SetData("checkedKeys", roleHasMenuIds).Json(ctx)
}

// Create adds a new menu.
// @Summary Add menu
// @Description Adds a new menu to the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreateMenuRequest true "Menu data"
// @Success 200 {object} response.Response "Success"
// @Router /system/menu [post]
func (c *MenuController) Create(ctx *gin.Context) {
	var param dto.CreateMenuRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateMenuValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if menu := c.MenuService.GetMenuByMenuName(param.MenuName); menu.MenuId > 0 {
		response.NewError().SetMsg("Failed to add menu " + param.MenuName + ", menu name already exists").Json(ctx)
		return
	}

	if err := c.MenuService.CreateMenu(dto.SaveMenu{
		MenuName:  param.MenuName,
		ParentId:  param.ParentId,
		OrderNum:  param.OrderNum,
		Path:      param.Path,
		Component: param.Component,
		Query:     param.Query,
		RouteName: param.RouteName,
		IsFrame:   &param.IsFrame,
		IsCache:   &param.IsCache,
		MenuType:  param.MenuType,
		Visible:   param.Visible,
		Perms:     param.Perms,
		Icon:      param.Icon,
		Status:    param.Status,
		CreateBy:  security.GetAuthUserName(ctx),
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update modifies an existing menu.
// @Summary Update menu
// @Description Modifies an existing menu in the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateMenuRequest true "Menu data"
// @Success 200 {object} response.Response "Success"
// @Router /system/menu [put]
func (c *MenuController) Update(ctx *gin.Context) {
	var param dto.UpdateMenuRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateMenuValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if menu := c.MenuService.GetMenuByMenuName(param.MenuName); menu.MenuId > 0 && menu.MenuId != param.MenuId {
		response.NewError().SetMsg("Failed to modify menu " + param.MenuName + ", menu name already exists").Json(ctx)
		return
	}

	if err := c.MenuService.UpdateMenu(dto.SaveMenu{
		MenuId:    param.MenuId,
		MenuName:  param.MenuName,
		ParentId:  param.ParentId,
		OrderNum:  param.OrderNum,
		Path:      param.Path,
		Component: param.Component,
		Query:     param.Query,
		RouteName: param.RouteName,
		IsFrame:   &param.IsFrame,
		IsCache:   &param.IsCache,
		MenuType:  param.MenuType,
		Visible:   param.Visible,
		Perms:     param.Perms,
		Icon:      param.Icon,
		Status:    param.Status,
		UpdateBy:  security.GetAuthUserName(ctx),
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Remove deletes a menu.
// @Summary Delete menu
// @Description Deletes a menu by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param menuId path int true "Menu ID"
// @Success 200 {object} response.Response "Success"
// @Router /system/menu/{menuId} [delete]
func (c *MenuController) Remove(ctx *gin.Context) {
	menuId, _ := strconv.Atoi(ctx.Param("menuId"))

	if c.MenuService.MenuHasChildren(menuId) {
		response.NewError().SetMsg("Sub-menu exists, deletion is not allowed").Json(ctx)
		return
	}

	if c.MenuService.MenuExistRole(menuId) {
		response.NewError().SetMsg("Menu has been assigned, deletion is not allowed").Json(ctx)
		return
	}

	if err := c.MenuService.DeleteMenu(menuId); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
