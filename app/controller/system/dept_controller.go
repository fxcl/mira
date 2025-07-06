package systemcontroller

import (
	"strconv"
	"strings"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/types/constant"
	"mira/common/utils"

	"github.com/gin-gonic/gin"
)

// DeptController handles department-related operations.
type DeptController struct {
	DeptService *service.DeptService
	UserService *service.UserService
}

// NewDeptController creates a new DeptController.
func NewDeptController(deptService *service.DeptService, userService *service.UserService) *DeptController {
	return &DeptController{
		DeptService: deptService,
		UserService: userService,
	}
}

// List retrieves the department list.
// @Summary Get department list
// @Description Retrieves the department list based on query parameters and user permissions.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.DeptListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=[]dto.DeptListResponse} "Success"
// @Router /system/dept/list [get]
func (c *DeptController) List(ctx *gin.Context) {
	var param dto.DeptListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	depts := c.DeptService.GetDeptList(param, security.GetAuthUserId(ctx))

	response.NewSuccess().SetData("data", depts).Json(ctx)
}

// ListExclude retrieves the department list, excluding a specific node and its children.
// @Summary Get department list (exclude node)
// @Description Retrieves the department list, excluding a specific department and its sub-departments.
// @Tags System
// @Accept json
// @Produce json
// @Param deptId path int true "Department ID to exclude"
// @Success 200 {object} response.Response{data=[]dto.DeptListResponse} "Success"
// @Router /system/dept/list/exclude/{deptId} [get]
func (c *DeptController) ListExclude(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))

	data := make([]dto.DeptListResponse, 0)

	depts := c.DeptService.GetDeptList(dto.DeptListRequest{}, security.GetAuthUserId(ctx))
	for _, dept := range depts {
		if dept.DeptId == deptId || utils.Contains(strings.Split(dept.Ancestors, ","), strconv.Itoa(deptId)) {
			continue
		}
		data = append(data, dept)
	}

	response.NewSuccess().SetData("data", data).Json(ctx)
}

// Detail retrieves the details of a specific department.
// @Summary Get department details
// @Description Retrieves the details of a department by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param deptId path int true "Department ID"
// @Success 200 {object} response.Response{data=dto.DeptDetailResponse} "Success"
// @Router /system/dept/{deptId} [get]
func (c *DeptController) Detail(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))

	dept := c.DeptService.GetDeptByDeptId(deptId)

	response.NewSuccess().SetData("data", dept).Json(ctx)
}

// Create adds a new department.
// @Summary Add department
// @Description Adds a new department to the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreateDeptRequest true "Department data"
// @Success 200 {object} response.Response "Success"
// @Router /system/dept [post]
func (c *DeptController) Create(ctx *gin.Context) {
	var param dto.CreateDeptRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateDeptValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dept := c.DeptService.GetDeptByDeptName(param.DeptName); dept.DeptId > 0 {
		response.NewError().SetMsg("Failed to add department " + param.DeptName + ", department name already exists").Json(ctx)
		return
	}

	// Splice ancestors, get the ancestor list of the parent
	parentDept := c.DeptService.GetDeptByDeptId(param.ParentId)
	if parentDept.Status == constant.EXCEPTION_STATUS {
		response.NewError().SetMsg("Department is disabled, adding is not allowed").Json(ctx)
		return
	}
	ancestors := parentDept.Ancestors + "," + strconv.Itoa(parentDept.DeptId)

	if err := c.DeptService.CreateDept(dto.SaveDept{
		ParentId:  param.ParentId,
		Ancestors: ancestors,
		DeptName:  param.DeptName,
		OrderNum:  param.OrderNum,
		Leader:    param.Leader,
		Phone:     param.Phone,
		Email:     param.Email,
		Status:    param.Status,
		CreateBy:  security.GetAuthUserName(ctx),
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update modifies an existing department.
// @Summary Update department
// @Description Modifies an existing department in the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateDeptRequest true "Department data"
// @Success 200 {object} response.Response "Success"
// @Router /system/dept [put]
func (c *DeptController) Update(ctx *gin.Context) {
	var param dto.UpdateDeptRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateDeptValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dept := c.DeptService.GetDeptByDeptName(param.DeptName); dept.DeptId > 0 && dept.DeptId != param.DeptId {
		response.NewError().SetMsg("Failed to modify department " + param.DeptName + ", department name already exists").Json(ctx)
		return
	}

	if dept := c.DeptService.GetDeptByDeptId(param.DeptId); dept.ParentId != param.ParentId && c.DeptService.DeptHasChildren(param.DeptId) {
		response.NewError().SetMsg("Sub-departments exist, cannot directly modify the parent department").Json(ctx)
		return
	}

	// Splice ancestors, get the ancestor list of the parent
	parentDept := c.DeptService.GetDeptByDeptId(param.ParentId)
	if parentDept.Status == constant.EXCEPTION_STATUS {
		response.NewError().SetMsg("Department is disabled, adding is not allowed").Json(ctx)
		return
	}
	ancestors := parentDept.Ancestors + "," + strconv.Itoa(parentDept.DeptId)

	if err := c.DeptService.UpdateDept(dto.SaveDept{
		DeptId:    param.DeptId,
		ParentId:  param.ParentId,
		Ancestors: ancestors,
		DeptName:  param.DeptName,
		OrderNum:  param.OrderNum,
		Leader:    param.Leader,
		Phone:     param.Phone,
		Email:     param.Email,
		Status:    param.Status,
		UpdateBy:  security.GetAuthUserName(ctx),
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Remove deletes a department.
// @Summary Delete department
// @Description Deletes a department by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param deptId path int true "Department ID"
// @Success 200 {object} response.Response "Success"
// @Router /system/dept/{deptId} [delete]
func (c *DeptController) Remove(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))

	if c.DeptService.DeptHasChildren(deptId) {
		response.NewError().SetMsg("Sub-departments exist, deletion is not allowed").Json(ctx)
		return
	}

	if c.UserService.UserHasDeptByDeptId(deptId) {
		response.NewError().SetMsg("Users exist in the department, deletion is not allowed").Json(ctx)
		return
	}

	if err := c.DeptService.DeleteDept(deptId); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
