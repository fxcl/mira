package dto

import "mira/anima/datetime"

// Department List
type DeptListResponse struct {
	DeptId     int               `json:"deptId"`
	ParentId   int               `json:"parentId"`
	Ancestors  string            `json:"ancestors"`
	DeptName   string            `json:"deptName"`
	OrderNum   int               `json:"orderNum"`
	Status     string            `json:"status"`
	CreateTime datetime.Datetime `json:"createTime"`
}

// Department List Tree
type DeptTreeListResponse struct {
	DeptListResponse
	Children []DeptTreeListResponse `json:"children"`
}

// Department Details
type DeptDetailResponse struct {
	DeptId     int               `json:"deptId"`
	ParentId   int               `json:"parentId"`
	Ancestors  string            `json:"ancestors"`
	DeptName   string            `json:"deptName"`
	OrderNum   int               `json:"orderNum"`
	Leader     string            `json:"leader"`
	Phone      string            `json:"phone"`
	Email      string            `json:"email"`
	Status     string            `json:"status"`
	CreateTime datetime.Datetime `json:"createTime"`
}

// Department Tree (User Management Tree)
type DeptTreeResponse struct {
	Id       int                `json:"id"`
	Label    string             `json:"label"`
	Children []DeptTreeResponse `json:"children" gorm:"-"`
	ParentId int                `json:"-"`
}
