package dto

import "mira/anima/datetime"

// Parameter List
type ConfigListResponse struct {
	ConfigId    int               `json:"configId"`
	ConfigName  string            `json:"configName"`
	ConfigKey   string            `json:"configKey"`
	ConfigValue string            `json:"configValue"`
	ConfigType  string            `json:"configType"`
	CreateTime  datetime.Datetime `json:"createTime"`
	Remark      string            `json:"remark"`
}

// Parameter Details
type ConfigDetailResponse struct {
	ConfigId    int    `json:"configId"`
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	Remark      string `json:"remark"`
}

// Parameter Export
type ConfigExportResponse struct {
	ConfigId    int    `excel:"name:Config ID;"`
	ConfigName  string `excel:"name:Config Name;"`
	ConfigKey   string `excel:"name:Config Key;"`
	ConfigValue string `excel:"name:Config Value;"`
	ConfigType  string `excel:"name:System Built-in;replace:Y_Yes,N_No;"`
}
