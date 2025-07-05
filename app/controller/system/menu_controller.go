package systemcontroller

import (
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuController struct{}

// Menu list
func (*MenuController) List(ctx *gin.Context) {
	var param dto.MenuListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	menus := (&service.MenuService{}).GetMenuList(param)

	response.NewSuccess().SetData("data", menus).Json(ctx)
}

// Menu details
func (*MenuController) Detail(ctx *gin.Context) {
	menuId, _ := strconv.Atoi(ctx.Param("menuId"))

	menu := (&service.MenuService{}).GetMenuByMenuId(menuId)

	response.NewSuccess().SetData("data", menu).Json(ctx)
}

// Get menu drop-down tree list
func (*MenuController) Treeselect(ctx *gin.Context) {
	menus := (&service.MenuService{}).MenuSelect()

	tree := (&service.MenuService{}).MenuSeleteToTree(menus, 0)

	response.NewSuccess().SetData("data", tree).Json(ctx)
}

// Load the corresponding role menu list tree
func (*MenuController) RoleMenuTreeselect(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))
	roleHasMenuIds := (&service.MenuService{}).GetMenuIdsByRoleId(roleId)

	menus := (&service.MenuService{}).MenuSelect()
	tree := (&service.MenuService{}).MenuSeleteToTree(menus, 0)

	response.NewSuccess().SetData("menus", tree).SetData("checkedKeys", roleHasMenuIds).Json(ctx)
}

// Add menu
func (*MenuController) Create(ctx *gin.Context) {
	var param dto.CreateMenuRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateMenuValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if menu := (&service.MenuService{}).GetMenuByMenuName(param.MenuName); menu.MenuId > 0 {
		response.NewError().SetMsg("Failed to add menu " + param.MenuName + ", menu name already exists").Json(ctx)
		return
	}

	if err := (&service.MenuService{}).CreateMenu(dto.SaveMenu{
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

// Update menu
func (*MenuController) Update(ctx *gin.Context) {
	var param dto.UpdateMenuRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateMenuValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if menu := (&service.MenuService{}).GetMenuByMenuName(param.MenuName); menu.MenuId > 0 && menu.MenuId != param.MenuId {
		response.NewError().SetMsg("Failed to modify menu " + param.MenuName + ", menu name already exists").Json(ctx)
		return
	}

	if err := (&service.MenuService{}).UpdateMenu(dto.SaveMenu{
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

// Delete menu
func (*MenuController) Remove(ctx *gin.Context) {
	menuId, _ := strconv.Atoi(ctx.Param("menuId"))

	if (&service.MenuService{}).MenuHasChildren(menuId) {
		response.NewError().SetMsg("Sub-menu exists, deletion is not allowed").Json(ctx)
		return
	}

	if (&service.MenuService{}).MenuExistRole(menuId) {
		response.NewError().SetMsg("Menu has been assigned, deletion is not allowed").Json(ctx)
		return
	}

	if err := (&service.MenuService{}).DeleteMenu(menuId); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
