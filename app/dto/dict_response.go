package dto

import "mira/anima/datetime"

// Dictionary Type List
type DictTypeListResponse struct {
	DictId     int               `json:"dictId"`
	DictName   string            `json:"dictName"`
	DictType   string            `json:"dictType"`
	Status     string            `json:"status"`
	CreateTime datetime.Datetime `json:"createTime"`
	Remark     string            `json:"remark"`
}

// Dictionary Type Details
type DictTypeDetailResponse struct {
	DictId   int    `json:"dictId"`
	DictName string `json:"dictName"`
	DictType string `json:"dictType"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// Dictionary Data List
type DictDataListResponse struct {
	DictCode   int               `json:"dictCode"`
	DictSort   int               `json:"dictSort"`
	DictLabel  string            `json:"dictLabel"`
	DictValue  string            `json:"dictValue"`
	DictType   string            `json:"dictType"`
	CssClass   string            `json:"cssClass"`
	ListClass  string            `json:"listClass"`
	IsDefault  string            `json:"isDefault"`
	Status     string            `json:"status"`
	CreateTime datetime.Datetime `json:"createTime"`
	Default    bool              `json:"default" gorm:"-"`
}

// Dictionary Data Details
type DictDataDetailResponse struct {
	DictCode  int    `json:"dictCode"`
	DictSort  int    `json:"dictSort"`
	DictLabel string `json:"dictLabel"`
	DictValue string `json:"dictValue"`
	DictType  string `json:"dictType"`
	CssClass  string `json:"cssClass"`
	ListClass string `json:"listClass"`
	IsDefault string `json:"isDefault"`
	Status    string `json:"status"`
	Default   bool   `json:"default" gorm:"-"`
}

// Dictionary Type Export
type DictTypeExportResponse struct {
	DictId   int    `excel:"name:Dict ID;"`
	DictName string `excel:"name:Dict Name;"`
	DictType string `excel:"name:Dict Type;"`
	Status   string `excel:"name:Status;replace:0_Normal,1_Disabled;"`
}

// Dictionary Data Export
type DictDataExportResponse struct {
	DictCode  int    `excel:"name:Dict Code;"`
	DictSort  int    `excel:"name:Dict Sort;"`
	DictLabel string `excel:"name:Dict Label;"`
	DictValue string `excel:"name:Dict Value;"`
	DictType  string `excel:"name:Dict Type;"`
	IsDefault string `excel:"name:Is Default;replace:Y_Yes,N_No;"`
	Status    string `excel:"name:Status;replace:0_Normal,1_Disabled;"`
}
