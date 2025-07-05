package dto

import "mira/anima/datetime"

// Save User
type SaveUser struct {
	UserId      int               `json:"userId"`
	DeptId      int               `json:"deptId"`
	UserName    string            `json:"userName"`
	NickName    string            `json:"nickName"`
	UserType    string            `json:"userType"`
	Email       string            `json:"email"`
	Phonenumber string            `json:"phonenumber"`
	Sex         string            `json:"sex"`
	Avatar      string            `json:"avatar"`
	Password    string            `json:"password"`
	LoginIP     string            `json:"loginIp"`
	LoginDate   datetime.Datetime `json:"loginDate"`
	Status      string            `json:"status"`
	CreateBy    string            `json:"createBy"`
	UpdateBy    string            `json:"updateBy"`
	Remark      string            `json:"remark"`
}

// User List
type UserListRequest struct {
	PageRequest
	UserName    string `query:"userName" form:"userName"`
	Phonenumber string `query:"phonenumber" form:"phonenumber"`
	Status      string `query:"status" form:"status"`
	DeptId      int    `query:"deptId" form:"deptId"`
	BeginTime   string `query:"params[beginTime]" form:"params[beginTime]"`
	EndTime     string `query:"params[endTime]" form:"params[endTime]"`
}

// Create User
type CreateUserRequest struct {
	DeptId      int    `json:"deptId"`
	UserName    string `json:"userName"`
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Sex         string `json:"sex"`
	Password    string `json:"password"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
	PostIds     []int  `json:"postIds"`
	RoleIds     []int  `json:"roleIds"`
}

// Update User
type UpdateUserRequest struct {
	UserId      int    `json:"userId"`
	DeptId      int    `json:"deptId"`
	UserName    string `json:"userName"`
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Sex         string `json:"sex"`
	Password    string `json:"password"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
	PostIds     []int  `json:"postIds"`
	RoleIds     []int  `json:"roleIds"`
}

// User Authorized Role
type AddUserAuthRoleRequest struct {
	UserId  int    `query:"userId" form:"userId"`
	RoleIds string `query:"roleIds" form:"roleIds"`
}

// Update Personal Information
type UpdateProfileRequest struct {
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Sex         string `json:"sex"`
}

// Update Personal Password
type UserProfileUpdatePwdRequest struct {
	OldPassword string `query:"oldPassword" form:"oldPassword"`
	NewPassword string `query:"newPassword" form:"newPassword"`
}

// User Import
type UserImportRequest struct {
	DeptId      int    `excel:"name:Dept ID;"`
	UserName    string `excel:"name:Login Name;"`
	NickName    string `excel:"name:User Name;"`
	Email       string `excel:"name:User Email;"`
	Phonenumber string `excel:"name:Phone Number;"`
	Sex         string `excel:"name:User Gender;replace:0_Male,1_Female,2_Unknown;"`
	Status      string `excel:"name:Account Status;replace:0_Normal,1_Disabled;"`
}
