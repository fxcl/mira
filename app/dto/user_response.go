package dto

import "mira/anima/datetime"

// User Authorization
type UserTokenResponse struct {
	UserId   int    `json:"userId"`
	DeptId   int    `json:"deptId"`
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	UserType string `json:"userType"`
	Password string `json:"-"`
	Status   string `json:"status"`
	DeptName string `json:"deptName"`
}

// User List
type UserListResponse struct {
	UserId      int               `json:"userId"`
	DeptId      int               `json:"deptId"`
	UserName    string            `json:"userName"`
	NickName    string            `json:"nickName"`
	Email       string            `json:"email"`
	Phonenumber string            `json:"phonenumber"`
	Sex         string            `json:"sex"`
	LoginIp     string            `json:"loginIp"`
	LoginDate   datetime.Datetime `json:"loginDate"`
	Status      string            `json:"status"`
	CreateTime  datetime.Datetime `json:"createTime"`
	Dept        struct {
		DeptId   int    `json:"deptId"`
		DeptName string `json:"deptName"`
		Leader   string `json:"leader"`
	} `json:"dept" gorm:"-"`
	DeptName string `json:"-"`
	Leader   string `json:"-"`
}

// User Details
type UserDetailResponse struct {
	UserId      int               `json:"userId"`
	DeptId      int               `json:"deptId"`
	UserName    string            `json:"userName"`
	NickName    string            `json:"nickName"`
	UserType    string            `json:"userType"`
	Email       string            `json:"email"`
	Phonenumber string            `json:"phonenumber"`
	Sex         string            `json:"sex"`
	Avatar      string            `json:"avatar"`
	Password    string            `json:"-"`
	LoginIP     string            `json:"loginIp"`
	LoginDate   datetime.Datetime `json:"loginDate"`
	Status      string            `json:"status"`
	CreateTime  datetime.Datetime `json:"createTime"`
	Admin       bool              `json:"admin" gorm:"-"`
}

// Authorized User Information
type AuthUserInfoResponse struct {
	UserDetailResponse
	Dept  DeptDetailResponse `json:"dept"`
	Roles []RoleListResponse `json:"roles"`
}

// User Export
type UserExportResponse struct {
	UserId      int    `excel:"name:User ID;"`
	UserName    string `excel:"name:Login Name;"`
	NickName    string `excel:"name:User Name;"`
	Email       string `excel:"name:User Email;"`
	Phonenumber string `excel:"name:Phone Number;"`
	Sex         string `excel:"name:User Gender;replace:0_Male,1_Female,2_Unknown;"`
	Status      string `excel:"name:Account Status;replace:0_Normal,1_Disabled;"`
	LoginIp     string `excel:"name:Last Login IP;"`
	LoginDate   string `excel:"name:Last Login Time;"`
	DeptName    string `excel:"name:Dept Name;"`
	DeptLeader  string `excel:"name:Dept Leader;"`
}
