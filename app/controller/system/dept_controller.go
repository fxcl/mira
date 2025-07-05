package systemcontroller

import (
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/types/constant"
	"mira/common/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type DeptController struct{}

// Department list
func (*DeptController) List(ctx *gin.Context) {
	var param dto.DeptListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	depts := (&service.DeptService{}).GetDeptList(param, security.GetAuthUserId(ctx))

	response.NewSuccess().SetData("data", depts).Json(ctx)
}

// Query department list (excluding nodes)
func (*DeptController) ListExclude(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))

	data := make([]dto.DeptListResponse, 0)

	depts := (&service.DeptService{}).GetDeptList(dto.DeptListRequest{}, security.GetAuthUserId(ctx))
	for _, dept := range depts {
		if dept.DeptId == deptId || utils.Contains(strings.Split(dept.Ancestors, ","), strconv.Itoa(deptId)) {
			continue
		}
		data = append(data, dept)
	}

	response.NewSuccess().SetData("data", data).Json(ctx)
}

// Get department details
func (*DeptController) Detail(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))

	dept := (&service.DeptService{}).GetDeptByDeptId(deptId)

	response.NewSuccess().SetData("data", dept).Json(ctx)
}

// Add department
func (*DeptController) Create(ctx *gin.Context) {
	var param dto.CreateDeptRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateDeptValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dept := (&service.DeptService{}).GetDeptByDeptName(param.DeptName); dept.DeptId > 0 {
		response.NewError().SetMsg("Failed to add department " + param.DeptName + ", department name already exists").Json(ctx)
		return
	}

	// Splice ancestors, get the ancestor list of the parent
	parentDept := (&service.DeptService{}).GetDeptByDeptId(param.ParentId)
	if parentDept.Status == constant.EXCEPTION_STATUS {
		response.NewError().SetMsg("Department is disabled, adding is not allowed").Json(ctx)
		return
	}
	ancestors := parentDept.Ancestors + "," + strconv.Itoa(parentDept.DeptId)

	if err := (&service.DeptService{}).CreateDept(dto.SaveDept{
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

// Update department
func (*DeptController) Update(ctx *gin.Context) {
	var param dto.UpdateDeptRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateDeptValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if dept := (&service.DeptService{}).GetDeptByDeptName(param.DeptName); dept.DeptId > 0 && dept.DeptId != param.DeptId {
		response.NewError().SetMsg("Failed to modify department " + param.DeptName + ", department name already exists").Json(ctx)
		return
	}

	if dept := (&service.DeptService{}).GetDeptByDeptId(param.DeptId); dept.ParentId != param.ParentId && (&service.DeptService{}).DeptHasChildren(param.DeptId) {
		response.NewError().SetMsg("Sub-departments exist, cannot directly modify the parent department").Json(ctx)
		return
	}

	// Splice ancestors, get the ancestor list of the parent
	parentDept := (&service.DeptService{}).GetDeptByDeptId(param.ParentId)
	if parentDept.Status == constant.EXCEPTION_STATUS {
		response.NewError().SetMsg("Department is disabled, adding is not allowed").Json(ctx)
		return
	}
	ancestors := parentDept.Ancestors + "," + strconv.Itoa(parentDept.DeptId)

	if err := (&service.DeptService{}).UpdateDept(dto.SaveDept{
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

// Delete department
func (*DeptController) Remove(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))

	if (&service.DeptService{}).DeptHasChildren(deptId) {
		response.NewError().SetMsg("Sub-departments exist, deletion is not allowed").Json(ctx)
		return
	}

	if (&service.UserService{}).UserHasDeptByDeptId(deptId) {
		response.NewError().SetMsg("Users exist in the department, deletion is not allowed").Json(ctx)
		return
	}

	if err := (&service.DeptService{}).DeleteDept(deptId); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
