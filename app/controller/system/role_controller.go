package systemcontroller

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/utils"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleController handles role-related operations.
type RoleController struct {
	RoleService *service.RoleService
	DeptService *service.DeptService
	UserService *service.UserService
}

// NewRoleController creates a new RoleController.
func NewRoleController(roleService *service.RoleService, deptService *service.DeptService, userService *service.UserService) *RoleController {
	return &RoleController{
		RoleService: roleService,
		DeptService: deptService,
		UserService: userService,
	}
}

// List retrieves a paginated list of roles.
// @Summary Get role list
// @Description Retrieves a paginated list of roles based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.RoleListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.RoleListResponse}} "Success"
// @Router /system/role/list [get]
func (c *RoleController) List(ctx *gin.Context) {
	var param dto.RoleListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	roles, total := c.RoleService.GetRoleList(param, true)

	response.NewSuccess().SetPageData(roles, total).Json(ctx)
}

// Detail retrieves the details of a specific role.
// @Summary Get role details
// @Description Retrieves the details of a role by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param roleId path int true "Role ID"
// @Success 200 {object} response.Response{data=dto.RoleDetailResponse} "Success"
// @Router /system/role/{roleId} [get]
func (c *RoleController) Detail(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))

	role, err := c.RoleService.GetRoleByRoleId(roleId)
	if err != nil {
		response.NewError().SetMsg(fmt.Sprintf("Failed to get role information: %v", err)).Json(ctx)
		return
	}

	response.NewSuccess().SetData("data", role).Json(ctx)
}

