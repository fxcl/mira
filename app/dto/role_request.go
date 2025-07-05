package dto

// Save Role
type SaveRole struct {
	RoleId            int    `json:"roleId"`
	RoleName          string `json:"roleName"`
	RoleKey           string `json:"roleKey"`
	RoleSort          int    `json:"roleSort"`
	DataScope         string `json:"dataScope"`
	MenuCheckStrictly *int   `json:"menuCheckStrictly"`
	DeptCheckStrictly *int   `json:"deptCheckStrictly"`
	Status            string `json:"status"`
	CreateBy          string `json:"createBy"`
	UpdateBy          string `json:"updateBy"`
	Remark            string `json:"remark"`
}

// Role List
type RoleListRequest struct {
	PageRequest
	RoleName  string `query:"roleName" form:"roleName"`
	RoleKey   string `query:"roleKey" form:"roleKey"`
	Status    string `query:"status" form:"status"`
	BeginTime string `query:"params[beginTime]" form:"params[beginTime]"`
	EndTime   string `query:"params[endTime]" form:"params[endTime]"`
}

// Create Role
type CreateRoleRequest struct {
	RoleName          string `json:"roleName"`
	RoleKey           string `json:"roleKey"`
	RoleSort          int    `json:"roleSort"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly"`
	Status            string `json:"status"`
	Remark            string `json:"remark"`
	MenuIds           []int  `json:"menuIds"`
}

// Update Role
type UpdateRoleRequest struct {
	RoleId            int    `json:"roleId"`
	RoleName          string `json:"roleName"`
	RoleKey           string `json:"roleKey"`
	RoleSort          int    `json:"roleSort"`
	DataScope         string `json:"dataScope"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly"`
	Status            string `json:"status"`
	Remark            string `json:"remark"`
	MenuIds           []int  `json:"menuIds"`
	DeptIds           []int  `json:"deptIds"`
}

// Query allocated user role list
type RoleAuthUserAllocatedListRequest struct {
	PageRequest
	RoleId      int    `query:"roleId" form:"roleId"`
	UserName    string `query:"userName" form:"userName"`
	Phonenumber string `query:"phonenumber" form:"phonenumber"`
}

// Batch select user authorization
type RoleAuthUserSelectAllRequest struct {
	RoleId  int    `query:"roleId" form:"roleId"`
	UserIds string `query:"userIds" form:"userIds"`
}

// Cancel user authorization
type RoleAuthUserCancelRequest struct {
	RoleId int `json:"roleId,string"`
	UserId int `json:"userId"`
}

// Batch cancel user authorization
type RoleAuthUserCancelAllRequest struct {
	RoleId  int    `query:"roleId" form:"roleId"`
	UserIds string `query:"userIds" form:"userIds"`
}
