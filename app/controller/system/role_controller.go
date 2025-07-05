package systemcontroller

import (
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/utils"
	"strconv"
	"time"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// Role list
func (*RoleController) List(ctx *gin.Context) {
	var param dto.RoleListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	roles, total := (&service.RoleService{}).GetRoleList(param, true)

	response.NewSuccess().SetPageData(roles, total).Json(ctx)
}

// Role details
func (*RoleController) Detail(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))

	role := (&service.RoleService{}).GetRoleByRoleId(roleId)

	response.NewSuccess().SetData("data", role).Json(ctx)
}

// Add role
func (*RoleController) Create(ctx *gin.Context) {
	var param dto.CreateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateRoleValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if role := (&service.RoleService{}).GetRoleByRoleName(param.RoleName); role.RoleId > 0 {
		response.NewError().SetMsg("Failed to add role " + param.RoleName + ", role name already exists").Json(ctx)
		return
	}

	if role := (&service.RoleService{}).GetRoleByRoleKey(param.RoleKey); role.RoleId > 0 {
		response.NewError().SetMsg("Failed to add role " + param.RoleName + ", permission character already exists").Json(ctx)
		return
	}

	menuCheckStrictly, deptCheckStrictly := 0, 0
	if param.MenuCheckStrictly {
		menuCheckStrictly = 1
	}
	if param.DeptCheckStrictly {
		deptCheckStrictly = 1
	}

	if err := (&service.RoleService{}).CreateRole(dto.SaveRole{
		RoleName:          param.RoleName,
		RoleKey:           param.RoleKey,
		RoleSort:          param.RoleSort,
		MenuCheckStrictly: &menuCheckStrictly,
		DeptCheckStrictly: &deptCheckStrictly,
		Status:            param.Status,
		CreateBy:          security.GetAuthUserName(ctx),
		Remark:            param.Remark,
	}, param.MenuIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update role
func (*RoleController) Update(ctx *gin.Context) {
	var param dto.UpdateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateRoleValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if role := (&service.RoleService{}).GetRoleByRoleName(param.RoleName); role.RoleId > 0 && role.RoleId != param.RoleId {
		response.NewError().SetMsg("Failed to modify role " + param.RoleName + ", role name already exists").Json(ctx)
		return
	}

	if role := (&service.RoleService{}).GetRoleByRoleKey(param.RoleKey); role.RoleId > 0 && role.RoleId != param.RoleId {
		response.NewError().SetMsg("Failed to modify role " + param.RoleName + ", permission character already exists").Json(ctx)
		return
	}

	menuCheckStrictly, deptCheckStrictly := 0, 0
	if param.MenuCheckStrictly {
		menuCheckStrictly = 1
	}
	if param.DeptCheckStrictly {
		deptCheckStrictly = 1
	}

	if err := (&service.RoleService{}).UpdateRole(dto.SaveRole{
		RoleId:            param.RoleId,
		RoleName:          param.RoleName,
		RoleKey:           param.RoleKey,
		RoleSort:          param.RoleSort,
		MenuCheckStrictly: &menuCheckStrictly,
		DeptCheckStrictly: &deptCheckStrictly,
		Status:            param.Status,
		UpdateBy:          security.GetAuthUserName(ctx),
		Remark:            param.Remark,
	}, param.MenuIds, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Delete role
func (*RoleController) Remove(ctx *gin.Context) {
	roleIds, err := utils.StringToIntSlice(ctx.Param("roleIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	roles := (&service.RoleService{}).GetRoleListByUserId(security.GetAuthUserId(ctx))

	for _, role := range roles {
		if err = validator.RemoveRoleValidator(roleIds, role.RoleId, role.RoleName); err != nil {
			response.NewError().SetMsg(err.Error()).Json(ctx)
			return
		}
	}

	if err = (&service.RoleService{}).DeleteRole(roleIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Change role status
func (*RoleController) ChangeStatus(ctx *gin.Context) {
	var param dto.UpdateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.ChangeRoleStatusValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.RoleService{}).UpdateRole(dto.SaveRole{
		RoleId:   param.RoleId,
		Status:   param.Status,
		UpdateBy: security.GetAuthUserName(ctx),
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Department tree
func (*RoleController) DeptTree(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))
	roleHasDeptIds := (&service.DeptService{}).GetDeptIdsByRoleId(roleId)

	depts := (&service.DeptService{}).DeptSelect()
	tree := (&service.DeptService{}).DeptSeleteToTree(depts, 0)

	response.NewSuccess().SetData("depts", tree).SetData("checkedKeys", roleHasDeptIds).Json(ctx)
}

// Assign data permissions
func (*RoleController) DataScope(ctx *gin.Context) {
	var param dto.UpdateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	deptCheckStrictly := 0
	if param.DeptCheckStrictly {
		deptCheckStrictly = 1
	}

	if err := (&service.RoleService{}).UpdateRole(dto.SaveRole{
		RoleId:            param.RoleId,
		DataScope:         param.DataScope,
		DeptCheckStrictly: &deptCheckStrictly,
		UpdateBy:          security.GetAuthUserName(ctx),
	}, nil, param.DeptIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Query the list of allocated user roles
func (*RoleController) RoleAuthUserAllocatedList(ctx *gin.Context) {
	var param dto.RoleAuthUserAllocatedListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	users, total := (&service.UserService{}).GetUserListByRoleId(param, security.GetAuthUserId(ctx), true)

	response.NewSuccess().SetPageData(users, total).Json(ctx)
}

// Query the list of unallocated user roles
func (*RoleController) RoleAuthUserUnallocatedList(ctx *gin.Context) {
	var param dto.RoleAuthUserAllocatedListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	users, total := (&service.UserService{}).GetUserListByRoleId(param, security.GetAuthUserId(ctx), false)

	response.NewSuccess().SetPageData(users, total).Json(ctx)
}

// Batch select user authorization
func (*RoleController) RoleAuthUserSelectAll(ctx *gin.Context) {
	var param dto.RoleAuthUserSelectAllRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	userIds, err := utils.StringToIntSlice(param.UserIds, ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = (&service.RoleService{}).AuthUserSelectAll(param.RoleId, userIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Cancel authorized user
func (*RoleController) RoleAuthUserCancel(ctx *gin.Context) {
	var param dto.RoleAuthUserCancelRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.RoleService{}).AuthUserDelete(param.RoleId, []int{param.UserId}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Batch cancel authorized users
func (*RoleController) RoleAuthUserCancelAll(ctx *gin.Context) {
	var param dto.RoleAuthUserCancelAllRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	userIds, err := utils.StringToIntSlice(param.UserIds, ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.RoleService{}).AuthUserDelete(param.RoleId, userIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Data export
func (*RoleController) Export(ctx *gin.Context) {
	var param dto.RoleListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.RoleExportResponse, 0)

	roles, _ := (&service.RoleService{}).GetRoleList(param, false)
	for _, role := range roles {
		list = append(list, dto.RoleExportResponse{
			RoleId:    role.RoleId,
			RoleName:  role.RoleName,
			RoleKey:   role.RoleKey,
			RoleSort:  role.RoleSort,
			DataScope: role.DataScope,
			Status:    role.Status,
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("role_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}
