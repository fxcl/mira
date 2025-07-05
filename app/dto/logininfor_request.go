package dto

import "mira/anima/datetime"

// Login Log List
type LogininforListRequest struct {
	PageRequest
	Ipaddr        string `query:"ipaddr" form:"ipaddr"`
	UserName      string `query:"userName" form:"userName"`
	Status        string `query:"status" form:"status"`
	BeginTime     string `query:"params[beginTime]" form:"params[beginTime]"`
	EndTime       string `query:"params[endTime]" form:"params[endTime]"`
	OrderByColumn string `query:"orderByColumn" form:"orderByColumn"`
	IsAsc         string `query:"isAsc" form:"isAsc"`
	OrderRule     string
}

// Save Login Log Information
type SaveLogininforRequest struct {
	UserName      string            `json:"userName"`
	Ipaddr        string            `json:"ipaddr"`
	LoginLocation string            `json:"loginLocation"`
	Browser       string            `json:"browser"`
	Os            string            `json:"os"`
	Status        string            `json:"status"`
	Msg           string            `json:"msg"`
	LoginTime     datetime.Datetime `json:"loginTime"`
}
