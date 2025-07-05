package dto

// Save Department
type SaveDept struct {
	DeptId    int    `json:"deptId"`
	ParentId  int    `json:"parentId"`
	Ancestors string `json:"ancestors"`
	DeptName  string `json:"deptName"`
	OrderNum  int    `json:"orderNum"`
	Leader    string `json:"leader"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	CreateBy  string `json:"createBy"`
	UpdateBy  string `json:"updateBy"`
}

// Department List
type DeptListRequest struct {
	DeptName string `query:"deptName" form:"deptName"`
	Status   string `query:"status" form:"status"`
}

// Create Department
type CreateDeptRequest struct {
	ParentId int    `json:"parentId"`
	DeptName string `json:"deptName"`
	OrderNum int    `json:"orderNum"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// Update Department
type UpdateDeptRequest struct {
	DeptId    int    `json:"deptId"`
	ParentId  int    `json:"parentId"`
	Ancestors string `json:"ancestors"`
	DeptName  string `json:"deptName"`
	OrderNum  int    `json:"orderNum"`
	Leader    string `json:"leader"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}
