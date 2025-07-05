package dto

import (
	"mira/anima/datetime"
)

// Role List
type RoleListResponse struct {
	RoleId            int               `json:"roleId"`
	RoleName          string            `json:"roleName"`
	RoleKey           string            `json:"roleKey"`
	RoleSort          int               `json:"roleSort"`
	DataScope         string            `json:"dataScope"`
	MenuCheckStrictly bool              `json:"menuCheckStrictly"`
	DeptCheckStrictly bool              `json:"deptCheckStrictly"`
	Status            string            `json:"status"`
	CreateTime        datetime.Datetime `json:"createTime"`
	Flag              bool              `json:"flag" gorm:"-"`
}

// Role Details
type RoleDetailResponse struct {
	RoleId            int    `json:"roleId"`
	RoleName          string `json:"roleName"`
	RoleKey           string `json:"roleKey"`
	RoleSort          int    `json:"roleSort"`
	DataScope         string `json:"dataScope"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly"`
	Status            string `json:"status"`
	Remark            string `json:"remark"`
}

// Role Export
type RoleExportResponse struct {
	RoleId    int    `excel:"name:Role ID;"`
	RoleName  string `excel:"name:Role Name;"`
	RoleKey   string `excel:"name:Role Permission;"`
	RoleSort  int    `excel:"name:Role Sort;"`
	DataScope string `excel:"name:Data Scope;replace:1_All Data Permissions,2_Custom Data Permissions,3_Department Data Permissions,4_Department and Below Data Permissions,5_Only Personal Data Permissions;"`
	Status    string `excel:"name:Role Status;replace:0_Normal,1_Disabled;"`
}