// Create adds a new role.
// @Summary Add role
// @Description Adds a new role to the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreateRoleRequest true "Role data"
// @Success 200 {object} response.Response "Success"
// @Router /system/role [post]
func (c *RoleController) Create(ctx *gin.Context) {
	var param dto.CreateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateRoleValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	role, err := c.RoleService.GetRoleByRoleName(param.RoleName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.NewError().SetMsg(fmt.Sprintf("Failed to validate role name: %v", err)).Json(ctx)
		return
	}
	if role.RoleId > 0 {
		response.NewError().SetMsg("Failed to add role " + param.RoleName + ", role name already exists").Json(ctx)
		return
	}

	role, err = c.RoleService.GetRoleByRoleKey(param.RoleKey)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.NewError().SetMsg(fmt.Sprintf("Failed to validate role permission character: %v", err)).Json(ctx)
		return
	}
	if role.RoleId > 0 {
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

	if err := c.RoleService.CreateRole(dto.SaveRole{
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

// Update modifies an existing role.
// @Summary Update role
// @Description Modifies an existing role in the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateRoleRequest true "Role data"
// @Success 200 {object} response.Response "Success"
// @Router /system/role [put]
func (c *RoleController) Update(ctx *gin.Context) {
	var param dto.UpdateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateRoleValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	role, err := c.RoleService.GetRoleByRoleName(param.RoleName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.NewError().SetMsg(fmt.Sprintf("Failed to validate role name: %v", err)).Json(ctx)
		return
	}
	if role.RoleId > 0 && role.RoleId != param.RoleId {
		response.NewError().SetMsg("Failed to modify role " + param.RoleName + ", role name already exists").Json(ctx)
		return
	}

	role, err = c.RoleService.GetRoleByRoleKey(param.RoleKey)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.NewError().SetMsg(fmt.Sprintf("Failed to validate role permission character: %v", err)).Json(ctx)
		return
	}
	if role.RoleId > 0 && role.RoleId != param.RoleId {
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

	if err := c.RoleService.UpdateRole(dto.SaveRole{
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

// Remove deletes one or more roles.
// @Summary Delete role
// @Description Deletes roles by their IDs.
// @Tags System
// @Accept json
// @Produce json
// @Param roleIds path string true "Role IDs, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /system/role/{roleIds} [delete]
func (c *RoleController) Remove(ctx *gin.Context) {
	roleIds, err := utils.StringToIntSlice(ctx.Param("roleIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	roles, err := c.RoleService.GetRoleListByUserId(security.GetAuthUserId(ctx))
	if err != nil {
		response.NewError().SetMsg(fmt.Sprintf("Failed to get user role list: %v", err)).Json(ctx)
		return
	}

	for _, role := range roles {
		if err = validator.RemoveRoleValidator(roleIds, role.RoleId, role.RoleName); err != nil {
			response.NewError().SetMsg(err.Error()).Json(ctx)
			return
		}
	}

	if err = c.RoleService.DeleteRole(roleIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// ChangeStatus changes the status of a role.
// @Summary Change role status
// @Description Changes the status (e.g., active/inactive) of a role.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateRoleRequest true "Role status data"
// @Success 200 {object} response.Response "Success"
// @Router /system/role/changeStatus [put]
func (c *RoleController) ChangeStatus(ctx *gin.Context) {
	var param dto.UpdateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.ChangeRoleStatusValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := c.RoleService.UpdateRole(dto.SaveRole{
		RoleId:   param.RoleId,
		Status:   param.Status,
		UpdateBy: security.GetAuthUserName(ctx),
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// DeptTree retrieves the department tree for a specific role.
// @Summary Get department tree for role
// @Description Retrieves the department tree and the department IDs associated with a specific role.
// @Tags System
// @Accept json
// @Produce json
// @Param roleId path int true "Role ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "Success, returns 'depts' and 'checkedKeys'"
// @Router /system/role/deptTree/{roleId} [get]
func (c *RoleController) DeptTree(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))
	roleHasDeptIds := c.DeptService.GetDeptIdsByRoleId(roleId)

	depts := c.DeptService.DeptSelect()
	tree := c.DeptService.DeptSeleteToTree(depts, 0)

	response.NewSuccess().SetData("depts", tree).SetData("checkedKeys", roleHasDeptIds).Json(ctx)
}

// DataScope assigns data permissions to a role.
// @Summary Assign data permissions
// @Description Assigns data permissions (e.g., department access) to a role.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateRoleRequest true "Data scope data"
// @Success 200 {object} response.Response "Success"
// @Router /system/role/dataScope [put]
func (c *RoleController) DataScope(ctx *gin.Context) {
	var param dto.UpdateRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	deptCheckStrictly := 0
	if param.DeptCheckStrictly {
		deptCheckStrictly = 1
	}

	if err := c.RoleService.UpdateRole(dto.SaveRole{
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

// RoleAuthUserAllocatedList retrieves a list of users allocated to a role.
// @Summary Query allocated user role list
// @Description Retrieves a paginated list of users who have been allocated to a specific role.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.RoleAuthUserAllocatedListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.UserListResponse}} "Success"
// @Router /system/role/authUser/allocatedList [get]
func (c *RoleController) RoleAuthUserAllocatedList(ctx *gin.Context) {
	var param dto.RoleAuthUserAllocatedListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	users, total := c.UserService.GetUserListByRoleId(param, security.GetAuthUserId(ctx), true)

	response.NewSuccess().SetPageData(users, total).Json(ctx)
}

// RoleAuthUserUnallocatedList retrieves a list of users not allocated to a role.
// @Summary Query unallocated user role list
// @Description Retrieves a paginated list of users who have not been allocated to a specific role.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.RoleAuthUserAllocatedListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.UserListResponse}} "Success"
// @Router /system/role/authUser/unallocatedList [get]
func (c *RoleController) RoleAuthUserUnallocatedList(ctx *gin.Context) {
	var param dto.RoleAuthUserAllocatedListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	users, total := c.UserService.GetUserListByRoleId(param, security.GetAuthUserId(ctx), false)

	response.NewSuccess().SetPageData(users, total).Json(ctx)
}

// RoleAuthUserSelectAll authorizes multiple users to a role.
// @Summary Batch select user authorization
// @Description Authorizes multiple users to a specific role in a single operation.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.RoleAuthUserSelectAllRequest true "Authorization data"
// @Success 200 {object} response.Response "Success"
// @Router /system/role/authUser/selectAll [put]
func (c *RoleController) RoleAuthUserSelectAll(ctx *gin.Context) {
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

	if err = c.RoleService.AuthUserSelectAll(param.RoleId, userIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// RoleAuthUserCancel cancels a user's authorization for a role.
// @Summary Cancel authorized user
// @Description Cancels a single user's authorization for a specific role.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.RoleAuthUserCancelRequest true "Cancellation data"
// @Success 200 {object} response.Response "Success"
// @Router /system/role/authUser/cancel [put]
func (c *RoleController) RoleAuthUserCancel(ctx *gin.Context) {
	var param dto.RoleAuthUserCancelRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := c.RoleService.AuthUserDelete(param.RoleId, []int{param.UserId}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// RoleAuthUserCancelAll cancels authorization for multiple users.
// @Summary Batch cancel authorized user
// @Description Cancels authorization for multiple users for a specific role.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.RoleAuthUserCancelAllRequest true "Batch cancellation data"
// @Success 200 {object} response.Response "Success"
// @Router /system/role/authUser/cancelAll [put]
func (c *RoleController) RoleAuthUserCancelAll(ctx *gin.Context) {
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

	if err := c.RoleService.AuthUserDelete(param.RoleId, userIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Export exports role data to an Excel file.
// @Summary Export roles
// @Description Exports role data to an Excel file based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.RoleListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /system/role/export [post]
func (c *RoleController) Export(ctx *gin.Context) {
	var param dto.RoleListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.RoleExportResponse, 0)

	roles, _ := c.RoleService.GetRoleList(param, false)
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
