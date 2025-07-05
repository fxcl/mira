package dto

import "mira/anima/datetime"

// Operation Log List
type OperLogListResponse struct {
	OperId        int               `json:"operId"`
	Title         string            `json:"title"`
	BusinessType  int               `json:"businessType"`
	Method        string            `json:"method"`
	RequestMethod string            `json:"requestMethod"`
	OperName      string            `json:"operName"`
	DeptName      string            `json:"deptName"`
	OperUrl       string            `json:"operUrl"`
	OperIp        string            `json:"operIp"`
	OperLocation  string            `json:"operLocation"`
	OperParam     string            `json:"operParam"`
	JsonResult    string            `json:"jsonResult"`
	Status        int               `json:"status"`
	ErrorMsg      string            `json:"errorMsg"`
	OperTime      datetime.Datetime `json:"operTime"`
	CostTime      int               `json:"costTime"`
}

// Operation Log Export
type OperLogExportResponse struct {
	OperId        int    `excel:"name:Operation ID;"`
	Title         string `excel:"name:Operation Module;"`
	BusinessType  int    `excel:"name:Business Type;replace:0_Other,1_Add,2_Update,3_Delete,4_Auth,5_Export,6_Import,7_Force,8_Gen Code,9_Clean;"`
	Method        string `excel:"name:Request Method;"`
	RequestMethod string `excel:"name:Request Mode;"`
	OperName      string `excel:"name:Operator;"`
	DeptName      string `excel:"name:Dept Name;"`
	OperUrl       string `excel:"name:Request URL;"`
	OperIp        string `excel:"name:Operator IP;"`
	OperLocation  string `excel:"name:Operator Location;"`
	OperParam     string `excel:"name:Request Params;"`
	JsonResult    string `excel:"name:Return Params;"`
	Status        int    `excel:"name:Operation Status;replace:0_Normal,1_Abnormal;"`
	ErrorMsg      string `excel:"name:Error Message;"`
	OperTime      string `excel:"name:Operation Time;"`
	CostTime      string `excel:"name:Cost Time;"`
}
